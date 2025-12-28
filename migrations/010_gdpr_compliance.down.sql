-- Rollback GDPR compliance migration

-- Drop PIAS
DROP TABLE IF EXISTS privacy_impact_assessments CASCADE;

-- Drop encryption functions
DROP FUNCTION IF EXISTS decrypt_sensitive_field(BYTEA, TEXT);
DROP FUNCTION IF EXISTS encrypt_sensitive_field(TEXT, TEXT);

-- Remove soft delete columns
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE marks DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE absences DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE messages DROP COLUMN IF EXISTS deleted_at;

-- Remove encrypted columns
ALTER TABLE users DROP COLUMN IF EXISTS tax_code_encrypted;
ALTER TABLE users DROP COLUMN IF EXISTS phone_encrypted;
ALTER TABLE users DROP COLUMN IF EXISTS address_encrypted;

-- Drop retention policies
DROP TABLE IF EXISTS data_retention_policies CASCADE;

-- Drop export requests
DROP TABLE IF EXISTS data_exports CASCADE;

-- Drop deletion requests
DROP TABLE IF EXISTS data_deletion_requests CASCADE;

-- Drop access logs (including partitions)
DROP TABLE IF EXISTS data_access_logs CASCADE;

-- Drop consents
DROP TABLE IF EXISTS user_consents CASCADE;

-- Drop pgcrypto extension (only if not used elsewhere)
-- DROP EXTENSION IF EXISTS pgcrypto;
