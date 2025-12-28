-- Enable Row Level Security on sensitive tables
ALTER TABLE schools ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE students ENABLE ROW LEVEL SECURITY;
ALTER TABLE marks ENABLE ROW LEVEL SECURITY;
ALTER TABLE absences ENABLE ROW LEVEL SECURITY;
ALTER TABLE documents ENABLE ROW LEVEL SECURITY;

-- Create helper function to get current user ID (mock implementation for Postgres, usually matches JWT claim)
-- In Supabase this is auth.uid(). In a custom setup we might use a session variable set by middleware.
-- For standard implementation we can use `current_setting('app.current_user_id', true)`
CREATE OR REPLACE FUNCTION current_user_id() RETURNS INTEGER AS $$
BEGIN
    RETURN current_setting('app.current_user_id', true)::INTEGER;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION current_user_role() RETURNS VARCHAR AS $$
BEGIN
    RETURN current_setting('app.current_user_role', true)::VARCHAR;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION current_school_id() RETURNS INTEGER AS $$
BEGIN
    RETURN current_setting('app.current_school_id', true)::INTEGER;
EXCEPTION WHEN OTHERS THEN
    RETURN NULL;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- --- POLICIES ---

-- 1. Schools
-- Public read for login/selection (or restricted by domain?)
CREATE POLICY "Schools are viewable by everyone" ON schools FOR SELECT USING (true);

-- 2. Users
-- Users can see themselves
CREATE POLICY "Users can see own profile" ON users FOR SELECT USING (id = current_user_id());
-- Staff can see users in their school
CREATE POLICY "Staff can see users in school" ON users FOR SELECT USING (
    school_id = current_school_id() 
    AND current_user_role() IN ('Admin', 'Dirigente', 'Segreteria', 'Insegnante')
);

-- 3. Students
-- Teacher sees students in their school (simplification, real: in their classes)
CREATE POLICY "Teachers see school students" ON students FOR SELECT USING (
    school_id = current_school_id()
    AND current_user_role() IN ('Admin', 'Dirigente', 'Segreteria', 'Insegnante')
);
-- Parents see their children
CREATE POLICY "Parents see own children" ON students FOR SELECT USING (
    (parent_1_id = current_user_id() OR parent_2_id = current_user_id())
);
-- Students see themselves
CREATE POLICY "Students see themselves" ON students FOR SELECT USING (id = (SELECT id FROM students WHERE user_id = current_user_id()));


-- 4. Marks
-- Teachers can insert/update marks for their school (or specific classes - tough in SQL without complex joins)
CREATE POLICY "Teachers manage marks" ON marks FOR ALL USING (
    school_id = current_school_id() 
    AND current_user_role() = 'Insegnante'
    -- Ideally check if teacher teaches subject/class
);

-- Parents/Students view marks
CREATE POLICY "View own/children marks" ON marks FOR SELECT USING (
    student_id IN (
        SELECT id FROM students WHERE 
        user_id = current_user_id() -- Student
        OR parent_1_id = current_user_id() -- Parent 1
        OR parent_2_id = current_user_id() -- Parent 2
    )
);

-- 5. Documents
-- Director/Secretary manage documents
CREATE POLICY "Admin manages documents" ON documents FOR ALL USING (
    school_id = current_school_id()
    AND current_user_role() IN ('Admin', 'Dirigente', 'Segreteria')
);

-- Viewable by owner student/parent
CREATE POLICY "View own documents" ON documents FOR SELECT USING (
    status = 'SIGNED' AND
    student_id IN (
        SELECT id FROM students WHERE 
        user_id = current_user_id()
        OR parent_1_id = current_user_id()
        OR parent_2_id = current_user_id()
    )
);
