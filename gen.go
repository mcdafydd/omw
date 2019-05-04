//+build generate

package main

import "github.com/zserge/lorca"

func main() {
	// generate manifests, or do other preparations for your assets.
	lorca.Embed("main", "./assets.go", "www/build/es6-bundled")
}
