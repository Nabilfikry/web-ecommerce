package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecommerce/backend-service/configs"
	"ecommerce/backend-service/internal/model"
	"ecommerce/backend-service/internal/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	PayOrder(ctx context.Context, buyerID, orderID uuid.UUID) (*PayOrderResponse, error)
	HandleNotification(ctx context.Context, headers http.Header, rawBody []byte) error
	SimulatePayment(ctx context.Context, req SimulatePaymentRequest) error
	GetPaymentStatus(ctx context.Context, invoiceNumber string) (*PaymentStatusResponse, error)
}

type service struct {
	repo       Repository
	dokuClient DokuClient
	cfg        *configs.Config
}

func NewService(repo Repository, dokuClient DokuClient, cfg *configs.Config) Service {
	return &service{
		repo:       repo,
		dokuClient: dokuClient,
		cfg:        cfg,
	}
}

func (s *service) PayOrder(ctx context.Context, buyerID, orderID uuid.UUID) (*PayOrderResponse, error) {
	ord, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Order not found", 404, nil)
		}
		return nil, err
	}

	if ord.BuyerID != buyerID {
		return nil, shared.NewAppError("FORBIDDEN", "You are not authorized to pay for this order", 403, nil)
	}

	if ord.PaymentStatus == model.PaymentStatusPaid {
		return nil, shared.NewAppError("ALREADY_PAID", "Order has already been paid", 400, nil)
	}

	if ord.Status == model.StatusCancelled || ord.Status == model.StatusRejected {
		return nil, shared.NewAppError("INVALID_STATUS", fmt.Sprintf("Cannot pay for %s order", ord.Status), 400, nil)
	}

	// If order already has an active invoice and payment url, reuse it
	if ord.InvoiceNumber != nil && *ord.InvoiceNumber != "" && ord.PaymentURL != nil && *ord.PaymentURL != "" {
		return &PayOrderResponse{
			OrderID:       ord.ID,
			InvoiceNumber: *ord.InvoiceNumber,
			Amount:        ord.TotalAmount,
			PaymentURL:    *ord.PaymentURL,
			PaymentStatus: ord.PaymentStatus,
		}, nil
	}

	// Generate unique invoice number: INV-YYYYMMDD-ORDID(first 8 chars)
	invoiceNumber := fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102"), ord.ID.String()[:8])

	// Build Line Items for DOKU
	lineItems := make([]DokuLineItem, len(ord.Items))
	for i, item := range ord.Items {
		lineItems[i] = DokuLineItem{
			Name:     item.ProductName,
			Price:    item.UnitPrice,
			Quantity: item.Quantity,
		}
	}

	buyerName := "Pembeli SimpleShop"
	buyerEmail := "buyer@example.com"
	if ord.Buyer != nil {
		buyerName = ord.Buyer.Name
		buyerEmail = ord.Buyer.Email
	}

	callbackURL := fmt.Sprintf("%s/%s", s.cfg.DokuCallbackURL, ord.ID.String())

	dokuReq := DokuCheckoutRequest{
		Order: DokuOrderPayload{
			InvoiceNumber:     invoiceNumber,
			Amount:            ord.TotalAmount,
			CallbackURL:       callbackURL,
			CallbackURLCancel: callbackURL,
			LineItems:         lineItems,
		},
		Payment: DokuPaymentConfig{
			PaymentDueDate: 60, // 60 minutes
		},
		Customer: DokuCustomerPayload{
			ID:      buyerID.String(),
			Name:    buyerName,
			Email:   buyerEmail,
			Phone:   ord.ShippingAddress.Phone,
			Address: fmt.Sprintf("%s, %s", ord.ShippingAddress.Address, ord.ShippingAddress.City),
		},
	}

	dokuResp, err := s.dokuClient.CreateCheckoutSession(ctx, dokuReq)
	if err != nil {
		return nil, shared.NewAppError("PAYMENT_GATEWAY_ERROR", fmt.Sprintf("Failed creating DOKU payment session: %v", err), 502, err)
	}

	paymentURL := dokuResp.Response.Payment.URL

	// Create payment record and update order
	paymentRec := &model.Payment{
		ID:            uuid.New(),
		OrderID:       ord.ID,
		InvoiceNumber: invoiceNumber,
		Amount:        ord.TotalAmount,
		Status:        "PENDING",
	}

	if err := s.repo.CreatePaymentRecord(ctx, paymentRec, invoiceNumber, paymentURL); err != nil {
		return nil, fmt.Errorf("failed saving payment record: %w", err)
	}

	return &PayOrderResponse{
		OrderID:       ord.ID,
		InvoiceNumber: invoiceNumber,
		Amount:        ord.TotalAmount,
		PaymentURL:    paymentURL,
		PaymentStatus: model.PaymentStatusUnpaid,
	}, nil
}

