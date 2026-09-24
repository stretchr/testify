#!/usr/bin/env bash

set -euo pipefail

unformatted_files="$(gofmt -l .)"
if [ -n "$unformatted_files" ]; then
  echo "Go code is not formatted:"
  gofmt -d .
  exit 1
fi

go run ./_readme-gofmt/main.go

go generate ./...
repository_status="$(git status --short)"
if [ -n "$repository_status" ]; then
  echo "Go generate output does not match commit."
  echo "Did you forget to run go generate ./... ?"
  exit 1
fi
