#!/usr/bin/env bash

set -e

# Run go vet on all packages. To exclude specific tests that are known to
# trigger vet warnings, use the 'novet' build tag. This is used in the
# following tests:
#
# require/requirements_testing_test.go:
#
#    The 'testing' tests test testify behavior 😜 against a real testing.T,
#    running tests in goroutines to capture Goexit behavior.
#    Such usage triggers would normally trigger vet warnings (SA2002).

go vet -tags novet ./...
