package printer

import (
	"fmt"
	"io"

	"github.com/dty1er/kubecolor/color"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/printer"
)

/*
   thankfully copied from https://github.com/goccy/go-yaml/blob/master/cmd/ycat/ycat.go and
   some modifications are applied by kubecolor author
*/

const escape = "\x1b"

func format(attr color.Color) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func colorProperty(c color.Color) func() *printer.Property {
	return func() *printer.Property {
		return &printer.Property{
			Prefix: color.Format(c),
			Suffix: color.Reset(),
		}
	}
}

func printYaml(w io.Writer, output []byte) error {
	tokens := lexer.Tokenize(string(output))

	p := printer.Printer{
		MapKey: colorProperty(KeyColor),
		Anchor: colorProperty(AnchorColor),
		Alias:  colorProperty(AliasColor),
		Bool:   colorProperty(BoolColor),
		String: colorProperty(StringColor),
		Number: colorProperty(NumberColor),
	}

	fmt.Fprintf(w, p.PrintTokens(tokens)+"\n")
	return nil
}
