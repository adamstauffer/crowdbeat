package main

import (
	"os"

	"github.com/adamstauffer/crowdbeat/cmd"

	_ "github.com/adamstauffer/crowdbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
