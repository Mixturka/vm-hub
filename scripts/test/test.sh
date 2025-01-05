#!/bin/bash

set -e

make test-env-up
make migrate-test

echo "Running tests..."
if ! go test ./... -v; then
  echo "Tests failed. Cleaning up..."
  make clean-tests
  exit 1
fi

echo "Tests passed. Cleaning up..."
make clean-tests
