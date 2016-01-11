package graph

import (
    "image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
)

type NamerObject struct {
    namer interfaces.Namer
}

func NamerObjectInit(namer interfaces.Namer) NamerObject {
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
    tCol Color
}

var _ NamedBox = (*NamedBoxObject)(nil)

func NamedBoxObjectInit(box image.Rectangle, nCol, hCol, sCol, fCol, tCol Color, pad image.Point, namer interfaces.Namer) NamedBoxObject {
    return NamedBoxObject{SelectableBoxInit(box, nCol, hCol, sCol, fCol, pad), NamerObjectInit(namer), tCol}
}

func (b NamedBoxObject) Draw(ctxt interface{}){
	b.DrawDefaultText(ctxt)
}

func (b NamedBoxObject) DrawDefaultText(ctxt interface{}){
	b.DrawDefault(ctxt)
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		x, y, _, _ := boxToDraw(&b, b.pad)
		context.SetSourceRGB(b.tCol.r, b.tCol.g, b.tCol.b)
		context.SetFontSize(float64(global.fontSize))
		tx, ty := float64(global.textX), float64(global.textY)
		context.MoveTo(x + tx, y + ty)
		context.ShowText(b.namer.Name())
	}
}



