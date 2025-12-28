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

    -- Constraints
    CONSTRAINT idx_users_email_school UNIQUE (email, school_id)
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_school_id ON users(school_id);

-- 4. User Sessions (Refresh Tokens)
CREATE TABLE user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL UNIQUE,
    user_agent VARCHAR(255),
    ip_address VARCHAR(45),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_sessions_token ON user_sessions(refresh_token);
