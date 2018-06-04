#!/bin/bash

PROGNAME=gcsbeat
GODIR=github.com/GoogleCloudPlatform/gcsbeat

set -e -x

echo "pwd: " $PWD
ls -lah

#export GOPATH=$PWD

echo "Environment Variables"
env

mkdir -p /go/src/github.com/GoogleCloudPlatform
ln -s $PWD/$PROGNAME /go/src/github.com/GoogleCloudPlatform/gcsbeat

#cd $GODIR
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
ls -lah
