package graph

import (
	"image"
	"github.com/gotk3/gotk3/cairo"
	//"github.com/axel-freesp/sge/freesp"
)

type Process struct {
	SelectableBox
	userObj ProcessObject
}

func ProcessNew(box image.Rectangle, userObj ProcessObject) *Process {
	return &Process{SelectableBoxInit(box), userObj}
}

var _ ProcessIf = (*Process)(nil)

func (p *Process) UserObj() ProcessObject {
	return p.userObj
}

/*
 *	Drawer interface
 */

func (p *Process) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		x, y, w, h := boxToDraw(p)
		var mode ColorMode
		if p.selected {
			mode = SelectedMode
		} else if p.highlighted {
			mode = HighlightMode
		} else {
			mode = NormalMode
		}
		drawBody(context, x, y, w, h, p.Name(), mode)
	}
}

/*
 *	BBoxer interface
 */

func (p *Process) BBox() image.Rectangle {
	p.box.Min = p.Position()
	p.box.Max = p.box.Min.Add(image.Point{global.processWidth, global.processHeight})
	return p.box
}

/*
 *	freesp.Positioner interface
 */

func (p *Process) Position() image.Point {
	return p.userObj.Position()
}

func (p *Process) SetPosition(pos image.Point) {
	p.userObj.SetPosition(pos)
}


/*
 *	Namer interface
 */

func (p *Process) Name() string {
	return p.userObj.Name()
}

func (p *Process) SetName(newName string) {
	p.userObj.SetName(newName)
}


