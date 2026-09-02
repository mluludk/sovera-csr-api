DROP POLICY IF EXISTS crawling_targets_org_isolation ON crawling_targets;
DROP INDEX IF EXISTS idx_crawling_targets_org;
ALTER TABLE crawling_targets RENAME COLUMN org_id TO tenant_id;
CREATE INDEX IF NOT EXISTS idx_crawling_targets_tenant ON crawling_targets (tenant_id);
CREATE POLICY crawling_targets_tenant_isolation ON crawling_targets
    FOR ALL
    USING (
        tenant_id IS NULL
        OR tenant_id = NULLIF(current_setting('app.current_tenant_id', true), '')::uuid
    );
