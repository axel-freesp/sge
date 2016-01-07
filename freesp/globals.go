package freesp

import (
	"github.com/axel-freesp/sge/tool"
)

var signalTypes map[string]*signalType
var nodeTypes map[string]*nodeType
var libraries map[string]*library
var ioTypes map[string]*iotype
var registeredNodeTypes tool.StringList
var registeredSignalTypes tool.StringList
var registeredIOTypes tool.StringList

func Init() {
	signalTypes = make(map[string]*signalType)
	nodeTypes = make(map[string]*nodeType)
	libraries = make(map[string]*library)
	ioTypes = make(map[string]*iotype)
	registeredNodeTypes = tool.StringListInit()
	registeredSignalTypes = tool.StringListInit()
	registeredIOTypes = tool.StringListInit()
}
