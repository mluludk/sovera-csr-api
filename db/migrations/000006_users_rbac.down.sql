-- Rollback migration 000006
DELETE FROM users WHERE email IN ('admin@laz.id', 'director@laz.id', 'fundraiser@laz.id');

ALTER TABLE deal_pipelines DROP COLUMN IF EXISTS pic_user_id;
DROP INDEX IF EXISTS idx_deals_pic_user;

ALTER TABLE users ALTER COLUMN role TYPE VARCHAR(50) USING role::TEXT;
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'FUNDRAISER';
ALTER TABLE users DROP COLUMN IF EXISTS is_active;

DROP TYPE IF EXISTS user_org_role;
