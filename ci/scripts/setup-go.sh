#!/bin/bash

PROGNAME=gcsbeat
GODIR=github.com/GoogleCloudPlatform/gcsbeat

set -e -x

echo "pwd: " $PWD
ls -lah

export GOPATH=$PWD

echo "Environment Variables"
env

mkdir -p $GODIR
cp -R ./$PROGNAME/* src/$GODIR


cd $GODIR
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
ls -lah
