package graph

import (
	"image"
	interfaces "github.com/axel-freesp/sge/interface"
)

/*
 * 	Graph items
 */

type NodeIf interface {
	interfaces.Namer
	interfaces.Porter
	BoxedSelecter
	Drawer
	UserObj() interfaces.NodeObject
}

type PortIf interface {
	BoxedSelecter
	Drawer
	UserObj() interfaces.PortObject
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
	interfaces.Namer
	BoxedSelecter
	Drawer
	UserObj() interfaces.ArchObject
}

type ProcessIf interface {
	interfaces.Namer
	BoxedSelecter
	Drawer
	UserObj() interfaces.ProcessObject
}

/*
 * 	Basic items
 */

type BBoxer interface {
    interfaces.Positioner
    interfaces.Shaper
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


