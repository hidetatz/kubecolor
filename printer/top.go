package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/formattedwriter"
)

type TopPrinter struct {
	Writer     io.Writer
	WithHeader bool

	isFirstLine bool
}

func (tp *TopPrinter) Print(outReader io.Reader) {
	scanner := bufio.NewScanner(outReader)
	w := formattedwriter.New(tp.Writer)
	tp.isFirstLine = true
	for scanner.Scan() {
		columns := tab.Split(scanner.Text(), -1)
		result := []string{}

		for i, column := range columns {
			result = append(result, color.Apply(column, DecideColor(column, i, tp.Palette(), tp.DecideColor)))
		}

		fmt.Fprintf(w, "%+v\n", strings.Join(result, "\t"))
		if tp.isFirstLine {
			tp.isFirstLine = false
		}
	}

	w.Flush()
}

func (tp *TopPrinter) isHeader() bool {
	return tp.WithHeader && tp.isFirstLine
}

func (tp *TopPrinter) DecideColor(msg string) (color.Color, bool) {
	if tp.isHeader() {
		return HeaderColor, true
	}

	return color.Color(0), false
}

func (tp *TopPrinter) Palette() []color.Color {
	return []color.Color{color.Green, color.Magenta, color.Cyan, color.Blue, color.White, color.Yellow}
}
