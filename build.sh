#!/bin/sh

# Build progressive web app for UI
(cd www && npm run build)

# Bundle assets.go for Lorca
go generate

echo "Now run the local build command for your operating system"