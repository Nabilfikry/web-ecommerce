-- 000006_create_payments_and_order_payment_fields.up.sql

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS payment_status VARCHAR(30) NOT NULL DEFAULT 'unpaid',
    ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(100),
    ADD COLUMN IF NOT EXISTS payment_url TEXT;

CREATE INDEX IF NOT EXISTS idx_orders_invoice_number ON orders(invoice_number);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    invoice_number VARCHAR(100) NOT NULL,
    amount NUMERIC(14,2) NOT NULL CHECK (amount >= 0),
    payment_method VARCHAR(50),
    channel VARCHAR(50),
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',
    doku_reference_id VARCHAR(100),
    raw_payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_invoice ON payments(invoice_number);
