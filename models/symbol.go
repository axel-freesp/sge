package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
)

// Public:

// Access the Pixbuf objects
func normalPixbuf(s freesp.Symbol) *gdk.Pixbuf {
	return normalTable[s]
}

func readonlyPixbuf(s freesp.Symbol) *gdk.Pixbuf {
	return readonlyTable[s]
}

type symbolConfig struct {
	symbol   freesp.Symbol
	filename string
}

var sConfig = []symbolConfig{
	{freesp.SymbolSignalType, "signal-type.png"},
	{freesp.SymbolInputPort, "inport-green.png"},
	{freesp.SymbolOutputPort, "outport-red.png"},
	{freesp.SymbolConnection, "link.png"},
	{freesp.SymbolImplElement, "test0.png"},
	{freesp.SymbolImplGraph, "test1.png"},
	{freesp.SymbolLibrary, "test0.png"},
	{freesp.SymbolInputPortType, "inport-green.png"},
	{freesp.SymbolOutputPortType, "outport-red.png"},
	{freesp.SymbolInputNode, "input.png"},
	{freesp.SymbolOutputNode, "output.png"},
	{freesp.SymbolProcessingNode, "node.png"},
	{freesp.SymbolNodeType, "node-type.png"},
	{freesp.SymbolSignalGraph, "test1.png"},
}

var normalTable, readonlyTable map[freesp.Symbol]*gdk.Pixbuf

func symbolInit(iconPath string) (err error) {
	normalTable = make(map[freesp.Symbol]*gdk.Pixbuf)
	readonlyTable = make(map[freesp.Symbol]*gdk.Pixbuf)
	for _, s := range sConfig {
		normalTable[s.symbol], err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/%s", iconPath, s.filename))
		if err != nil {
			err = fmt.Errorf("symbolInit error loading %s: %s\n", s.filename, err)
			return
		}
		err = readonlyInit(iconPath, s)
		if err != nil {
			err = fmt.Errorf("symbolInit error: %s\n", err)
			return
		}
	}
	return
}

func readonlyInit(iconPath string, config symbolConfig) (err error) {
	pixbuf, err := gdk.PixbufCopy(normalTable[config.symbol])
	if err != nil {
		err = fmt.Errorf("symbolInit error loading %s: %s\n", config.filename, err)
		return
	}
	readonlyTable[config.symbol] = pixbuf
	data := pixbuf.GetPixels()
	pixelSize := 3
	if pixbuf.GetHasAlpha() {
		pixelSize = 4
	}
	for y := 0; y < pixbuf.GetHeight(); y++ {
		rowIdx := y * pixbuf.GetRowstride()
		for x := 0; x < pixbuf.GetWidth(); x++ {
			pixel := func(i int) *byte { return &data[rowIdx+x*pixelSize+i] }
			r, g, b := int(*pixel(0)), int(*pixel(1)), int(*pixel(2))
			r = (r + 255) / 2
			g = (g + 255) / 2
			b = (b + 255) / 2
			*pixel(0), *pixel(1), *pixel(2) = byte(r), byte(g), byte(b)
		}
	}

	return
}

/*
 * What a pity: alpha colors come wrong!
 *
func readonlyInit(iconPath string, config symbolConfig) (err error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", iconPath, config.filename))
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		err = fmt.Errorf("readonlyInit error: image.Decode failed: %s", err)
		return
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	data := make([]byte, width*height*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		j := width * 4 * y
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[j+4*x+0] = byte(r)
			data[j+4*x+1] = byte(g)
			data[j+4*x+2] = byte(b)
			data[j+4*x+3] = byte(a)
		}
	}

	readonlyTable[config.symbol], err = gdk.PixbufNewFromBytes(data, gdk.COLORSPACE_RGB, true, 8, width, height, width*4)
	if err != nil {
		err = fmt.Errorf("readonlyInit error: PixbufNewFromXpmData failed: %s", err)
	}
	return
}
*/
