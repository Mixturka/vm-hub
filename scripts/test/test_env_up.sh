#!/bin/sh

echo "Starting test database environment..."
docker compose -f ./test/docker/docker-compose.yml up -d postgres_test

until docker exec vm-hub-postgres-test-db pg_isready -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done