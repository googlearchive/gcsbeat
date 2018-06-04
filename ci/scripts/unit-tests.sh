#!/bin/bash

source $(dirname $0)/setup-go.sh

go test -cover ./... > test_coverage.txt
mkdir coverage-results
mv test_coverage.txt coverage-results/

echo "Checking for style guidelines"
make check

echo "Testsuite"
make testsuite
