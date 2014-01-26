#!/bin/bash
TARGET="/opt/dtop/"
pushd ../
export GOPATH=`pwd`
go build eu.dominiek/dtop
mkdir -p $TARGET
cp -rf ./static $TARGET
mv dtop $TARGET
unset GOPATH
popd
