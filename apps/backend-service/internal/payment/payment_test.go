package payment_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"ecommerce/backend-service/configs"
	"ecommerce/backend-service/internal/auth"
	"ecommerce/backend-service/internal/cart"
	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/order"
	"ecommerce/backend-service/internal/payment"
	"ecommerce/backend-service/internal/product"
	"github.com/google/uuid"
)

func TestDokuSignatureVerification(t *testing.T) {
	clientID := "MCH-TEST-123"
	secretKey := "my-secret-key-456"
	requestID := "REQ-UUID-789"
	timestamp := "2026-11-07T12:00:00Z"
	target := "/checkout/v1/payment"

	body := []byte(`{"order":{"invoice_number":"INV-001","amount":150000}}`)
	digest := payment.ComputeDigest(body)

	sig := payment.GenerateSignature(clientID, requestID, timestamp, target, digest, secretKey)
	if sig == "" {
		t.Fatal("expected non-empty signature")
	}

	// Should verify true
	if !payment.VerifySignature(clientID, requestID, timestamp, target, digest, sig, secretKey) {
		t.Fatal("expected valid signature to verify true")
	}

	// Should verify false with wrong secret
	if payment.VerifySignature(clientID, requestID, timestamp, target, digest, sig, "wrong-secret") {
		t.Fatal("expected signature with wrong secret to verify false")
	}

	// Should verify false with tampered body
	tamperedDigest := payment.ComputeDigest([]byte(`{"tampered":true}`))
	if payment.VerifySignature(clientID, requestID, timestamp, target, tamperedDigest, sig, secretKey) {
		t.Fatal("expected tampered digest to verify false")
	}
}

func setupTestDB(t *testing.T) *configs.Config {
	cfg := configs.LoadConfig()
	_, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	return cfg
}

