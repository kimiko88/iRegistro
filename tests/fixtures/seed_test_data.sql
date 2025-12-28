-- Seed Data for Testing

-- 1. School
INSERT INTO schools (id, name, code, city) 
VALUES (1, 'Test School', 'TEST001', 'Test City');

-- 2. Users
INSERT INTO users (id, school_id, email, password_hash, role, first_name, last_name) VALUES
(1, 1, 'admin@test.com', '$2a$10$test', 'Admin', 'Admin', 'User'),
(2, 1, 'teacher1@test.com', '$2a$10$test', 'Teacher', 'John', 'Doe'),
(3, 1, 'student1@test.com', '$2a$10$test', 'Student', 'Jane', 'Smith');

-- 3. Academic
INSERT INTO campuses (id, school_id, name) VALUES (1, 1, 'Main Campus');
INSERT INTO curriculums (id, school_id, code, name) VALUES (1, 1, 'SCI', 'Science');
INSERT INTO classes (id, school_id, year, section, academic_year) VALUES (1, 1, 1, 'A', '2024/2025');
INSERT INTO subjects (id, school_id, name) VALUES (1, 1, 'Mathematics');

-- 4. Enrollment
INSERT INTO students (id, school_id, user_id, first_name, last_name, tax_code) 
VALUES (1, 1, 3, 'Jane', 'Smith', 'TAXCODE123');

INSERT INTO class_enrollments (school_id, student_id, class_id, academic_year)
VALUES (1, 1, 1, '2024/2025');
