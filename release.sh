#!/bin/bash

version="$1"

[ -z "$version" ] && echo "Usage: ./release.sh <version>" && exit 1;

mkdir -p bin
rm -fr bin/*

build() {
	echo "building OS=$1 ARCH=$2"
	GOOS=$1 GOARCH=$2 go build -o bin/mdj-diff-$1-$2-$version
}

build darwin arm64
build darwin amd64
build linux amd64
build windows amd64
