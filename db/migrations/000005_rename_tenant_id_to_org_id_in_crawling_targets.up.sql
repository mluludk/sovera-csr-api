-- Migration 000005: Rename tenant_id to org_id in crawling_targets for naming convention consistency

-- 1. Drop old RLS policy referencing tenant_id
DROP POLICY IF EXISTS crawling_targets_tenant_isolation ON crawling_targets;

-- 2. Rename column tenant_id → org_id
ALTER TABLE crawling_targets RENAME COLUMN tenant_id TO org_id;

-- 3. Drop old index and recreate with correct name
DROP INDEX IF EXISTS idx_crawling_targets_tenant;
CREATE INDEX IF NOT EXISTS idx_crawling_targets_org ON crawling_targets (org_id);

-- 4. Recreate RLS policy using org_id
CREATE POLICY crawling_targets_org_isolation ON crawling_targets
    FOR ALL
    USING (
        org_id IS NULL
        OR org_id = NULLIF(current_setting('app.current_tenant_id', true), '')::uuid
    );
