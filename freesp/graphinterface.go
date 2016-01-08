package freesp

import (
	"image"
)

type Positioner interface {
	Position() image.Point
	SetPosition(image.Point)
}

type Shaper interface {
	Shape() image.Point
	SetShape(image.Point)
}

type Namer interface {
	Name() string
	SetName(string)
}

type Porter interface {
	NumInPorts() int
	NumOutPorts() int
	InPortIndex(portname string) int
	OutPortIndex(portname string) int
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

type Directioner interface {
	Direction() PortDirection
	SetDirection(PortDirection)
}
