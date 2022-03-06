#!/usr/bin/env bash

echo "" > coverage.txt 

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
