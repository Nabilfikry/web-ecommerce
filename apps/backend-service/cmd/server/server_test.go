package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ecommerce/backend-service/internal/auth"
	"ecommerce/backend-service/internal/cart"
	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/middleware"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/order"
	"ecommerce/backend-service/internal/product"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
)

func setupTestServer(t *testing.T) *httptest.Server {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	db, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	jwtSecret := "test-jwt-secret-for-e2e-testing-32-chars-long"
	rfSecret := "test-rf-secret-for-e2e-testing-32-chars-long"

	authRepo := auth.NewRepository(db)
	prodRepo := product.NewRepository(db)
	cartRepo := cart.NewRepository(db)
	orderRepo := order.NewRepository(db)

	authSvc := auth.NewService(authRepo, jwtSecret, rfSecret)
	prodSvc := product.NewService(prodRepo)
	cartSvc := cart.NewService(cartRepo, prodRepo)
	orderSvc := order.NewService(orderRepo)

	authHandler := auth.NewHandler(authSvc)
	prodHandler := product.NewHandler(prodSvc)
	cartHandler := cart.NewHandler(cartSvc)
	orderHandler := order.NewHandler(orderSvc)

	mux := http.NewServeMux()
	authMw := middleware.Auth(jwtSecret)
	buyerMw := middleware.RequireRole("buyer")
	sellerMw := middleware.RequireRole("seller")

	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Service: "backend-service",
		})
	}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/v1/health", healthHandler)

	authHandler.RegisterRoutes(mux, authMw)
	prodHandler.RegisterRoutes(mux, authMw, sellerMw)
	cartHandler.RegisterRoutes(mux, authMw, buyerMw)
	orderHandler.RegisterRoutes(mux, authMw, buyerMw, sellerMw)

	var rootHandler http.Handler = mux
	rootHandler = middleware.CORS("*")(rootHandler)
	rootHandler = middleware.Recovery(rootHandler)

	return httptest.NewServer(rootHandler)
}

