package color

import (
	"fmt"
)

type Color int

const escape = "\x1b"

var (
	colorsForDarkBackground = []Color{
		Cyan,
		Green,
		Magenta,
		White,
		Yellow,
	}

	colorsForLightBackground = []Color{
		Cyan,
		Green,
		Magenta,
		Black,
		Yellow,
		Blue,
	}
)

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func (c Color) sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.sequence(), val, escape)
}

func GetColors(dark bool) []Color {
	if dark {
		return colorsForDarkBackground
	}

	return colorsForLightBackground
}
