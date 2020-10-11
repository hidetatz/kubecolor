package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
)

type TopPrinter struct {
	Writer     io.Writer
	WithHeader bool

	isFirstLine bool
}

func (tp *TopPrinter) Print(outReader io.Reader) {
	tp.isFirstLine = true
	scanner := bufio.NewScanner(outReader)
	w := tp.Writer
	for scanner.Scan() {
		line := scanner.Text()
		columns := tab.Split(line, -1)
		spacesIndices := tab.FindAllStringIndex(line, -1)

		if len(columns) == len(spacesIndices)-1 {
			// It should not come here.
			panic("unexpected format of get. this must be a bug of kubecolor")
		}

		for i, column := range columns {
			c := tp.DecideColor(i, column)
			// Write colored column
			fmt.Fprintf(w, "%s", color.Apply(column, c))
			// Write spaces based on actual output
			// When writing the most left column, no extra spaces needed.
			if i <= len(spacesIndices)-1 {
				spacesIndex := spacesIndices[i]
				fmt.Fprintf(w, "%s", tp.toSpaces(spacesIndex[1]-spacesIndex[0]))
			}
		}

		fmt.Fprintf(w, "\n")

		if tp.isFirstLine {
			tp.isFirstLine = false
		}
	}
}

func (tp *TopPrinter) isHeader() bool {
	return tp.WithHeader && tp.isFirstLine
}

func (tp *TopPrinter) toSpaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

func (tp *TopPrinter) Palette() []color.Color {
	return []color.Color{color.Cyan, color.Magenta, color.Green, color.White, color.Blue}
}

func (tp *TopPrinter) DecideColor(index int, column string) color.Color {
	if tp.isHeader() {
		return HeaderColor
	}

	colors := tp.Palette()
	if index >= len(colors) {
		index = index % len(colors)
	}

	return colors[index]
}
