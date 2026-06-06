ALTER TABLE webhook_deliveries
    DROP COLUMN IF EXISTS attempt,
    DROP COLUMN IF EXISTS next_retry_at,
    DROP COLUMN IF EXISTS error_msg;
