#!/bin/sh

(cd www && npm run build)
go generate

./build-linux.sh
./build-macos.sh
./build-windows.bat
