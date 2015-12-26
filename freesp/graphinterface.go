package freesp

import(
	"image"
)

type Positioner interface {
	Position() image.Point
	SetPosition(image.Point)
}

type Namer interface {
	Name() string
}

type Porter interface {
	NumInPorts() int
	NumOutPorts() int
}

