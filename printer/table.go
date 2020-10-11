package printer

import (
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
)

// printTableFormat prints a line to w in kubectl "table" Format.
// Table format is something like:
// ---------------------------------------------------------
// NAME                     READY   STATUS    RESTARTS   AGE
// nginx-6799fc88d8-dnmv5   1/1     Running   0          31h
// nginx-6799fc88d8-m8pbc   1/1     Running   0          31h
// nginx-6799fc88d8-qdf9b   1/1     Running   0          31h
// nginx-8spn9              1/1     Running   0          31h
// nginx-dplns              1/1     Running   0          31h
// nginx-lpv5x              1/1     Running   0          31h
// ---------------------------------------------------------
// This function requires a line and tries to colorize it by each column.
// If dark is true, use colors which are readable in dark-backgrounded environment, else,
// it uses colors for light-backgrounded env.
// This function doesn't respect if the line is "header", so
// if you want to specify a special color for header, you must not pass the line
// to this function.
// deciderFn is a function to return context-specific color to be used to decorate a column.
// If the function returned ok=true, then returned color will be used for the column.
// If it returned ok=false, then default configurated color will be used.
// If deciderFn is null, then this function uses the default configurated color.
func printLineAsTableFormat(w io.Writer, line string, dark bool, deciderFn func(index int, column string) (color.Color, bool)) {
	columns := spaces.Split(line, -1)
	spacesIndices := spaces.FindAllStringIndex(line, -1)

	if len(columns) == len(spacesIndices)-1 {
		// It should not come here.
		panic("kubecolor: unexpected format as table. this must be a bug of kubecolor")
	}

	for i, column := range columns {
		c := decideColorForTable(i, column, dark)
		if deciderFn != nil {
			if cc, ok := deciderFn(i, column); ok {
				c = cc // prior injected deciderFn result
			}
		}
		// Write colored column
		fmt.Fprintf(w, "%s", color.Apply(column, c))
		// Write spaces based on actual output
		// When writing the most left column, no extra spaces needed.
		if i <= len(spacesIndices)-1 {
			spacesIndex := spacesIndices[i]
			fmt.Fprintf(w, "%s", toSpaces(spacesIndex[1]-spacesIndex[0]))
		}
	}

	fmt.Fprintf(w, "\n")
}

func decideColorForTable(index int, column string, dark bool) color.Color {
	colors := color.GetColors(dark)
	if index >= len(colors) {
		index = index % len(colors)
	}

	return colors[index]
}
