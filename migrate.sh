#!/bin/bash

# Load environment variables
source .env

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Running database migrations..."
echo "Database URL: $DB_URL"

# Create database if it doesn't exist
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || true

# Run SQL migration directly
psql $DB_URL -f internal/adapter/sqlc/schema/001_users.sql

echo "Migration completed successfully!"

# Generate SQLC code
echo "Generating SQLC code..."
sqlc generate

echo "SQLC generation completed!"