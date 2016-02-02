package graph

import (
	//"log"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/tool"
	"github.com/gotk3/gotk3/cairo"
	"image"
	"log"
)

type Node struct {
	NamedBoxObject
	userObj      bh.NodeIf
	path         string
	ports        []*Port
	selectedPort int
}

var _ NodeIf = (*Node)(nil)
var _ ContainerChild = (*Node)(nil)

func NodeNew(pos image.Point, n bh.NodeIf, path string) (ret *Node) {
	dy := NumericOption(PortDY)
	box := image.Rect(pos.X, pos.Y, pos.X+global.nodeWidth, pos.Y+global.nodeHeight+numPorts(n)*dy)
	config := DrawConfig{ColorInit(ColorOption(NodeNormal)),
		ColorInit(ColorOption(NodeHighlight)),
		ColorInit(ColorOption(NodeSelected)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{global.padX, global.padY}}
	ret = &Node{NamedBoxObjectInit(box, config, n), n, path, nil, -1}
	ret.RegisterOnHighlight(func(hit bool, pos image.Point) bool {
		return ret.onHighlight(hit, pos)
	})
	ret.RegisterOnSelect(func() bool {
		return ret.onSelect()
	}, func() bool {
		return ret.onDeselect()
	})
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

func (n Node) UserObj() bh.NodeIf {
	return n.userObj
}

func (n Node) InPorts() []bh.PortIf {
	return n.userObj.InPorts()
}

func (n Node) OutPorts() []bh.PortIf {
	return n.userObj.OutPorts()
}

func (n *Node) SelectPort(port bh.PortIf) {
	var index int
	if port.Direction() == gr.InPort {
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

func (n Node) GetSelectedPort(ownId bh.NodeIdIf) (port bh.PortIf, ok bool) {
	if n.selectedPort == -1 {
		return
	}
	ok = true
	port = n.ports[n.selectedPort].userObj
	return
}

func (n Node) NumInPorts() int {
	return len(n.userObj.InPorts())
}

func (n Node) NumOutPorts() int {
	return len(n.userObj.OutPorts())
}

func (n Node) InPortIndex(portName string) int {
	return n.userObj.InPortIndex(portName)
}

func (n Node) OutPortIndex(portName string) int {
	return n.userObj.OutPortIndex(portName)
}

func (n Node) InPort(idx int) (p BBoxer) {
	if idx >= n.NumInPorts() {
		log.Panicf("FIXME: Node.InPort(%d): index out of range.\n")
	}
	p = n.ports[idx]
	return
}

func (n Node) OutPort(idx int) (p BBoxer) {
	if idx >= n.NumOutPorts() {
		log.Panicf("FIXME: Node.OutPort(%d): index out of range.\n")
	}
	p = n.ports[idx+n.NumInPorts()]
	return
}

func (n Node) InPortByName(name string) (ret BoxedSelecter, ok bool) {
	for i := 0; i < len(n.UserObj().InPorts()); i++ {
		p := n.ports[i]
		if p.UserObj().Name() == name {
			ok = true
			ret = p
			return
		}
	}
	return
}

func (n Node) OutPortByName(name string) (ret BoxedSelecter, ok bool) {
	for i := 0; i < len(n.UserObj().OutPorts()); i++ {
		p := n.ports[i+len(n.UserObj().InPorts())]
		if p.UserObj().Name() == name {
			ok = true
			ret = p
			return
		}
	}
	return
}

func (n *Node) Expand() {
}

func (n *Node) Collapse() {
}

func (n *Node) SelectNode(ownId, selectId bh.NodeIdIf) (modified bool, node NodeIf) {
	if ownId.String() == selectId.String() {
		n.highlighted = true
		node = n
		modified = n.Select()
	} else {
		n.highlighted = false
		modified = n.Deselect()
	}
	return
}

func (n *Node) GetHighlightedNode(ownId bh.NodeIdIf) (id bh.NodeIdIf, ok bool) {
	if n.IsHighlighted() {
		id, ok = ownId, true
	}
	return
}

func (n *Node) GetSelectedNode(ownId bh.NodeIdIf) (id bh.NodeIdIf, ok bool) {
	if n.IsSelected() {
		id, ok = ownId, true
	}
	return
}

//
//		ContainerChild interface
//
func (n Node) Layout() (box image.Rectangle) {
	return n.BBox()
}

//
//      freesp.Positioner interface
//

var _ gr.Positioner = (*Node)(nil)
var _ gr.Positioner = (*Port)(nil)

// (overwrite BBoxObject default implementation)
func (n *Node) SetPosition(pos image.Point) {
	shift := pos.Sub(n.Position())
	//log.Printf("Node.SetPosition: path=%s, mode=%v, pos=%v\n", n.path, gr.PositionModeNormal, pos)
	n.userObj.SetPathModePosition(n.path, gr.PositionModeNormal, pos)
	n.BBoxDefaultSetPosition(pos)
	for _, p := range n.ports {
		p.SetPosition(p.Position().Add(shift))
	}
}

//
//      Drawer interface
//

var _ Drawer = (*Node)(nil)
var _ Drawer = (*Port)(nil)

func (n Node) Draw(ctxt interface{}) {
	n.NamedBoxDefaultDraw(ctxt)
	switch ctxt.(type) {
	case *cairo.Context:
		context := ctxt.(*cairo.Context)
		context.SetLineWidth(1)
		for _, p := range n.ports {
			p.Draw(context)
		}
	}
}

//
//	Selecter interface
//

var _ Selecter = (*Node)(nil)
var _ Selecter = (*Port)(nil)

func (n *Node) onSelect() (modified bool) {
	for i, p := range n.ports {
		if i == n.selectedPort {
			m := p.Select()
			modified = modified || m
		} else {
			m := p.Deselect()
			modified = modified || m
		}
	}
	return
}

func (n *Node) onDeselect() (modified bool) {
	for _, p := range n.ports {
		m := p.Deselect()
		modified = modified || m
	}
	return
}

//
//	Highlighter interface
//

var _ Highlighter = (*Node)(nil)
var _ Highlighter = (*Port)(nil)

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

func numPorts(n Porter) int {
	npi := len(n.InPorts())
	npo := len(n.OutPorts())
	return tool.MaxInt(npi, npo)
}
