#!/bin/bash

source $(dirname $0)/setup-go.sh

go test -cover ./... > test_coverage.txt

cat test_coverage.txt
