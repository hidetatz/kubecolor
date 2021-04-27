package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dty1er/kubecolor/color"
)

type ApplyPrinter struct {
	DarkBackground bool
}

// kubectl apply
// deployment.apps/foo unchanged
// deployment.apps/bar created
// deployment.apps/quux configured
func (ap *ApplyPrinter) Print(r io.Reader, w io.Writer) {
	const (
		applyActionCreated    = "created"
		applyActionConfigured = "configured"
		applyActionUnchanged  = "unchanged"
	)

	colors := map[string]color.Color{
		applyActionCreated:    color.Green,
		applyActionConfigured: color.Yellow,
		applyActionUnchanged:  color.Magenta,
	}

	colorize := func(line, action string, wr io.Writer) {
		arg := strings.TrimSuffix(line, " "+action)
		fmt.Fprintf(w, "%s %s\n", arg, color.Apply(action, colors[action]))
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasSuffix(line, " "+applyActionCreated):
			colorize(line, applyActionCreated, w)
		case strings.HasSuffix(line, " "+applyActionConfigured):
			colorize(line, applyActionConfigured, w)
		case strings.HasSuffix(line, " "+applyActionUnchanged):
			colorize(line, applyActionUnchanged, w)
		default:
			fmt.Fprintf(w, "%s\n", color.Apply(line, color.Green))
		}
	}
}
