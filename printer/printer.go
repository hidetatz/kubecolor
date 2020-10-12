package printer

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

var (
	// Preset of colors for background
	// Please use them when you just need random colors
	colorsForDarkBackground = []color.Color{
		color.Cyan,
		color.Green,
		color.Magenta,
		color.White,
		color.Yellow,
	}

	colorsForLightBackground = []color.Color{
		color.Cyan,
		color.Green,
		color.Magenta,
		color.Black,
		color.Yellow,
		color.Blue,
	}

	// colors to be recommended to be used for some context
	// e.g. Json, Yaml, kubectl-describe format etc.

	// colors which look good in dark-backgrounded environment
	KeyColorForDark    = color.White
	StringColorForDark = color.Cyan
	BoolColorForDark   = color.Green
	NumberColorForDark = color.Magenta
	NullColorForDark   = color.Yellow
	HeaderColorForDark = color.White // for plain table

	// colors which look good in light-backgrounded environment
	KeyColorForLight    = color.Black
	StringColorForLight = color.Blue
	BoolColorForLight   = color.Green
	NumberColorForLight = color.Magenta
	NullColorForLight   = color.Yellow
	HeaderColorForLight = color.Black // for plain table
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

func colorByValue(val string, dark bool) color.Color {
	if val == "null" || val == "<none>" || val == "<unknown>" {
		return NullColor
	}

	if val == "true" || val == "false" {
		return BoolColor
	}

	if _, err := strconv.Atoi(val); err == nil {
		return NumberColor
	}

	return StringColor
}
