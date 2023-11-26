// Package main is the entrypoint of the application
package main

import (
	"fmt"
	"os"

	"github.com/leonardoce/backup-filepath/cmd/filepath_adapter"
)

func main() {
	err := filepath_adapter.Cmd().Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
