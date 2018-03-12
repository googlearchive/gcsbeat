package main

import (
	"os"

	"github.com/GoogleCloudPlatform/gcsbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
