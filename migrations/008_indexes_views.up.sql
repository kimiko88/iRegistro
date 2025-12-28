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
