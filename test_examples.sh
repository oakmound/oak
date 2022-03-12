#!/usr/bin/env bash

examples=$(find ./examples | grep main.go$)
for ex in $examples 
do
    echo $ex
    dir=$(dirname $ex)
    cd $dir
    timeout 10 go run $(basename $ex)
    retVal=$?
    echo "exit status" $retVal
    if [ $retVal -ne 124 ]; then
        exit 1
    fi
    cd ../..
done