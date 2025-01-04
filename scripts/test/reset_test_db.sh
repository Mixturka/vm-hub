#!/bin/sh

echo "Resetting the database to its original state..."
MIGRATE_CMD="docker run --rm --network=docker_app_network -v $(pwd)/internal/infrustructure/database/postgres/migrations:/migrations migrate/migrate"
$MIGRATE_CMD -path=/migrations -database "postgres://postgres:postgres@postgres_test:5432/postgres?sslmode=disable" -verbose down -all
