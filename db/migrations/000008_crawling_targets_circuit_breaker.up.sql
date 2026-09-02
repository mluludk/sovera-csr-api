-- Migration 000008: Add Circuit Breaker and Fault Tolerance columns to crawling_targets

ALTER TABLE crawling_targets
    ADD COLUMN IF NOT EXISTS consecutive_failures INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS last_http_status INT,
    ADD COLUMN IF NOT EXISTS last_error_message TEXT,
    ADD COLUMN IF NOT EXISTS health_status VARCHAR(50) NOT NULL DEFAULT 'HEALTHY';

-- Index for filtering healthy / degraded / disabled targets
CREATE INDEX IF NOT EXISTS idx_crawling_targets_health ON crawling_targets(health_status);
