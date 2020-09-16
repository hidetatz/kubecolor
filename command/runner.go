package command

import (
	"context"
	"fmt"
	"os"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/printer"
)

func Run(args []string) error {
	ctx := context.Background()
	output, err := kubectl.Execute(ctx, args)

	if err != nil {
		fmt.Fprintf(os.Stderr, color.Apply(err.Error(), color.Red))
		return nil
	}

	subcommandInfo, err := kubectl.InspectSubcommandInfo(args)
	if err != nil {
		printer.PrintPlain(output)
		return nil
	}

	printer.Print(output, subcommandInfo)
	return nil
}
