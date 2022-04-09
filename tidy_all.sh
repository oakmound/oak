#!/usr/bin/env bash

go mod tidy 

cd examples/fallback-font 
go mod tidy 

cd ../clipboard 
go mod tidy 

cd ../svg
go mod tidy 