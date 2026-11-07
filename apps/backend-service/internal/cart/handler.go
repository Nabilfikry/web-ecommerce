package cart

import (
	"encoding/json"
	"errors"
	"net/http"

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

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler, buyerMiddleware func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/cart", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleGetCart))))
	mux.Handle("POST /api/v1/cart/items", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleAddItem))))
	mux.Handle("PUT /api/v1/cart/items/{id}", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleUpdateItem))))
	mux.Handle("DELETE /api/v1/cart/items/{id}", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleRemoveItem))))
	mux.Handle("DELETE /api/v1/cart", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleClearCart))))
}

func (h *Handler) handleGetCart(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	cart, err := h.service.GetCart(r.Context(), authUser.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, cart, "")
}

func (h *Handler) handleAddItem(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	cart, err := h.service.AddItem(r.Context(), authUser.ID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, cart, "Item added to cart")
}

func (h *Handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid cart item ID UUID format", "BAD_REQUEST", nil)
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	cart, err := h.service.UpdateItem(r.Context(), authUser.ID, itemID, req.Quantity)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, cart, "Cart item updated")
}

func (h *Handler) handleRemoveItem(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid cart item ID UUID format", "BAD_REQUEST", nil)
		return
	}

	cart, err := h.service.RemoveItem(r.Context(), authUser.ID, itemID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, cart, "Cart item removed")
}

func (h *Handler) handleClearCart(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	cart, err := h.service.ClearCart(r.Context(), authUser.ID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, cart, "Cart cleared successfully")
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		shared.Error(w, appErr.StatusCode, appErr.Message, appErr.Code, nil)
		return
	}
	shared.Error(w, http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR", nil)
}
