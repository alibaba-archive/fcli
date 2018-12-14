#!/bin/bash

files=( $(find . -name "*.go" | grep -v "^./vendor/") )

badFiles=()
for f in "${files[@]}"; do
    echo "---------------------- $f ----------------------"
    gofmt -s -l "$f"
    res="$?"
    if [ "$res" -ne 0 ]; then
        badFiles+=( "$f" )
    fi
done

echo "--------------------------------------------"
if [ ${#badFiles[@]} -eq 0 ]; then
    echo 'Congratulations!  All Go source files are properly formatted.'
else
    {
        echo "These files are not properly gofmt'd:"
        for f in "${badFiles[@]}"; do
            echo " - $f"
        done
        echo
        echo 'Please reformat the above files using "gofmt -s -w" and commit the result.'
        echo

    }
    false
fi
