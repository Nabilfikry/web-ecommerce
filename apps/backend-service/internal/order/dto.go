package order

import (
	"time"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
)

type CheckoutRequest struct {
	ShippingAddress model.ShippingAddress `json:"shipping_address"`
	Notes           *string               `json:"notes,omitempty"`
}

type UpdateStatusRequest struct {
	Status model.OrderStatus `json:"status"`
}

type CancelOrderRequest struct {
	Reason *string `json:"reason,omitempty"`
}

type UserSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email,omitempty"`
}

type OrderItemResponse struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	UnitPrice   float64   `json:"unit_price"`
	Quantity    int       `json:"quantity"`
	Subtotal    float64   `json:"subtotal"`
}

type OrderResponse struct {
	ID              uuid.UUID             `json:"id"`
	BuyerID         uuid.UUID             `json:"buyer_id"`
	Buyer           *UserSummary          `json:"buyer,omitempty"`
	SellerID        uuid.UUID             `json:"seller_id"`
	Seller          *UserSummary          `json:"seller,omitempty"`
	TotalAmount     float64               `json:"total_amount"`
	Status          model.OrderStatus     `json:"status"`
	ShippingAddress model.ShippingAddress `json:"shipping_address"`
	Notes           *string               `json:"notes,omitempty"`
	Items           []OrderItemResponse   `json:"items"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

type ListOrdersParams struct {
	Page     int
	Limit    int
	Status   string
	BuyerID  *uuid.UUID
	SellerID *uuid.UUID
}

