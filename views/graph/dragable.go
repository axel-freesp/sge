package graph

import (
	//"fmt"
	//"log"
	"image"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/cairo"
)

type GObject interface {
	freesp.Namer
	freesp.Positioner
	freesp.Porter
}

type ColorMode int

const (
	NormalMode       ColorMode = iota
	HighlightMode
	SelectedMode
	NumColorMode
)

type Dragable interface {
	GObject
	UserObj() GObject
	BBox() image.Rectangle
	Draw(context *cairo.Context, mode ColorMode)
	Check(pos image.Point) (r image.Rectangle, ok bool)
}

type DragableObject struct {
	box      image.Rectangle
	userObj  GObject
}

var _ GObject = (*DragableObject)(nil)

func (drg *DragableObject) Name() string {
	return (*drg).userObj.Name()
}

func (drg *DragableObject) Position() image.Point {
	return (*drg).userObj.Position()
}

func (drg *DragableObject) SetPosition(pos image.Point) {
	(*drg).userObj.SetPosition(pos)
	size := drg.box.Size()
	drg.box.Min = pos
	drg.box.Max = drg.box.Min.Add(size)
}

func (drg *DragableObject) NumInPorts() int {
	return drg.userObj.NumInPorts()
}

func (drg *DragableObject) NumOutPorts() int {
	return drg.userObj.NumOutPorts()
}

func (drg *DragableObject) UserObj() GObject {
	return drg.userObj
}

func (drg *DragableObject) BBox() image.Rectangle {
	return drg.box
}

func hit(drg Dragable, p image.Point) bool {
	return drg.BBox().Overlaps(image.Rectangle{p, p})
}

func HitObj(drg []Dragable, p image.Point) (idx int, ok bool) {
	for i, o := range drg {
		if hit(o, p) {
			return i, true
		}
	}
	return -1, false
}

// Check if a given object (index idx) would overlap with any other
// if we move it to p coordinates:
func Overlaps(drg []Dragable, idx int, p image.Point) bool {
	box := drg[idx].BBox()
	newBox := box.Add(p.Sub(box.Min))
	for i, d := range drg {
		if i != idx && newBox.Overlaps(d.BBox()) {
			return true
		}
	}
	return false
}
