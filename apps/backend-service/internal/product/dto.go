package product

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    *string `json:"image_url"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock"`
}

type SellerSummary struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email,omitempty"`
}

type ProductResponse struct {
	ID          uuid.UUID      `json:"id"`
	SellerID    uuid.UUID      `json:"seller_id"`
	Seller      *SellerSummary `json:"seller,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	Stock       int            `json:"stock"`
	ImageURL    *string        `json:"image_url"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type StockUpdateResponse struct {
	ID    uuid.UUID `json:"id"`
	Stock int       `json:"stock"`
}

type ListProductsParams struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	SellerID *uuid.UUID
}

