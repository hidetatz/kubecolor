package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
)

type WithFuncPrinter struct {
	Fn func(line string) color.Color
}

func (wp *WithFuncPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c := wp.Fn(line)
		fmt.Fprintf(w, "%s\n", color.Apply(line, c))
	}
}
