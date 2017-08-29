#!/usr/bin/env bash
set -e
echo "" > coverage.txt
# go test -coverprofile=profile.out -covermode=atomic .
# if [ -f profile.out ]; then
#     cat profile.out >> coverage.txt
#     rm profile.out
# fi
go test -coverprofile=profile.out -covermode=atomic ./shape
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./oakerr
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./timing
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./physics
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
go test -coverprofile=profile.out -covermode=atomic ./event
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
go test -coverprofile=profile.out -covermode=atomic ./collision
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