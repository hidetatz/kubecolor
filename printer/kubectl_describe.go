package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dty1er/kubecolor/color"
)

// DescribePrinter is a specific printer to print kubectl describe format.
type DescribePrinter struct {
	DarkBackground bool
}

func (dp *DescribePrinter) Print(r io.Reader, w io.Writer) {
	basicIndentWidth := 2 // according to kubectl describe format
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		// Split a line by spaces to colorize and render them
		// For example:
		// e.g. 1-----------------
		// Status:         Running
		// -----------------------
		// spacesIndices: [[7, 15]] // <- where spaces locate
		// columns: ["Status:", "Running"]
		//
		// e.g. 2--------------------------------------------
		//     Ports:          10001/TCP, 5000/TCP, 18000/TCP
		// --------------------------------------------------
		// spacesIndices: [[0, 3], [10, 19]] // <- where spaces locate
		// columns: ["Ports:", "10001/TCP, 5000/TCP, 18000/TCP"]
		//
		// So now, we know where to render which column.
		spacesIndices := spaces.FindAllStringIndex(line, -1)
		columns := spaces.Split(line, -1)
		// when the line has indent (spaces on left), the first item will be
		// just a "" and we don't need it so remove
		if len(columns) > 0 {
			if columns[0] == "" {
				columns = columns[1:]
			}
		}

		// First, identify if there is an indent
		indentCnt := findIndent(line)
		indent := toSpaces(indentCnt)
		if indentCnt > 0 {
			// when an indent exists, removes it because it's already captured by "indent" var
			spacesIndices = spacesIndices[1:]
		}

		// First, write the first value assuming it's a key
		keyColor := getColorByKeyIndent(indentCnt, basicIndentWidth, dp.DarkBackground)
		if strings.HasSuffix(columns[0], ":") {
			// trailing ":" should not be colorized
			fmt.Fprintf(w, "%s%s:", indent, color.Apply(strings.TrimRight(columns[0], ":"), keyColor))
		} else {
			fmt.Fprintf(w, "%s%s", indent, color.Apply(columns[0], keyColor))
		}

		if len(columns) == 1 {
			fmt.Fprint(w, "\n")
			continue
		}

		// Then, write values
		// In this for loop, we write spaces and a value
		// For example, if the line looked like:
		// e.g.--------------------------------------
		//     Key:<- spaces1 ->Value1<- spaces2 ->Value2<- spaces3 ->Value3
		// ------------------------------------------
		// In each loop, we write spaces+value
		// So, each iteration will look like
		// when i == 0 then continue // because indent and key is already written
		// when i == 1 then Write spaces1 and Value1
		// when i == 2 then Write spaces2 and Value1
		// when i == 3 then Write spaces3 and Value1
		for i, column := range columns {
			if i == 0 {
				continue
			}

			spacesPos := spacesIndices[i-1]
			spacesCnt := spacesPos[1] - spacesPos[0]
			fmt.Fprintf(w, "%s%s", toSpaces(spacesCnt), color.Apply(column, getColorByValueType(column, dp.DarkBackground)))
		}

		fmt.Fprint(w, "\n")
	}
}
