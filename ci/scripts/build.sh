#!/bin/bash

source ./setup-go.sh

echo "Building Source"
make

echo "Building Releases"
make release
