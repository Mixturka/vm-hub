#!/bin/sh
echo "Running database migrations..."
echo $POSTGRES_MIGRATIONS_PATH

MIGRATE_CMD="docker run --rm --network=docker_app_network -v $POSTGRES_MIGRATIONS_PATH:/migrations migrate/migrate"

if ! $MIGRATE_CMD -path=/migrations -database "postgres://postgres:postgres@postgres_test:5432/postgres?sslmode=disable" -verbose up; then
  echo "Migration failed. Cleaning up..."
  make test-env-down
  exit 1
fi