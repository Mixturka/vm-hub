#!/bin/bash

set -e

make test-env-up

echo "Running tests..."
if ! go test ./... -count=1 -v; then
  echo "Tests failed. Cleaning up..."
  make clean-tests
  exit 1
fi

echo "Tests passed. Cleaning up..."
make clean-tests
