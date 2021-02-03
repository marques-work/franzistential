package main

import (
	"github.com/marques-work/franzistential/cli"
	"github.com/marques-work/franzistential/logging"
)

func main() {
	if err := cli.Execute(); err != nil {
		logging.Die("Failed with error: %v", err)
	}
}
