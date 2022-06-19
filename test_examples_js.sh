#!/usr/bin/env bash

set -e

examples=$(find ./examples | grep main.go$)
root=$(pwd)
for ex in $examples
do
    echo "$ex"
    dir=$(dirname "$ex")
    # excluded because screenopts explicitly demonstrates
    # desktop-specific features
    if [[ "$dir" == "./examples/screenopts" ]]; then
      continue
    fi
    # excluded because text includes a find-font dependency
    # that does not compile in js
    if [[ "$dir" == "./examples/text" ]]; then
      continue
    fi
    cd "$dir"
    GOOS=js GOARCH=wasm go build .
    cd "$root"
done