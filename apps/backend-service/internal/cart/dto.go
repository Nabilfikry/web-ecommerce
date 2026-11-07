package cart

import (
	"time"

	"ecommerce/backend-service/internal/product"
	"github.com/google/uuid"
)

type AddItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity"`
}

type CartItemResponse struct {
	ID        uuid.UUID               `json:"id"`
	CartID    uuid.UUID               `json:"cart_id"`
	ProductID uuid.UUID               `json:"product_id"`
	Product   product.ProductResponse `json:"product"`
	Quantity  int                     `json:"quantity"`
	Subtotal  float64                 `json:"subtotal"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

type CartResponse struct {
	ID          uuid.UUID          `json:"id"`
	BuyerID     uuid.UUID          `json:"buyer_id"`
	Items       []CartItemResponse `json:"items"`
	TotalItems  int                `json:"total_items"`
	TotalAmount float64            `json:"total_amount"`
}
