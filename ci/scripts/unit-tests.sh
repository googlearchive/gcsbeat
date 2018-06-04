#!/bin/bash

source $(dirname $0)/setup-go.sh

go test -cover ./... > test_coverage.txt
mkdir coverage-results
mv test_coverage.txt coverage-results/

make ci
