#!/usr/bin/env bash

set -euo pipefail

# If GOMOD is defined we are running with Go Modules enabled, either
# automatically or via the GO111MODULE=on environment variable. Codegen only
# works with modules, so skip generation if modules is not in use.
gomod="$(go env GOMOD)"
if [[ -z "$gomod" ]]; then
  echo "Skipping go generate because modules not enabled and required"
  exit 0
fi

go generate ./...
repository_status="$(git status --short)"
if [ -n "$repository_status" ]; then
  echo "Go generate had not been run"
  git status --short
  git diff
  exit 1
fi
