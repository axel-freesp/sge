package interfaces

import (
	"fmt"
	"image"
)

//
//		Context (implemented by *sge.Global)
//
type Context interface {
	SelectNode(NodeObject)          // single click selection
	EditNode(NodeObject)            // double click selection
	SelectPort(PortObject)          // single click selection
	SelectConnect(ConnectionObject) // single click selection
	SelectArch(ArchObject)
	SelectProcess(ProcessObject)
	SelectChannel(ChannelObject)
	SelectMapElement(MapElemObject)
}

/*
 * 	Graph user objects
 */

type GraphElement interface{}

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
	MapElemObject(NodeObject) (MapElemObject, bool)
}

type MapElemObject interface {
	Positioner
	NodeObject() NodeObject
	ProcessObject() (ProcessObject, bool)
}

type PlatformObject interface {
	ArchObjects() []ArchObject
}

type ArchObject interface {
	Namer
	ModePositioner
	ProcessObjects() []ProcessObject
	IOTypeObjects() []IOTypeObject
	PortObjects() []ArchPortObject
}

type ArchPortObject interface {
	ModePositioner
	Channel() ChannelObject
}

type IOTypeObject interface {
	Namer
	IOModer
}

type ProcessObject interface {
	Namer
	ModePositioner
	ArchObject() ArchObject
	InChannelObjects() []ChannelObject
	OutChannelObjects() []ChannelObject
}

type ChannelObject interface {
	ModePositioner
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

type ModePositioner interface {
	ModePosition(PositionMode) image.Point
	SetModePosition(PositionMode, image.Point)
}

type PositionMode int

const (
	PositionModeNormal PositionMode = iota
	PositionModeMapping
)

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
