package printer

import (
	"io"
	"regexp"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

var (
	KeyColor    = color.White
	StringColor = color.Cyan
	BoolColor   = color.Green
	NumberColor = color.Magenta

	// for json
	NullColor = color.Yellow

	// for yaml
	AnchorColor = color.Magenta
	AliasColor  = color.Yellow

	// for plain table
	HeaderColor = color.White
)

var tab = regexp.MustCompile("\\s{2,}")

func Print(outReader io.Reader, writer io.Writer, subcommandInfo *kubectl.SubcommandInfo) {
	withHeader := !subcommandInfo.NoHeader
	switch subcommandInfo.Subcommand.String() {
	case "top":
		// kubectl top supports only pod or node
		// if subcommandInfo.Target != kubectl.Pod && subcommandInfo.Target != kubectl.Node {
		// 	PrintPlain(outReader, writer)
		// 	return
		// }

		printer := &TopPrinter{Writer: writer, WithHeader: withHeader}
		printer.Print(outReader)

	case "get":
		printer := &GetPrinter{Writer: writer, WithHeader: withHeader, FormatOpt: subcommandInfo.FormatOption}
		printer.Print(outReader)

	case "describe":
		printer := &DescribePrinter{Writer: writer}
		printer.Print(outReader)

	default:
		PrintPlain(outReader, writer)
	}
}

func DecideColor(column string, i int, palette []color.Color, decider func(column string) (color.Color, bool)) color.Color {
	if c, decided := decider(column); decided {
		return c
	}

	if i >= len(palette) {
		i = i % len(palette)
	}

	return palette[i]
}
