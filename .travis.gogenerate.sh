#!/bin/bash

GO111MODULE=on go generate ./...
if [ -n "$(git diff)" ]; then
  echo "Go generate had not been run"
  git diff
  exit 1
fi
