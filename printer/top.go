package printer

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/formattedwriter"
)

var colorsTop = []color.Color{color.Green, color.Blue, color.Magenta, color.Cyan}

// PrintTop prints the output of kubectl top command.
func PrintTop(output []byte, withHeader bool) {
	if string(output) == "" {
		return
	}

	w := formattedwriter.New(os.Stdout)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	index := 0
	isHeader := func() bool { return index == 0 }
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
				if i >= len(colorsTop) {
					i = i % len(colorsTop)
				}
				c = colorsTop[i]
			}

			colorized := color.Apply(column, c)
			result = append(result, colorized)
		}

		fmt.Fprintf(w, "%+v\n", strings.Join(result, "\t"))
		index++
	}
	w.Flush()
}
