package payment

import (
	"encoding/json"
	"time"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
)

type PayOrderResponse struct {
	OrderID       uuid.UUID           `json:"order_id"`
	InvoiceNumber string              `json:"invoice_number"`
	Amount        float64             `json:"amount"`
	PaymentURL    string              `json:"payment_url"`
	PaymentStatus model.PaymentStatus `json:"payment_status"`
}

type SimulatePaymentRequest struct {
	InvoiceNumber string `json:"invoice_number"`
	Status        string `json:"status"` // SUCCESS, EXPIRED, FAILED
}

type PaymentStatusResponse struct {
	ID              uuid.UUID `json:"id"`
	OrderID         uuid.UUID `json:"order_id"`
	InvoiceNumber   string    `json:"invoice_number"`
	Amount          float64   `json:"amount"`
	PaymentMethod   string    `json:"payment_method,omitempty"`
	Channel         string    `json:"channel,omitempty"`
	Status          string    `json:"status"`
	DokuReferenceID string    `json:"doku_reference_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// DOKU Checkout v1 Request Payload
type DokuLineItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type DokuOrderPayload struct {
	InvoiceNumber     string         `json:"invoice_number"`
	Amount            float64        `json:"amount"`
	CallbackURL       string         `json:"callback_url"`
	CallbackURLCancel string         `json:"callback_url_cancel"`
	LineItems         []DokuLineItem `json:"line_items"`
}

type DokuPaymentConfig struct {
	PaymentDueDate int `json:"payment_due_date"` // in minutes
}

type DokuCustomerPayload struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

type DokuCheckoutRequest struct {
	Order    DokuOrderPayload    `json:"order"`
	Payment  DokuPaymentConfig   `json:"payment"`
	Customer DokuCustomerPayload `json:"customer"`
}

// DOKU Checkout v1 Response Payload
type DokuCheckoutResponse struct {
	Response struct {
		Order struct {
			InvoiceNumber string `json:"invoice_number"`
		} `json:"order"`
		Payment struct {
			URL string `json:"url"`
		} `json:"payment"`
	} `json:"response"`
	Message []string `json:"message,omitempty"`
	Error   any      `json:"error,omitempty"`
}

// DOKU Checkout v1 Webhook Notification Payload
type DokuNotificationPayload struct {
	Order struct {
		InvoiceNumber string  `json:"invoice_number"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency,omitempty"`
	} `json:"order"`
	Transaction struct {
		Status        string  `json:"status"` // SUCCESS, EXPIRED, FAILED
		Date          string  `json:"date,omitempty"`
		Amount        float64 `json:"amount"`
		InvoiceNumber string  `json:"invoice_number"`
		JournalNum    string  `json:"journal_num,omitempty"`
		ApprovalCode  string  `json:"approval_code,omitempty"`
	} `json:"transaction"`
	Payment struct {
		PaymentMethodType string `json:"payment_method_type,omitempty"`
		PaidDate          string `json:"paid_date,omitempty"`
	} `json:"payment"`
	Customer struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"customer,omitempty"`
	Raw json.RawMessage `json:"-"`
}

