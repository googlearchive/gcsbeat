#!/bin/bash

source $(dirname $0)/setup-go.sh

echo "Building Source"
make

echo "Building Releases"
make release
