-- Enable UUID extension if needed, though we use SERIAL/INTEGER for IDs as requested.
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Schools Table
CREATE TABLE schools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE, -- Meccanographic code
    address TEXT,
    city VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    website VARCHAR(255),
    logo_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_schools_code ON schools(code);

-- 2. User Roles Enum (optional, or use VARCHAR with check)
CREATE TYPE user_role AS ENUM (
    'SuperAdmin', 
    'Admin', 
    'Principal', 
    'Teacher', 
    'Parent', 
    'Student', 
    'Secretary'
);

-- 3. Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    school_id INTEGER REFERENCES schools(id), -- Nullable for SuperAdmin
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL, -- using VARCHAR matching ENUM for flexibility or actual ENUM
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    tax_code VARCHAR(16), -- Codice Fiscale
    
    -- Security
    two_fa_enabled BOOLEAN DEFAULT FALSE,
    two_fa_secret VARCHAR(100),
    failed_logins INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    reset_token_hash VARCHAR(255),
    reset_token_exp TIMESTAMP WITH TIME ZONE,
    
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    -- Auth/Identity
    auth_method VARCHAR(50) DEFAULT 'LOCAL', -- LOCAL, SPID, CIE
    sp_id_provider VARCHAR(100),
    cie_serial_number VARCHAR(100),
    external_id VARCHAR(100),
    last_auth_at TIMESTAMP WITH TIME ZONE,

    -- Constraints
    CONSTRAINT idx_users_email_school UNIQUE (email, school_id)
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_school_id ON users(school_id);

-- 4. User Sessions (Refresh Tokens)
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent VARCHAR(255),
    ip_address VARCHAR(45),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_token ON user_sessions(refresh_token);
-- 1. Campuses (Sedi)
CREATE TABLE campuses (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL, -- e.g. "Sede Centrale"
    address VARCHAR(255),
    contact_email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_campuses_school_id ON campuses(school_id);

-- 2. Curriculums (Indirizzi di studio)
CREATE TABLE curriculums (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL, -- e.g. "INF", "LIC-SCI"
    name VARCHAR(100) NOT NULL, -- e.g. "Informatica", "Liceo Scientifico"
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_curriculums_school_id ON curriculums(school_id);

-- 3. Classes
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    campus_id INTEGER REFERENCES campuses(id),
    curriculum_id INTEGER REFERENCES curriculums(id),
    
    year INTEGER NOT NULL, -- e.g. 1, 2, 3, 4, 5
    section VARCHAR(5) NOT NULL, -- e.g. "A", "B"
    academic_year VARCHAR(20) NOT NULL, -- e.g. "2024/2025"
    
    room VARCHAR(50),
    coordinator_id INTEGER REFERENCES users(id), -- Teacher ID
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT uq_class_definition UNIQUE (school_id, year, section, academic_year)
);

CREATE INDEX idx_classes_school_year ON classes(school_id, academic_year);
CREATE INDEX idx_classes_curriculum ON classes(curriculum_id);

-- 4. Subjects (Materie)
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    short_name VARCHAR(20), -- e.g. "MAT", "ITA"
    department VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_subjects_school_id ON subjects(school_id);

-- 5. Class Subject Assignments (Which teacher teaches what in which class)
CREATE TABLE class_subject_assignments (
    id SERIAL PRIMARY KEY,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    teacher_id INTEGER REFERENCES users(id),
    
    hours_per_week INTEGER DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT uq_class_subject UNIQUE (class_id, subject_id, teacher_id) -- Or just class+subject? Sometimes split?
);

CREATE INDEX idx_assignments_class_id ON class_subject_assignments(class_id);
CREATE INDEX idx_assignments_teacher_id ON class_subject_assignments(teacher_id);
-- 1. Students Table
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id), -- Link to User account for login
    
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE,
    place_of_birth VARCHAR(100),
    tax_code VARCHAR(16) NOT NULL, -- Codice Fiscale
    gender VARCHAR(1), -- M/F
    
    address VARCHAR(255),
    city VARCHAR(100),
    email VARCHAR(255), -- Personal email
    phone VARCHAR(50),
    
    -- Parent Links (Often implicit via separate table, but simple fields first)
    parent_1_id INTEGER REFERENCES users(id),
    parent_2_id INTEGER REFERENCES users(id),
    
    enrollment_date DATE DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'ACTIVE', -- ACTIVE, TRANSFERRED, GRADUATED
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT uq_student_tax_code UNIQUE (school_id, tax_code)
);

