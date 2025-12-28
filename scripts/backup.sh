#!/bin/bash

# Configuration
DB_CONTAINER="iregistro_db"
DB_USER="postgres"
DB_NAME="iregistro"
BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="backup_$DATE.sql.gz"

# Create backup dir
mkdir -p $BACKUP_DIR

echo "Starting backup of $DB_NAME..."

# Dump and gzip
docker exec $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME | gzip > "$BACKUP_DIR/$FILENAME"

if [ $? -eq 0 ]; then
  echo "Backup successful: $FILENAME"
  
  # Optional: Upload to S3 (requires aws-cli installed)
  # aws s3 cp "$BACKUP_DIR/$FILENAME" s3://my-bucket/
  
  # Cleanup old backups (Keep 7 days)
  find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -exec rm {} \;
else
  echo "Backup failed!"
  exit 1
fi
