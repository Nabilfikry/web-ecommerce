-- 000006_create_payments_and_order_payment_fields.down.sql

DROP TABLE IF EXISTS payments;

ALTER TABLE orders
    DROP COLUMN IF EXISTS payment_status,
    DROP COLUMN IF EXISTS invoice_number,
    DROP COLUMN IF EXISTS payment_url;

