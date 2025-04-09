#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$1
    echo "  Creating database '$database'"
    for i in {1..10}; do
        if psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -c "CREATE DATABASE $database"; then
            psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -c "GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER"
            return 0
        fi
        echo "  Attempt $i failed, retrying..."
        sleep 2
    done
    echo "  Failed to create database '$database' after 10 attempts"
    exit 1
}

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to become available..."
for i in {1..30}; do
    if pg_isready -U "$POSTGRES_USER" -h localhost; then
        break
    fi
    sleep 1
    if [ $i -eq 30 ]; then
        echo "PostgreSQL never became available"
        exit 1
    fi
done
echo "PostgreSQL is ready"

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
    echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
    for db in $(echo "$POSTGRES_MULTIPLE_DATABASES" | tr ',' ' '); do
        if ! psql -U "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "$db"; then
            create_user_and_database "$db"
        else
            echo "  Database '$db' already exists - skipping creation"
        fi
    done
    echo "Multiple databases created"
fi