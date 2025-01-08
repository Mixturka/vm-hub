#!/bin/bash

set -e

echo "Starting test database environment..."

docker compose -f ./test/docker/docker-compose.yml up -d postgres_test redis_test

until docker exec vm-hub-postgres-test-db pg_isready -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

echo "PostgreSQL is ready."
echo "Waiting for Redis..."
until docker exec vm-hub-redis-test-db redis-cli ping | grep -q PONG; do
  echo "Redis is not ready yet. Retrying in 2 seconds..."
  sleep 2
done
echo "Redis is ready."

echo "All test services are up and running!"