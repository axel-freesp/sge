package graph

import (
	"image"
	"github.com/axel-freesp/sge/freesp"
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
}

type BBoxer interface {
	BBox() image.Rectangle
}

type Selecter interface {
	Select() bool
	Deselect() bool
	IsSelected() bool
	CheckHit(image.Point) (hit, modified bool)
}

type Drawer interface {
	Draw(context interface{})
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

func (b *SelectableObject) IsSelected() bool {
	return b.selected
}

func (b *SelectableObject) CheckHit(p image.Point) (hit, modified bool) {
	test := image.Rectangle{p, p}
	hit = b.BBox().Overlaps(test)
	modified = b.highlighted != hit
	b.highlighted = hit
	return
}


// Check if the selected object would overlap with any other
// if we move it to p coordinates:
func Overlaps(drg []NodeIf, p image.Point) bool {
	var box image.Rectangle
	for _, n := range drg {
		if n.IsSelected() {
			box = n.BBox()
			break
		}
	}
	newBox := box.Add(p.Sub(box.Min))
	for _, n := range drg {
		if !n.IsSelected() && newBox.Overlaps(n.BBox()) {
			return true
		}
	}
	return false
}
