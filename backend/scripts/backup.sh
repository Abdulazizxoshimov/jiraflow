#!/usr/bin/env bash
# PostgreSQL backup — dumps the DB, optionally uploads to S3/MinIO.
# Usage: ./backup.sh
# Env vars:
#   DATABASE_URL        — postgres connection string (required)
#   BACKUP_DIR          — local backup directory (default: /backups)
#   S3_BUCKET           — s3://bucket/prefix to sync backups (optional)
#   RETENTION_DAYS      — how many days to keep local backups (default: 7)

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FILE="${BACKUP_DIR}/jiraflow-${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

echo "[backup] dumping database → ${FILE}"
pg_dump "$DATABASE_URL" | gzip > "$FILE"
echo "[backup] done ($(du -sh "$FILE" | cut -f1))"

# Upload to S3/MinIO if configured
if [[ -n "${S3_BUCKET:-}" ]]; then
    echo "[backup] uploading to ${S3_BUCKET}"
    aws s3 cp "$FILE" "${S3_BUCKET}/$(basename "$FILE")" --no-progress
    echo "[backup] upload complete"
fi

# Remove backups older than RETENTION_DAYS
echo "[backup] pruning backups older than ${RETENTION_DAYS} days"
find "$BACKUP_DIR" -name "jiraflow-*.sql.gz" -mtime "+${RETENTION_DAYS}" -delete

echo "[backup] finished at ${TIMESTAMP}"
