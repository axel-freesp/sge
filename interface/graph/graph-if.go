package graph

import (
	"image"
)

type XmlCreator interface {
	CreateXml() (buf []byte, err error)
}

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

type PositionMode string

const (
	PositionModeNormal   PositionMode = "normal"
	PositionModeMapping  PositionMode = "mapping"
	PositionModeExpanded PositionMode = "expanded"
)

type Shaper interface {
	Shape() image.Point
	SetShape(image.Point)
}

type Directioner interface {
	Direction() PortDirection
	SetDirection(PortDirection)
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

type IOModer interface {
	IOMode() IOMode
	SetIOMode(IOMode)
}

type IOMode string

const (
	IOModeShmem IOMode = "Shared Memory"
	IOModeAsync IOMode = "Asynchronous"
	IOModeSync  IOMode = "Isochronous"
)
