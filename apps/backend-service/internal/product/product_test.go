package product

import (
	"context"
	"os"
	"testing"

	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, Repository) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=admin password=admin dbname=ecommerce_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}
	db, err := database.InitDB(dsn)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	return db, NewRepository(db)
}

func createTestSeller(t *testing.T, db *gorm.DB) *model.User {
	seller := &model.User{
		ID:           uuid.New(),
		Name:         "Seller Test",
		Email:        "seller_" + uuid.New().String()[:8] + "@test.com",
		PasswordHash: "hashed",
		Role:         model.RoleSeller,
	}
	if err := db.Create(seller).Error; err != nil {
		t.Fatalf("failed to create test seller: %v", err)
	}
	return seller
}

func TestProductService_CRUD(t *testing.T) {
	db, repo := setupTestDB(t)
	svc := NewService(repo)
	ctx := context.Background()

	seller1 := createTestSeller(t, db)
	seller2 := createTestSeller(t, db)

	// 1. Validation error on Create
	_, err := svc.CreateProduct(ctx, seller1.ID, CreateProductRequest{
		Name:  "AB", // too short
		Price: 1000,
		Stock: 5,
	})
	if err == nil {
		t.Errorf("expected error for short product name, got nil")
	}

	_, err = svc.CreateProduct(ctx, seller1.ID, CreateProductRequest{
		Name:  "Valid Product Name",
		Price: -100, // invalid price
		Stock: 5,
	})
	if err == nil {
		t.Errorf("expected error for negative price, got nil")
	}

	// 2. Successful Create
	created, err := svc.CreateProduct(ctx, seller1.ID, CreateProductRequest{
		Name:        "Mechanical Keyboard 75%",
		Description: "Wireless mechanical keyboard with RGB",
		Price:       750000,
		Stock:       15,
	})
	if err != nil {
		t.Fatalf("failed to create valid product: %v", err)
	}

	if created.Name != "Mechanical Keyboard 75%" {
		t.Errorf("expected name %s, got %s", "Mechanical Keyboard 75%", created.Name)
	}
	if created.Stock != 15 {
		t.Errorf("expected stock 15, got %d", created.Stock)
	}

	// 3. Get By ID
	fetched, err := svc.GetProductByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to get product by ID: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, fetched.ID)
	}

	// 4. Update Product by Seller 2 (Unauthorized / Forbidden)
	newName := "Hacked Product Name"
	_, err = svc.UpdateProduct(ctx, seller2.ID, created.ID, UpdateProductRequest{
		Name: &newName,
	})
	if err == nil {
		t.Errorf("expected forbidden error for unauthorized seller update, got nil")
	}

	// 5. Update Product by Owner Seller 1
	updatedPrice := 790000.0
	updated, err := svc.UpdateProduct(ctx, seller1.ID, created.ID, UpdateProductRequest{
		Price: &updatedPrice,
	})
	if err != nil {
		t.Fatalf("failed to update product: %v", err)
	}
	if updated.Price != 790000 {
		t.Errorf("expected price 790000, got %f", updated.Price)
	}

	// 6. Update Stock
	stockRes, err := svc.UpdateStock(ctx, seller1.ID, created.ID, 50)
	if err != nil {
		t.Fatalf("failed to update stock: %v", err)
	}
	if stockRes.Stock != 50 {
		t.Errorf("expected stock 50, got %d", stockRes.Stock)
	}

	// 7. List Products
	list, meta, err := svc.ListProducts(ctx, ListProductsParams{
		Page:   1,
		Limit:  10,
		Search: "Keyboard",
	})
	if err != nil {
		t.Fatalf("failed to list products: %v", err)
	}
	if len(list) == 0 {
		t.Errorf("expected list to contain at least 1 product, got 0")
	}
	if meta.Total == 0 {
		t.Errorf("expected total > 0, got %d", meta.Total)
	}

	// 8. Delete Product (Soft delete)
	if err := svc.DeleteProduct(ctx, seller1.ID, created.ID); err != nil {
		t.Fatalf("failed to delete product: %v", err)
	}

	// 9. Fetch deleted product should return 404
	_, err = svc.GetProductByID(ctx, created.ID)
	if err == nil {
		t.Errorf("expected 404 for deleted product, got nil")
	}
}
