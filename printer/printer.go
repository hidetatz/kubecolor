package printer

import (
	"io"
	"regexp"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

var (
	KeyColor    = color.White
	StringColor = color.Cyan
	BoolColor   = color.Green
	NumberColor = color.Magenta
	NullColor   = color.Yellow

	// for plain table
	HeaderColor = color.White
)

var spaces = regexp.MustCompile("\\s{2,}")

func Print(r io.Reader, w io.Writer, subcommandInfo *kubectl.SubcommandInfo, darkBackground bool) {
	withHeader := !subcommandInfo.NoHeader
	switch subcommandInfo.Subcommand {
	case kubectl.Top:
		printer := &TopPrinter{Writer: w, WithHeader: withHeader, DarkBackground: darkBackground}
		printer.Print(r)

	case kubectl.Get:
		printer := &GetPrinter{Writer: w, WithHeader: withHeader, FormatOpt: subcommandInfo.FormatOption, DarkBackground: darkBackground}
		printer.Print(r)

	case kubectl.Describe:
		printer := &DescribePrinter{Writer: w}
		printer.Print(r)

	default:
		PrintPlain(r, w)
	}
}

func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}
