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

func RemoveRegisteredNodeType(nt NodeTypeIf) {
	name := nt.TypeName()
	delete(nodeTypes, name)
	_, ok := registeredNodeTypes.Find(name)
	if ok {
		registeredNodeTypes.Remove(name)
	}
}

func RemoveRegisteredSignalType(st SignalType) {
	name := st.TypeName()
	delete(signalTypes, name)
	_, ok := registeredSignalTypes.Find(name)
	if ok {
		registeredSignalTypes.Remove(name)
	}
}

func RemoveRegisteredIOType(iot IOTypeIf) {
	name := iot.Name()
	delete(ioTypes, name)
	_, ok := registeredIOTypes.Find(name)
	if ok {
		registeredIOTypes.Remove(name)
	}
}
