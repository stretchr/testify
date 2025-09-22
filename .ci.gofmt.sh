#!/usr/bin/env bash

if [ -n "$(go fmt ./...)" ]; then
  echo "Go code is not formatted:"
  go fmt ./...
  exit 1
fi

go generate ./...
if [ -n "$(git status -s -uno)" ]; then
  echo "Go generate output does not match commit."
  echo "Did you forget to run go generate ./... ?"
  exit 1
fi
