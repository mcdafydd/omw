//go:generate go run -tags generate gen.go

package main

import (
	"log"

	"github.com/mcdafydd/omw/cmd"
)

func main() {
	cmd.Execute()
	log.Println("exiting...")
}


