package payment

import (
	"context"
	"encoding/json"
	"time"

	"ecommerce/backend-service/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	DB() *gorm.DB
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	GetOrderByInvoice(ctx context.Context, invoiceNumber string) (*model.Order, error)
	GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*model.Payment, error)
	CreatePaymentRecord(ctx context.Context, p *model.Payment, invoiceNumber, paymentURL string) error
	ProcessPaymentSuccess(ctx context.Context, invoiceNumber, paymentMethod, dokuRef string, raw json.RawMessage) error
	ProcessPaymentExpired(ctx context.Context, invoiceNumber string, raw json.RawMessage) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) DB() *gorm.DB {
	return r.db
}

func (r *repository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	var ord model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Buyer").
		Preload("Seller").
		Where("id = ?", orderID).
		First(&ord).Error
	if err != nil {
		return nil, err
	}
	return &ord, nil
}

func (r *repository) GetOrderByInvoice(ctx context.Context, invoiceNumber string) (*model.Order, error) {
	var ord model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Buyer").
		Preload("Seller").
		Where("invoice_number = ?", invoiceNumber).
		First(&ord).Error
	if err != nil {
		return nil, err
	}
	return &ord, nil
}

func (r *repository) GetPaymentByInvoice(ctx context.Context, invoiceNumber string) (*model.Payment, error) {
	var p model.Payment
	err := r.db.WithContext(ctx).
		Where("invoice_number = ?", invoiceNumber).
		Order("created_at DESC").
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) CreatePaymentRecord(ctx context.Context, p *model.Payment, invoiceNumber, paymentURL string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order with payment session info
		if err := tx.Model(&model.Order{}).Where("id = ?", p.OrderID).Updates(map[string]any{
			"invoice_number": invoiceNumber,
			"payment_url":    paymentURL,
			"payment_status": model.PaymentStatusUnpaid,
			"updated_at":     time.Now(),
		}).Error; err != nil {
			return err
		}

		// Insert payment record
		return tx.Create(p).Error
	})
}

func (r *repository) ProcessPaymentSuccess(ctx context.Context, invoiceNumber, paymentMethod, dokuRef string, raw json.RawMessage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ord model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("invoice_number = ?", invoiceNumber).
			First(&ord).Error; err != nil {
			return err
		}

		// If already paid, idempotent return
		if ord.PaymentStatus == model.PaymentStatusPaid {
			return nil
		}

		// Update order status: paid and confirmed (or pending fulfillment)
		updateFields := map[string]any{
			"payment_status": model.PaymentStatusPaid,
			"updated_at":     time.Now(),
		}
		if ord.Status == model.StatusPending {
			updateFields["status"] = model.StatusConfirmed
		}

		if err := tx.Model(&model.Order{}).Where("id = ?", ord.ID).Updates(updateFields).Error; err != nil {
			return err
		}

		// Update payment record
		return tx.Model(&model.Payment{}).Where("invoice_number = ?", invoiceNumber).Updates(map[string]any{
			"status":            "SUCCESS",
			"payment_method":    paymentMethod,
			"doku_reference_id": dokuRef,
			"raw_payload":       raw,
			"updated_at":        time.Now(),
		}).Error
	})
}

func (r *repository) ProcessPaymentExpired(ctx context.Context, invoiceNumber string, raw json.RawMessage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ord model.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			Where("invoice_number = ?", invoiceNumber).
			First(&ord).Error; err != nil {
			return err
		}

		// If already paid, do not expire
		if ord.PaymentStatus == model.PaymentStatusPaid {
			return nil
		}

		// Mark order as cancelled and expired
		if err := tx.Model(&model.Order{}).Where("id = ?", ord.ID).Updates(map[string]any{
			"payment_status": model.PaymentStatusExpired,
			"status":         model.StatusCancelled,
			"updated_at":     time.Now(),
		}).Error; err != nil {
			return err
		}

		// Restore product stocks
		for _, item := range ord.Items {
			if err := tx.Model(&model.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Update payment record
		return tx.Model(&model.Payment{}).Where("invoice_number = ?", invoiceNumber).Updates(map[string]any{
			"status":      "EXPIRED",
			"raw_payload": raw,
			"updated_at":  time.Now(),
		}).Error
	})
}

