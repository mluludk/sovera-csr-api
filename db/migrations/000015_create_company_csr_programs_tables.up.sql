-- Migration 000014 Up: Create company_csr_programs & company_csr_program_focuses tables

CREATE TABLE IF NOT EXISTS company_csr_programs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    program_type VARCHAR(50),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50) DEFAULT 'ACTIVE',
    budget_amount NUMERIC(15,2),
    impact_summary TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS company_csr_program_focuses (
    program_id UUID NOT NULL REFERENCES company_csr_programs(id) ON DELETE CASCADE,
    focus_id UUID NOT NULL REFERENCES csr_focuses(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (program_id, focus_id)
);

CREATE INDEX IF NOT EXISTS idx_company_csr_programs_company_id ON company_csr_programs(company_id);
CREATE INDEX IF NOT EXISTS idx_company_csr_programs_status ON company_csr_programs(status);
CREATE INDEX IF NOT EXISTS idx_company_csr_program_focuses_program_id ON company_csr_program_focuses(program_id);
CREATE INDEX IF NOT EXISTS idx_company_csr_program_focuses_focus_id ON company_csr_program_focuses(focus_id);
