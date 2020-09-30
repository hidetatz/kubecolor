package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/formattedwriter"
	"github.com/dty1er/kubecolor/kubectl"
)

var colorsGet = []color.Color{color.Blue, color.Magenta, color.Cyan}

func PrintGet(output []byte, target kubectl.Target, withHeader bool, formatOpt kubectl.FormatOption) {
	if string(output) == "" {
		return
	}

	if formatOpt == kubectl.Json {
		err := printJson(os.Stdout, output)
		if err != nil {
			PrintPlain(output)
		}
		return
	}

	if formatOpt == kubectl.Yaml {
		err := printYaml(os.Stdout, output)
		if err != nil {
			PrintPlain(output)
		}
		return
	}

	switch target {
	case kubectl.Pod:
		printGetPod(output, withHeader)
	}
}

func printGetPod(output []byte, withHeader bool) {
	w := formattedwriter.New(os.Stdout)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	index := 0
	isHeader := func() bool { return withHeader && index == 0 }

	for scanner.Scan() {
		line := scanner.Text()
		columns := tab.Split(line, -1)

		result := []string{}
		for i, column := range columns {
			// only header line
			var c color.Color
			if isHeader() {
				c = HeaderColor
			} else {
				if i >= len(colorsGet) {
					i = i % len(colorsGet)
				}
				c = colorsGet[i]
			}

			colorized := Colorize(column, c)
			result = append(result, colorized)
		}

		fmt.Fprintf(w, "%+v\n", strings.Join(result, "\t"))
		index++
	}
	w.Flush()
}

func Colorize(msg string, defaultColor color.Color) string {
	// for status of something
	switch msg {
	case "Running", "Succeeded", "Completed":
		return color.Apply(msg, color.Green)
	case "CrashLoopBackOff":
		return color.Apply(msg, color.Red)

		// more status?
	}

	// When Readiness is "n/n" then green
	// When Readiness is "n/m" then yellow
	if strings.Count(msg, "/") == 1 {
		arr := strings.Split(msg, "/")
		ready := arr[0]
		total := arr[1]
		if ready != total {
			return color.Apply(msg, color.Yellow)
		} else {
			return color.Apply(msg, color.Green)
		}
	}

	return color.Apply(msg, defaultColor)
}
