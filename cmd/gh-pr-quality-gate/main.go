package main

import (
	"fmt"
	"os"

	"github.com/Ysoseri1224/gh-pr-quality-gate/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr, version); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
