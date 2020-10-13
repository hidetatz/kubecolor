package printer

import (
	"bufio"
	"fmt"
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

// Print reads data in r then colorize them and write it in w.
// subcommandInfo is necessary because this function decides how to decorate the information by it.
// When darkBackground is true, colors which looks readable in darkBackground.
// For example, When darkBackground is true, this function will never colorized the data in black.
// If this function does not know how to colorize the data, it just prints them out in plain.
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
		printer := &DescribePrinter{Writer: w, DarkBackground: darkBackground}
		printer.Print(r)

	default:
		PrintPlain(r, w)
	}
}

// PrintPlaing reads r then writes it to w without any decorations.
func PrintPlain(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", scanner.Text())
	}
}

// PrintWithColor reads r then writes it to w in given color.
func PrintWithColor(r io.Reader, w io.Writer, c color.Color) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", color.Apply(scanner.Text(), c))
	}
}

// PrintErrorOrWarning reads r then writes it to w in error or warning color.
// if the line has "error" prefix, it will be error color, otherwise it will be printed in warning color.
func PrintErrorOrWarning(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.ToLower(line), "error") {
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Red))
		} else {
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Yellow))
		}
	}
}

// toSpaces returns repeated spaces whose length is n.
func toSpaces(n int) string {
	return strings.Repeat(" ", n)
}

// getColorByKeyIndent returns a color based on the given indent.
// When you want to change key color based on indent depth (e.g. Json, Yaml), use this function
func getColorByKeyIndent(indent int, basicIndentWidth int, dark bool) color.Color {
	switch indent / basicIndentWidth % 2 {
	case 1:
		if dark {
			return color.White
		}
		return color.Black
	default:
		return color.Yellow
	}
}

// getColorByValueType returns a color by value.
// This is intended to be used to colorize any structured data e.g. Json, Yaml.
func getColorByValueType(val string, dark bool) color.Color {
	if val == "null" || val == "<none>" || val == "<unknown>" {
		if dark {
			return NullColorForDark
		}
		return NullColorForLight
	}

	if val == "true" || val == "false" {
		if dark {
			return BoolColorForDark
		}
		return BoolColorForLight
	}

	if _, err := strconv.Atoi(val); err == nil {
		if dark {
			return NumberColorForDark
		}
		return NumberColorForLight
	}

	if dark {
		return StringColorForDark
	}
	return StringColorForLight
}

// getColorsByBackground returns a preset of colors depending on given background color
func getColorsByBackground(dark bool) []color.Color {
	if dark {
		return colorsForDarkBackground
	}

	return colorsForLightBackground
}

// getHeaderColorByBackground returns a defined color for Header (not actual data) by the background color
func getHeaderColorByBackground(dark bool) color.Color {
	if dark {
		return HeaderColorForDark
	}

	return HeaderColorForLight
}

// findIndent returns a length of indent (spaces at left) in the given line
func findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}
