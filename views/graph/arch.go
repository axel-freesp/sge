package graph

import (
	"image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
)

type Arch struct {
	SelectableBox
	userObj interfaces.ArchObject
}

func ArchNew(box image.Rectangle, userObj interfaces.ArchObject) *Arch {
	return &Arch{SelectableBoxInit(box), userObj}
}

var _ ArchIf = (*Arch)(nil)

func (a *Arch) UserObj() interfaces.ArchObject {
	return a.userObj
}

/*
 *	Drawer interface
 */

func (a *Arch) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		context.Stroke()
	}
}

/*
 *	freesp.Positioner interface
 */

func (a *Arch) Position() image.Point {
	return a.userObj.Position()
}

func (a *Arch) SetPosition(pos image.Point) {
	a.userObj.SetPosition(pos)
}

/*
 *	BBoxer interface
 */

func (a *Arch) BBox() image.Rectangle {
	a.box.Min = a.Position()
	a.box.Max = a.box.Min.Add(a.Shape())
	return a.box
}

/*
 *  freesp.Shaper API
 */

func (a *Arch) Shape() image.Point {
	return a.userObj.Shape()
}

func (a *Arch) SetShape(shape image.Point) {
	a.userObj.SetShape(shape)
}

/*
 *	Namer interface
 */

func (a *Arch) Name() string {
	return a.userObj.Name()
}

func (a *Arch) SetName(newName string) {
	a.userObj.SetName(newName)
}


