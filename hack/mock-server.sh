#!/bin/bash

npm install -g anyproxy

mkdir -p ~/.fcli

cp ./.config.yaml ~/.fcli/config.yaml

anyproxy  --rule ./hack/integrity-test/rule.js  &
