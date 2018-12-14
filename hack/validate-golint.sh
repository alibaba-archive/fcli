#!/bin/bash

files=( $(find . -name "*.go" | grep -v "^./vendor/") )

badFiles=()
for f in "${files[@]}"; do
    echo "---------------------- $f ----------------------"
    res=$( golint "$f" )
    if [ "$res" ]; then
        echo "$res"
        badFiles+=( "$f" )
    fi
done

echo "--------------------------------------------"
if [ ${#badFiles[@]} -eq 0 ]; then
    echo 'Congratulations!  All Go source files have been linted.'
else
    {
        echo "These files are not properly golint'd:"
        for f in "${badFiles[@]}"; do
            echo " - $f"
        done
        echo
        echo 'Please fix the above errors. You can test via "golint" and commit the result.'
        echo

    }
    false
fi
