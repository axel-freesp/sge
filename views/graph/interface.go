package graph

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	"image"
)

/*
 * 	Graph items
 */

type NodeIf interface {
	Porter
	NamedBox
	UserObj() bh.NodeIf
	Expand()
	Collapse()
	InPort(int) BBoxer
	OutPort(int) BBoxer
	SelectPort(port bh.PortIf)
	GetSelectedPort() (ok bool, port bh.PortIf)
}

type Porter interface {
	InPorts() []bh.PortIf
	OutPorts() []bh.PortIf
	InPortIndex(portname string) int
	OutPortIndex(portname string) int
}

type PortIf interface {
	BoxedSelecter
	UserObj() bh.PortIf
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
	UserObj() pf.ArchIf
	IsLinked(name string) bool
	ChannelPort(pf.ChannelIf) ArchPortIf
	SelectProcess(pf.ProcessIf) ProcessIf
	GetSelectedProcess() (ok bool, pr pf.ProcessIf, p ProcessIf)
	SelectChannel(ch pf.ChannelIf)
	GetSelectedChannel() (ok bool, pr pf.ChannelIf)
}

type ArchPortIf interface {
	BoxedSelecter
}

type ProcessIf interface {
	NamedBox
	UserObj() pf.ProcessIf
	ArchObject() ArchIf
	SetArchObject(ArchIf)
	SelectChannel(ch pf.ChannelIf)
	GetSelectedChannel() (ok bool, pr pf.ChannelIf)
}

/*
 * 	Basic items
 */

type NamedBox interface {
	gr.Namer
	BoxedSelecter
}

type BoxedSelecter interface {
	BBoxer
	Selecter
	Highlighter
	Drawer
	CheckHit(image.Point) (hit, modified bool)
}

type BBoxer interface {
	gr.Positioner
	gr.Shaper
	BBox() image.Rectangle
}

type Drawer interface {
	Draw(context interface{})
}

type Selecter interface {
	Select() bool
	Deselect() bool
	IsSelected() bool
	RegisterOnSelect(inSelect, onDeselect OnSelectionFunc)
}
type OnSelectionFunc func()

type Highlighter interface {
	IsHighlighted() bool
	RegisterOnHighlight(onHighlight OnHighlightFunc)
	DoHighlight(state bool, pos image.Point) (modified bool)
}
type OnHighlightFunc func(hit bool, pos image.Point) (modified bool)
