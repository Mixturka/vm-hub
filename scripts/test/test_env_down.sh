#!/bin/bash

echo "Tearing down test database environment..."
docker compose -f ./test/docker/docker-compose.yml down -v