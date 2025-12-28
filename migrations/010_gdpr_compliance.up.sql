-- Enable pgcrypto extension for encryption
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- GDPR Consent Management
CREATE TABLE user_consents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    consent_type VARCHAR(50) NOT NULL CHECK (consent_type IN ('COMMUNICATIONS', 'PHOTOS', 'BIOMETRIC', 'DATA_PROCESSING')),
    given_at TIMESTAMP DEFAULT NOW(),
    granted BOOLEAN DEFAULT FALSE,
    ip_address INET,
    user_agent TEXT,
    revoked_at TIMESTAMP,
    revoked_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, consent_type)
);

CREATE INDEX idx_user_consents_user_id ON user_consents(user_id);
CREATE INDEX idx_user_consents_type ON user_consents(consent_type);

-- GDPR Data Access Audit Log
CREATE TABLE data_access_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accessed_by_user_id UUID NOT NULL REFERENCES users(id),
    accessed_user_id UUID REFERENCES users(id),
    resource_type VARCHAR(50) NOT NULL CHECK (resource_type IN ('MARKS', 'ABSENCES', 'MESSAGES', 'DOCUMENTS', 'PROFILE', 'COLLOQUIUMS')),
    resource_id UUID,
    accessed_at TIMESTAMP DEFAULT NOW(),
    ip_address INET,
    user_agent TEXT,
    purpose VARCHAR(200),
    action VARCHAR(20) CHECK (action IN ('READ', 'WRITE', 'DELETE', 'EXPORT'))
);

CREATE INDEX idx_data_access_logs_accessed_by ON data_access_logs(accessed_by_user_id);
CREATE INDEX idx_data_access_logs_accessed_user ON data_access_logs(accessed_user_id);
CREATE INDEX idx_data_access_logs_accessed_at ON data_access_logs(accessed_at);

-- Partition by month for performance
CREATE TABLE data_access_logs_2024_01 PARTITION OF data_access_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- GDPR Data Deletion Requests  
CREATE TABLE data_deletion_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    requested_at TIMESTAMP DEFAULT NOW(),
    scheduled_deletion_at TIMESTAMP,
    reason TEXT,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'COMPLETED', 'CANCELLED')),
    processed_by_user_id UUID REFERENCES users(id),
    completed_at TIMESTAMP,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_data_deletion_requests_user_id ON data_deletion_requests(user_id);
CREATE INDEX idx_data_deletion_requests_status ON data_deletion_requests(status);

-- GDPR Data Export Requests
CREATE TABLE data_exports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    requested_at TIMESTAMP DEFAULT NOW(),
    format VARCHAR(10) DEFAULT 'JSON' CHECK (format IN ('JSON', 'CSV', 'XML')),
    file_path TEXT,
    file_size BIGINT,
    expiry_date TIMESTAMP DEFAULT NOW() + INTERVAL '30 days',
    downloaded_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PROCESSING', 'READY', 'EXPIRED', 'ERROR')),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_data_exports_user_id ON data_exports(user_id);
CREATE INDEX idx_data_exports_status ON data_exports(status);

-- Data Retention Policy Tracking
CREATE TABLE data_retention_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data_type VARCHAR(50) NOT NULL UNIQUE,
    retention_days INTEGER NOT NULL,
    description TEXT,
    legal_basis TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default retention policies (Italian school law)
INSERT INTO data_retention_policies (data_type, retention_days, description, legal_basis) VALUES
('MARKS', 1825, 'Student marks - 5 years', 'Italian School Law'),
('ABSENCES', 1825, 'Student absences - 5 years', 'Italian School Law'),
('AUDIT_LOGS', 365, 'Security audit logs - 1 year', 'GDPR Art. 32'),
('MESSAGES', 730, 'Internal messaging - 2 years', 'School Policy'),
('TEMP_FILES', 7, 'Temporary files - 7 days', 'Data Minimization'),
('DOCUMENTS', 3650, 'Official documents - 10 years', 'Italian Archive Law');

-- Add encryption for sensitive fields (using pgcrypto)
-- This migration adds encrypted columns for existing sensitive data

-- Add encrypted columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS tax_code_encrypted BYTEA;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_encrypted BYTEA;
ALTER TABLE users ADD COLUMN IF NOT EXISTS address_encrypted BYTEA;

-- Add deleted_at for soft delete
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE marks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE absences ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
ALTER TABLE messages ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

-- Create function for automatic encryption
CREATE OR REPLACE FUNCTION encrypt_sensitive_field(data TEXT, key TEXT)
RETURNS BYTEA AS $$
BEGIN
    RETURN pgp_sym_encrypt(data, key);
END;
$$ LANGUAGE plpgsql;

-- Create function for automatic decryption
CREATE OR REPLACE FUNCTION decrypt_sensitive_field(data BYTEA, key TEXT)
RETURNS TEXT AS $$
BEGIN
    RETURN pgp_sym_decrypt(data, key);
END;
$$ LANGUAGE plpgsql;

-- Privacy impact assessment log
CREATE TABLE privacy_impact_assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feature_name VARCHAR(100) NOT NULL,
    assessed_by_user_id UUID REFERENCES users(id),
    assessment_date TIMESTAMP DEFAULT NOW(),
    risk_level VARCHAR(20) CHECK (risk_level IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    mitigation_measures TEXT,
    approved BOOLEAN DEFAULT FALSE,
    approved_by_user_id UUID REFERENCES users(id),
    approved_at TIMESTAMP
);

COMMENT ON TABLE user_consents IS 'GDPR Art. 7 - Conditions for consent';
COMMENT ON TABLE data_access_logs IS 'GDPR Art. 30 - Records of processing activities';
COMMENT ON TABLE data_deletion_requests IS 'GDPR Art. 17 - Right to erasure';
COMMENT ON TABLE data_exports IS 'GDPR Art. 20 - Right to data portability';
