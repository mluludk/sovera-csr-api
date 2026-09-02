-- Rollback migration 000008
DROP INDEX IF EXISTS idx_crawling_targets_health;

ALTER TABLE crawling_targets
    DROP COLUMN IF EXISTS consecutive_failures,
    DROP COLUMN IF EXISTS last_http_status,
    DROP COLUMN IF EXISTS last_error_message,
    DROP COLUMN IF EXISTS health_status;
