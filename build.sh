#!/bin/bash

echo "Compiling scope"

mkdir -p ./tmp

cd ..
CGO_ENABLED=1 GOPATH=`pwd`/go GOARCH=arm GOARM=7 CXX=arm-linux-gnueabihf-g++ PKG_CONFIG_LIBDIR=/usr/lib/arm-linux-gnueabihf/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig CC=arm-linux-gnueabihf-gcc go build -o ./build/tmp/falcon/falcon.bhdouglass_falcon -ldflags '-extld=arm-linux-gnueabihf-g++' ./src
cd build

echo "Moving files into place"

cp ../click/manifest.json . # Allow clickable to find the manifest easily
cp ../click/* ./tmp/
cp ../images/* ./tmp/falcon
cp ../src/*.ini ./tmp/falcon