CREATE INDEX idx_students_school_id ON students(school_id);
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_students_tax_code ON students(tax_code);

-- 2. Class Enrollments (Students in Classes)
CREATE TABLE class_enrollments (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    class_id INTEGER NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    
    academic_year VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL DEFAULT CURRENT_DATE,
    end_date DATE, -- Null if currently enrolled
    
    active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT uq_enrollment_period UNIQUE (student_id, class_id, academic_year)
);

CREATE INDEX idx_enrollments_class_id ON class_enrollments(class_id);
CREATE INDEX idx_enrollments_student_id ON class_enrollments(student_id);
CREATE INDEX idx_enrollments_year ON class_enrollments(academic_year);
-- 1. Marks (Partitioned by Year/Date Range)
-- Since 'year' isn't a date, we use mark_date for partitioning
CREATE TABLE marks (
    id SERIAL, -- Partitions don't support global serial PK constraints easily in some PG versions, but PG10+ is usually okay strictly speaking it's "id" + keys.
               -- Often better to use UUID or composite key for partitioning.
               -- For simplicity we use BIGSERIAL and PRIMARY KEY (id, mark_date) for partitioning requirement.
    school_id INTEGER NOT NULL, -- FK References Check is tricky on partitions, usually ignored or manually managed
    student_id INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    teacher_id INTEGER,
    
    mark_date DATE NOT NULL,
    value DECIMAL(4, 2) NOT NULL, -- e.g. 8.50
    type VARCHAR(50), -- Oral, Written, Practical
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Partition Key must be part of PK
    PRIMARY KEY (id, mark_date)
) PARTITION BY RANGE (mark_date);

-- Indexes (Local)
CREATE INDEX idx_marks_student_subject ON marks(student_id, subject_id, mark_date);
CREATE INDEX idx_marks_class_date ON marks(class_id, mark_date);
CREATE INDEX idx_marks_teacher ON marks(teacher_id);

-- Create Initial Partitions
CREATE TABLE marks_y2023 PARTITION OF marks
    FOR VALUES FROM ('2023-09-01') TO ('2024-09-01');

CREATE TABLE marks_y2024 PARTITION OF marks
    FOR VALUES FROM ('2024-09-01') TO ('2025-09-01');
    
CREATE TABLE marks_y2025 PARTITION OF marks
    FOR VALUES FROM ('2025-09-01') TO ('2026-09-01');


-- 2. Absences (Partitioned by Date)
CREATE TABLE absences (
    id SERIAL,
    school_id INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    class_id INTEGER NOT NULL,
    
    absence_date DATE NOT NULL,
    start_hour INTEGER, -- e.g. 1st hour
    end_hour INTEGER,
    type VARCHAR(50), -- ABSENCE, LATE, EARLY_EXIT
    justified BOOLEAN DEFAULT FALSE,
    justification_note TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    PRIMARY KEY (id, absence_date)
) PARTITION BY RANGE (absence_date);

CREATE INDEX idx_absences_student_date ON absences(student_id, absence_date);
CREATE INDEX idx_absences_class_date ON absences(class_id, absence_date);

-- Partitions
CREATE TABLE absences_y2023 PARTITION OF absences
    FOR VALUES FROM ('2023-09-01') TO ('2024-09-01');

CREATE TABLE absences_y2024 PARTITION OF absences
    FOR VALUES FROM ('2024-09-01') TO ('2025-09-01');

CREATE TABLE absences_y2025 PARTITION OF absences
    FOR VALUES FROM ('2025-09-01') TO ('2026-09-01');

-- 3. Substitutions (Regular table, low volume)
CREATE TABLE teacher_substitutions (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    original_teacher_id INTEGER REFERENCES users(id),
    substitute_teacher_id INTEGER REFERENCES users(id),
    class_id INTEGER REFERENCES classes(id),
    
    date DATE NOT NULL,
    hour INTEGER NOT NULL,
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- 1. Notifications
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- ABSENCE, MARK, MESSAGE, SYSTEM
    title VARCHAR(255) NOT NULL,
    message TEXT,
    data JSONB, -- Link details, e.g. { "mark_id": 123 }
    
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id) WHERE is_read = FALSE;

