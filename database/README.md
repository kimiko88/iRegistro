# iRegistro Database Schema

This directory contains the database migration files and seed data for the iRegistro application (PostgreSQL).

## Structure
- `migrations/`: Contains `golang-migrate` compatible SQL files.
    - `001`: Schools, Users, Roles.
    - `002`: Academic Structure (Campuses, Classes, Curriculums, Subjects).
    - `003`: Students and Enrollments.
    - `004`: Core Features (Marks, Absences - Partitioned).
    - `005`: Communication (Messages, Chats, Colloquiums).
    - `006`: Documents (Reports, Signatures, PCTO).
    - `007`: System (Audit Logs - Partitioned, Exports).
    - `008`: Views (Performance & Statistics).
- `seeds/`: Contains `mock_data.sql` for populating the development database.
- `schema.sql`: Aggregated schema DDL.

## Setup Instructions

### Prerequisites
- PostgreSQL 12+ (for Partitioning support)
- `golang-migrate` tool (optional, for running migrations)

### Running Migrations
```bash
# Using migrate tool
migrate -source file://migrations -database postgres://user:pass@localhost:5432/iregistro?sslmode=disable up

# Or manually applying files 001 to 008 in order using psql
psql -U user -d iregistro -f schema.sql
```

### Seeding Data
```bash
psql -U user -d iregistro -f seeds/mock_data.sql
```

## Features
- **Partitioning**: `marks` and `absences` tables are partitioned by RANGE (Date/Year) for performance scaling.
- **Search**: `messages` table uses GIN index on `to_tsvector` for full-text search.
- **Audit**: `audit_logs` are partitioned by month.
