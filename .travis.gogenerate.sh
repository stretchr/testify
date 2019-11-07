#!/bin/bash

echo ".travis.gogenerate.sh"

go build -o codegen ./_codegen
go get github.com/ernesto-jimenez/gogen/imports
go generate ./...
if [ -n "$(git diff assert/ mock/ require/ suite/)" ]; then
  echo "Go generate had not been run"
  git diff
  exit 1
else
  echo "Go generate had been run"
fi

