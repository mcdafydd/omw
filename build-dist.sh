#!/bin/bash

# Uses xgo to build binaries for windows, mac, and linux amd64
go get -u github.com/karalabe/xgo
docker build -t omw-build/latest .
xgo -image outofmyway-build/latest -targets=darwin/amd64,windows/amd64 -branch master -dest build/ -v -x -race github.com/mcdafydd/omw

