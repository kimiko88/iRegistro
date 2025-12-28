-- Test Users for iRegistro
-- Password for all users: "password123"
-- Password hash is bcrypt hash of "password123"

-- Default School
INSERT INTO schools (id, name, code, address, city, email)
VALUES (1, 'Test School', 'TEST001', 'Via Roma 1', 'Roma', 'info@testschool.it')
ON CONFLICT (code) DO NOTHING;

-- Admin User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    1,
    'admin@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',  -- password123
    'Admin',
    'Test',
    'Admin',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Director User  
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    2,
    'director@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',
    'Maria',
    'Rossi',
    'Dirigente',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Teacher User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    3,
    'teacher@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',
    'Giovanni',
    'Bianchi',
    'Insegnante',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Secretary User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    4,
    'secretary@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',
    'Laura',
    'Verdi',
    'Segreteria',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Parent User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    5,
    'parent@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',
    'Paolo',
    'Neri',
    'Genitore',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Student User
INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, created_at, updated_at)
VALUES (
    6,
    'student@test.it',
    '$2a$14$xgqsXa.yBpG8CXmaKymrzO/8vFf7DHLrfGwSiPuY1l3Am1D5m3y3K',
    'Marco',
    'Blu',
    'Studente',
    1,
    NOW(),
    NOW()
) ON CONFLICT (email, school_id) DO UPDATE SET password_hash = EXCLUDED.password_hash;

-- Reset sequences to avoid "duplicate key" errors on next insert
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('schools_id_seq', (SELECT MAX(id) FROM schools));

SELECT 'Test users created successfully! Login with password: password123' AS status;
