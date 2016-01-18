package graph

import (
	"image"
	interfaces "github.com/axel-freesp/sge/interface"
)

/*
 * 	Graph items
 */

type NodeIf interface {
	interfaces.Porter
	NamedBox
	UserObj() interfaces.NodeObject
}

type PortIf interface {
	BoxedSelecter
	UserObj() interfaces.PortObject
}

type ConnectIf interface {
	BoxedSelecter
	IsLinked(nodeName string) bool
	From() NodeIf
	To() NodeIf
	FromId() int
	ToId() int
}

type ArchIf interface {
	NamedBox
	Processes() []ProcessIf
	UserObj() interfaces.ArchObject
	IsLinked(name string) bool
	ChannelPort(interfaces.ChannelObject) ArchPortIf
	SelectProcess(interfaces.ProcessObject) ProcessIf
	GetSelectedProcess() (ok bool, pr interfaces.ProcessObject, p ProcessIf)
	SelectChannel(ch interfaces.ChannelObject)
	GetSelectedChannel() (ok bool, pr interfaces.ChannelObject)
}

type ArchPortIf interface {
	BoxedSelecter
}

type ProcessIf interface {
	NamedBox
	UserObj() interfaces.ProcessObject
	ArchObject() ArchIf
	SetArchObject(ArchIf)
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

type Drawer interface {
	Draw(context interface{})
}

type BoxedSelecter interface {
	BBoxer
	Selecter
    Highlighter
	Drawer
	CheckHit(image.Point) (hit, modified bool)
}

type NamedBox interface {
	interfaces.Namer
	BoxedSelecter
}