func (s *service) HandleNotification(ctx context.Context, headers http.Header, rawBody []byte) error {
	// Verify DOKU Webhook HMAC-SHA256 signature
	if !s.dokuClient.VerifyWebhookSignature(headers, rawBody, "/api/v1/payments/notifications/doku") {
		return shared.NewAppError("INVALID_SIGNATURE", "DOKU webhook signature verification failed", 401, nil)
	}

	var payload DokuNotificationPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return shared.NewAppError("BAD_REQUEST", "Invalid webhook JSON body", 400, err)
	}

	invoiceNumber := payload.Order.InvoiceNumber
	if invoiceNumber == "" {
		invoiceNumber = payload.Transaction.InvoiceNumber
	}
	if invoiceNumber == "" {
		return shared.NewAppError("BAD_REQUEST", "Invoice number missing in notification payload", 400, nil)
	}

	status := strings.ToUpper(payload.Transaction.Status)
	paymentMethod := payload.Payment.PaymentMethodType
	dokuRef := payload.Transaction.JournalNum

	switch status {
	case "SUCCESS":
		return s.repo.ProcessPaymentSuccess(ctx, invoiceNumber, paymentMethod, dokuRef, json.RawMessage(rawBody))
	case "EXPIRED", "FAILED":
		return s.repo.ProcessPaymentExpired(ctx, invoiceNumber, json.RawMessage(rawBody))
	default:
		// Other statuses (PENDING, etc.) can simply be acknowledged
		return nil
	}
}

func (s *service) SimulatePayment(ctx context.Context, req SimulatePaymentRequest) error {
	status := strings.ToUpper(req.Status)
	if status == "" {
		status = "SUCCESS"
	}

	mockPayload := DokuNotificationPayload{}
	mockPayload.Order.InvoiceNumber = req.InvoiceNumber
	mockPayload.Transaction.InvoiceNumber = req.InvoiceNumber
	mockPayload.Transaction.Status = status
	mockPayload.Transaction.JournalNum = fmt.Sprintf("MOCK-TRX-%d", time.Now().Unix())
	mockPayload.Payment.PaymentMethodType = "MOCK_DOKU_QRIS"

	raw, _ := json.Marshal(mockPayload)

	switch status {
	case "SUCCESS":
		return s.repo.ProcessPaymentSuccess(ctx, req.InvoiceNumber, "MOCK_DOKU_QRIS", mockPayload.Transaction.JournalNum, json.RawMessage(raw))
	case "EXPIRED", "FAILED":
		return s.repo.ProcessPaymentExpired(ctx, req.InvoiceNumber, json.RawMessage(raw))
	default:
		return shared.NewAppError("BAD_REQUEST", fmt.Sprintf("Unsupported simulation status: %s", status), 400, nil)
	}
}

func (s *service) GetPaymentStatus(ctx context.Context, invoiceNumber string) (*PaymentStatusResponse, error) {
	p, err := s.repo.GetPaymentByInvoice(ctx, invoiceNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewAppError("NOT_FOUND", "Payment not found", 404, nil)
		}
		return nil, err
	}

	return &PaymentStatusResponse{
		ID:              p.ID,
		OrderID:         p.OrderID,
		InvoiceNumber:   p.InvoiceNumber,
		Amount:          p.Amount,
		PaymentMethod:   p.PaymentMethod,
		Channel:         p.Channel,
		Status:          p.Status,
		DokuReferenceID: p.DokuReferenceID,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}, nil
}

