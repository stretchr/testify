#!/bin/bash

# Verify that the code snippets in README.md are formatted.
# The tool https://github.com/hougesen/mdsf is used.

if [ -n "$(mdsf verify --config .mdsf.json README.md 2>&1)" ]; then
  echo "Go code in the README.md is not formatted."
  echo "Did you forget to run 'mdsf format --config .mdsf.json README.md'?"
  exit 1
fi
