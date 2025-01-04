#!/bin/bash

set -e

echo "Starting test database environment..."

IMAGE_NAME="migrate_local"

# Check if the image exists, and build it if not
if [[ "$(docker images -q "$IMAGE_NAME" 2> /dev/null)" == "" ]]; then
  echo "Image $IMAGE_NAME not found, building..."
  docker build -t "$IMAGE_NAME" -f ./test/docker/Dockerfile.migrate ./test/docker
else
  echo "Image $IMAGE_NAME found, skipping build..."
fi

# Start the test database using docker-compose
docker compose -f ./test/docker/docker-compose.yml up -d postgres_test

# Wait for PostgreSQL to be ready
until docker exec vm-hub-postgres-test-db pg_isready -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

echo "PostgreSQL is ready."
