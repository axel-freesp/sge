package freesp

import (
	"github.com/axel-freesp/sge/tool"
)

var signalTypes map[string]*signalType
var nodeTypes map[string]*nodeType
var portTypes map[string]*portType
var libraries map[string]*library
var registeredNodeTypes tool.StringList
var registeredSignalTypes tool.StringList

func Init() {
	signalTypes = make(map[string]*signalType)
	nodeTypes = make(map[string]*nodeType)
	portTypes = make(map[string]*portType)
	libraries = make(map[string]*library)
	registeredNodeTypes = tool.StringListInit()
	registeredSignalTypes = tool.StringListInit()
}

