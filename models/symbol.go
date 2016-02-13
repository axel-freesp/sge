package models

import (
	"fmt"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/gotk3/gotk3/gdk"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type method int

const (
	gdkReadFile = iota
	ioReadFile
)

const imageMethod = gdkReadFile

// Access the Pixbuf objects
func normalPixbuf(s tr.Symbol) *gdk.Pixbuf {
	return normalTable[s]
}

func readonlyPixbuf(s tr.Symbol) *gdk.Pixbuf {
	return readonlyTable[s]
}

type symbolConfig struct {
	symbol   tr.Symbol
	filename string
}

const suffix = "png"

func makeFilename(s string) string {
	return fmt.Sprintf("%s.%s", s, suffix)
}

var sConfig = []symbolConfig{
	{tr.SymbolSignalType, makeFilename("signal-type")},
	{tr.SymbolInputPort, makeFilename("inputPic")},
	{tr.SymbolOutputPort, makeFilename("outputPic")},
	{tr.SymbolConnection, makeFilename("linkPic")},
	{tr.SymbolImplElement, makeFilename("implementationPic")},
	{tr.SymbolImplGraph, makeFilename("implGraphPic")},
	{tr.SymbolLibrary, makeFilename("libraryPic")},
	{tr.SymbolInputPortType, makeFilename("inputPic")},
	{tr.SymbolOutputPortType, makeFilename("outputPic")},
	{tr.SymbolInputNode, makeFilename("inputNodePic")},
	{tr.SymbolOutputNode, makeFilename("outputNodePic")},
	{tr.SymbolProcessingNode, makeFilename("processingNodePic")},
	{tr.SymbolNodeType, makeFilename("nodeTypePic")},
	{tr.SymbolSignalGraph, makeFilename("signalGraphPic")},
	{tr.SymbolPlatform, makeFilename("platformPic")},
	{tr.SymbolArch, makeFilename("archPic")},
	{tr.SymbolProcess, makeFilename("processPic")},
	{tr.SymbolIOType, makeFilename("signal-type")},
	{tr.SymbolInChannel, makeFilename("inputPic")},
	{tr.SymbolOutChannel, makeFilename("outputPic")},
	{tr.SymbolMappings, makeFilename("mappingsPic")},
	{tr.SymbolMapped, makeFilename("mappedPic")},
	{tr.SymbolUnmapped, makeFilename("unmappedPic")},
}

var normalTable, readonlyTable map[tr.Symbol]*gdk.Pixbuf

func symbolInit(iconPath string) (err error) {
	normalTable = make(map[tr.Symbol]*gdk.Pixbuf)
	readonlyTable = make(map[tr.Symbol]*gdk.Pixbuf)
	for _, s := range sConfig {
		normalTable[s.symbol], err = normalInit(iconPath, s)
		if err != nil {
			err = fmt.Errorf("symbolInit error loading %s: normalInit failed: %s\n", s.filename, err)
			return
		}
		readonlyTable[s.symbol], err = readonlyInit(iconPath, s)
		if err != nil {
			err = fmt.Errorf("symbolInit error loading %s: readonlyInit failed: %s\n", s.filename, err)
			return
		}
	}
	return
}

// Make all pixels brighter, keep alpha
func readonlyInit(iconPath string, config symbolConfig) (pixbuf *gdk.Pixbuf, err error) {
	pixbuf, err = gdk.PixbufCopy(normalTable[config.symbol])
	if err != nil {
		err = fmt.Errorf("symbolInit error loading %s: %s\n", config.filename, err)
		return
	}
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

func normalInit(iconPath string, config symbolConfig) (img *gdk.Pixbuf, err error) {
	switch imageMethod {
	case gdkReadFile:
		img, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/%s", iconPath, config.filename))
	case ioReadFile:
		var file *os.File
		file, err = os.Open(fmt.Sprintf("%s/%s", iconPath, config.filename))
		if err != nil {
			return
		}
		defer file.Close()
		var decImg image.Image
		decImg, _, err = image.Decode(file)
		if err != nil {
			err = fmt.Errorf("readonlyInit error: image.Decode failed: %s", err)
			return
		}
		//img, err = gdk.PixbufNewFromBytes(ImageRgbaToGdkColorspace(decImg))
		_ = decImg
		if err != nil {
			err = fmt.Errorf("readonlyInit error: PixbufNewFromBytes failed: %s", err)
		}
	default:
	}
	return
}

func ImageRgbaToGdkColorspace(img image.Image) (data []byte,
	colorspace gdk.Colorspace,
	hasAlpha bool,
	bitsPerSample, width, height, stride int) {
	bounds := img.Bounds()
	width = bounds.Max.X - bounds.Min.X
	height = bounds.Max.Y - bounds.Min.Y
	colorspace = gdk.COLORSPACE_RGB
	hasAlpha = true
	bitsPerSample = 8
	stride = width * 4
	data = make([]byte, width*height*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		j := stride * y
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			switch a {
			case 0xFFFF:
				data[j+4*x+0] = byte(r >> 8)
				data[j+4*x+1] = byte(g >> 8)
				data[j+4*x+2] = byte(b >> 8)
				data[j+4*x+3] = byte(a >> 8)
			case 0:
				data[j+4*x+0] = 0
				data[j+4*x+1] = 0
				data[j+4*x+2] = 0
				data[j+4*x+3] = 0
			default:
				data[j+4*x+0] = byte((r << 8) / a)
				data[j+4*x+1] = byte((g << 8) / a)
				data[j+4*x+2] = byte((b << 8) / a)
				data[j+4*x+3] = byte(a >> 8)
			}
		}
	}
	return
}
