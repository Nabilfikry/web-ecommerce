package cart

import (
	"context"
	"errors"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetOrCreateCart(ctx context.Context, buyerID uuid.UUID) (*model.Cart, error)
	GetCartWithItems(ctx context.Context, buyerID uuid.UUID) (*model.Cart, error)
	GetCartItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) (*model.CartItem, error)
	GetCartItemByProduct(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (*model.CartItem, error)
	AddItem(ctx context.Context, item *model.CartItem) error
	UpdateItem(ctx context.Context, item *model.CartItem) error
	DeleteItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetOrCreateCart(ctx context.Context, buyerID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.WithContext(ctx).Where("buyer_id = ?", buyerID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = model.Cart{
			ID:      uuid.New(),
			BuyerID: buyerID,
		}
		if createErr := r.db.WithContext(ctx).Create(&cart).Error; createErr != nil {
			return nil, createErr
		}
		return &cart, nil
	}
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *gormRepository) GetCartWithItems(ctx context.Context, buyerID uuid.UUID) (*model.Cart, error) {
	cart, err := r.GetOrCreateCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	var fullCart model.Cart
	err = r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("cart_items.created_at ASC")
		}).
		Preload("Items.Product").
		Preload("Items.Product.Seller").
		First(&fullCart, "id = ?", cart.ID).Error
	if err != nil {
		return nil, err
	}
	return &fullCart, nil
}

func (r *gormRepository) GetCartItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("id = ? AND cart_id = ?", itemID, cartID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository) GetCartItemByProduct(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository) AddItem(ctx context.Context, item *model.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *gormRepository) UpdateItem(ctx context.Context, item *model.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *gormRepository) DeleteItem(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.CartItem{}, "id = ? AND cart_id = ?", itemID, cartID).Error
}

func (r *gormRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.CartItem{}, "cart_id = ?", cartID).Error
}

