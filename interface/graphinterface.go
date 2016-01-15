package interfaces

import (
	"fmt"
	"image"
)

/*
 * 	Graph user objects
 */

type GraphElement interface {
}

type GraphObject interface {
	NodeObjects() []NodeObject
}

type NodeObject interface {
	Namer
	Porter
	Positioner
}

type PortObject interface {
	Namer
	Directioner
	NodeObject() NodeObject
	ConnectionObjects() []PortObject
	ConnectionObject(PortObject) ConnectionObject
}

type ConnectionObject interface {
	FromObject() PortObject
	ToObject() PortObject
}

type MappingObject interface {
	GraphObject() GraphObject
	PlatformObject() PlatformObject
	MappedObject(NodeObject) (ProcessObject, bool)
}

type PlatformObject interface {
	ArchObjects() []ArchObject
}

type ArchObject interface {
	Namer
	Positioner
	Shaper
	ProcessObjects() []ProcessObject
	IOTypeObjects() []IOTypeObject
	PortObjects() []ArchPortObject
}

type ArchPortObject interface {
	Positioner
	Channel() ChannelObject
}

type IOTypeObject interface {
	Namer
	Positioner
	IOModer
}

type ProcessObject interface {
	Namer
	Positioner
	ArchObject() ArchObject
	InChannelObjects() []ChannelObject
	OutChannelObjects() []ChannelObject
}

type ChannelObject interface {
	Positioner
	Directioner
	IOTypeObject() IOTypeObject
	ProcessObject() ProcessObject
	LinkObject() ChannelObject
	ArchPortObject() ArchPortObject
}

/*
 * 	Basic types
 */

type Namer interface {
	Name() string
	SetName(string)
}

type Positioner interface {
	Position() image.Point
	SetPosition(image.Point)
}

type Shaper interface {
	Shape() image.Point
	SetShape(image.Point)
}

type Porter interface {
	GetInPorts() []PortObject
	GetOutPorts() []PortObject
	InPortIndex(portname string) int
	OutPortIndex(portname string) int
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

var _ fmt.Stringer = PortDirection(false)

func (d PortDirection) String() (s string) {
	if d == InPort {
		s = "Input"
	} else {
		s = "Output"
	}
	return
}

type Directioner interface {
	Direction() PortDirection
	SetDirection(PortDirection)
}

type IOMode string

const (
	IOModeShmem IOMode = "Shared Memory"
	IOModeAsync IOMode = "Asynchronous"
	IOModeSync  IOMode = "Isochronous"
)

type IOModer interface {
	IOMode() IOMode
	SetIOMode(IOMode)
}
