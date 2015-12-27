package graph

import (
	"image/color"
)

// Index arguments for NumericOption()
const (
	NodeWidth = iota
	NodeHeight
	NodePadX
	NodePadY
	NodeTextX
	NodeTextY
	PortW
	PortH
	PortX0
	PortY0
	PortDY
	FontSize
)

func NumericOption(index int) int {
	return options.numericOptions[index].val
}

// Index arguments for ColorOption()
const (
	Background = iota
	Normal
	Highlight
	Selected
	InputPort
	OutputPort
	BoxFrame
	Text
	HighlightInPort
	HighlightOutPort
)

func ColorOption(index int) (r, g, b float64) {
	c := options.colorOptions[index].val
	r = float64(c.R) / 255.0
	g = float64(c.G) / 255.0
	b = float64(c.B) / 255.0
	return
}

// Index arguments for StringOption()
const (
	FontPath = iota
)

func StringOption(index int) string {
	return options.stringOptions[index].val
}

/*
 *	Private
 */

type optionNumeric struct {
	label string
	val int
}

type optionColor struct {
	label string
	val color.RGBA
}

type optionString struct {
	label string
	val string
}

type gOptions struct {
	numericOptions []optionNumeric
	colorOptions  []optionColor
	stringOptions []optionString
}

// Default options: hardcoded, read-only
var defaultOptions = gOptions{
	[]optionNumeric{
		{"Node Width", 100},
		{"Node Height", 32},
		{"Node PadX", 5},
		{"Node PadY", 2},
		{"Node TextX", 10},
		{"Node TextY", 14},
		{"Port W", 8},
		{"Port H", 8},
		{"Port X0", -3},
		{"Port Y0", 24},
		{"Port DY", 12},
		{"Font Size", 12},
	},
	[]optionColor{
		{"Background", color.RGBA{240, 240, 240, 0xff}},
		{"Normal", color.RGBA{255, 204, 146, 0xff}},
		{"Highlight", color.RGBA{255, 234, 170, 0xff}},
		{"Selected", color.RGBA{221, 255, 190, 0xff}},
		{"InputPort", color.RGBA{255, 60, 60, 0xff}},
		{"OutputPort", color.RGBA{60, 60, 255, 0xff}},
		{"BoxFrame", color.RGBA{0, 0, 0, 0xff}},
		{"Text", color.RGBA{0, 0, 0, 0xff}},
		{"HighlightInPort", color.RGBA{255, 255, 120, 0xff}},
		{"HighlightOutPort", color.RGBA{255, 255, 120, 0xff}},
	},
	[]optionString{	// actually not needed anymore:
		{"FontPath", "/usr/share/fonts/truetype"},
	},
}

// Options, initialized by default options
var options = defaultOptions

