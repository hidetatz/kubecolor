package printer

import (
	"regexp"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

var (
	KeyColor    = color.White
	StringColor = color.Green
	BoolColor   = color.Cyan
	NumberColor = color.Blue

	// for json
	NullColor = color.Yellow

	// for yaml
	AnchorColor = color.Magenta
	AliasColor  = color.Yellow

	// for plain table
	HeaderColor = color.White
)

var tab = regexp.MustCompile("\\s{2,}")

func Print(output []byte, subcommandInfo *kubectl.SubcommandInfo) {
	withHeader := !subcommandInfo.NoHeader
	switch subcommandInfo.Subcommand.String() {
	case "top":
		// kubectl top supports only pod or node
		if subcommandInfo.Target != kubectl.Pod && subcommandInfo.Target != kubectl.Node {
			PrintPlain(output)
			return
		}

		PrintTop(output, withHeader)

	case "get":
		PrintGet(output, subcommandInfo.Target, withHeader, subcommandInfo.FormatOption)

	case "describe":
		// TODO implement
		PrintPlain(output)

	default:
		PrintPlain(output)
	}
}
