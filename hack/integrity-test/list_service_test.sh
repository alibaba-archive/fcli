#!/bin/bash

expect="{
  \"Services\": [
    \"demo\"
  ],
  \"NextToken\": null
}"

result=$(go run main.go service list)

if ! [ "${expect}" == "${result}" ]; then
    echo
    echo "list_service_test error: "
    echo "expect: ${expect}"
    echo "result: ${result}"
    false
fi