func TestPaymentLifecycle_SuccessAndStockRestoration(t *testing.T) {
	cfg := setupTestDB(t)
	db, _ := database.InitDB(cfg.DatabaseURL)

	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret, cfg.RefreshTokenSecret)
	prodRepo := product.NewRepository(db)
	prodSvc := product.NewService(prodRepo)
	cartRepo := cart.NewRepository(db)
	cartSvc := cart.NewService(cartRepo, prodRepo)
	ordRepo := order.NewRepository(db)
	ordSvc := order.NewService(ordRepo)

	dokuClient := payment.NewDokuClient(cfg.DokuClientID, cfg.DokuSecretKey, cfg.DokuIsProduction, true)
	payRepo := payment.NewRepository(db)
	paySvc := payment.NewService(payRepo, dokuClient, cfg)

	ctx := context.Background()

	// 1. Create Seller & Product
	sellerRes, err := authSvc.Register(ctx, auth.RegisterRequest{
		Name:     "Payment Seller",
		Email:    "seller_pay_" + uuid.New().String()[:8] + "@shop.com",
		Password: "Password123",
		Role:     model.RoleSeller,
	})
	if err != nil {
		t.Fatalf("register seller failed: %v", err)
	}

	initialStock := 10
	prodRes, err := prodSvc.CreateProduct(ctx, sellerRes.ID, product.CreateProductRequest{
		Name:  "Mechanical Numpad",
		Price: 250000,
		Stock: initialStock,
	})
	if err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	// 2. Create Buyer & Order
	buyerRes, err := authSvc.Register(ctx, auth.RegisterRequest{
		Name:     "Payment Buyer",
		Email:    "buyer_pay_" + uuid.New().String()[:8] + "@shop.com",
		Password: "Password123",
		Role:     model.RoleBuyer,
	})
	if err != nil {
		t.Fatalf("register buyer failed: %v", err)
	}

	orderQty := 3
	_, err = cartSvc.AddItem(ctx, buyerRes.ID, cart.AddItemRequest{
		ProductID: prodRes.ID,
		Quantity:  orderQty,
	})
	if err != nil {
		t.Fatalf("add to cart failed: %v", err)
	}

	checkoutRes, err := ordSvc.Checkout(ctx, buyerRes.ID, order.CheckoutRequest{
		ShippingAddress: model.ShippingAddress{
			RecipientName: "Payment Buyer",
			Phone:         "0811223344",
			Address:       "Jl. Merdeka 10",
			City:          "Bandung",
			PostalCode:    "40111",
		},
	})
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}

	// Verify stock decremented
	pAfterOrder, _ := prodSvc.GetProductByID(ctx, prodRes.ID)
	if pAfterOrder.Stock != initialStock-orderQty {
		t.Fatalf("expected stock %d, got %d", initialStock-orderQty, pAfterOrder.Stock)
	}

	// 3. Initiate Payment
	payRes, err := paySvc.PayOrder(ctx, buyerRes.ID, checkoutRes.ID)
	if err != nil {
		t.Fatalf("pay order failed: %v", err)
	}
	if payRes.InvoiceNumber == "" || payRes.PaymentURL == "" {
		t.Fatalf("expected non-empty invoice number and payment URL, got: %+v", payRes)
	}

	// 4. Simulate Payment Webhook SUCCESS
	webhookPayload := payment.DokuNotificationPayload{}
	webhookPayload.Order.InvoiceNumber = payRes.InvoiceNumber
	webhookPayload.Transaction.InvoiceNumber = payRes.InvoiceNumber
	webhookPayload.Transaction.Status = "SUCCESS"
	webhookPayload.Transaction.JournalNum = "JRN-12345"
	webhookPayload.Payment.PaymentMethodType = "VIRTUAL_ACCOUNT_BCA"

	rawPayload, _ := json.Marshal(webhookPayload)
	headers := make(http.Header)
	headers.Set("Client-Id", "MCH-TEST")
	headers.Set("Request-Id", "REQ-1")
	headers.Set("Request-Timestamp", time.Now().UTC().Format(time.RFC3339))
	headers.Set("Signature", "HMACSHA256=test")

	err = paySvc.HandleNotification(ctx, headers, rawPayload)
	if err != nil {
		t.Fatalf("handle notification failed: %v", err)
	}

	// Verify order status updated to confirmed and payment_status to paid
	updatedOrder, err := ordSvc.GetOrderByID(ctx, buyerRes.ID, checkoutRes.ID)
	if err != nil {
		t.Fatalf("get order failed: %v", err)
	}
	if updatedOrder.PaymentStatus != model.PaymentStatusPaid {
		t.Fatalf("expected payment_status 'paid', got %s", updatedOrder.PaymentStatus)
	}
	if updatedOrder.Status != model.StatusConfirmed {
		t.Fatalf("expected status 'confirmed', got %s", updatedOrder.Status)
	}

	// 5. Test EXPIRED Webhook with another order to test stock restoration
	_, _ = cartSvc.AddItem(ctx, buyerRes.ID, cart.AddItemRequest{
		ProductID: prodRes.ID,
		Quantity:  2,
	})
	checkout2, err := ordSvc.Checkout(ctx, buyerRes.ID, order.CheckoutRequest{
		ShippingAddress: model.ShippingAddress{
			RecipientName: "Payment Buyer",
			Phone:         "0811223344",
			Address:       "Jl. Merdeka 10",
			City:          "Bandung",
			PostalCode:    "40111",
		},
	})
	if err != nil {
		t.Fatalf("checkout 2 failed: %v", err)
	}

	payRes2, err := paySvc.PayOrder(ctx, buyerRes.ID, checkout2.ID)
	if err != nil {
		t.Fatalf("pay order 2 failed: %v", err)
	}

	// Trigger EXPIRED via simulation
	err = paySvc.SimulatePayment(ctx, payment.SimulatePaymentRequest{
		InvoiceNumber: payRes2.InvoiceNumber,
		Status:        "EXPIRED",
	})
	if err != nil {
		t.Fatalf("simulate payment expired failed: %v", err)
	}

	// Verify order 2 is cancelled and expired
	ord2After, _ := ordSvc.GetOrderByID(ctx, buyerRes.ID, checkout2.ID)
	if ord2After.PaymentStatus != model.PaymentStatusExpired {
		t.Fatalf("expected payment_status 'expired', got %s", ord2After.PaymentStatus)
	}
	if ord2After.Status != model.StatusCancelled {
		t.Fatalf("expected status 'cancelled', got %s", ord2After.Status)
	}

	// Verify stock was restored by 2
	pAfterExpired, _ := prodSvc.GetProductByID(ctx, prodRes.ID)
	if pAfterExpired.Stock != pAfterOrder.Stock {
		t.Fatalf("expected stock to be restored to %d, got %d", pAfterOrder.Stock, pAfterExpired.Stock)
	}
}