-- 2. Conversations
CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    subject VARCHAR(255),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE conversation_participants (
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_read_at TIMESTAMP WITH TIME ZONE,
    
    PRIMARY KEY (conversation_id, user_id)
);

-- 3. Messages
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id INTEGER REFERENCES users(id), -- Nullable for system messages
    body TEXT NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_conversation ON messages(conversation_id);
-- FTS Index on Body
CREATE INDEX idx_messages_search ON messages USING GIN (to_tsvector('italian', body));

-- 4. Attachments
CREATE TABLE message_attachments (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    file_type VARCHAR(100),
    file_size INTEGER,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 5. Colloquiums (Parent-Teacher Meetings)
CREATE TABLE colloquium_slots (
    id SERIAL PRIMARY KEY,
    school_id INTEGER NOT NULL REFERENCES schools(id),
    teacher_id INTEGER NOT NULL REFERENCES users(id),
    
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_bookings INTEGER DEFAULT 1,
    
    location VARCHAR(255), -- Room or Link
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_colloquium_teacher ON colloquium_slots(school_id, teacher_id, date);

CREATE TABLE colloquium_bookings (
    id SERIAL PRIMARY KEY,
    slot_id INTEGER NOT NULL REFERENCES colloquium_slots(id) ON DELETE CASCADE,
    parent_id INTEGER NOT NULL REFERENCES users(id),
    student_id INTEGER REFERENCES students(id), -- Which child is discussed
    
    status VARCHAR(50) DEFAULT 'CONFIRMED', -- CONFIRMED, CANCELLED
    notes TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
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
-- 1. Student Current Marks View
-- Calculates average per subject for current year
CREATE OR REPLACE VIEW student_current_marks AS
SELECT 
    m.student_id,
    m.subject_id,
    s.name AS subject_name,
    AVG(m.value) AS average_mark,
    COUNT(m.id) AS mark_count,
    MAX(m.mark_date) AS last_mark_date
FROM marks m
JOIN subjects s ON m.subject_id = s.id
WHERE m.mark_date >= DATE_TRUNC('year', CURRENT_DATE) -- Simplification for "Current Year"
GROUP BY m.student_id, m.subject_id, s.name;

-- 2. Class Statistics View
CREATE OR REPLACE VIEW class_statistics AS
SELECT 
    c.id AS class_id,
    c.year,
    c.section,
    COUNT(DISTINCT e.student_id) AS total_students,
    AVG(m.value) AS class_average,
    COUNT(DISTINCT a.id) AS total_absences
FROM classes c
LEFT JOIN class_enrollments e ON c.id = e.class_id AND e.active = TRUE
LEFT JOIN marks m ON c.id = m.class_id
LEFT JOIN absences a ON c.id = a.class_id
GROUP BY c.id, c.year, c.section;

-- 3. Teacher Workload View
CREATE OR REPLACE VIEW teacher_workload AS
SELECT 
    u.id AS teacher_id,
    u.first_name,
    u.last_name,
    COUNT(DISTINCT csa.class_id) AS total_classes,
    SUM(csa.hours_per_week) AS total_hours
FROM users u
JOIN class_subject_assignments csa ON u.id = csa.teacher_id
WHERE u.role = 'Teacher'
GROUP BY u.id, u.first_name, u.last_name;

-- 4. Colloquium Availability View
CREATE OR REPLACE VIEW colloquium_availability AS
SELECT 
    cs.id AS slot_id,
    cs.teacher_id,
    cs.date,
    cs.start_time,
    cs.end_time,
    (cs.max_bookings - COUNT(cb.id)) AS available_spots
FROM colloquium_slots cs
LEFT JOIN colloquium_bookings cb ON cs.id = cb.slot_id AND cb.status = 'CONFIRMED'
WHERE cs.date >= CURRENT_DATE
GROUP BY cs.id, cs.teacher_id, cs.date, cs.start_time, cs.end_time, cs.max_bookings
HAVING (cs.max_bookings - COUNT(cb.id)) > 0;
