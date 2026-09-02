DROP POLICY IF EXISTS crawling_targets_tenant_isolation ON crawling_targets;
ALTER TABLE crawling_targets DISABLE ROW LEVEL SECURITY;
DROP INDEX IF EXISTS idx_crawling_targets_tenant;
ALTER TABLE crawling_targets DROP COLUMN IF EXISTS tenant_id;
