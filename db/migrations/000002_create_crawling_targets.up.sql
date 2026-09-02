-- Migration: Create crawling_targets and crawling_logs tables
CREATE TABLE IF NOT EXISTS crawling_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_name VARCHAR(255) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    target_url TEXT NOT NULL,
    check_interval_hours INT DEFAULT 24,
    last_scraped_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_crawling_targets_active_next_run 
ON crawling_targets (is_active, next_run_at) 
WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS crawling_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_id UUID REFERENCES crawling_targets(id) ON DELETE SET NULL,
    task_id VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    http_status_code INT,
    error_message TEXT,
    execution_time_ms INT,
    content_hash VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed initial target data
INSERT INTO crawling_targets (source_name, source_type, target_url, check_interval_hours, is_active)
VALUES 
  ('IDX Sustainability Report Telkom Indonesia', 'PDF_DOCUMENT', 'https://www.idx.co.id/StaticData/NewsAndAnnouncement/SR_TLKM_2025.pdf', 24, true),
  ('Mandiri CSR Program Release', 'NEWS_ARTICLE', 'https://www.bankmandiri.co.id/csr-updates', 12, true),
  ('Astra International Sustainability Overview', 'PDF_DOCUMENT', 'https://www.astra.co.id/Sustainability/Report-2025.pdf', 24, true)
ON CONFLICT DO NOTHING;
