package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
)

type TopPrinter struct {
	Writer         io.Writer
	WithHeader     bool
	DarkBackground bool

	isFirstLine bool
}

func (tp *TopPrinter) Print(outReader io.Reader) {
	tp.isFirstLine = true
	scanner := bufio.NewScanner(outReader)
	for scanner.Scan() {
		line := scanner.Text()
		if tp.isHeader() {
			fmt.Fprintf(tp.Writer, "%s\n", color.Apply(line, HeaderColor))
			tp.isFirstLine = false
			continue
		}

		printLineAsTableFormat(tp.Writer, line, tp.DarkBackground, nil)
	}
}

func (tp *TopPrinter) isHeader() bool {
	return tp.WithHeader && tp.isFirstLine
}
