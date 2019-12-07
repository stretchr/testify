#!/bin/bash

if [[ "$TRAVIS_GO_VERSION" =~ ^1\.(12|13)(\..*)?$ ]]; then
  exit 0
fi

GO111MODULE=on go generate ./...
if [ -n "$(git diff)" ]; then
  echo "Go generate had not been run"
  git diff
  exit 1
fi
