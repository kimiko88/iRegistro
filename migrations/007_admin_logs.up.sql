-- 1. Audit Logs (Partitioned by Month)
CREATE TABLE audit_logs (
    id BIGSERIAL, -- BIGSERIAL for logs
    school_id INTEGER NOT NULL,
    user_id INTEGER, -- Nullable for system actions
    action VARCHAR(100) NOT NULL, -- e.g. "USER.LOGIN", "MARK.CREATE"
    entity VARCHAR(100), -- "users", "marks"
    entity_id VARCHAR(50), 
    
    changes JSONB, -- { "old": ..., "new": ... }
    ip_address VARCHAR(50),
    user_agent VARCHAR(255),
    
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

CREATE INDEX idx_audit_school_time ON audit_logs(school_id, timestamp);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_changes ON audit_logs USING GIN (changes);

-- Partitions
CREATE TABLE audit_logs_y2024m01 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE audit_logs_y2024m02 PARTITION OF audit_logs
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
-- ... typically created by cron job

-- 2. Data Exports
CREATE TABLE data_exports (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    requester_id INTEGER NOT NULL REFERENCES users(id),
    
    type VARCHAR(50), -- "GDPR", "ACADEMIC"
    status VARCHAR(50) DEFAULT 'PENDING',
    file_url VARCHAR(255),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- 3. User Imports
CREATE TABLE user_imports (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    requester_id INTEGER NOT NULL REFERENCES users(id),
    
    filename VARCHAR(255),
    total_records INTEGER,
    processed_records INTEGER,
    error_log JSONB,
    status VARCHAR(50) DEFAULT 'PROCESSING',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Settings (Key-Value)
CREATE TABLE school_settings (
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    description TEXT,
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (school_id, key)
);
