# Self-Hosting Guide

This guide describes how to deploy iRegistro using Docker Compose on a Linux server.

## Prerequisites
- Docker Engine & Docker Compose
- 4GB RAM minimum recommended
- 20GB Disk Space

## Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/k/iRegistro.git
   cd iRegistro
   ```

2. **Configure Environment**
   Copy the example configuration:
   ```bash
   cp .env.example .env
   # Edit .env and set secure passwords!
   ```

3. **Start Services**
   ```bash
   docker-compose up -d
   ```
   This starts:
   - PostgreSQL (Port 5432)
   - MinIO (Port 9000/9001)
   - pgAdmin (Port 5050)
   - Backup Service

4. **Initialize Database**
   Migrations in `migrations/` are auto-mapped but you may need to run them manually or check logs if `entrypoint-initdb` runs them (standard postgres image behavior for `.sql` files, but we have `golang-migrate` format).
   
   Recommended: Use `golang-migrate` locally or use the consolidated `schema.sql`.
   ```bash
   # Copy schema to container
   docker cp schema.sql iregistro_db:/tmp/schema.sql
   # Execute
   docker exec -it iregistro_db psql -U postgres -d iregistro -f /tmp/schema.sql
   ```

## Backups
- **Automated**: The `postgres_backup` container creates daily backups in `./backups`.
- **Manual**: Run `./scripts/backup.sh`.
- **Restore**: Run `./scripts/restore.sh <file>`.

## Storage (MinIO)
- Access Console: `http://localhost:9001`
- Login: `minioadmin` / `minioadmin` (or values in `.env`)
- Create buckets: `documents`, `attachments`
