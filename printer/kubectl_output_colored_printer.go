package printer

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

// KubectlOutputColoredPrinter is a printer to print data depending on
// which kubectl subcommand is executed.
type KubectlOutputColoredPrinter struct {
	SubcommandInfo    *kubectl.SubcommandInfo
	DarkBackground    bool
	Recursive         bool
	ObjFreshThreshold time.Duration
}

// Print reads r then write it to w, its format is based on kubectl subcommand.
// If given subcommand is not supported by the printer, it prints data in Green.
func (kp *KubectlOutputColoredPrinter) Print(r io.Reader, w io.Writer) {
	withHeader := !kp.SubcommandInfo.NoHeader

	var printer Printer = &SingleColoredPrinter{Color: color.Green} // default in green

	switch kp.SubcommandInfo.Subcommand {
	case kubectl.Top, kubectl.APIResources:
		printer = NewTablePrinter(withHeader, kp.DarkBackground, nil)

	case kubectl.APIVersions:
		printer = NewTablePrinter(false, kp.DarkBackground, nil) // api-versions always doesn't have header

	case kubectl.Get:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.None, kp.SubcommandInfo.FormatOption == kubectl.Wide:
			printer = NewTablePrinter(
				withHeader,
				kp.DarkBackground,
				func(_ int, column string) (color.Color, bool) {
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

					// Object age when fresh then green
					if checkIfObjFresh(column, kp.ObjFreshThreshold) {
						return color.Green, true
					}

					return 0, false
				},
			)
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{DarkBackground: kp.DarkBackground}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{DarkBackground: kp.DarkBackground}
		}

	case kubectl.Describe:
		printer = &DescribePrinter{
			DarkBackground: kp.DarkBackground,
			TablePrinter:   NewTablePrinter(false, kp.DarkBackground, nil),
		}
	case kubectl.Explain:
		printer = &ExplainPrinter{
			DarkBackground: kp.DarkBackground,
			Recursive:      kp.Recursive,
		}
	case kubectl.Version:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{DarkBackground: kp.DarkBackground}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{DarkBackground: kp.DarkBackground}
		case kp.SubcommandInfo.Short:
			printer = &VersionShortPrinter{
				DarkBackground: kp.DarkBackground,
			}
		default:
			printer = &VersionPrinter{
				DarkBackground: kp.DarkBackground,
			}
		}
	case kubectl.Options:
		printer = &OptionsPrinter{
			DarkBackground: kp.DarkBackground,
		}
	case kubectl.Apply:
		switch {
		case kp.SubcommandInfo.FormatOption == kubectl.Json:
			printer = &JsonPrinter{DarkBackground: kp.DarkBackground}
		case kp.SubcommandInfo.FormatOption == kubectl.Yaml:
			printer = &YamlPrinter{DarkBackground: kp.DarkBackground}
		default:
			printer = &ApplyPrinter{DarkBackground: kp.DarkBackground}
		}
	}

	if kp.SubcommandInfo.Help {
		printer = &SingleColoredPrinter{Color: color.Yellow}
	}

	printer.Print(r, w)
}

func checkIfObjFresh(value string, threshold time.Duration) bool {
	// decode HumanDuration from k8s.io/apimachinery/pkg/util/duration
	durationRegex := regexp.MustCompile(`^(?P<years>\d+y)?(?P<days>\d+d)?(?P<hours>\d+h)?(?P<minutes>\d+m)?(?P<seconds>\d+s)?$`)
	matches := durationRegex.FindStringSubmatch(value)
	if len(matches) > 0 {
		years := parseInt64(matches[1])
		days := parseInt64(matches[2])
		hours := parseInt64(matches[3])
		minutes := parseInt64(matches[4])
		seconds := parseInt64(matches[5])
		objAgeSeconds := years*365*24*3600 + days*24*3600 + hours*3600 + minutes*60 + seconds
		objAgeDuration, err := time.ParseDuration(fmt.Sprintf("%ds", objAgeSeconds))
		if err != nil {
			return false
		}
		if objAgeDuration < threshold {
			return true
		}
	}
	return false
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}