func TestEndToEndEcommerceWorkflow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	client := server.Client()
	randID := uuid.New().String()[:6]

	// 1. Health check
	resp, err := client.Get(server.URL + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("health check failed: %v, status: %d", err, resp.StatusCode)
	}

	// 2. Register Seller
	sellerEmail := "seller_" + randID + "@shop.com"
	regSellerBody, _ := json.Marshal(auth.RegisterRequest{
		Name:     "Seller E2E",
		Email:    sellerEmail,
		Password: "Password123",
		Role:     model.RoleSeller,
	})
	resp, err = client.Post(server.URL+"/api/v1/auth/register", "application/json", bytes.NewReader(regSellerBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to register seller: %v, status: %d", err, resp.StatusCode)
	}

	// 3. Login Seller
	loginSellerBody, _ := json.Marshal(auth.LoginRequest{
		Email:    sellerEmail,
		Password: "Password123",
	})
	resp, err = client.Post(server.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginSellerBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to login seller: %v, status: %d", err, resp.StatusCode)
	}

	var sellerLoginResp shared.ApiResponse[auth.LoginData]
	_ = json.NewDecoder(resp.Body).Decode(&sellerLoginResp)
	sellerToken := sellerLoginResp.Data.AccessToken

	// 4. Create Product with Seller Token
	prodBody, _ := json.Marshal(product.CreateProductRequest{
		Name:        "Ergonomic Mechanical Keyboard",
		Description: "Wireless Gateron Pro switches with PBT keycaps",
		Price:       650000,
		Stock:       10,
	})
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/products", bytes.NewReader(prodBody))
	req.Header.Set("Authorization", "Bearer "+sellerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create product: %v, status: %d", err, resp.StatusCode)
	}

	var prodResp shared.ApiResponse[product.ProductResponse]
	_ = json.NewDecoder(resp.Body).Decode(&prodResp)
	productID := prodResp.Data.ID

	// 5. Register Buyer
	buyerEmail := "buyer_" + randID + "@customer.com"
	regBuyerBody, _ := json.Marshal(auth.RegisterRequest{
		Name:     "Buyer E2E",
		Email:    buyerEmail,
		Password: "Password123",
		Role:     model.RoleBuyer,
	})
	resp, err = client.Post(server.URL+"/api/v1/auth/register", "application/json", bytes.NewReader(regBuyerBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to register buyer: %v, status: %d", err, resp.StatusCode)
	}

	// 6. Login Buyer
	loginBuyerBody, _ := json.Marshal(auth.LoginRequest{
		Email:    buyerEmail,
		Password: "Password123",
	})
	resp, err = client.Post(server.URL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBuyerBody))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to login buyer: %v, status: %d", err, resp.StatusCode)
	}

	var buyerLoginResp shared.ApiResponse[auth.LoginData]
	_ = json.NewDecoder(resp.Body).Decode(&buyerLoginResp)
	buyerToken := buyerLoginResp.Data.AccessToken

	// 7. Add Product to Buyer's Cart
	addCartBody, _ := json.Marshal(cart.AddItemRequest{
		ProductID: productID,
		Quantity:  2,
	})
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/cart/items", bytes.NewReader(addCartBody))
	req.Header.Set("Authorization", "Bearer "+buyerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to add item to cart: %v, status: %d", err, resp.StatusCode)
	}

	// 8. Buyer Checkout
	checkoutBody, _ := json.Marshal(order.CheckoutRequest{
		ShippingAddress: model.ShippingAddress{
			RecipientName: "Buyer E2E",
			Phone:         "08123456789",
			Address:       "Jl. Sudirman No. 12",
			City:          "Jakarta",
			PostalCode:    "10110",
		},
	})
	req, _ = http.NewRequest("POST", server.URL+"/api/v1/orders", bytes.NewReader(checkoutBody))
	req.Header.Set("Authorization", "Bearer "+buyerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed to checkout order: %v, status: %d", err, resp.StatusCode)
	}

	var orderCreatedResp shared.ApiResponse[order.OrderResponse]
	_ = json.NewDecoder(resp.Body).Decode(&orderCreatedResp)
	orderID := orderCreatedResp.Data.ID

	if orderCreatedResp.Data.TotalAmount != 1300000 {
		t.Errorf("expected total amount 1,300,000, got %f", orderCreatedResp.Data.TotalAmount)
	}
	if orderCreatedResp.Data.Status != model.StatusPending {
		t.Errorf("expected status pending, got %s", orderCreatedResp.Data.Status)
	}

	// 9. Seller Confirms Order (pending -> confirmed)
	confirmBody, _ := json.Marshal(order.UpdateStatusRequest{
		Status: model.StatusConfirmed,
	})
	req, _ = http.NewRequest("PATCH", server.URL+"/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(confirmBody))
	req.Header.Set("Authorization", "Bearer "+sellerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to confirm order: %v, status: %d", err, resp.StatusCode)
	}

	// 10. Seller Ships Order (confirmed -> processing -> shipped)
	procBody, _ := json.Marshal(order.UpdateStatusRequest{Status: model.StatusProcessing})
	req, _ = http.NewRequest("PATCH", server.URL+"/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(procBody))
	req.Header.Set("Authorization", "Bearer "+sellerToken)
	req.Header.Set("Content-Type", "application/json")
	_, _ = client.Do(req)

	shipBody, _ := json.Marshal(order.UpdateStatusRequest{Status: model.StatusShipped})
	req, _ = http.NewRequest("PATCH", server.URL+"/api/v1/orders/"+orderID.String()+"/status", bytes.NewReader(shipBody))
	req.Header.Set("Authorization", "Bearer "+sellerToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to ship order: %v, status: %d", err, resp.StatusCode)
	}

	// 11. Buyer Completes Order (shipped -> completed)
	req, _ = http.NewRequest("PATCH", server.URL+"/api/v1/orders/"+orderID.String()+"/complete", nil)
	req.Header.Set("Authorization", "Bearer "+buyerToken)
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to complete order: %v, status: %d", err, resp.StatusCode)
	}

	var completedResp shared.ApiResponse[order.OrderResponse]
	_ = json.NewDecoder(resp.Body).Decode(&completedResp)
	if completedResp.Data.Status != model.StatusCompleted {
		t.Errorf("expected final order status completed, got %s", completedResp.Data.Status)
	}
}
