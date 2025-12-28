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
