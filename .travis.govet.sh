#!/bin/bash

set -e

go_vet_in_dir() {
    cd $1
    go vet ./...
}

go_vet_in_dir .
go_vet_in_dir ./v2
