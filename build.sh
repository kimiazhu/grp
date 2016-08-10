#!/usr/bin/env bash

rm -rf "bin/grp_linux"
export set GOOS=linux
export set GOARCH=amd64
go build -o bin/grp_linux

rm -rf "bin/grp_darwin"
export set GOOS=darwin
export set GOARCH=amd64
go build -o bin/grp_darwin

rm -rf "bin/grp.exe"
export set GOOS=windows
export set GOARCH=amd64
go build -o bin/grp.exe
