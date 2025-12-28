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
