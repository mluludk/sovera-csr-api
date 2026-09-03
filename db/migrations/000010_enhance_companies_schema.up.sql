-- Migration 000009 Up: Enhance companies table schema with enterprise corporate fields

ALTER TABLE companies
ADD COLUMN IF NOT EXISTS legal_name VARCHAR(255),
ADD COLUMN IF NOT EXISTS industry_id VARCHAR(100),
ADD COLUMN IF NOT EXISTS company_type VARCHAR(50) DEFAULT 'SWASTA',
ADD COLUMN IF NOT EXISTS website TEXT,
ADD COLUMN IF NOT EXISTS linkedin_url TEXT,
ADD COLUMN IF NOT EXISTS headquarters VARCHAR(255),
ADD COLUMN IF NOT EXISTS employee_range VARCHAR(50),
ADD COLUMN IF NOT EXISTS revenue_range VARCHAR(50),
ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS ticker VARCHAR(20),
ADD COLUMN IF NOT EXISTS parent_company_id UUID REFERENCES companies(id) ON DELETE SET NULL;

-- Populate website from website_url
UPDATE companies SET website = website_url WHERE website IS NULL AND website_url IS NOT NULL;

-- Populate ticker from stock_code
UPDATE companies SET ticker = stock_code WHERE ticker IS NULL AND stock_code IS NOT NULL;

-- Populate is_public = true for public listed companies
UPDATE companies SET is_public = TRUE WHERE stock_code IS NOT NULL OR ticker IS NOT NULL;

-- Populate company_type = 'BUMN' for BUMN companies
UPDATE companies SET company_type = 'BUMN' WHERE industry_sector LIKE '%BUMN%';
UPDATE companies SET company_type = 'SWASTA_TBK' WHERE (stock_code IS NOT NULL OR ticker IS NOT NULL) AND company_type != 'BUMN';

-- Populate parent_company_id for subsidiaries
UPDATE companies SET parent_company_id = (SELECT id FROM companies WHERE slug = 'pertamina') WHERE slug IN ('pertamina-geothermal-energy', 'pertamina-ep', 'pertamina-hulu-energi', 'pertamina-patra-niaga', 'kilang-pertamina-internasional');
UPDATE companies SET parent_company_id = (SELECT id FROM companies WHERE slug = 'pln') WHERE slug IN ('pln-energi-gas');
UPDATE companies SET parent_company_id = (SELECT id FROM companies WHERE slug = 'telkom-indonesia') WHERE slug IN ('telkomsel');
UPDATE companies SET parent_company_id = (SELECT id FROM companies WHERE slug = 'astra-international') WHERE slug IN ('toyota-astra-motor', 'asuransi-astra', 'astra-agro-lestari');

CREATE INDEX IF NOT EXISTS idx_companies_company_type ON companies(company_type);
CREATE INDEX IF NOT EXISTS idx_companies_is_public ON companies(is_public);
CREATE INDEX IF NOT EXISTS idx_companies_ticker ON companies(ticker);
CREATE INDEX IF NOT EXISTS idx_companies_parent_company_id ON companies(parent_company_id);
