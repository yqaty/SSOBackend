package main

import (
	"os"

	"github.com/UniqueStudio/UniqueSSOBackend/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
