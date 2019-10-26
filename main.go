//go:generate statik -src=./www/build/es5-bundled

package main

import (
	"github.com/mcdafydd/omw/cmd"
)

func main() {
	cmd.Execute()
}
