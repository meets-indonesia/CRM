#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$1
    echo "  Creating database '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE $database;
        GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
}

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to become available..."
while ! pg_isready -U "$POSTGRES_USER" -h localhost; do
    sleep 1
done
echo "PostgreSQL is ready"

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
    echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
    for db in $(echo "$POSTGRES_MULTIPLE_DATABASES" | tr ',' ' '); do
        # Check if database already exists
        if ! psql -U "$POSTGRES_USER" -lqt | cut -d \| -f 1 | grep -qw "$db"; then
            create_user_and_database "$db"
        else
            echo "  Database '$db' already exists - skipping creation"
        fi
    done
    echo "Multiple databases created"
fi