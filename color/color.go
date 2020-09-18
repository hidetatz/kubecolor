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
