package printer

import (
	"fmt"
	"io"

	"github.com/hokaccha/go-prettyjson"

	fcolor "github.com/fatih/color"
)

func printJson(w io.Writer, output []byte) error {
	f := &prettyjson.Formatter{
		KeyColor:        fcolor.New(fcolor.Attribute(KeyColor.Sequence())),
		StringColor:     fcolor.New(fcolor.Attribute(StringColor.Sequence())),
		BoolColor:       fcolor.New(fcolor.Attribute(BoolColor.Sequence())),
		NumberColor:     fcolor.New(fcolor.Attribute(NumberColor.Sequence())),
		NullColor:       fcolor.New(fcolor.Attribute(NullColor.Sequence())),
		StringMaxLength: 0,
		DisabledColor:   false,
		Indent:          2,
		Newline:         "\n",
	}
	colorized, err := f.Format(output)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, string(colorized)+"\n") // go-prettyjson trims leading space
	return nil
}
