-- Migration 000010 Up: Cleanup legacy company columns (website_url & stock_code)

-- 1. Ensure all data in website and ticker are synced from legacy columns
UPDATE companies SET website = website_url WHERE website IS NULL AND website_url IS NOT NULL;
UPDATE companies SET ticker = stock_code WHERE ticker IS NULL AND stock_code IS NOT NULL;

-- 2. Drop legacy columns
ALTER TABLE companies DROP COLUMN IF EXISTS website_url;
ALTER TABLE companies DROP COLUMN IF EXISTS stock_code;
