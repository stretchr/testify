#!/bin/bash
set -e

run_tests() {
    go test -v -race ./...
}

run_tests

cd ./v2
run_tests
