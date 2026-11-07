package cart

import (
	"context"
	"os"
	"testing"

	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, Repository, product.Repository) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db, NewRepository(db), product.NewRepository(db)
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

func createTestProduct(t *testing.T, db *gorm.DB, sellerID uuid.UUID, stock int) *model.Product {
	prod := &model.Product{
		ID:          uuid.New(),
		SellerID:    sellerID,
		Name:        "Test Product " + uuid.New().String()[:6],
		Description: "A great test item",
		Price:       250000,
		Stock:       stock,
		IsActive:    true,
	}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}
	return prod
}

func TestCartService_Operations(t *testing.T) {
	db, repo, prodRepo := setupTestDB(t)
	svc := NewService(repo, prodRepo)
	ctx := context.Background()

	seller := createTestUser(t, db, model.RoleSeller)
	buyer := createTestUser(t, db, model.RoleBuyer)
	prod := createTestProduct(t, db, seller.ID, 10)

	// 1. Initial empty cart
	cartRes, err := svc.GetCart(ctx, buyer.ID)
	if err != nil {
		t.Fatalf("failed to get empty cart: %v", err)
	}
	if cartRes.TotalItems != 0 || cartRes.TotalAmount != 0 {
		t.Errorf("expected empty cart, got %d items", cartRes.TotalItems)
	}

	// 2. Add item to cart exceeding stock -> error
	_, err = svc.AddItem(ctx, buyer.ID, AddItemRequest{
		ProductID: prod.ID,
		Quantity:  15,
	})
	if err == nil {
		t.Errorf("expected error when adding quantity exceeding stock, got nil")
	}

	// 3. Add item to cart successfully
	cartRes, err = svc.AddItem(ctx, buyer.ID, AddItemRequest{
		ProductID: prod.ID,
		Quantity:  2,
	})
	if err != nil {
		t.Fatalf("failed to add item to cart: %v", err)
	}
	if cartRes.TotalItems != 2 {
		t.Errorf("expected 2 total items, got %d", cartRes.TotalItems)
	}
	if cartRes.TotalAmount != 500000 {
		t.Errorf("expected 500000 total amount, got %f", cartRes.TotalAmount)
	}

	itemID := cartRes.Items[0].ID

	// 4. Update item quantity in cart
	cartRes, err = svc.UpdateItem(ctx, buyer.ID, itemID, 3)
	if err != nil {
		t.Fatalf("failed to update item quantity: %v", err)
	}
	if cartRes.TotalItems != 3 {
		t.Errorf("expected 3 total items, got %d", cartRes.TotalItems)
	}

	// 5. Remove item from cart
	cartRes, err = svc.RemoveItem(ctx, buyer.ID, itemID)
	if err != nil {
		t.Fatalf("failed to remove item from cart: %v", err)
	}
	if cartRes.TotalItems != 0 {
		t.Errorf("expected 0 items after remove, got %d", cartRes.TotalItems)
	}

	// 6. Add again and Clear Cart
	_, _ = svc.AddItem(ctx, buyer.ID, AddItemRequest{
		ProductID: prod.ID,
		Quantity:  1,
	})
	cartRes, err = svc.ClearCart(ctx, buyer.ID)
	if err != nil {
		t.Fatalf("failed to clear cart: %v", err)
	}
	if cartRes.TotalItems != 0 {
		t.Errorf("expected 0 items after clear, got %d", cartRes.TotalItems)
	}
}

