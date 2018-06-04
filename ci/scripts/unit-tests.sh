#!/bin/bash

source ./setup-go.sh

go test -cover ./... > test_coverage.txt
mv test_coverage.txt $GOPATH/coverage-results/.

make ci
