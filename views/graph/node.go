package graph

import (
	//"fmt"
	//"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	"github.com/axel-freesp/sge/freesp"
)

type Global struct {
	padX, padY,
	portX0, portY0, portDY,
	portW, portH,
	textX, textY, fontSize int
}

var global = Global{
	padX: NumericOption(NodePadX),
	padY: NumericOption(NodePadY),
	portX0: NumericOption(PortX0),
	portY0: NumericOption(PortY0),
	portDY: NumericOption(PortDY),
	portW: NumericOption(PortW),
	portH: NumericOption(PortH),
	textX: NumericOption(NodeTextX),
	textY: NumericOption(NodeTextY),
	fontSize: NumericOption(FontSize),
}

type Node struct {
	DragableObject
	highlightInport, highlightOutport int
}

var _ Dragable = (*Node)(nil)

func NodeNew(box image.Rectangle, userObj freesp.Node) *Node {
	return &Node{DragableObject{box, userObj}, -1, -1}
}

func (n *Node) Draw(context *cairo.Context, mode ColorMode){
	x, y, w, h := boxToDraw(n)
	x0, y0, dy, pw, ph := portGeom()
	drawBody(context, x, y, w, h, n.Name(), mode)
	context.SetLineWidth(1)
	for i := 0; i < n.NumInPorts(); i++ {
		if i == n.highlightInport {
			context.SetSourceRGB(ColorOption(HighlightInPort))
		} else {
			context.SetSourceRGB(ColorOption(InputPort))
		}
		context.Rectangle(x + x0, y + y0 + float64(i)*dy, pw, ph)
		context.FillPreserve()
		context.SetSourceRGB(ColorOption(BoxFrame))
		context.Stroke()
	}
	for i := 0; i < n.NumOutPorts(); i++ {
		if i == n.highlightOutport {
			context.SetSourceRGB(ColorOption(HighlightOutPort))
		} else {
			context.SetSourceRGB(ColorOption(OutputPort))
		}
		context.Rectangle(x + w - pw - x0, y + y0 + float64(i)*dy, pw, ph)
		context.FillPreserve()
		context.SetSourceRGB(ColorOption(BoxFrame))
		context.Stroke()
	}
}

func (n *Node) Check(pos image.Point) (rect image.Rectangle, ok bool) {
	box := n.BBox()
	dx, dy := global.padX, global.padY
	w0, h0 := box.Size().X, box.Size().Y
	x := box.Min.X + dx
	y := box.Min.Y + dy
	w, _ := w0 - 2*dx, h0 - 2*dy
	x0 := NumericOption(PortX0)
	y0, dy := global.portY0, global.portDY
	pw, ph := global.portW, global.portH

	d := image.Point{pw, ph}
	p := image.Rectangle{pos, pos}

	oldHiIn := n.highlightInport
	oldHiOut := n.highlightOutport
	n.highlightInport = -1
	n.highlightOutport = -1
	var oldRect image.Rectangle

	for i := 0; i < n.NumInPorts(); i++ {
		p0 := image.Point{x + x0, y + y0 + i*dy}
		p1 := p0.Add(d)
		r := image.Rectangle{p0, p1}
		if i == oldHiIn {
			oldRect = r
		}
		if p.Overlaps(r) {
			ok = true
			rect = r
			n.highlightInport = i
		}
	}
	for i := 0; i < n.NumOutPorts(); i++ {
		p0 := image.Point{x + w - pw - x0, y + y0 + i*dy}
		p1 := p0.Add(d)
		r := image.Rectangle{p0, p1}
		if i == oldHiOut {
			oldRect = r
		}
		if p.Overlaps(r){
			ok = true
			rect = r
			n.highlightOutport = i
		}
	}
	if oldRect.Size().X > 0 {
		if rect.Size().X > 0 {
			rect = rect.Union(oldRect)
		} else {
			rect = oldRect
		}
		ok = true
	}
	return
}

func size(drg Dragable) (w, h float64) {
	box := drg.BBox()
	r := box.Size()
	w = float64(r.X)
	h = float64(r.Y)
	return
}

func boxToDraw(drg Dragable) (x, y, w, h float64) {
	dx, dy := global.padX, global.padY
	w0, h0 := size(drg)
	x = float64(drg.BBox().Min.X) + float64(dx)
	y = float64(drg.BBox().Min.Y) + float64(dy)
	w, h = w0 - float64(2*dx), h0 - float64(2*dy)
	return
}

func portGeom() (x0, y0, dy, w, h float64) {
	x0 = float64(global.portX0)
	y0, dy = float64(global.portY0), float64(global.portDY)
	w, h = float64(global.portW), float64(global.portH)
	return
}

func drawBody(context *cairo.Context, x, y, w, h float64, name string, mode ColorMode) {
	context.SetSourceRGB(chooseColor(mode))
	context.Rectangle(x, y, w, h)
	context.FillPreserve()
	context.SetSourceRGB(ColorOption(BoxFrame))
	context.SetLineWidth(2)
	context.Stroke()
	context.SetSourceRGB(ColorOption(Text))
	context.SetFontSize(float64(global.fontSize))
	tx, ty := float64(global.textX), float64(global.textY)
	context.MoveTo(x + tx, y + ty)
	context.ShowText(name)
}

func chooseColor(mode ColorMode) (r, g, b float64) {
	switch mode {
	case NormalMode:
		r, g, b = ColorOption(Normal)
	case HighlightMode:
		r, g, b = ColorOption(Highlight)
	case SelectedMode:
		r, g, b = ColorOption(Selected)
	}
	return
}



