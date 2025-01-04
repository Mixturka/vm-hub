#!/bin/bash
echo "Running database migrations..."

MIGRATE_CMD="docker run --rm --network=docker_app_network -v $POSTGRES_MIGRATIONS_PATH:/migrations migrate/migrate"

if ! $MIGRATE_CMD -path=/migrations -database $TEST_POSTGRES_URL -verbose up; then
  echo "Migration failed. Cleaning up..."
  make test-env-down
  exit 1
fi