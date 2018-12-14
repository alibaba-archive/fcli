#!/bin/bash

TEST_PATH=./...

pkg_list=$(go list -e ${TEST_PATH})

for pkg in $pkg_list
do
    echo "------------------${pkg}------------------"
    go test -v -coverprofile cover.out $pkg
done