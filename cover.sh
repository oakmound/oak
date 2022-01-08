#!/usr/bin/env bash

echo "" > coverage.txt

TAGS="${OAK_COVER_TAGS:-nooswindow}"

examples=$(find ./examples | grep main.go$)
for ex in $examples 
do
    echo $ex
    dir=$(dirname $ex)
    cd $dir
    timeout 10 go run --tags="${TAGS}" $(basename $ex)
    retVal=$?
    echo "exit status" $retVal
    if [ $retVal -ne 124 ]; then
        exit 1
    fi
    cd ../..
done

if [[ ! -z ${OAK_COVER_EXAMPLES_ONLY} ]]; then
    exit 0  
fi 

set -e
go test -coverprofile=profile.out -covermode=atomic ./alg
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./alg/intgeom
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./alg/floatgeom
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./collision
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./collision/ray
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./debugstream
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./dlog
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./event
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./fileutil
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./mouse
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic --tags=nooswindow .
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./oakerr
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./physics
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./render
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./render/mod
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./render/particle
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./scene
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./shape
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./dlog
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./timing
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
