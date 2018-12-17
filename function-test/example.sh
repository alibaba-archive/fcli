#!/bin/bash

expect_res='fcli: function compute command line tools

Usage:
  fcli [flags]
  fcli [command]

Available Commands:
  alias           alias related operation
  config          Configure the fcli
  function        function related operation
  help            Help about any command
  service         service related operation
  shell           interactive shell
  sls             sls related operation
  trigger         trigger related operation
  version         fcli version information

Flags:
  -h, --help   help for fcli

Use "fcli [command] --help" for more information about a command.'

res=`go run main.go`


if [ "$res" = "$expect_res" ]; then
    echo 'Congratulations!  example.sh passed.'
else
    echo 'example.sh fail'
    false
fi