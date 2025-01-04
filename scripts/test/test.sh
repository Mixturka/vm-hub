#!/bin/bash
make test-env-up
make migrate-test

echo "Running tests..."
go test ./... -v

make clean-tests