package freesp

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/axel-freesp/sge/tool"
)

var signalTypes map[string]bh.SignalTypeIf
var nodeTypes map[string]bh.NodeTypeIf
var libraries map[string]bh.LibraryIf
var ioTypes map[string]pf.IOTypeIf
var registeredNodeTypes tool.StringList
var registeredSignalTypes tool.StringList
var registeredIOTypes tool.StringList

func Init() {
	signalTypes = make(map[string]bh.SignalTypeIf)
	nodeTypes = make(map[string]bh.NodeTypeIf)
	// TODO: is this redundant with global.libraries?
	libraries = make(map[string]bh.LibraryIf)
	ioTypes = make(map[string]pf.IOTypeIf)
	registeredNodeTypes = tool.StringListInit()
	registeredSignalTypes = tool.StringListInit()
	registeredIOTypes = tool.StringListInit()
}

// TODO: Turn into interface
func GetRegisteredNodeTypes() []string {
	return registeredNodeTypes.Strings()
}

func GetRegisteredSignalTypes() []string {
	return registeredSignalTypes.Strings()
}

func GetNodeTypeByName(typeName string) (nType bh.NodeTypeIf, ok bool) {
	nType, ok = nodeTypes[typeName]
	return
}

func GetSignalTypeByName(typeName string) (sType bh.SignalTypeIf, ok bool) {
	sType, ok = signalTypes[typeName]
	return
}

func GetLibraryByName(filename string) (lib bh.LibraryIf, ok bool) {
	lib, ok = libraries[filename]
	return
}

func GetRegisteredIOTypes() []string {
	return registeredIOTypes.Strings()
}

func GetIOTypeByName(name string) (ioType pf.IOTypeIf, ok bool) {
	ioType, ok = ioTypes[name]
	return
}

func RegisterSignalType(st bh.SignalTypeIf) {
	signalTypes[st.TypeName()] = st
	registeredSignalTypes.Append(st.TypeName())
}

func RegisterNodeType(nt bh.NodeTypeIf) {
	nodeTypes[nt.TypeName()] = nt
	registeredNodeTypes.Append(nt.TypeName())
}

func RegisterLibrary(lib bh.LibraryIf) {
	libraries[lib.Filename()] = lib
}

func RegisterIOType(iot pf.IOTypeIf) {
	ioTypes[iot.Name()] = iot
	registeredIOTypes.Append(iot.Name())
}

func RemoveRegisteredLibrary(lib bh.LibraryIf) {
	delete(libraries, lib.Filename())
}

func RemoveRegisteredNodeType(nt bh.NodeTypeIf) {
	name := nt.TypeName()
	delete(nodeTypes, name)
	_, ok := registeredNodeTypes.Find(name)
	if ok {
		registeredNodeTypes.Remove(name)
	}
}

func RemoveRegisteredSignalType(st bh.SignalTypeIf) {
	name := st.TypeName()
	delete(signalTypes, name)
	_, ok := registeredSignalTypes.Find(name)
	if ok {
		registeredSignalTypes.Remove(name)
	}
}

func RemoveRegisteredIOType(iot pf.IOTypeIf) {
	name := iot.Name()
	delete(ioTypes, name)
	_, ok := registeredIOTypes.Find(name)
	if ok {
		registeredIOTypes.Remove(name)
	}
}
