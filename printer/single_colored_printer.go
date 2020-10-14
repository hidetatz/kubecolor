package printer

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
)

type SingleColoredPrinter struct {
	Color color.Color
}

func (sp *SingleColoredPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", color.Apply(scanner.Text(), sp.Color))
	}
}
