package printer

import (
	"io"
	"regexp"
)

var spaces = regexp.MustCompile("\\s{2,}")

// Printer can print something
type Printer interface {
	Print(r io.Reader, w io.Writer)
}
