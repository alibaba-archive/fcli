#!/bin/bash

files=( $(find ./hack/integrity-test -name "*.sh") )

for f in "${files[@]}"; do
    bash ${f}
done