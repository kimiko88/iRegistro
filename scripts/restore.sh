#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./restore.sh <backup_file.sql.gz>"
    exit 1
fi

BACKUP_FILE=$1
DB_CONTAINER="iregistro_db"
DB_USER="postgres"
DB_NAME="iregistro"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "File not found: $BACKUP_FILE"
    exit 1
fi

echo "WARNING: This will overwrite the database $DB_NAME. Continue? (y/N)"
read -r response
if [[ ! "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
then
    echo "Aborted."
    exit 0
fi

echo "Restoring from $BACKUP_FILE..."

# Unzip and pipe to psql
gunzip -c "$BACKUP_FILE" | docker exec -i $DB_CONTAINER psql -U $DB_USER -d $DB_NAME

if [ $? -eq 0 ]; then
    echo "Restore completed successfully."
else
    echo "Restore failed."
    exit 1
fi
