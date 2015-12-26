package graph

import (
	//"fmt"
	"image"
	"github.com/axel-freesp/sge/freesp"
)

type GObject interface {
	freesp.Namer
	freesp.Positioner
	freesp.Porter
}

type ColorMode int

const (
	Normal       ColorMode = 0
	Highlight    ColorMode = 1
	Selected     ColorMode = 2
	NumColorMode ColorMode = 3
)

type DragableObject struct {
	//CtxMenuer
	//CtxMenuProvider
	//Menu     ctxMenu
	box      image.Rectangle
	image    []*image.RGBA
	userObj  GObject
	//doCtxMenu func(*DragableObject, int)
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
}

func (drg *DragableObject) NumInPorts() int {
	return drg.userObj.NumInPorts()
}

func (drg *DragableObject) NumOutPorts() int {
	return drg.userObj.NumOutPorts()
}

/*
func (drg *DragableObject) CtxMenuOpen(gc draw2d.GraphicContext, pos image.Point) image.Rectangle {
	return drg.Menu.CtxMenuOpen(gc, pos)
}

func (drg *DragableObject) CtxMenuMouse(gc draw2d.GraphicContext, pos image.Point) (image.Rectangle, bool) {
	return drg.Menu.CtxMenuMouse(gc, pos)
}

func (drg *DragableObject) CtxMenuChoose() image.Rectangle {
	return drg.Menu.CtxMenuChoose()
}

// CtxMenuer interface
func (drg *DragableObject) CtxMenu(id int) {
	drg.doCtxMenu(drg, id)
}
*/

func (drg *DragableObject) BBox() image.Rectangle {
	return drg.box
}

func (drg *DragableObject) Init(box image.Rectangle,
			userObj GObject) {
			/*,
			DrawNode func(*draw2dimg.GraphicContext, DragableObject, ColorMode),
			doCtxMenu func(*DragableObject, int)*/
	drg.userObj = userObj
	drg.box = box
	drg.image = make([]*image.RGBA, NumColorMode)
	//drg.doCtxMenu = doCtxMenu
	//drg.Menu.Init(drg)
	for _ = range drg.image {
		/*
		drg.image[i] = image.NewRGBA(box.Sub(box.Min))
		gc := draw2dimg.NewGraphicContext(drg.image[i])
		DrawNode(gc, *drg, ColorMode(i))
		*/
	}
	return
}

func (drg *DragableObject) hit(p image.Point) bool {
	return drg.box.Overlaps(image.Rectangle{p, p})
}

func HitObj(drg []DragableObject, p image.Point) (idx int, ok bool) {
	for i, o := range drg {
		if o.hit(p) {
			return i, true
		}
	}
	return -1, false
}

// Check if a given object (index idx) would overlap with any other
// if we move it to p coordinates:
func Overlaps(drg []DragableObject, idx int, p image.Point) bool {
	newBox := drg[idx].box.Add(p.Sub(drg[idx].box.Min))
	for i, d := range drg {
		if i != idx && newBox.Overlaps(d.box) {
			return true
		}
	}
	return false
}
