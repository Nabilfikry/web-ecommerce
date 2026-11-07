package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ecommerce/backend-service/internal/middleware"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler, sellerMiddleware func(http.Handler) http.Handler) {
	// Public Routes
	mux.HandleFunc("GET /api/v1/products", h.handleListProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", h.handleGetProduct)

	// Seller Protected Routes
	mux.Handle("POST /api/v1/products", authMiddleware(sellerMiddleware(http.HandlerFunc(h.handleCreateProduct))))
	mux.Handle("PUT /api/v1/products/{id}", authMiddleware(sellerMiddleware(http.HandlerFunc(h.handleUpdateProduct))))
	mux.Handle("DELETE /api/v1/products/{id}", authMiddleware(sellerMiddleware(http.HandlerFunc(h.handleDeleteProduct))))
	mux.Handle("PATCH /api/v1/products/{id}/stock", authMiddleware(sellerMiddleware(http.HandlerFunc(h.handleUpdateStock))))
}

func (h *Handler) handleListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	search := q.Get("search")
	sort := q.Get("sort")

	var sellerID *uuid.UUID
	if sIDStr := q.Get("seller_id"); sIDStr != "" {
		if parsed, err := uuid.Parse(sIDStr); err == nil {
			sellerID = &parsed
		}
	}

	products, meta, err := h.service.ListProducts(r.Context(), ListProductsParams{
		Page:     page,
		Limit:    limit,
		Search:   search,
		Sort:     sort,
		SellerID: sellerID,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.SuccessWithMeta(w, products, meta)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid product ID UUID format", "BAD_REQUEST", nil)
		return
	}

	prod, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, prod, "")
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	prod, err := h.service.CreateProduct(r.Context(), authUser.ID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusCreated, prod, "Product created successfully")
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid product ID UUID format", "BAD_REQUEST", nil)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	prod, err := h.service.UpdateProduct(r.Context(), authUser.ID, id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, prod, "Product updated successfully")
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid product ID UUID format", "BAD_REQUEST", nil)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), authUser.ID, id); err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"}, "Product deleted successfully")
}

func (h *Handler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid product ID UUID format", "BAD_REQUEST", nil)
		return
	}

	var req UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	stockRes, err := h.service.UpdateStock(r.Context(), authUser.ID, id, req.Stock)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, stockRes, "Stock updated successfully")
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		shared.Error(w, appErr.StatusCode, appErr.Message, appErr.Code, nil)
		return
	}
	shared.Error(w, http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR", nil)
}

