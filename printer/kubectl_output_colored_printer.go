package printer

import (
	"io"
	"strconv"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

// KubectlOutputColoredPrinter is a printer to print data depending on
// which kubectl subcommand is executed.
type KubectlOutputColoredPrinter struct {
	SubcommandInfo *kubectl.SubcommandInfo
	DarkBackground bool
}

// Print reads r then write it to w, its format is based on kubectl subcommand.
// If given subcommand is not supported by the printer, it prints data in Green.
func (kp *KubectlOutputColoredPrinter) Print(r io.Reader, w io.Writer) {
	withHeader := !kp.SubcommandInfo.NoHeader

	var printer Printer = &SingleColoredPrinter{Color: color.Green} // default in green

	switch kp.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		printer = &TablePrinter{WithHeader: withHeader, DarkBackground: kp.DarkBackground}

	case kubectl.Get:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.None, kp.SubcommandInfo.FormatOption == kubectl.Wide:
			printer = &TablePrinter{
				WithHeader:     withHeader,
				DarkBackground: kp.DarkBackground,
				ColorDeciderFn: func(_ int, column string) (color.Color, bool) {
					if column == "CrashLoopBackOff" {
						return color.Red, true
					}

					// When Readiness is "n/m" then yellow
					if strings.Count(column, "/") == 1 {
						if arr := strings.Split(column, "/"); arr[0] != arr[1] {
							_, e1 := strconv.Atoi(arr[0])
							_, e2 := strconv.Atoi(arr[1])
							if e1 == nil && e2 == nil { // check both is number
								return color.Yellow, true
							}
						}

					}

					return 0, false
				},
			}
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{DarkBackground: kp.DarkBackground}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{DarkBackground: kp.DarkBackground}
		}

	case kubectl.Describe:
		printer = &DescribePrinter{DarkBackground: kp.DarkBackground}
	}

	printer.Print(r, w)
}
