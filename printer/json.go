package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/hidetatz/kubecolor/color"
)

type JsonPrinter struct {
	DarkBackground bool
	PlainHierarchy bool
}

func (jp *JsonPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		jp.printLineAsJsonFormat(line, w, jp.DarkBackground)
	}
}

func (jp *JsonPrinter) printLineAsJsonFormat(line string, w io.Writer, dark bool) {
	indentCnt := findIndent(line)
	indent := toSpaces(indentCnt)
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
		// fmt.Fprintf(w, "%s", toSpaces(indentCnt))
		fmt.Fprintf(w, "%s%s", indent, trimmedLine)
		fmt.Fprintf(w, "\n")
		return
	}

	// when coming here, the line must be one of:
	// "key": {
	// "key": [
	// "key": value
	// "key": value,
	// value,
	// value
	splitted := strings.SplitN(trimmedLine, ": ", 2) // if key contains ": " this works in a wrong way but it's unlikely to happen

	if len(splitted) == 1 {
		// when coming here, it will be a value in an array
		fmt.Fprintf(w, "%s%s\n", indent, jp.toColorizedJsonValue(splitted[0], dark))
		return
	}

	key := splitted[0]
	val := splitted[1]

	fmt.Fprintf(w, "%s%s: %s\n", indent, jp.toColorizedJsonKey(key, indentCnt, 4, dark), jp.toColorizedJsonValue(val, dark))
}

// toColorizedJsonKey returns colored json key
func (jp *JsonPrinter) toColorizedJsonKey(key string, indentCnt, basicWidth int, dark bool) string {
	hasColon := strings.HasSuffix(key, ":")
	// remove colon and double quotations although they might not exist actually
	key = strings.TrimRight(key, ":")
	doubleQuoteTrimmed := strings.TrimRight(strings.TrimLeft(key, `"`), `"`)

	format := `"%s"`
	if hasColon {
		format += ":"
	}

	return fmt.Sprintf(format, color.Apply(doubleQuoteTrimmed, getColorByKeyIndent(indentCnt, basicWidth, dark, jp.PlainHierarchy)))
}

// toColorizedJsonValue returns colored json value.
// This function checks it trailing comma and double quotation exist
// then colorize the given value considering them.
func (jp *JsonPrinter) toColorizedJsonValue(value string, dark bool) string {
	if value == "{" {
		return "{"
	}

	if value == "[" {
		return "["
	}

	if value == "{}," {
		return "{},"
	}

	if value == "{}" {
		return "{}"
	}

	hasComma := strings.HasSuffix(value, ",")
	// remove comma and double quotations although they might not exist actually
	value = strings.TrimRight(value, ",")

	isString := strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)
	doubleQuoteTrimmedValue := strings.TrimRight(strings.TrimLeft(value, `"`), `"`)

	var format string
	switch {
	case hasComma && isString:
		format = `"%s",`
	case hasComma:
		format = `%s,`
	case isString:
		format = `"%s"`
	default:
		format = `%s`
	}

	return fmt.Sprintf(format, color.Apply(doubleQuoteTrimmedValue, getColorByValueType(value, dark)))
}
