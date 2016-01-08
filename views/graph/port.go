package graph

import (
	"image"
	"github.com/gotk3/gotk3/cairo"
	"github.com/axel-freesp/sge/freesp"
)

type Port struct {
	SelectableBox
	userObj PortObject
}

func PortNew(box image.Rectangle, userObj PortObject) *Port {
	return &Port{SelectableBoxInit(box), userObj}
}

var _ PortIf = (*Port)(nil)

func (p *Port) UserObj() PortObject {
	return p.userObj
}

/*
 *	Drawer interface
 */

func (p *Port) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		var r, g, b float64
		switch p.userObj.(freesp.Directioner).Direction() {
		case freesp.InPort:
			if p.selected {
				r, g, b = ColorOption(SelectInPort)
			} else if p.highlighted {
				r, g, b = ColorOption(HighlightInPort)
			} else {
				r, g, b = ColorOption(InputPort)
			}
		case freesp.OutPort:
			if p.selected {
				r, g, b = ColorOption(SelectOutPort)
			} else if p.highlighted {
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

/*
 *	freesp.Positioner interface
 */

func (p *Port) Position() image.Point {
	return p.BBox().Min
}

func (p *Port) SetPosition(pos image.Point) {
	size := p.BBox().Size()
	p.box.Min = pos
	p.box.Max = p.box.Min.Add(size)
}


