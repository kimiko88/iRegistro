-- Mock Data for iRegistro Development

-- 1. School
INSERT INTO schools (name, code, city, email)
VALUES ('Liceo Classico Alessandro Manzoni', 'LCMKZ12345', 'Milano', 'segreteria@manzoni.edu.it');

-- 2. Users (50 Total as requested: 1 Admin, 1 Dir, 1 Sec, 15 Teachers, 30 Parents, 2 Students explicit)
-- Passwords are 'hashed_secret' placeholder
INSERT INTO users (school_id, email, password_hash, role, first_name, last_name) VALUES
((SELECT id FROM schools LIMIT 1), 'admin@manzoni.edu.it', '$2a$10$xyz...', 'Admin', 'Mario', 'Rossi'),
((SELECT id FROM schools LIMIT 1), 'dirigente@manzoni.edu.it', '$2a$10$xyz...', 'Dirigente', 'Giulia', 'Bianchi'),
((SELECT id FROM schools LIMIT 1), 'segreteria@manzoni.edu.it', '$2a$10$xyz...', 'Segreteria', 'Anna', 'Verdi');

-- Teachers (15)
INSERT INTO users (school_id, email, password_hash, role, first_name, last_name)
SELECT 
    (SELECT id FROM schools LIMIT 1), 
    'docente' || generate_series || '@manzoni.edu.it', 
    '$2a$10$xyz...', 
    'Insegnante', 
    'Docente', 
    'Numero' || generate_series
FROM generate_series(1, 15);

-- 3. Academic Structure
-- Campuses
INSERT INTO campuses (school_id, name) VALUES ((SELECT id FROM schools LIMIT 1), 'Sede Centrale');

-- Curriculums
INSERT INTO curriculums (school_id, code, name) VALUES ((SELECT id FROM schools LIMIT 1), 'CLASSICO', 'Liceo Classico');

-- Classes (5: 1A, 2A, 3A, 4A, 5A)
INSERT INTO classes (school_id, campus_id, curriculum_id, year, section, academic_year)
SELECT 
    (SELECT id FROM schools LIMIT 1),
    (SELECT id FROM campuses LIMIT 1),
    (SELECT id FROM curriculums LIMIT 1),
    s, 'A', '2024/2025'
FROM generate_series(1, 5) AS s;

-- Subjects
INSERT INTO subjects (school_id, name, short_name) VALUES
((SELECT id FROM schools LIMIT 1), 'Italiano', 'ITA'),
((SELECT id FROM schools LIMIT 1), 'Latino', 'LAT'),
((SELECT id FROM schools LIMIT 1), 'Greco', 'GRE'),
((SELECT id FROM schools LIMIT 1), 'Matematica', 'MAT'),
((SELECT id FROM schools LIMIT 1), 'Storia', 'STO'),
((SELECT id FROM schools LIMIT 1), 'Filosofia', 'FIL'),
((SELECT id FROM schools LIMIT 1), 'Inglese', 'ING'),
((SELECT id FROM schools LIMIT 1), 'Scienze', 'SCI');

-- Assignments (Randomly assign teachers to classes/subjects)
INSERT INTO class_subject_assignments (class_id, subject_id, teacher_id)
SELECT 
    c.id, 
    s.id, 
    (SELECT id FROM users WHERE role = 'Insegnante' ORDER BY RANDOM() LIMIT 1)
FROM classes c
CROSS JOIN subjects s
WHERE s.short_name IN ('ITA', 'MAT', 'STO'); -- Assign core subjects to all classes

-- 4. Students (150 total, ~30 per class)
INSERT INTO students (school_id, first_name, last_name, tax_code, date_of_birth)
SELECT 
    (SELECT id FROM schools LIMIT 1),
    'Studente', 
    'Cognome' || generate_series, 
    'TAXCODE' || generate_series,
    '2008-01-01'::date + (generate_series || ' days')::interval
FROM generate_series(1, 150);

-- Enrollments (Distribute 150 students into 5 classes)
INSERT INTO class_enrollments (school_id, student_id, class_id, academic_year)
SELECT 
    (SELECT id FROM schools LIMIT 1),
    s.id,
    (SELECT id FROM classes WHERE year = ((s.id % 5) + 1) LIMIT 1),
    '2024/2025'
FROM students s;

-- 5. Marks (200 Test Marks)
INSERT INTO marks (school_id, student_id, class_id, subject_id, teacher_id, mark_date, value, type)
SELECT 
    (SELECT id FROM schools LIMIT 1),
    s.id,
    e.class_id,
    (SELECT id FROM subjects ORDER BY RANDOM() LIMIT 1),
    (SELECT id FROM users WHERE role = 'Insegnante' ORDER BY RANDOM() LIMIT 1),
    '2024-10-01'::date + ((random() * 60)::int || ' days')::interval,
    (random() * 8 + 2)::decimal(4,2), -- Marks 2.00 to 10.00
    CASE WHEN random() > 0.5 THEN 'Written' ELSE 'Oral' END
FROM students s
JOIN class_enrollments e ON s.id = e.student_id
ORDER BY RANDOM()
LIMIT 200;

-- 6. Absences (100 Test Absences)
INSERT INTO absences (school_id, student_id, class_id, absence_date, type, justified)
SELECT 
    (SELECT id FROM schools LIMIT 1),
    s.id,
    e.class_id,
    '2024-10-01'::date + ((random() * 60)::int || ' days')::interval,
    'ABSENCE',
    (random() > 0.5)
FROM students s
JOIN class_enrollments e ON s.id = e.student_id
ORDER BY RANDOM()
LIMIT 100;
