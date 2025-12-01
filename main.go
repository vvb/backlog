package main

import (
	"fmt"
	"os"

	"github.com/vvb/backlog/cmd"
)

func main() {
	// If no arguments are provided, default to "backlog list -i" (interactive list view)
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "list", "-i")
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
