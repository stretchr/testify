#!/usr/bin/env bash

set -euo pipefail

if [ -n "$(gofmt -l .)" ]; then
  echo "Go code is not formatted:"
  gofmt -d .
  exit 1
fi

go run ./_readme-gofmt/main.go

go generate ./...
if [ -n "$(git status -s -uno)" ]; then
  echo "Go generate output does not match commit."
  echo "Did you forget to run go generate ./... ?"
  exit 1
fi
