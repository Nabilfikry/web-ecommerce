package payment

import (
	"encoding/json"
	"errors"
	"io"
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

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
	authMiddleware func(http.Handler) http.Handler,
	buyerMiddleware func(http.Handler) http.Handler,
) {
	// Buyer initiates payment for an order
	mux.Handle("POST /api/v1/orders/{id}/pay", authMiddleware(buyerMiddleware(http.HandlerFunc(h.handlePayOrder))))

	// DOKU Webhook Notification Endpoint (Public, verified via HMAC-SHA256 signature)
	mux.HandleFunc("POST /api/v1/payments/notifications/doku", h.handleNotification)

	// Development & Demo payment simulation helper
	mux.HandleFunc("POST /api/v1/payments/simulate", h.handleSimulate)

	// Payment status inquiry
	mux.HandleFunc("GET /api/v1/payments/status/{invoice}", h.handleGetPaymentStatus)
}

func (h *Handler) handlePayOrder(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.service.PayOrder(r.Context(), authUser.ID, orderID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, resp, "DOKU payment session created successfully")
}

func (h *Handler) handleNotification(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		shared.Error(w, http.StatusBadRequest, "Failed to read request body", "BAD_REQUEST", nil)
		return
	}

	if err := h.service.HandleNotification(r.Context(), r.Header, rawBody); err != nil {
		h.handleError(w, err)
		return
	}

	// DOKU expects 200 OK with success confirmation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Notification received and processed",
	})
}

func (h *Handler) handleSimulate(w http.ResponseWriter, r *http.Request) {
	var req SimulatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.Error(w, http.StatusBadRequest, "Invalid request body JSON", "BAD_REQUEST", nil)
		return
	}

	if req.InvoiceNumber == "" {
		shared.Error(w, http.StatusBadRequest, "invoice_number is required", "BAD_REQUEST", nil)
		return
	}

	if err := h.service.SimulatePayment(r.Context(), req); err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, map[string]string{
		"invoice_number": req.InvoiceNumber,
		"status":         req.Status,
	}, "Payment simulated successfully")
}

func (h *Handler) handleGetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	invoiceNumber := r.PathValue("invoice")
	if invoiceNumber == "" {
		shared.Error(w, http.StatusBadRequest, "Invoice number is required", "BAD_REQUEST", nil)
		return
	}

	resp, err := h.service.GetPaymentStatus(r.Context(), invoiceNumber)
	if err != nil {
		h.handleError(w, err)
		return
	}

	shared.Success(w, http.StatusOK, resp, "")
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		shared.Error(w, appErr.StatusCode, appErr.Message, appErr.Code, nil)
		return
	}

	shared.Error(w, http.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR", nil)
}
