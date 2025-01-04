#!/bin/bash

IMAGE_NAME="migrate"

echo "Starting test database environment..."

if [[ "$(docker images -q $IMAGE_NAME 2> /dev/null)" == "" ]]; then
  echo "Image $IMAGE_NAME not found, building..."
  docker compose -f ./test/docker/docker-compose.yml build migrate
else
  echo "Image $IMAGE_NAME found, skipping build..."
fi

docker compose -f ./test/docker/docker-compose.yml up -d postgres_test

until docker exec vm-hub-postgres-test-db pg_isready -U postgres; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done