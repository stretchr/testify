#!/bin/bash

echo ".travis.govet.sh"

cd "$(dirname $0)"
DIRS=". assert require mock suite _codegen"
set -e
for subdir in $DIRS; do
  pushd $subdir
  go vet
  popd
done

