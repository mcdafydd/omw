//go:generate go run -tags generate gen.go

package main

import (
	"github.com/mcdafydd/omw/cmd"
)

func main() {
	cmd.Execute()
}
