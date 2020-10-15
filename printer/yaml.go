package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dty1er/kubecolor/color"
)

type YamlPrinter struct {
	DarkBackground bool

	inString bool
}

func (yp *YamlPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	yp.inString = false
	for scanner.Scan() {
		line := scanner.Text()
		inString := printLineAsYamlFormat(line, w, yp.DarkBackground, yp.inString)
		yp.inString = inString
	}
}

func printLineAsYamlFormat(line string, w io.Writer, dark, inString bool) bool {
	indentCnt := findIndent(line) // can be 0
	indent := toSpaces(indentCnt) // so, can be empty

	trimmedLine := strings.TrimLeft(line, " ")

	isArrayValue := strings.HasPrefix(trimmedLine, "- ")
	prefix := ""
	if isArrayValue {
		prefix = "- "
		trimmedLine = strings.TrimLeft(trimmedLine, "- ")
	}

	splitted := strings.SplitN(trimmedLine, ": ", 2) // assuming key does not contain ": " while value might do

	if len(splitted) == 2 && !inString {
		key := color.Apply(splitted[0], getColorByKeyIndent(indentCnt+len(prefix), 2, dark))
		val := color.Apply(splitted[1], getColorByValueType(splitted[1], dark))

		fmt.Fprintf(w, "%s%s%s: %s\n", indent, prefix, key, val)

		isValueStringNotClosed := (strings.HasPrefix(splitted[1], "'") && !strings.HasSuffix(splitted[1], "'")) ||
			(strings.HasPrefix(splitted[1], `"`) && !strings.HasSuffix(splitted[1], `"`))
		return isValueStringNotClosed
	}

	isYamlKey := isValYamlKey(trimmedLine, inString)
	if isYamlKey {
		c := getColorByKeyIndent(indentCnt+len(prefix), 2, dark)
		fmt.Fprintf(w, "%s%s%s:\n", indent, prefix, color.Apply(strings.TrimRight(trimmedLine, ":"), c))
		return false
	}

	fmt.Fprintf(w, "%s%s%s\n", indent, prefix, color.Apply(trimmedLine, getColorByValueType(trimmedLine, dark)))

	strClosed := strings.HasSuffix(trimmedLine, "'") || strings.HasSuffix(trimmedLine, `"`)

	if isArrayValue {
		return false
	}

	if !inString {
		return false
	}

	return !strClosed
}

func isValYamlKey(s string, inString bool) bool {
	// key must end with :
	if !strings.HasSuffix(s, ":") {
		return false
	}

	// even if it ends with :, if it's in a string value, it's not a key
	if inString {
		return false
	}

	return true
}
