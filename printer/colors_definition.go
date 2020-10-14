package printer

import "github.com/dty1er/kubecolor/color"

var (
	// Preset of colors for background
	// Please use them when you just need random colors
	colorsForDarkBackground = []color.Color{
		color.Cyan,
		color.Green,
		color.Magenta,
		color.White,
		color.Yellow,
	}

	colorsForLightBackground = []color.Color{
		color.Cyan,
		color.Green,
		color.Magenta,
		color.Black,
		color.Yellow,
		color.Blue,
	}

	// colors to be recommended to be used for some context
	// e.g. Json, Yaml, kubectl-describe format etc.

	// colors which look good in dark-backgrounded environment
	KeyColorForDark    = color.White
	StringColorForDark = color.Cyan
	BoolColorForDark   = color.Green
	NumberColorForDark = color.Magenta
	NullColorForDark   = color.Yellow
	HeaderColorForDark = color.White // for plain table

	// colors which look good in light-backgrounded environment
	KeyColorForLight    = color.Black
	StringColorForLight = color.Blue
	BoolColorForLight   = color.Green
	NumberColorForLight = color.Magenta
	NullColorForLight   = color.Yellow
	HeaderColorForLight = color.Black // for plain table
)
