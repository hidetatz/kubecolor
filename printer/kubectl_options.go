package printer

import (
	"bufio"
	"io"
)

type OptionsPrinter struct {
	DarkBackground bool
}

func (op *OptionsPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		_ = line
	}
}
