-- Migration 000011 Up: Create company_csr_profiles table

CREATE TABLE IF NOT EXISTS company_csr_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID UNIQUE NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    has_csr BOOLEAN DEFAULT TRUE NOT NULL,
    csr_department_name VARCHAR(255),
    csr_email_public VARCHAR(255),
    csr_focus TEXT[] DEFAULT '{}',
    budget_range VARCHAR(100),
    proposal_acceptance VARCHAR(50) DEFAULT 'OPEN',
    website_source TEXT,
    last_verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_company_csr_profiles_company_id ON company_csr_profiles(company_id);
CREATE INDEX IF NOT EXISTS idx_company_csr_profiles_has_csr ON company_csr_profiles(has_csr);
CREATE INDEX IF NOT EXISTS idx_company_csr_profiles_proposal_acceptance ON company_csr_profiles(proposal_acceptance);
