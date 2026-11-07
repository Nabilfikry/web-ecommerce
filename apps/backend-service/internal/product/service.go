package product

import (
	"context"
	"errors"
	"math"
	"strings"

	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	CreateProduct(ctx context.Context, sellerID uuid.UUID, req CreateProductRequest) (*ProductResponse, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*ProductResponse, error)
	UpdateProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID) error
	UpdateStock(ctx context.Context, sellerID uuid.UUID, id uuid.UUID, stock int) (*StockUpdateResponse, error)
	ListProducts(ctx context.Context, params ListProductsParams) ([]ProductResponse, *shared.ApiMeta, error)
}

type productService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, sellerID uuid.UUID, req CreateProductRequest) (*ProductResponse, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 || len(name) > 150 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Product name must be between 3 and 150 characters", 400, nil)
	}

	if req.Price <= 0 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Product price must be greater than 0", 400, nil)
	}

	if req.Stock < 0 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Product stock cannot be negative", 400, nil)
	}

	prod := &model.Product{
		ID:          uuid.New(),
		SellerID:    sellerID,
		Name:        name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if err := s.repo.CreateProduct(ctx, prod); err != nil {
		return nil, err
	}

	return s.toProductResponse(prod), nil
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*ProductResponse, error) {
	prod, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Product not found", 404, nil)
		}
		return nil, err
	}
	return s.toProductResponse(prod), nil
}

func (s *productService) UpdateProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	prod, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Product not found", 404, nil)
		}
		return nil, err
	}

	if prod.SellerID != sellerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to update this product", 403, nil)
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 || len(name) > 150 {
			return nil, shared.NewAppError("VALIDATION_ERROR", "Product name must be between 3 and 150 characters", 400, nil)
		}
		prod.Name = name
	}

	if req.Description != nil {
		prod.Description = *req.Description
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, shared.NewAppError("VALIDATION_ERROR", "Product price must be greater than 0", 400, nil)
		}
		prod.Price = *req.Price
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, shared.NewAppError("VALIDATION_ERROR", "Product stock cannot be negative", 400, nil)
		}
		prod.Stock = *req.Stock
	}

	if req.ImageURL != nil {
		prod.ImageURL = req.ImageURL
	}

	if req.IsActive != nil {
		prod.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateProduct(ctx, prod); err != nil {
		return nil, err
	}

	return s.toProductResponse(prod), nil
}

func (s *productService) DeleteProduct(ctx context.Context, sellerID uuid.UUID, id uuid.UUID) error {
	prod, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.NewAppError("NOT_FOUND", "Product not found", 404, nil)
		}
		return err
	}

	if prod.SellerID != sellerID {
		return shared.NewAppError("FORBIDDEN", "You are not authorized to delete this product", 403, nil)
	}

	return s.repo.SoftDeleteProduct(ctx, id)
}

func (s *productService) UpdateStock(ctx context.Context, sellerID uuid.UUID, id uuid.UUID, stock int) (*StockUpdateResponse, error) {
	if stock < 0 {
		return nil, shared.NewAppError("VALIDATION_ERROR", "Product stock cannot be negative", 400, nil)
	}

	prod, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Product not found", 404, nil)
		}
		return nil, err
	}

	if prod.SellerID != sellerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to update this product's stock", 403, nil)
	}

	if err := s.repo.UpdateStock(ctx, id, stock); err != nil {
		return nil, err
	}

	return &StockUpdateResponse{
		ID:    id,
		Stock: stock,
	}, nil
}

func (s *productService) ListProducts(ctx context.Context, params ListProductsParams) ([]ProductResponse, *shared.ApiMeta, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 12
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	products, total, err := s.repo.ListProducts(ctx, params)
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

	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = *s.toProductResponse(&p)
	}

	return resp, meta, nil
}

func (s *productService) toProductResponse(p *model.Product) *ProductResponse {
	resp := &ProductResponse{
		ID:          p.ID,
		SellerID:    p.SellerID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageURL:    p.ImageURL,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}

	if p.Seller != nil {
		resp.Seller = &SellerSummary{
			ID:    p.Seller.ID,
			Name:  p.Seller.Name,
			Email: p.Seller.Email,
		}
	}

	return resp
}

