package cart

import (
	"context"
	"errors"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/product"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	GetCart(ctx context.Context, buyerID uuid.UUID) (*CartResponse, error)
	AddItem(ctx context.Context, buyerID uuid.UUID, req AddItemRequest) (*CartResponse, error)
	UpdateItem(ctx context.Context, buyerID uuid.UUID, itemID uuid.UUID, quantity int) (*CartResponse, error)
	RemoveItem(ctx context.Context, buyerID uuid.UUID, itemID uuid.UUID) (*CartResponse, error)
	ClearCart(ctx context.Context, buyerID uuid.UUID) (*CartResponse, error)
}

type cartService struct {
	repo        Repository
	productRepo product.Repository
}

func NewService(repo Repository, productRepo product.Repository) Service {
	return &cartService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *cartService) GetCart(ctx context.Context, buyerID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.GetCartWithItems(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	return s.toCartResponse(cart), nil
}

func (s *cartService) AddItem(ctx context.Context, buyerID uuid.UUID, req AddItemRequest) (*CartResponse, error) {
	if req.Quantity <= 0 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Quantity must be greater than 0", 400, nil)
	}

	prod, err := s.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Product not found", 404, nil)
		}
		return nil, err
	}

	if !prod.IsActive {
		return nil, shared.NewAppError("PRODUCT_INACTIVE", "Product is currently inactive", 400, nil)
	}

	cart, err := s.repo.GetOrCreateCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	existingItem, err := s.repo.GetCartItemByProduct(ctx, cart.ID, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	targetQuantity := req.Quantity
	if existingItem != nil {
		targetQuantity += existingItem.Quantity
	}

	if targetQuantity > prod.Stock {
		return nil, shared.NewAppError("INSUFFICIENT_STOCK", "Requested quantity exceeds available product stock", 400, nil)
	}

	if existingItem != nil {
		existingItem.Quantity = targetQuantity
		if err := s.repo.UpdateItem(ctx, existingItem); err != nil {
			return nil, err
		}
	} else {
		newItem := &model.CartItem{
			ID:        uuid.New(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := s.repo.AddItem(ctx, newItem); err != nil {
			return nil, err
		}
	}

	return s.GetCart(ctx, buyerID)
}

func (s *cartService) UpdateItem(ctx context.Context, buyerID uuid.UUID, itemID uuid.UUID, quantity int) (*CartResponse, error) {
	if quantity <= 0 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Quantity must be greater than 0", 400, nil)
	}

	cart, err := s.repo.GetOrCreateCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	item, err := s.repo.GetCartItem(ctx, cart.ID, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Cart item not found", 404, nil)
		}
		return nil, err
	}

	if quantity > item.Product.Stock {
		return nil, shared.NewAppError("INSUFFICIENT_STOCK", "Requested quantity exceeds available product stock", 400, nil)
	}

	item.Quantity = quantity
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, buyerID)
}

func (s *cartService) RemoveItem(ctx context.Context, buyerID uuid.UUID, itemID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.GetOrCreateCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DeleteItem(ctx, cart.ID, itemID); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, buyerID)
}

func (s *cartService) ClearCart(ctx context.Context, buyerID uuid.UUID) (*CartResponse, error) {
	cart, err := s.repo.GetOrCreateCart(ctx, buyerID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ClearCart(ctx, cart.ID); err != nil {
		return nil, err
	}

	return s.GetCart(ctx, buyerID)
}

func (s *cartService) toCartResponse(cart *model.Cart) *CartResponse {
	totalItems := 0
	totalAmount := 0.0

	itemsResp := make([]CartItemResponse, 0, len(cart.Items))
	for _, it := range cart.Items {
		subtotal := float64(it.Quantity) * it.Product.Price
		totalItems += it.Quantity
		totalAmount += subtotal

		sellerSummary := &product.SellerSummary{}
		if it.Product.Seller != nil {
			sellerSummary.ID = it.Product.Seller.ID
			sellerSummary.Name = it.Product.Seller.Name
			sellerSummary.Email = it.Product.Seller.Email
		} else {
			sellerSummary.ID = it.Product.SellerID
		}

		itemResp := CartItemResponse{
			ID:        it.ID,
			CartID:    it.CartID,
			ProductID: it.ProductID,
			Product: product.ProductResponse{
				ID:          it.Product.ID,
				SellerID:    it.Product.SellerID,
				Seller:      sellerSummary,
				Name:        it.Product.Name,
				Description: it.Product.Description,
				Price:       it.Product.Price,
				Stock:       it.Product.Stock,
				ImageURL:    it.Product.ImageURL,
				IsActive:    it.Product.IsActive,
				CreatedAt:   it.Product.CreatedAt,
				UpdatedAt:   it.Product.UpdatedAt,
			},
			Quantity:  it.Quantity,
			Subtotal:  subtotal,
			CreatedAt: it.CreatedAt,
			UpdatedAt: it.UpdatedAt,
		}
		itemsResp = append(itemsResp, itemResp)
	}

	return &CartResponse{
		ID:          cart.ID,
		BuyerID:     cart.BuyerID,
		Items:       itemsResp,
		TotalItems:  totalItems,
		TotalAmount: totalAmount,
	}
}
