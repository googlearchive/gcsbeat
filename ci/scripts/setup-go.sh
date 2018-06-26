#!/bin/bash

PROGNAME=gcsbeat
GODIR=github.com/GoogleCloudPlatform/gcsbeat

mkdir -p $GOPATH/src/github.com/GoogleCloudPlatform
ln -s $PWD/$PROGNAME $GOPATH/src/$GODIR

cd $GOPATH/src/$GODIR
echo "Gopath is: " $GOPATH
echo "pwd is: " $PWD
ls -lah

