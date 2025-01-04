#!/bin/bash

echo "Resetting the database to its original state..."
echo $TEST_POSTGRES_URL
MIGRATE_CMD="docker run --rm --network=docker_app_network -v $POSTGRES_MIGRATIONS_PATH:/migrations migrate/migrate"
$MIGRATE_CMD -path=/migrations -database $TEST_POSTGRES_URL -verbose down -all
