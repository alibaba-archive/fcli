#!/usr/bin/env bash
set -e

version=`cat version/VERSION`
export VERSION=${version}

bash hack/.go-autogen

# build binary.
xgo -go 1.8.3 -out bundles/fcli --targets="linux/amd64,darwin/amd64,windows/386" $GOPATH/src/github.com/aliyun/fcli

# package binary.
cd bundles
cp fcli-darwin* fcli && zip fcli-v${version}-darwin-amd64.zip fcli && rm -rf fcli
cp fcli-linux* fcli && zip fcli-v${version}-linux-amd64.zip fcli && rm -rf fcli
cp fcli-windows* fcli.exe && zip fcli-v${version}-win-386.zip fcli.exe && rm -rf fcli.exe
