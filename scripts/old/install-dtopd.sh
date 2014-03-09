#!/bin/bash
TARGET="/opt/dtop/"

RED='\e[0;31m'
GREEN='\e[0;32m'
NC='\e[0m'
OK=" ->$GREEN OK$NC"
FAILED=" ->$RED FAILED$NC"

pushd ../ > /dev/null
export GOPATH=`pwd`
echo "building dtop ... "
go build eu.dominiek/dtop > /dev/null
BUILD_EXIT=$?
test 0 -eq $BUILD_EXIT && echo -e $OK
test 0 -ne $BUILD_EXIT && echo -e $FAILED && popd && exit 1
echo "copying to $TARGET ... "
mkdir -p $TARGET
cp -rf ./static $TARGET
mv dtop $TARGET
MV_EXIT=$?
unset GOPATH
test 0 -eq $BUILD_EXIT && echo -e $OK
test 0 -ne $BUILD_EXIT && echo -e $FAILED && popd && exit 1
popd > /dev/null

pushd ../scripts > /dev/null
echo "registering dtopd daemon ... "
sudo cp dtopd -f /etc/init.d/
sudo update-rc.d -f dtopd remove > /dev/null 
sudo update-rc.d dtopd defaults > /dev/null
sudo service dtopd start > /dev/null
START_EXIT=$?
test 0 -eq $BUILD_EXIT && echo -e $OK
test 0 -ne $BUILD_EXIT && echo -e $FAILED && popd && exit 1
popd > /dev/null