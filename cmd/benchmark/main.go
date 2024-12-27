package main

import (
	"fmt"
	"os"

	"github.com/ingo-eichhorst/arch-bench/internal/adapters/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
