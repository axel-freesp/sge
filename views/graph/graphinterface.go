package graph

import (
    "image"
)

// user objects

type NodeObject interface {
	Namer
	Porter
	Positioner
}

type PortObject interface {
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

// graph items

type NodeIf interface {
	Namer
	Porter
	BoxedSelecter
	Drawer
	UserObj() NodeObject
	GetOutPort(index int) PortObject
	GetInPort(index int) PortObject
}

type PortIf interface {
	BoxedSelecter
	Drawer
	UserObj() PortObject
}

type ConnectIf interface {
	BoxedSelecter
	Drawer
	IsLinked(nodeName string) bool
	From() NodeIf
	To() NodeIf
	FromId() int
	ToId() int
}

type ArchIf interface {
	Namer
	BoxedSelecter
	Drawer
	UserObj() ArchObject
}

type ProcessIf interface {
	Namer
	BoxedSelecter
	Drawer
	UserObj() ProcessObject
}

// Basic types

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

type BBoxer interface {
    Positioner
    Shaper
	BBox() image.Rectangle
}

type Selecter interface {
	Select() bool
	Deselect() bool
	IsSelected() bool
}

type Highlighter interface {
	IsHighlighted() bool
}

type BoxedSelecter interface {
	BBoxer
	Selecter
    Highlighter
	CheckHit(image.Point) (hit, modified bool)
}


type Drawer interface {
	Draw(context interface{})
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

