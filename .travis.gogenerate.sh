#!/bin/bash

# If GOMOD is defined we are running with Go Modules enabled, either
# automatically or via the GO111MODULE=on environment variable. Codegen only
# works with modules, so skip generation if modules is not in use.
if [[ -z "$(go env GOMOD)" ]]; then
  echo "Skipping go generate because modules not enabled and required"
  exit 0
fi

go_generate_in_dir() {
  cd $1
  go generate ./...
  if [ -n "$(git status -s -uno)" ]; then
    echo "Go generate output does not match commit."
    echo "Did you forget to run go generate ./... ?"
    exit 1
  fi
}

go_generate_in_dir .
go_generate_in_dir ./v2

