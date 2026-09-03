-- Migration 000013 Up: Create company_csr_focuses junction table

CREATE TABLE IF NOT EXISTS company_csr_focuses (
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    focus_id UUID NOT NULL REFERENCES csr_focuses(id) ON DELETE CASCADE,
    priority SMALLINT,
    confidence NUMERIC(5,2),
    source_id UUID,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (company_id, focus_id)
);

CREATE INDEX IF NOT EXISTS idx_company_csr_focuses_company_id ON company_csr_focuses(company_id);
CREATE INDEX IF NOT EXISTS idx_company_csr_focuses_focus_id ON company_csr_focuses(focus_id);
