package freesp

import (
	"image"
)

type Positioner interface {
	Position() image.Point
	SetPosition(image.Point)
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

type Directioner interface {
	Direction() PortDirection
	SetDirection(PortDirection)
}

type Filenamer interface {
	Filename() string
	SetFilename(string)
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
}
