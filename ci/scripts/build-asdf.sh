#!/bin/bash

SCRIPTPATH=$(dirname $0)


source $SCRIPTPATH/install-dependencies.sh
source $(dirname $0)/setup-go.sh

echo "Building Source"
make
