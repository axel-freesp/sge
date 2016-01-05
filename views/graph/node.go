package graph

import (
	//"fmt"
	//"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/tool"
)

type ColorMode int

const (
	NormalMode       ColorMode = iota
	HighlightMode
	SelectedMode
	NumColorMode
)

var (
	nodeWidth  = NumericOption(NodeWidth)
	nodeHeight = NumericOption(NodeHeight)
)

type Node struct {
	SelectableObject
	userObj  NodeObject
	ports []*Port
	selectedPort int
}

var _ NodeIf = (*Node)(nil)

func NodeNew(pos image.Point, n freesp.Node) (ret *Node) {
	dy := NumericOption(PortDY)
	box := image.Rect(pos.X, pos.Y, pos.X + nodeWidth, pos.Y + nodeHeight + numPorts(n)*dy)
	ret = &Node{SelectableObject{box, false, false}, n, nil, -1}
	portBox := image.Rect(0, 0, global.portW, global.portH)
	portBox = portBox.Add(box.Min)
	shiftIn := image.Point{global.padX + global.portX0, global.padY + global.portY0}
	shiftOut := image.Point{box.Size().X - global.padX - global.portW - global.portX0, global.padY + global.portY0}
	b := portBox.Add(shiftIn)
	for _, p := range n.InPorts() {
		p := PortNew(b, p)
		ret.ports = append(ret.ports, p)
		b = b.Add(image.Point{0, global.portDY})
	}
	b = portBox.Add(shiftOut)
	for _, p := range n.OutPorts() {
		p := PortNew(b, p)
		ret.ports = append(ret.ports, p)
		b = b.Add(image.Point{0, global.portDY})
	}
	return
}

func (n *Node) UserObj() NodeObject {
	return n.userObj
}

func (n *Node) GetOutPort(index int) PortObject {
	return n.ports[index + n.NumInPorts()].UserObj()
}

func (n *Node) GetInPort(index int) PortObject {
	return n.ports[index].UserObj()
}


func (n *Node) SelectPort(userObj PortObject) {
	var index int
	if userObj.Direction() == freesp.InPort {
		index = n.InPortIndex(userObj.Name())
	} else {
		index = n.OutPortIndex(userObj.Name())
		if index >= 0 {
			index += n.NumInPorts()
		}
	}
	for i, p := range n.ports {
		if i == index {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	n.selectedPort = index
}

func (n *Node) GetSelectedPort() (ok bool, port PortObject) {
	if n.selectedPort == -1 {
		return
	}
	ok = true
	port = n.ports[n.selectedPort].userObj
	return
}

func (n *Node) NumInPorts() int {
	return n.userObj.NumInPorts()
}

func (n *Node) NumOutPorts() int {
	return n.userObj.NumOutPorts()
}

func (n *Node) InPortIndex(portName string) int {
	return n.userObj.InPortIndex(portName)
}

func (n *Node) OutPortIndex(portName string) int {
	return n.userObj.OutPortIndex(portName)
}


/*
 *      freesp.Namer implementation
 */

func (n *Node) Name() string {
	return n.userObj.Name()
}

func (n *Node) SetName(newName string) {
	n.userObj.SetName(newName)
}

var _ freesp.Namer  = (*Node)(nil)

/*
 *      freesp.Positioner interface
 */

var _ freesp.Positioner  = (*Node)(nil)
var _ freesp.Positioner  = (*Port)(nil)

// (overwrite DragableObject implementation)
func (n *Node) SetPosition(pos image.Point) {
	shift := pos.Sub(n.Position())
	(*n).userObj.SetPosition(pos)
	n.box = n.box.Add(shift)
	for _, p := range n.ports {
		p.SetPosition(p.Position().Add(shift))
	}
}


/*
 *      Drawer interface
 */

var _ Drawer  = (*Node)(nil)
var _ Drawer  = (*Port)(nil)

func (n *Node) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		x, y, w, h := boxToDraw(n)
		var mode ColorMode
		if n.selected {
			mode = SelectedMode
		} else if n.highlighted {
			mode = HighlightMode
		} else {
			mode = NormalMode
		}
		drawBody(context, x, y, w, h, n.Name(), mode)
		context.SetLineWidth(1)
		for _, p := range n.ports {
			p.Draw(context)
		}
	}
}



/*
 *	Selecter interface (overwrite DragableObject implementation)
 */

var _ Selecter  = (*Node)(nil)
var _ Selecter  = (*Port)(nil)

func (n *Node) Select() (selected bool) {
	selected = n.selected
	n.selected = true
	for i, p := range n.ports {
		if i == n.selectedPort {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	return
}

func (n *Node) Deselect() (selected bool) {
	selected = n.selected
	n.selected = false
	for _, p := range n.ports {
		p.Deselect()
	}
	return
}

func (n *Node) CheckHit(pos image.Point) (hit, modified bool) {
	test := image.Rectangle{pos, pos}
	hit = n.BBox().Overlaps(test)
	modified = n.highlighted != hit
	n.highlighted = hit
	if hit {
		n.selectedPort = -1
		for i, p := range n.ports {
			phit, mod := p.CheckHit(pos)
			if phit {
				n.selectedPort = i
			}
			modified = modified || mod
		}
	}
	return
}


/*
 *	Private functions
 */

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

func size(b BBoxer) (w, h float64) {
	box := b.BBox()
	r := box.Size()
	w = float64(r.X)
	h = float64(r.Y)
	return
}

func boxToDraw(b BBoxer) (x, y, w, h float64) {
	dx, dy := global.padX, global.padY
	w0, h0 := size(b)
	a := b.BBox().Min
	x = float64(a.X) + float64(dx)
	y = float64(a.Y) + float64(dy)
	w = w0 - float64(2*dx)
	h = h0 - float64(2*dy)
	return
}

func numPorts(n freesp.Porter) int {
	npi := n.NumInPorts()
	npo := n.NumOutPorts()
	return tool.MaxInt(npi, npo)
}



