package printer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
)

type GetPrinter struct {
	Writer     io.Writer
	WithHeader bool
	FormatOpt  kubectl.FormatOption

	isFirstLine bool
	inString    bool
}

func (gp *GetPrinter) Print(outReader io.Reader) {
	gp.isFirstLine = true
	scanner := bufio.NewScanner(outReader)
	for scanner.Scan() {
		line := scanner.Text()
		switch gp.FormatOpt {
		case kubectl.Json:
			gp.PrintJson(line)
		case kubectl.Yaml:
			gp.PrintYaml(line)
		default:
			gp.PrintTable(line)
		}
	}
}

func (gp *GetPrinter) PrintTable(line string) {
	w := gp.Writer
	columns := tab.Split(line, -1)
	spacesIndices := tab.FindAllStringIndex(line, -1)

	// The format of kubectl get is like
	// NAME     READY  STATUS
	// pod-a    1/1    Running
	// pod-b-2  1/1    Running
	// pod-c    1/1    Running
	// Spaces must locate between each column, so we validate it by this check
	if len(columns) == len(spacesIndices)-1 {
		// It should not come here.
		panic("unexpected format of get. this must be a bug of kubecolor")
	}

	for i, column := range columns {
		c := gp.DecideColor(i, column)
		// Write colored column
		fmt.Fprintf(w, "%s", color.Apply(column, c))
		// Write spaces based on actual output
		// When writing the most left column, no extra spaces needed.
		if i <= len(spacesIndices)-1 {
			spacesIndex := spacesIndices[i]
			fmt.Fprintf(w, "%s", gp.toSpaces(spacesIndex[1]-spacesIndex[0]))
		}
	}

	fmt.Fprintf(w, "\n")

	if gp.isFirstLine {
		gp.isFirstLine = false
	}
}

func (gp *GetPrinter) PrintJson(line string) {
	w := gp.Writer
	indentCnt := gp.findIndent(line)
	trimmedLine := strings.TrimLeft(line, " ")

	if strings.HasPrefix(trimmedLine, "{") ||
		strings.HasPrefix(trimmedLine, "}") ||
		strings.HasPrefix(trimmedLine, "]") {
		// when coming here, it must not be starting with key.
		// that patterns are:
		// {
		// }
		// },
		// ]
		// ],
		// note: it must not be "[" because it will be always after key
		// in this case, just write it without color
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, "%s", trimmedLine)
		fmt.Fprintf(w, "\n")
		return
	}

	// when coming here:
	// "key": {
	// "key": [
	// "key": value
	// "key": value,
	trimmed := strings.SplitN(trimmedLine, ": ", 2) // if key contains ": " this works in a wrong way but it's unlikely to happen

	if len(trimmed) == 1 {
		// when coming here, it will be a value in an array
		if strings.HasSuffix(trimmed[0], ",") {
			// when coming here, it must be `value,`
			ss := strings.TrimRight(trimmed[0], ",") // this is a value; it might be double-quoted or not
			if strings.HasPrefix(ss, `"`) && strings.HasSuffix(ss, `"`) {
				ss = strings.TrimLeft(ss, `"`)
				ss = strings.TrimRight(ss, `"`)
				fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
				fmt.Fprintf(w, `"`)
				fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByIndent(indentCnt)))
				fmt.Fprintf(w, `",`)
				fmt.Fprintf(w, "\n")
			} else {
				fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
				fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
				fmt.Fprintf(w, "\n")
			}
		} else {
			ss := trimmed[0]
			// when coming here, it must be `value`
			if strings.HasPrefix(ss, `"`) && strings.HasSuffix(ss, `"`) {
				ss = strings.TrimLeft(ss, `"`)
				ss = strings.TrimRight(ss, `"`)
				fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
				fmt.Fprintf(w, `"`)
				fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByIndent(indentCnt)))
				fmt.Fprintf(w, `"`)
				fmt.Fprintf(w, "\n")
			} else {
				fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
				fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
				fmt.Fprintf(w, "\n")
			}
		}
		return
	}

	key := trimmed[0]
	key = strings.TrimLeft(key, `"`)
	key = strings.TrimRight(key, `"`)

	if strings.HasSuffix(trimmedLine, "{") {
		// trim double quotation and colon, bracket
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, `"`)
		fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
		fmt.Fprintf(w, `": {`)
		fmt.Fprintf(w, "\n")
	} else if strings.HasSuffix(trimmedLine, "[") {
		// trim double quotation and colon, bracket
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, `"`)
		fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
		fmt.Fprintf(w, `": [`)
		fmt.Fprintf(w, "\n")
	} else if strings.HasSuffix(trimmed[1], ",") {
		// when coming here, it must be `"key": value,`
		ss := strings.TrimRight(trimmed[1], ",") // this is a value; it might be double-quoted or not
		if strings.HasPrefix(ss, `"`) && strings.HasSuffix(ss, `"`) {
			ss = strings.TrimLeft(ss, `"`)
			ss = strings.TrimRight(ss, `"`)
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, `"`)
			fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
			fmt.Fprintf(w, `": "`)
			fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
			fmt.Fprintf(w, `",`)
			fmt.Fprintf(w, "\n")
		} else {
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, `"`)
			fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
			fmt.Fprintf(w, `": `)
			fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
			fmt.Fprintf(w, `,`)
			fmt.Fprintf(w, "\n")
		}
	} else {
		// when coming here, it must be `"key": value`
		ss := trimmed[1]
		if strings.HasPrefix(ss, `"`) && strings.HasSuffix(ss, `"`) {
			ss = strings.TrimLeft(ss, `"`)
			ss = strings.TrimRight(ss, `"`)
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, `"`)
			fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
			fmt.Fprintf(w, `": "`)
			fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
			fmt.Fprintf(w, `"`)
			fmt.Fprintf(w, "\n")
		} else {
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, `"`)
			fmt.Fprintf(w, "%s", color.Apply(key, gp.colorByIndent(indentCnt)))
			fmt.Fprintf(w, `": `)
			fmt.Fprintf(w, "%s", color.Apply(ss, gp.colorByValue(ss)))
			fmt.Fprintf(w, "\n")
		}
	}
}

