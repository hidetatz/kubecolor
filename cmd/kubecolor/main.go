package main

import (
	"fmt"
	"os"

	"github.com/dty1er/kubecolor/command"
)

func main() {
	err := command.Run(os.Args[1:], false)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
