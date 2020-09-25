package color

import (
	"fmt"
)

type Color int

const escape = "\x1b"

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	HiBlack Color = iota + 90
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

func Format(c Color) string {
	return fmt.Sprintf("%s[%dm", escape, c.Sequence())
}

func Reset() string {
	return Format(Color(0))
}

func (c Color) Sequence() int {
	return int(c)
}

func Apply(val string, c Color) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c.Sequence(), val, escape)
}
