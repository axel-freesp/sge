package interfaces

import (
	"fmt"
    "image"
)

/*
 * 	Graph user objects
 */

type NodeObject interface {
	Namer
	Porter
	Positioner
}

type PortObject interface {
	Namer
	Directioner
	NodeObject() NodeObject
	ConnectionObject(PortObject) ConnectionObject
}

type ArchObject interface {
	Namer
	Positioner
	Shaper
}

type ProcessObject interface {
	Namer
	Positioner
}

type ConnectionObject interface {
	FromObject() PortObject
	ToObject() PortObject
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

