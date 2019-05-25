#!/bin/sh

# Build progressive web app for UI
(cd www && npm run build)

# Bundle assets.go for Lorca
go generate

echo "Now run the local build command for your operating system"

# xgo --targets=linux/amd64,darwin/amd64,windows/amd64 -branch master -dest build/ github.com/mcdafydd/omw
