#!/bin/bash


source $(dirname $0)/setup-go.sh

go test -cover ./... > test_coverage.txt
mv test_coverage.txt $GOPATH/coverage-results/.

make ci
