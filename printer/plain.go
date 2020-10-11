package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dty1er/kubecolor/color"
)

func PrintPlain(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", scanner.Text())
	}
}

func PrintWithColor(r io.Reader, w io.Writer, c color.Color) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n", color.Apply(scanner.Text(), c))
	}
}

func PrintErrorOrWarning(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.ToLower(line), "error") {
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Red))
		} else {
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Yellow))
		}
	}
}
