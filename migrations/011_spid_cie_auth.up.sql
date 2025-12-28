-- Add SPID and CIE authentication tracking to users table

-- Add authentication method tracking
ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_method VARCHAR(20) DEFAULT 'email' 
    CHECK (auth_method IN ('email', 'spid', 'cie'));

-- Add SPID provider tracking (Aruba, Infocert, Lepida, Poste, Tim, Register, SielteID, Intesa)
ALTER TABLE users ADD COLUMN IF NOT EXISTS spid_provider VARCHAR(50);

-- Add CIE tracking
ALTER TABLE users ADD COLUMN IF NOT EXISTS cie_serial_number VARCHAR(50);

-- Add last authentication timestamp
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_auth_at TIMESTAMP;

-- Add external identity ID (SPID unique identifier)
ALTER TABLE users ADD COLUMN IF NOT EXISTS external_id VARCHAR(255);

-- Create index for external ID lookups
CREATE INDEX IF NOT EXISTS idx_users_external_id ON users(external_id);
CREATE INDEX IF NOT EXISTS idx_users_auth_method ON users(auth_method);

-- Create authentication audit log
CREATE TABLE IF NOT EXISTS auth_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    auth_method VARCHAR(20) NOT NULL,
    provider VARCHAR(50),
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    failure_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_audit_user_id ON auth_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_audit_created_at ON auth_audit_logs(created_at);

COMMENT ON COLUMN users.auth_method IS 'Authentication method: email, spid, or cie';
COMMENT ON COLUMN users.spid_provider IS 'SPID identity provider (e.g., Poste, Aruba, Infocert)';
COMMENT ON COLUMN users.cie_serial_number IS 'CIE card serial number for verification';
COMMENT ON COLUMN users.external_id IS 'External identity ID from SPID or CIE';
COMMENT ON TABLE auth_audit_logs IS 'Audit log for authentication attempts (GDPR compliance)';
