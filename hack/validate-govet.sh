#!/bin/bash

go vet .
res=$?

echo "--------------------------------------------"
if [ "${res}" -ne 0  ]; then
    echo
    echo 'Please fix the above errors. You can test via "go vet" and commit the result.'
    echo
    false
else
    echo 'Congratulations!  All Go source files have been vetted.'
fi
