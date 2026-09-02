DROP INDEX IF EXISTS idx_institution_programs_cluster;
DROP INDEX IF EXISTS idx_organizations_org_type;

ALTER TABLE institution_programs
DROP COLUMN IF EXISTS esg_pillar,
DROP COLUMN IF EXISTS target_sdgs,
DROP COLUMN IF EXISTS primary_cluster;

ALTER TABLE organizations
DROP COLUMN IF EXISTS org_type;
