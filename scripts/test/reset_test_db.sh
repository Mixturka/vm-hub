#!/bin/bash

set -e

echo "Resetting the database to its original state..."

if [[ -z "$TEST_POSTGRES_URL" ]]; then
  echo "Error: TEST_POSTGRES_URL is not set."
  exit 1
fi

MIGRATE_CMD="docker run --rm --network=docker_app_network -v $POSTGRES_MIGRATIONS_PATH:/migrations migrate_local migrate"

$MIGRATE_CMD -path=/migrations -database "$TEST_POSTGRES_URL" -verbose down -all
