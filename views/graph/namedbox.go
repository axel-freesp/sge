package graph

import (
	"github.com/axel-freesp/sge/interface/graph"
	"github.com/gotk3/gotk3/cairo"
	"image"
)

type NamerObject struct {
	namer graph.Namer
}

func NamerObjectInit(namer graph.Namer) NamerObject {
	return NamerObject{namer}
}

func (n NamerObject) Name() string {
	return n.namer.Name()
}

func (n *NamerObject) SetName(newName string) {
	n.namer.SetName(newName)
}

type NamedBoxObject struct {
	SelectableBox
	NamerObject
}

var _ NamedBox = (*NamedBoxObject)(nil)

func NamedBoxObjectInit(box image.Rectangle, config DrawConfig, namer graph.Namer) NamedBoxObject {
	return NamedBoxObject{SelectableBoxInit(box, config), NamerObjectInit(namer)}
}

func (b NamedBoxObject) Draw(ctxt interface{}) {
	b.NamedBoxDefaultDraw(ctxt)
}

func (b NamedBoxObject) NamedBoxDefaultDraw(ctxt interface{}) {
	b.SelectorDefaultDraw(ctxt)
	switch ctxt.(type) {
	case *cairo.Context:
		context := ctxt.(*cairo.Context)
		x, y, _, _ := boxToDraw(&b, b.config.pad)
		context.SetSourceRGB(b.config.tCol.r, b.config.tCol.g, b.config.tCol.b)
		context.SetFontSize(float64(global.fontSize))
		tx, ty := float64(global.textX), float64(global.textY)
		context.MoveTo(x+tx, y+ty)
		context.ShowText(b.namer.Name())
	}
}
