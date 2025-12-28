-- Rollback SPID and CIE authentication tracking

-- Drop audit log table
DROP TABLE IF EXISTS auth_audit_logs CASCADE;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_external_id;
DROP INDEX IF EXISTS idx_users_auth_method;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS external_id;
ALTER TABLE users DROP COLUMN IF EXISTS last_auth_at;
ALTER TABLE users DROP COLUMN IF EXISTS cie_serial_number;
ALTER TABLE users DROP COLUMN IF EXISTS spid_provider;
ALTER TABLE users DROP COLUMN IF EXISTS auth_method;
