package main

import (
	"errors"
	"os"

	"github.com/dty1er/kubecolor/command"
)

func main() {
	err := command.Run(os.Args[1:])
	if err != nil {
		var ke *command.KubectlError
		if errors.As(err, &ke) {
			os.Exit(ke.ExitCode)
		}
		os.Exit(1)
	}
}
