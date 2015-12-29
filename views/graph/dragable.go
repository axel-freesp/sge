package graph

import (
	//"fmt"
	//"log"
	"image"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/cairo"
)

// user objects

type NodeObject interface {
	freesp.Namer
	freesp.Positioner
	freesp.Porter
}

type PortObject freesp.Port


// graph items

type NodeIf interface {
	freesp.Namer
	freesp.Positioner
	freesp.Porter // own interface??
	SelectableBox
	Drawer
	UserObj() NodeObject
	GetOutPort(index int) PortObject
	GetInPort(index int) PortObject
}

type PortIf interface {
	freesp.Positioner
	SelectableBox
	Drawer
	UserObj() PortObject
}

type ConnectIf interface {
	SelectableBox
	Drawer
	IsLinked(nodeName string) bool
	From() NodeIf
	To() NodeIf
	FromId() int
	ToId() int
}

// Grouping

type SelectableBox interface {
	BBoxer
	Selecter
	Hitter
}

type Dragable interface {
	freesp.Namer
    SelectableBox
}


type BBoxer interface {
	BBox() image.Rectangle
}

type Selecter interface {
	Select() bool
	Deselect() bool
}

type Hitter interface {
	CheckHit(image.Point) bool
}

type Drawer interface {
	Draw(*cairo.Context)
}


// Implementations

type SelectableObject struct {
	box      image.Rectangle
	selected, highlighted bool
}

var _ SelectableBox = (*SelectableObject)(nil)

/*
 *      Default Positioner implementation
 */

func (b *SelectableObject) Position() image.Point {
	return b.BBox().Min
}

func (b *SelectableObject) SetPosition(pos image.Point) {
	size := b.box.Size()
	b.box.Min = pos
	b.box.Max = b.box.Min.Add(size)
}

/*
 *      Default BBox implementation
 */

func (b *SelectableObject) BBox() image.Rectangle {
	return b.box
}

/*
 *      Default Selecter implementation
 */

var _ Selecter  = (*SelectableObject)(nil)

func (b *SelectableObject) Select() (selected bool) {
	selected = b.selected
	b.selected = true
	return
}

func (b *SelectableObject) Deselect() (selected bool) {
	selected = b.selected
	b.selected = false
	return
}

func (b *SelectableObject) CheckHit(p image.Point) (ok bool) {
	test := image.Rectangle{p, p}
	ok = b.BBox().Overlaps(test)
	b.highlighted = ok
	return
}


// Check if a given object (index idx) would overlap with any other
// if we move it to p coordinates:
func Overlaps(drg []NodeIf, idx int, p image.Point) bool {
	box := drg[idx].BBox()
	newBox := box.Add(p.Sub(box.Min))
	for i, d := range drg {
		if i != idx && newBox.Overlaps(d.BBox()) {
			return true
		}
	}
	return false
}
