package graph

import (
	//"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	"github.com/axel-freesp/sge/tool"
	interfaces "github.com/axel-freesp/sge/interface"
)

type ColorMode int

const (
	NormalMode       ColorMode = iota
	HighlightMode
	SelectedMode
	NumColorMode
)

type Node struct {
	SelectableBox
	userObj  interfaces.NodeObject
	ports []*Port
	selectedPort int
}

var _ NodeIf = (*Node)(nil)

func NodeNew(pos image.Point, n interfaces.NodeObject) (ret *Node) {
	dy := NumericOption(PortDY)
	box := image.Rect(pos.X, pos.Y, pos.X + global.nodeWidth, pos.Y + global.nodeHeight + numPorts(n)*dy)
	ret = &Node{SelectableBoxInit(box), n, nil, -1}
	ret.RegisterOnHighlight(func(hit bool, pos image.Point) bool{
		return ret.onHighlight(hit, pos)
	})
	ret.RegisterOnSelect(func(){
		ret.onSelect()
	}, func(){
		ret.onDeselect()
	})
	portBox := image.Rect(0, 0, global.portW, global.portH)
	portBox = portBox.Add(box.Min)
	shiftIn := image.Point{global.padX + global.portX0, global.padY + global.portY0}
	shiftOut := image.Point{box.Size().X - global.padX - global.portW - global.portX0, global.padY + global.portY0}
	b := portBox.Add(shiftIn)
	for _, p := range n.GetInPorts() {
		p := PortNew(b, p)
		ret.ports = append(ret.ports, p)
		b = b.Add(image.Point{0, global.portDY})
	}
	b = portBox.Add(shiftOut)
	for _, p := range n.GetOutPorts() {
		p := PortNew(b, p)
		ret.ports = append(ret.ports, p)
		b = b.Add(image.Point{0, global.portDY})
	}
	return
}

func (n Node) UserObj() interfaces.NodeObject {
	return n.userObj
}

func (n Node) GetInPorts() []interfaces.PortObject {
	return n.userObj.GetInPorts()
}

func (n Node) GetOutPorts() []interfaces.PortObject {
	return n.userObj.GetOutPorts()
}

func (n *Node) SelectPort(port interfaces.PortObject) {
	var index int
	if port.Direction() == interfaces.InPort {
		index = n.InPortIndex(port.Name())
	} else {
		index = n.OutPortIndex(port.Name())
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

func (n Node) GetSelectedPort() (ok bool, port interfaces.PortObject) {
	if n.selectedPort == -1 {
		return
	}
	ok = true
	port = n.ports[n.selectedPort].userObj
	return
}

func (n Node) NumInPorts() int {
	return len(n.userObj.GetInPorts())
}

func (n Node) NumOutPorts() int {
	return len(n.userObj.GetOutPorts())
}

func (n Node) InPortIndex(portName string) int {
	return n.userObj.InPortIndex(portName)
}

func (n Node) OutPortIndex(portName string) int {
	return n.userObj.OutPortIndex(portName)
}


//
//      freesp.Namer implementation
//

func (n Node) Name() string {
	return n.userObj.Name()
}

func (n Node) SetName(newName string) {
	n.userObj.SetName(newName)
}

var _ interfaces.Namer  = (*Node)(nil)


//
//      freesp.Positioner interface
//

var _ interfaces.Positioner  = (*Node)(nil)
var _ interfaces.Positioner  = (*Port)(nil)

// (overwrite BBoxObject default implementation)
func (n *Node) SetPosition(pos image.Point) {
	shift := pos.Sub(n.Position())
	n.userObj.SetPosition(pos)
	n.DefaultSetPosition(pos)
	for _, p := range n.ports {
		p.SetPosition(p.Position().Add(shift))
	}
}


//
//      Drawer interface
//

var _ Drawer  = (*Node)(nil)
var _ Drawer  = (*Port)(nil)

func (n Node) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		x, y, w, h := boxToDraw(&n)
		var mode ColorMode
		if n.IsSelected() {
			mode = SelectedMode
		} else if n.IsHighlighted() {
			mode = HighlightMode
		} else {
			mode = NormalMode
		}
		nodeDrawBody(context, x, y, w, h, n.Name(), mode)
		context.SetLineWidth(1)
		for _, p := range n.ports {
			p.Draw(context)
		}
	}
}


//
//	Selecter interface
//

var _ Selecter  = (*Node)(nil)
var _ Selecter  = (*Port)(nil)

func (n *Node) onSelect() (selected bool) {
	for i, p := range n.ports {
		if i == n.selectedPort {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	return
}

func (n *Node) onDeselect() (selected bool) {
	for _, p := range n.ports {
		p.Deselect()
	}
	return
}


//
//	Highlighter interface
//

var _ Highlighter  = (*Node)(nil)
var _ Highlighter  = (*Port)(nil)

func (n *Node) onHighlight(hit bool, pos image.Point) (modified bool) {
	n.selectedPort = -1
	if hit {
		for i, p := range n.ports {
			phit, mod := p.CheckHit(pos)
			if phit {
				n.selectedPort = i
			}
			modified = modified || mod
		}
	} else {
		for _, p := range n.ports {
			p.CheckHit(pos)
		}
	}
	return
}


//
//	Private functions
//

func nodeDrawBody(context *cairo.Context, x, y, w, h float64, name string, mode ColorMode) {
	context.SetSourceRGB(nodeChooseColor(mode))
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

func nodeChooseColor(mode ColorMode) (r, g, b float64) {
	switch mode {
	case NormalMode:
		r, g, b = ColorOption(NodeNormal)
	case HighlightMode:
		r, g, b = ColorOption(NodeHighlight)
	case SelectedMode:
		r, g, b = ColorOption(NodeSelected)
	}
	return
}

func numPorts(n interfaces.Porter) int {
	npi := len(n.GetInPorts())
	npo := len(n.GetOutPorts())
	return tool.MaxInt(npi, npo)
}



