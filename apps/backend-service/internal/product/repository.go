package product

import (
	"context"
	"strings"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	UpdateProduct(ctx context.Context, product *model.Product) error
	SoftDeleteProduct(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) error
	ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, int64, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *gormRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var prod model.Product
	err := r.db.WithContext(ctx).
		Preload("Seller").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&prod).Error
	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *gormRepository) UpdateProduct(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *gormRepository) SoftDeleteProduct(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}

func (r *gormRepository) UpdateStock(ctx context.Context, id uuid.UUID, stock int) error {
	return r.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("stock", stock).Error
}

func (r *gormRepository) ListProducts(ctx context.Context, params ListProductsParams) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("deleted_at IS NULL")

	if params.SellerID != nil {
		query = query.Where("seller_id = ?", *params.SellerID)
	} else {
		// Public catalog only shows active products
		query = query.Where("is_active = ?", true)
	}

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply Sorting
	switch params.Sort {
	case "price_asc":
		query = query.Order("price ASC, created_at DESC")
	case "price_desc":
		query = query.Order("price DESC, created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply Pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Preload("Seller").Offset(offset).Limit(params.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
