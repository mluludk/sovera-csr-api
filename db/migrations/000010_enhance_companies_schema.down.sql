-- Migration 000009 Down: Revert enterprise corporate fields addition from companies

ALTER TABLE companies
DROP COLUMN IF EXISTS parent_company_id,
DROP COLUMN IF EXISTS ticker,
DROP COLUMN IF EXISTS is_public,
DROP COLUMN IF EXISTS revenue_range,
DROP COLUMN IF EXISTS employee_range,
DROP COLUMN IF EXISTS headquarters,
DROP COLUMN IF EXISTS linkedin_url,
DROP COLUMN IF EXISTS website,
DROP COLUMN IF EXISTS company_type,
DROP COLUMN IF EXISTS industry_id,
DROP COLUMN IF EXISTS legal_name;
