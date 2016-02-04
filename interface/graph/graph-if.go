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
	Positioner
	ModePosition(PositionMode) image.Point
	SetModePosition(PositionMode, image.Point)
	SetActiveMode(PositionMode)
	ActiveMode() PositionMode
}

type PathModePositioner interface {
	ModePositioner
	PathModePosition(string, PositionMode) image.Point
	SetPathModePosition(string, PositionMode, image.Point)
	PathList() []string
	SetActivePath(string)
	ActivePath() string
}

type Expander interface {
	Expanded() bool
	SetExpanded(bool)
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
