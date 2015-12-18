package models

import (
	"fmt"
	"github.com/gotk3/gotk3/gdk"
)

// Public:
type Symbol int

const (
	SymbolInputPort Symbol = iota
	SymbolOutputPort
	SymbolSignalType
	SymbolConnection
	SymbolImplElement
	SymbolImplGraph
	SymbolLibrary
	SymbolInputPortType
	SymbolOutputPortType
	SymbolInputNode
	SymbolOutputNode
	SymbolProcessingNode
	SymbolNodeType
	SymbolSignalGraph
)

// Access the Pixbuf objects
func symbolPixbuf(s Symbol) *gdk.Pixbuf {
	return symbolTable[s]
}

type symbolConfig struct {
	symbol   Symbol
	filename string
}

var sConfig = []symbolConfig{
	{SymbolSignalType, "signal-type.png"},
	{SymbolInputPort, "inport-green.png"},
	{SymbolOutputPort, "outport-red.png"},
	{SymbolConnection, "link.png"},
	{SymbolImplElement, "test0.png"},
	{SymbolImplGraph, "test1.png"},
	{SymbolLibrary, "test0.png"},
	{SymbolInputPortType, "inport-green.png"},
	{SymbolOutputPortType, "outport-red.png"},
	{SymbolInputNode, "input.png"},
	{SymbolOutputNode, "output.png"},
	{SymbolProcessingNode, "node.png"},
	{SymbolNodeType, "node-type.png"},
	{SymbolSignalGraph, "test1.png"},
}

var symbolTable map[Symbol]*gdk.Pixbuf

func symbolInit(iconPath string) (err error) {
	symbolTable = make(map[Symbol]*gdk.Pixbuf)
	for _, s := range sConfig {
		symbolTable[s.symbol], err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/%s", iconPath, s.filename))
		if err != nil {
			err = fmt.Errorf("%ssymbolInit error loading %s: %s\n", err, s.filename, err)
		}
	}
	return
}
