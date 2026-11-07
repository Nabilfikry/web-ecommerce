package order

import (
	"context"
	"os"
	"testing"

	"ecommerce/backend-service/internal/cart"
	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, Repository, cart.Repository, product.Repository) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db, NewRepository(db), cart.NewRepository(db), product.NewRepository(db)
}

func createTestUser(t *testing.T, db *gorm.DB, role model.UserRole) *model.User {
	user := &model.User{
		ID:           uuid.New(),
		Name:         "User " + string(role),
		Email:        "user_" + uuid.New().String()[:8] + "@test.com",
		PasswordHash: "hashed",
		Role:         role,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func createTestProduct(t *testing.T, db *gorm.DB, sellerID uuid.UUID, price float64, stock int) *model.Product {
	prod := &model.Product{
		ID:          uuid.New(),
		SellerID:    sellerID,
		Name:        "Product " + uuid.New().String()[:6],
		Description: "High quality test item",
		Price:       price,
		Stock:       stock,
		IsActive:    true,
	}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}
	return prod
}

func TestOrderService_CheckoutAndLifecycle(t *testing.T) {
	db, repo, cartRepo, prodRepo := setupTestDB(t)
	cartSvc := cart.NewService(cartRepo, prodRepo)
	svc := NewService(repo)
	ctx := context.Background()

	seller1 := createTestUser(t, db, model.RoleSeller)
	seller2 := createTestUser(t, db, model.RoleSeller)
	buyer := createTestUser(t, db, model.RoleBuyer)

	prodA := createTestProduct(t, db, seller1.ID, 300000, 10)
	prodB := createTestProduct(t, db, seller1.ID, 150000, 5)
	prodOtherSeller := createTestProduct(t, db, seller2.ID, 50000, 20)

	// 1. Multiple sellers in cart -> checkout should fail (BR-010)
	_, _ = cartSvc.AddItem(ctx, buyer.ID, cart.AddItemRequest{
		ProductID: prodA.ID,
		Quantity:  1,
	})
	_, _ = cartSvc.AddItem(ctx, buyer.ID, cart.AddItemRequest{
		ProductID: prodOtherSeller.ID,
		Quantity:  1,
	})

	shippingAddr := model.ShippingAddress{
		RecipientName: "Budi",
		Phone:         "0812345678",
		Address:       "Jl. Dago No. 10",
		City:          "Bandung",
		PostalCode:    "40135",
	}

	_, err := svc.Checkout(ctx, buyer.ID, CheckoutRequest{
		ShippingAddress: shippingAddr,
	})
	if err == nil {
		t.Errorf("expected MULTIPLE_SELLERS_DETECTED error for multi-seller checkout, got nil")
	}

	// Clear cart and setup single-seller order
	_, _ = cartSvc.ClearCart(ctx, buyer.ID)
	_, _ = cartSvc.AddItem(ctx, buyer.ID, cart.AddItemRequest{
		ProductID: prodA.ID,
		Quantity:  2, // 2 * 300,000 = 600,000
	})
	_, _ = cartSvc.AddItem(ctx, buyer.ID, cart.AddItemRequest{
		ProductID: prodB.ID,
		Quantity:  1, // 1 * 150,000 = 150,000
	})

	// 2. Successful Checkout
	notes := "Mohon packing rapi"
	orderRes, err := svc.Checkout(ctx, buyer.ID, CheckoutRequest{
		ShippingAddress: shippingAddr,
		Notes:           &notes,
	})
	if err != nil {
		t.Fatalf("failed to checkout: %v", err)
	}

	if orderRes.Status != model.StatusPending {
		t.Errorf("expected initial status pending, got %s", orderRes.Status)
	}
	if orderRes.TotalAmount != 750000 {
		t.Errorf("expected total amount 750000, got %f", orderRes.TotalAmount)
	}
	if len(orderRes.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(orderRes.Items))
	}

	// Verify stock was decremented
	var updatedProdA model.Product
	_ = db.First(&updatedProdA, "id = ?", prodA.ID)
	if updatedProdA.Stock != 8 {
		t.Errorf("expected prodA stock to be 8, got %d", updatedProdA.Stock)
	}

	// Verify cart is now empty
	cartAfter, _ := cartSvc.GetCart(ctx, buyer.ID)
	if cartAfter.TotalItems != 0 {
		t.Errorf("expected cart to be empty after checkout, got %d items", cartAfter.TotalItems)
	}

	// 3. Status Transitions by Seller
	// Pending -> Confirmed
	confirmed, err := svc.UpdateStatus(ctx, seller1.ID, orderRes.ID, model.StatusConfirmed)
	if err != nil {
		t.Fatalf("failed to confirm order: %v", err)
	}
	if confirmed.Status != model.StatusConfirmed {
		t.Errorf("expected confirmed status, got %s", confirmed.Status)
	}

	// Confirmed -> Processing
	processing, err := svc.UpdateStatus(ctx, seller1.ID, orderRes.ID, model.StatusProcessing)
	if err != nil {
		t.Fatalf("failed to process order: %v", err)
	}
	if processing.Status != model.StatusProcessing {
		t.Errorf("expected processing status, got %s", processing.Status)
	}

	// Processing -> Shipped
	shipped, err := svc.UpdateStatus(ctx, seller1.ID, orderRes.ID, model.StatusShipped)
	if err != nil {
		t.Fatalf("failed to ship order: %v", err)
	}
	if shipped.Status != model.StatusShipped {
		t.Errorf("expected shipped status, got %s", shipped.Status)
	}

	// Shipped -> Completed by Buyer
	completed, err := svc.CompleteOrder(ctx, buyer.ID, orderRes.ID)
	if err != nil {
		t.Fatalf("failed to complete order: %v", err)
	}
	if completed.Status != model.StatusCompleted {
		t.Errorf("expected completed status, got %s", completed.Status)
	}
}

func TestOrderService_CancelAndStockRestoration(t *testing.T) {
	db, repo, cartRepo, prodRepo := setupTestDB(t)
	cartSvc := cart.NewService(cartRepo, prodRepo)
	svc := NewService(repo)
	ctx := context.Background()

	seller := createTestUser(t, db, model.RoleSeller)
	buyer := createTestUser(t, db, model.RoleBuyer)
	prod := createTestProduct(t, db, seller.ID, 100000, 5)

	_, _ = cartSvc.AddItem(ctx, buyer.ID, cart.AddItemRequest{
		ProductID: prod.ID,
		Quantity:  3,
	})

	orderRes, err := svc.Checkout(ctx, buyer.ID, CheckoutRequest{
		ShippingAddress: model.ShippingAddress{
			RecipientName: "Buyer",
			Phone:         "08123",
			Address:       "Address",
			City:          "City",
			PostalCode:    "12345",
		},
	})
	if err != nil {
		t.Fatalf("failed checkout: %v", err)
	}

	// Cancel Order by Buyer
	cancelled, err := svc.CancelOrder(ctx, buyer.ID, orderRes.ID)
	if err != nil {
		t.Fatalf("failed to cancel order: %v", err)
	}
	if cancelled.Status != model.StatusCancelled {
		t.Errorf("expected status cancelled, got %s", cancelled.Status)
	}

	// Verify product stock was restored back to 5
	var restoredProd model.Product
	_ = db.First(&restoredProd, "id = ?", prod.ID)
	if restoredProd.Stock != 5 {
		t.Errorf("expected restored stock 5, got %d", restoredProd.Stock)
	}
}

