//go:generate statik -src=./www/dist

package main

import (
	"github.com/mcdafydd/omw/cmd"
)

func main() {
	cmd.Execute()
}
