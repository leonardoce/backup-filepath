// Package main is the entrypoint of the application
package main

import (
	"fmt"
	"os"

	"github.com/leonardoce/backup-filepath/cmd/injector"
)

func main() {
	err := injector.Cmd().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
