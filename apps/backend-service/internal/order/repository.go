package order

import (
	"context"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	DB() *gorm.DB
	GetOrderByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	ListOrders(ctx context.Context, params ListOrdersParams) ([]model.Order, int64, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) DB() *gorm.DB {
	return r.db
}

func (r *gormRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var ord model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Buyer").
		Preload("Seller").
		First(&ord, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ord, nil
}

func (r *gormRepository) ListOrders(ctx context.Context, params ListOrdersParams) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Order{})

	if params.BuyerID != nil {
		query = query.Where("buyer_id = ?", *params.BuyerID)
	}
	if params.SellerID != nil {
		query = query.Where("seller_id = ?", *params.SellerID)
	}
	if params.Status != "" && params.Status != "all" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	err := query.
		Preload("Items").
		Preload("Buyer").
		Preload("Seller").
		Order("created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

