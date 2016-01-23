package freesp

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	pf "github.com/axel-freesp/sge/interface/platform"
)

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

func GetSignalTypeByName(typeName string) (sType bh.SignalType, ok bool) {
	sType, ok = signalTypes[typeName]
	return
}

func GetRegisteredIOTypes() []string {
	return registeredIOTypes.Strings()
}

func GetIOTypeByName(name string) (ioType pf.IOTypeIf, ok bool) {
	ioType, ok = ioTypes[name]
	return
}
