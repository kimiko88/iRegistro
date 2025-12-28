-- 1. Documents (Generic: Reports, PDP, PFI, etc.)
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- PAGELLA, PDP, PFI, 15_MAGGIO
    
    data JSONB, -- Flexible content (grades list, PDP goals, etc.)
    file_url VARCHAR(255), -- If generated PDF is stored
    
    student_id INTEGER REFERENCES students(id),
    class_id INTEGER REFERENCES classes(id),
    
    created_by INTEGER REFERENCES users(id),
    academic_year VARCHAR(20),
    status VARCHAR(50) DEFAULT 'DRAFT', -- DRAFT, SIGNED, ARCHIVED
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_documents_school_type ON documents(school_id, type);
CREATE INDEX idx_documents_student ON documents(student_id);

-- 2. Signatures
CREATE TABLE document_signatures (
    id SERIAL PRIMARY KEY,
    document_id INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    signer_id INTEGER NOT NULL REFERENCES users(id),
    signature_timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    signer_name VARCHAR(100), -- Snapshot of name at signing
    ip_address VARCHAR(50),
    is_valid BOOLEAN DEFAULT TRUE
);

-- 3. PCTO (Alternanza Scuola Lavoro)
CREATE TABLE pcto_projects (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    companies JSONB, -- Partner companies details
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE pcto_assignments (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES pcto_projects(id) ON DELETE CASCADE,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    company_id VARCHAR(100), -- ID inside the companies JSON
    
    start_date DATE,
    end_date DATE,
    hours_planned INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'ACTIVE',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE pcto_hours (
    id SERIAL PRIMARY KEY,
    assignment_id INTEGER NOT NULL REFERENCES pcto_assignments(id) ON DELETE CASCADE,
    
    week VARCHAR(20), -- "2024-W12"
    hours_worked INTEGER NOT NULL,
    description TEXT,
    tutor_notes TEXT,
    status VARCHAR(50) DEFAULT 'PENDING',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 4. Orientation
CREATE TABLE orientation_activities (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    date DATE,
    hours INTEGER,
    target_grade INTEGER -- e.g. 5th year
);

CREATE TABLE orientation_participations (
    id SERIAL PRIMARY KEY,
    activity_id INTEGER NOT NULL REFERENCES orientation_activities(id) ON DELETE CASCADE,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    hours_earned INTEGER,
    evaluation VARCHAR(255),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. PDP/PGF Specifics (Optional - using Documents table predominantly, but sometimes relational)
-- Creating views for convenience
CREATE OR REPLACE VIEW pdp_documents AS
SELECT * FROM documents WHERE type = 'PDP';
