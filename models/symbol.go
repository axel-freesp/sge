package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
)

// Public:

// Access the Pixbuf objects
func symbolPixbuf(s freesp.Symbol) *gdk.Pixbuf {
	return symbolTable[s]
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

var symbolTable map[freesp.Symbol]*gdk.Pixbuf

func symbolInit(iconPath string) (err error) {
	symbolTable = make(map[freesp.Symbol]*gdk.Pixbuf)
	for _, s := range sConfig {
		symbolTable[s.symbol], err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/%s", iconPath, s.filename))
		if err != nil {
			err = fmt.Errorf("%ssymbolInit error loading %s: %s\n", err, s.filename, err)
		}
	}
	return
}
