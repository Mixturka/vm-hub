#!/bin/bash

set -e

echo "Running database migrations..."

if [[ -z "$TEST_POSTGRES_URL" ]]; then
  echo "Error: TEST_POSTGRES_URL is not set."
  exit 1
fi

MIGRATE_CMD="docker run --rm --network=docker_app_network -v $POSTGRES_MIGRATIONS_PATH:/migrations migrate_local migrate"

if ! $MIGRATE_CMD -path=/migrations -database "$TEST_POSTGRES_URL" -verbose up; then
  echo "Migration failed. Cleaning up..."
  make test-env-down
  exit 1
fi

echo "Migrations completed successfully."
