package graph

import (
	"image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
)

type Port struct {
	SelectableBox
	userObj interfaces.PortObject
}

func PortNew(box image.Rectangle, userObj interfaces.PortObject) *Port {
	return &Port{SelectableBoxInit(box), userObj}
}

var _ PortIf = (*Port)(nil)

func (p Port) UserObj() interfaces.PortObject {
	return p.userObj
}

//
//	Drawer interface
//

func (p Port) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		var r, g, b float64
		switch p.userObj.Direction() {
		case interfaces.InPort:
			if p.IsSelected() {
				r, g, b = ColorOption(SelectInPort)
			} else if p.IsHighlighted() {
				r, g, b = ColorOption(HighlightInPort)
			} else {
				r, g, b = ColorOption(InputPort)
			}
		case interfaces.OutPort:
			if p.IsSelected() {
				r, g, b = ColorOption(SelectOutPort)
			} else if p.IsHighlighted() {
				r, g, b = ColorOption(HighlightOutPort)
			} else {
				r, g, b = ColorOption(OutputPort)
			}
		}
		x := float64(p.BBox().Min.X)
		y := float64(p.BBox().Min.Y)
		w := float64(p.BBox().Size().X)
		h := float64(p.BBox().Size().Y)
		context.SetSourceRGB(r, g, b)
		context.Rectangle(x, y, w, h)
		context.FillPreserve()
		context.SetSourceRGB(ColorOption(BoxFrame))
		context.Stroke()
	}
}

