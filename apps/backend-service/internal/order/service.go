package order

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service interface {
	Checkout(ctx context.Context, buyerID uuid.UUID, req CheckoutRequest) (*OrderResponse, error)
	GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error)
	ListOrders(ctx context.Context, params ListOrdersParams) ([]OrderResponse, *shared.ApiMeta, error)
	CancelOrder(ctx context.Context, buyerID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error)
	UpdateStatus(ctx context.Context, sellerID uuid.UUID, orderID uuid.UUID, targetStatus model.OrderStatus) (*OrderResponse, error)
	CompleteOrder(ctx context.Context, buyerID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error)
}

type orderService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &orderService{repo: repo}
}

func (s *orderService) Checkout(ctx context.Context, buyerID uuid.UUID, req CheckoutRequest) (*OrderResponse, error) {
	// Validate Shipping Address
	addr := req.ShippingAddress
	if strings.TrimSpace(addr.RecipientName) == "" || strings.TrimSpace(addr.Phone) == "" ||
		strings.TrimSpace(addr.Address) == "" || strings.TrimSpace(addr.City) == "" ||
		strings.TrimSpace(addr.PostalCode) == "" {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Shipping address fields (recipient_name, phone, address, city, postal_code) are required", 400, nil)
	}

	orderID := uuid.New()
	var createdOrder *model.Order

	err := s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Get buyer cart
		var cart model.Cart
		if err := tx.Where("buyer_id = ?", buyerID).First(&cart).Error; err != nil {
			return shared.NewAppError("NOT_FOUND", "Cart not found for buyer", 404, err)
		}

		// 2. Lock cart items
		var cartItems []model.CartItem
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("cart_id = ?", cart.ID).
			Order("created_at ASC").
			Find(&cartItems).Error; err != nil {
			return err
		}

		if len(cartItems) == 0 {
			return shared.NewAppError("EMPTY_CART", "Shopping cart is empty", 400, nil)
		}

		var sellerID uuid.UUID
		var totalAmount float64
		orderItems := make([]model.OrderItem, 0, len(cartItems))

		// 3. Lock products and validate stock & single-seller rule (BR-010)
		for i, item := range cartItems {
			var prod model.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("id = ? AND deleted_at IS NULL", item.ProductID).
				First(&prod).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return shared.NewAppError("NOT_FOUND", fmt.Sprintf("Product %s not found or deleted", item.ProductID), 404, err)
				}
				return err
			}

			if !prod.IsActive {
				return shared.NewAppError("PRODUCT_INACTIVE", fmt.Sprintf("Product '%s' is inactive", prod.Name), 400, nil)
			}

			// BR-010: Ensure single seller per checkout
			if i == 0 {
				sellerID = prod.SellerID
			} else if prod.SellerID != sellerID {
				return shared.NewAppError("MULTIPLE_SELLERS_DETECTED", "All items in a single checkout must be from the same seller", 400, nil)
			}

			// Validate sufficient stock
			if prod.Stock < item.Quantity {
				return shared.NewAppError("INSUFFICIENT_STOCK", fmt.Sprintf("Insufficient stock for '%s' (available: %d, requested: %d)", prod.Name, prod.Stock, item.Quantity), 400, nil)
			}

			subtotal := float64(item.Quantity) * prod.Price
			totalAmount += subtotal

			orderItems = append(orderItems, model.OrderItem{
				ID:          uuid.New(),
				OrderID:     orderID,
				ProductID:   prod.ID,
				ProductName: prod.Name,
				UnitPrice:   prod.Price,
				Quantity:    item.Quantity,
				Subtotal:    subtotal,
			})

			// 4. Atomically decrement stock
			if err := tx.Model(&model.Product{}).
				Where("id = ?", prod.ID).
				Update("stock", prod.Stock-item.Quantity).Error; err != nil {
				return err
			}
		}

		// 5. Create Order
		ord := &model.Order{
			ID:              orderID,
			BuyerID:         buyerID,
			SellerID:        sellerID,
			TotalAmount:     totalAmount,
			Status:          model.StatusPending,
			ShippingAddress: addr,
			Notes:           req.Notes,
			Items:           orderItems,
		}

		if err := tx.Create(ord).Error; err != nil {
			return err
		}

		// 6. Delete all cart items
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		createdOrder = ord
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, buyerID, createdOrder.ID)
}

func (s *orderService) GetOrderByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error) {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Order not found", 404, nil)
		}
		return nil, err
	}

	// Verify authorization: only buyer or seller of this order can view
	if ord.BuyerID != userID && ord.SellerID != userID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to view this order", 403, nil)
	}

	return s.toOrderResponse(ord), nil
}

