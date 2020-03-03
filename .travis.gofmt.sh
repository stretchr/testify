#!/bin/bash

check_gofmt_in_dir() {
  if [ -n "$(gofmt -l $1)" ]; then
    echo "Go code is not formatted:"
    gofmt -d .
    exit 1
  fi
}

check_gofmt_in_dir .
check_gofmt_in_dir ./v2
