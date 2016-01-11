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
	Processes() []ProcessIf
	UserObj() interfaces.ArchObject
	IsLinked(name string) bool
	ChannelPort(interfaces.ChannelObject) ArchPortIf
	SelectProcess(interfaces.ProcessObject)
	GetSelectedProcess() (ok bool, pr interfaces.ProcessObject)
	SelectChannel(ch interfaces.ChannelObject)
	GetSelectedChannel() (ok bool, pr interfaces.ChannelObject)
}

type ArchPortIf interface {
	BoxedSelecter
	Drawer
}

type ProcessIf interface {
	interfaces.Namer
	BoxedSelecter
	Drawer
	UserObj() interfaces.ProcessObject
	ArchObject() ArchIf
	SelectChannel(ch interfaces.ChannelObject)
	GetSelectedChannel() (ok bool, pr interfaces.ChannelObject)
}

/*
 * 	Basic items
 */

type BBoxer interface {
    interfaces.Positioner
    interfaces.Shaper
	BBox() image.Rectangle
}

type OnSelectionFunc func()

type Selecter interface {
	Select() bool
	Deselect() bool
	IsSelected() bool
	RegisterOnSelect(inSelect, onDeselect OnSelectionFunc)
}

type OnHighlightFunc func(hit bool, pos image.Point) (modified bool)

type Highlighter interface {
	IsHighlighted() bool
	RegisterOnHighlight(onHighlight OnHighlightFunc)
	DoHighlight(state bool, pos image.Point) (modified bool)
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


