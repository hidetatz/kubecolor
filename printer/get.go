package printer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/formattedwriter"
	"github.com/dty1er/kubecolor/kubectl"
)

type GetPrinter struct {
	Writer io.Writer
	// Target     kubectl.Target
	WithHeader bool
	FormatOpt  kubectl.FormatOption

	isFirstLine bool
}

func (gp *GetPrinter) Print(outReader io.Reader) {
	w := formattedwriter.New(gp.Writer)

	gp.isFirstLine = true
	scanner := bufio.NewScanner(outReader)
	lines := []byte{}
	for scanner.Scan() {
		line := scanner.Text()

		if gp.FormatOpt == kubectl.Json || gp.FormatOpt == kubectl.Yaml {
			// kubectl get can specify json or yaml. Json and Yaml Printers require complete Json or Yaml so
			// make a buffer to store all the output
			lines = append(lines, []byte(line+"\n")...)
			// in this case, print will be done outside of scan
			continue
		}

		columns := tab.Split(line, -1)
		result := []string{}
		for i, column := range columns {
			c := DecideColor(column, i, gp.Palette(), gp.DecideColor)
			result = append(result, color.Apply(column, c))
		}

		fmt.Fprintf(w, "%+v\n", strings.Join(result, "\t"))

		if gp.isFirstLine {
			gp.isFirstLine = false
		}
	}

	switch gp.FormatOpt {
	case kubectl.Json:
		printJson(os.Stdout, lines)
	case kubectl.Yaml:
		printYaml(os.Stdout, lines)
	default:
		w.Flush()
	}
}

func (gp *GetPrinter) isHeader() bool {
	return gp.WithHeader && gp.isFirstLine
}

func (gp *GetPrinter) Palette() []color.Color {
	return []color.Color{color.Cyan, color.Magenta, color.Green, color.White, color.Blue}
}

func (gp *GetPrinter) DecideColor(column string) (color.Color, bool) {
	if gp.isHeader() {
		return HeaderColor, true
	}

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

	return color.Color(0), false
}
