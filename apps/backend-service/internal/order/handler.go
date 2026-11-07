package order

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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.Handler) http.Handler,
	buyerMiddleware func(http.Handler) http.Handler,
	sellerMiddleware func(http.Handler) http.Handler,
) {
	// Checkout (Buyer only)
	mux.Handle("POST /api/v1/orders", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleCheckout))))

	// List & Detail (Both Buyer and Seller, scoped by role in handler)
	mux.Handle("GET /api/v1/orders", authMiddleware(http.HandlerFunc(h.handleListOrders)))
	mux.Handle("GET /api/v1/orders/{id}", authMiddleware(http.HandlerFunc(h.handleGetOrderDetail)))

	// Cancel & Complete (Buyer only)
	mux.Handle("PATCH /api/v1/orders/{id}/cancel", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleCancelOrder))))
	mux.Handle("PATCH /api/v1/orders/{id}/complete", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handleCompleteOrder))))

	// Status Update (Seller only)
	mux.Handle("PATCH /api/v1/orders/{id}/status", authMiddleware(sellerMiddleware(http.HandlerFunc(h.handleUpdateStatus))))
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	res, err := h.service.Checkout(r.Context(), authUser.ID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusCreated, res, "Order created successfully")
}

func (h *Handler) handleListOrders(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	status := q.Get("status")

	params := ListOrdersParams{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	if authUser.Role == "seller" {
		params.SellerID = &authUser.ID
	} else {
		params.BuyerID = &authUser.ID
	}

	orders, meta, err := h.service.ListOrders(r.Context(), params)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.SuccessWithMeta(w, orders, meta)
}

func (h *Handler) handleGetOrderDetail(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid order ID UUID format", "BAD_REQUEST", nil)
		return
	}

	ord, err := h.service.GetOrderByID(r.Context(), authUser.ID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, ord, "")
}

func (h *Handler) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid order ID UUID format", "BAD_REQUEST", nil)
		return
	}

	ord, err := h.service.CancelOrder(r.Context(), authUser.ID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, ord, "Order cancelled successfully. Product stock restored.")
}

func (h *Handler) handleUpdateStatus(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid order ID UUID format", "BAD_REQUEST", nil)
		return
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	ord, err := h.service.UpdateStatus(r.Context(), authUser.ID, orderID, req.Status)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, ord, "Order status updated successfully")
}

func (h *Handler) handleCompleteOrder(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.GetAuthUser(r.Context())
	if !ok || authUser == nil {
		shared.Error(w, http.StatusUnauthorized, "User context not found", "UNAUTHORIZED", nil)
		return
	}

	idStr := r.PathValue("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid order ID UUID format", "BAD_REQUEST", nil)
		return
	}

	ord, err := h.service.CompleteOrder(r.Context(), authUser.ID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, ord, "Order confirmed as completed. Thank you!")
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		shared.Error(w, appErr.StatusCode, appErr.Message, appErr.Code, nil)
		return
	}
	shared.Error(w, http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR", nil)
}
