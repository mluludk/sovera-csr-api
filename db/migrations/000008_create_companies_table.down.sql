-- Migration 000008 Down: Drop companies table and remove foreign keys

ALTER TABLE public_corporate_signals DROP COLUMN IF EXISTS company_id;
ALTER TABLE crawling_targets DROP COLUMN IF EXISTS company_id;
DROP TABLE IF EXISTS companies;
