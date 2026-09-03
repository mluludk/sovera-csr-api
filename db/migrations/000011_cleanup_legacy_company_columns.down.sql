-- Migration 000010 Down: Restore legacy company columns (website_url & stock_code)

ALTER TABLE companies ADD COLUMN IF NOT EXISTS website_url TEXT;
ALTER TABLE companies ADD COLUMN IF NOT EXISTS stock_code VARCHAR(20);

UPDATE companies SET website_url = website WHERE website_url IS NULL AND website IS NOT NULL;
UPDATE companies SET stock_code = ticker WHERE stock_code IS NULL AND ticker IS NOT NULL;