func (s *orderService) ListOrders(ctx context.Context, params ListOrdersParams) ([]OrderResponse, *shared.ApiMeta, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	orders, total, err := s.repo.ListOrders(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	meta := &shared.ApiMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	resp := make([]OrderResponse, len(orders))
	for i, o := range orders {
		resp[i] = *s.toOrderResponse(&o)
	}

	return resp, meta, nil
}

func (s *orderService) CancelOrder(ctx context.Context, buyerID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error) {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Order not found", 404, nil)
		}
		return nil, err
	}

	if ord.BuyerID != buyerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to cancel this order", 403, nil)
	}

	if ord.Status != model.StatusPending {
		return nil, shared.NewAppError("INVALID_STATUS", fmt.Sprintf("Cannot cancel order with status '%s'. Only 'pending' orders can be cancelled", ord.Status), 400, nil)
	}

	// In ACID transaction: update order status to cancelled and restore product stock
	err = s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).Where("id = ?", ord.ID).Update("status", model.StatusCancelled).Error; err != nil {
			return err
		}

		for _, item := range ord.Items {
			if err := tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, buyerID, orderID)
}

func (s *orderService) UpdateStatus(ctx context.Context, sellerID uuid.UUID, orderID uuid.UUID, targetStatus model.OrderStatus) (*OrderResponse, error) {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Order not found", 404, nil)
		}
		return nil, err
	}

	if ord.SellerID != sellerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to update this order's status", 403, nil)
	}

	// Validate allowed state transitions
	isValidTransition := false
	switch ord.Status {
	case model.StatusPending:
		isValidTransition = (targetStatus == model.StatusConfirmed || targetStatus == model.StatusRejected)
	case model.StatusConfirmed:
		isValidTransition = (targetStatus == model.StatusProcessing)
	case model.StatusProcessing:
		isValidTransition = (targetStatus == model.StatusShipped)
	}

	if !isValidTransition {
		return nil, shared.NewAppError("INVALID_STATUS_TRANSITION", fmt.Sprintf("Cannot transition order from '%s' to '%s'", ord.Status, targetStatus), 400, nil)
	}

	err = s.repo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Order{}).Where("id = ?", ord.ID).Update("status", targetStatus).Error; err != nil {
			return err
		}

		// If rejected, automatically restore stock
		if targetStatus == model.StatusRejected {
			for _, item := range ord.Items {
				if err := tx.Model(&model.Product{}).
					Where("id = ?", item.ProductID).
					Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, sellerID, orderID)
}

func (s *orderService) CompleteOrder(ctx context.Context, buyerID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error) {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Order not found", 404, nil)
		}
		return nil, err
	}

	if ord.BuyerID != buyerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to complete this order", 403, nil)
	}

	if ord.Status != model.StatusShipped {
		return nil, shared.NewAppError("INVALID_STATUS", fmt.Sprintf("Cannot complete order with status '%s'. Status must be 'shipped'", ord.Status), 400, nil)
	}

	if err := s.repo.DB().WithContext(ctx).Model(&model.Order{}).Where("id = ?", ord.ID).Update("status", model.StatusCompleted).Error; err != nil {
		return nil, err
	}

	return s.GetOrderByID(ctx, buyerID, orderID)
}

func (s *orderService) toOrderResponse(ord *model.Order) *OrderResponse {
	resp := &OrderResponse{
		ID:              ord.ID,
		BuyerID:         ord.BuyerID,
		SellerID:        ord.SellerID,
		TotalAmount:     ord.TotalAmount,
		Status:          ord.Status,
		ShippingAddress: ord.ShippingAddress,
		Notes:           ord.Notes,
		CreatedAt:       ord.CreatedAt,
		UpdatedAt:       ord.UpdatedAt,
	}

	if ord.Buyer != nil {
		resp.Buyer = &UserSummary{
			ID:    ord.Buyer.ID,
			Name:  ord.Buyer.Name,
			Email: ord.Buyer.Email,
		}
	}

	if ord.Seller != nil {
		resp.Seller = &UserSummary{
			ID:    ord.Seller.ID,
			Name:  ord.Seller.Name,
			Email: ord.Seller.Email,
		}
	}

	items := make([]OrderItemResponse, len(ord.Items))
	for i, it := range ord.Items {
		items[i] = OrderItemResponse{
			ID:          it.ID,
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			UnitPrice:   it.UnitPrice,
			Quantity:    it.Quantity,
			Subtotal:    it.Subtotal,
		}
	}
	resp.Items = items

	return resp
}
