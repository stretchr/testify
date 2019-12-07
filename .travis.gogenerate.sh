#!/bin/bash

# If GOMOD is defined we are running with Go Modules enabled, either
# automatically or via the GO111MODULE=on environment variable. If modules is
# enabled we skip generation because at the moment the codegen only works
# without modules.
if [[ -z "$(go env GOMOD)" ]]; then
  echo "Skipping go generate because modules are in use"
  exit 0
fi

go get github.com/ernesto-jimenez/gogen/imports
go generate ./...
if [ -n "$(git diff)" ]; then
  echo "Go generate had not been run"
  git diff
  exit 1
fi
