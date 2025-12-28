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