func (gp *GetPrinter) PrintYaml(line string) {
	w := gp.Writer
	indentCnt := gp.findIndent(line)
	trimmedLine := strings.TrimLeft(line, " ")

	if strings.HasPrefix(trimmedLine, "-") {
		// when coming here, it must be "- key: value" or "- value"
		trimmed := strings.TrimLeft(trimmedLine, "- ")
		if strings.Contains(trimmed, ": ") && unicode.IsLetter(rune(trimmed[0])) {
			// when coming here, it must be "- key: value"
			ss := strings.SplitN(trimmed, ": ", 2) // assuming key must not contain ": " while value might do
			k, v := ss[0], ss[1]
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, "- ")
			fmt.Fprintf(w, "%s", color.Apply(k, gp.colorByIndent(indentCnt+2))) // add length of "- "
			fmt.Fprintf(w, ": ")
			fmt.Fprintf(w, "%s", color.Apply(v, gp.colorByValue(v)))
			fmt.Fprintf(w, "\n")
		} else {
			// when coming here, it must be "- value"
			fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
			fmt.Fprintf(w, "- ")
			fmt.Fprintf(w, "%s", color.Apply(trimmed, gp.colorByValue(trimmed)))
			fmt.Fprintf(w, "\n")
		}
		return
	}

	// when coming here, "key:" or "key: value" or "value"
	if strings.Contains(trimmedLine, ": ") && unicode.IsLetter(rune(trimmedLine[0])) {
		// when coming here, it must be "key: value"
		ss := strings.SplitN(trimmedLine, ": ", 2) // assuming key must not contain ": " while value might do
		k, v := ss[0], ss[1]
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, "%s", color.Apply(k, gp.colorByIndent(indentCnt)))
		fmt.Fprintf(w, ": ")
		fmt.Fprintf(w, "%s", color.Apply(v, gp.colorByValue(v)))
		fmt.Fprintf(w, "\n")
	} else if strings.HasSuffix(trimmedLine, ":") && unicode.IsLetter(rune(trimmedLine[0])) {
		// when coming here, it must be "key:"
		trimmed := strings.TrimRight(trimmedLine, ":")
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, "%s", color.Apply(trimmed, gp.colorByIndent(indentCnt)))
		fmt.Fprintf(w, ":")
		fmt.Fprintf(w, "\n")
	} else {
		// when coming here, it must be just a "value"
		// when a string was too long, the line can be broken and come here
		fmt.Fprintf(w, "%s", gp.toSpaces(indentCnt))
		fmt.Fprintf(w, "%s", color.Apply(trimmedLine, gp.colorByValue(trimmedLine)))
		fmt.Fprintf(w, "\n")
	}
}

func (gp *GetPrinter) colorByIndent(indent int) color.Color {
	switch indent / 4 % 2 {
	case 1:
		return color.White
	default:
		return color.Yellow
	}
}

func (gp *GetPrinter) findIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func (gp *GetPrinter) toSpaces(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}

func (gp *GetPrinter) isHeader() bool {
	return gp.WithHeader && gp.isFirstLine
}

func (gp *GetPrinter) Palette() []color.Color {
	return []color.Color{color.Cyan, color.Magenta, color.Green, color.White, color.Blue}
}

func (gp *GetPrinter) DecideColor(index int, column string) color.Color {
	if gp.isHeader() {
		return HeaderColor
	}

	if column == "CrashLoopBackOff" {
		return color.Red
	}

	// When Readiness is "n/m" then yellow
	if strings.Count(column, "/") == 1 {
		if arr := strings.Split(column, "/"); arr[0] != arr[1] {
			_, e1 := strconv.Atoi(arr[0])
			_, e2 := strconv.Atoi(arr[1])
			if e1 == nil && e2 == nil { // check both is number
				return color.Yellow
			}
		}

	}

	colors := gp.Palette()
	if index >= len(colors) {
		index = index % len(colors)
	}

	return colors[index]
}

func (gp *GetPrinter) colorByValue(val string) color.Color {
	if val == "null" {
		return NullColor
	}

	if val == "true" || val == "false" {
		return BoolColor
	}

	if _, err := strconv.Atoi(val); err == nil {
		return NumberColor
	}

	return StringColor
}
