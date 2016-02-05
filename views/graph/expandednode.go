package graph

import (
	"fmt"
	freesp "github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	"github.com/gotk3/gotk3/cairo"
	"image"
	"log"
)

type PortConnector struct {
	LineObject
	port1, port2 BoxedSelecter
}

func PortConnectorNew(port1, port2 BoxedSelecter) *PortConnector {
	b1, b2 := port1.BBox(), port2.BBox()
	p1, p2 := b1.Min.Add(b1.Size().Div(2)), b2.Min.Add(b2.Size().Div(2))
	return &PortConnector{LineObjectInit(p1, p2), port1, port2}
}

func (p *PortConnector) Points() (p1, p2 image.Point) {
	b1, b2 := p.port1.BBox(), p.port2.BBox()
	p.p1, p.p2 = b1.Min.Add(b1.Size().Div(2)), b2.Min.Add(b2.Size().Div(2))
	p1, p2 = p.p1, p.p2
	return
}

func (p PortConnector) Draw(ctxt interface{}) {
	switch ctxt.(type) {
	case *cairo.Context:
		context := ctxt.(*cairo.Context)
		var color Color
		if p.IsSelected() {
			color = ColorInit(ColorOption(SelectLine))
		} else if p.IsHighlighted() {
			color = ColorInit(ColorOption(HighlightLine))
		} else {
			color = ColorInit(ColorOption(NormalLine))
		}
		p1, p2 := p.Points()
		context.SetSourceRGB(color.r, color.g, color.b)
		DrawArrow(context, p1, p2)
	}
}

type ExpandedNode struct {
	Container
	userObj     bh.NodeIf
	positioner  gr.Positioner
	connections []ConnectIf
	portconn    []*PortConnector
}

var _ NodeIf = (*ExpandedNode)(nil)

const (
	expandedPortWidth  = 10
	expandedPortHeight = 10
)

func ExpandedNodeNew(getPositioner GetPositioner, userObj bh.NodeIf, path string) (ret *ExpandedNode) {
	positioner := getPositioner(userObj, path)
	pos := positioner.Position()
	config := DrawConfig{ColorInit(ColorOption(NormalExpandedNode)),
		ColorInit(ColorOption(HighlightExpandedNode)),
		ColorInit(ColorOption(SelectExpandedNode)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{global.padX, global.padY}}
	cconfig := ContainerConfig{expandedPortWidth, expandedPortHeight, 120, 80}
	// Add children
	var g bh.SignalGraphTypeIf
	nt := userObj.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			g = impl.Graph()
			break
		}
	}
	var children []ContainerChild
	if g != nil {
		empty := image.Point{}
		first := image.Point{16, 32}
		shift := image.Point{16, 16}
		for i, n := range g.ProcessingNodes() {
			var ch ContainerChild
			var nPath string
			var mode gr.PositionMode
			if len(path) > 0 {
				nPath = fmt.Sprintf("%s/%s", path, userObj.Name())
			} else {
				nPath = userObj.Name()
			}
			n.SetActivePath(nPath)
			if n.Expanded() {
				mode = gr.PositionModeExpanded
			} else {
				mode = gr.PositionModeNormal
			}
			n.SetActiveMode(mode)
			log.Printf("ExpandedNodeNew TODO: position of child nodes...\n")
			chpos := n.PathModePosition(nPath, mode)
			if chpos == empty {
				chpos = pos.Add(first.Add(shift.Mul(i)))
			}
			if n.Expanded() {
				ch = ExpandedNodeNew(getPositioner, n, nPath)
			} else {
				ch = NodeNew(getPositioner, n, nPath)
			}
			children = append(children, ch)
		}
	}
	ret = &ExpandedNode{ContainerInit(children, config, userObj, cconfig),
		userObj, positioner, nil, nil}
	ret.ContainerInit()
	empty := image.Point{}
	config = DrawConfig{ColorInit(ColorOption(InputPort)),
		ColorInit(ColorOption(HighlightInPort)),
		ColorInit(ColorOption(SelectInPort)),
		ColorInit(ColorOption(BoxFrame)),
		Color{},
		image.Point{}}
	for i, p := range userObj.InPorts() {
		pos := p.ModePosition(gr.PositionModeExpanded)
		if pos == empty {
			pos = ret.CalcInPortPos(i)
		}
		positioner := gr.ModePositionerProxyNew(p, gr.PositionModeExpanded)
		ret.AddPort(config, p, positioner)
	}
	config = DrawConfig{ColorInit(ColorOption(OutputPort)),
		ColorInit(ColorOption(HighlightOutPort)),
		ColorInit(ColorOption(SelectOutPort)),
		ColorInit(ColorOption(BoxFrame)),
		Color{},
		image.Point{}}
	for i, p := range userObj.OutPorts() {
		pos := p.ModePosition(gr.PositionModeExpanded)
		if pos == empty {
			pos = ret.CalcOutPortPos(i)
		}
		positioner := gr.ModePositionerProxyNew(p, gr.PositionModeExpanded)
		ret.AddPort(config, p, positioner)
	}
	for _, n := range g.ProcessingNodes() {
		from, ok := ret.ChildByName(n.Name())
		if !ok {
			log.Printf("ExpandedNodeNew error: node %s not found\n", n.Name())
			continue
		}
		for _, p := range n.OutPorts() {
			fromId := from.OutPortIndex(p.Name())
			for _, c := range p.Connections() {
				to, ok := ret.ChildByName(c.Node().Name())
				if ok {
					toId := to.InPortIndex(c.Name())
					ret.connections = append(ret.connections, ConnectionNew(from, to, fromId, toId))
				} else {
					portname, ok := c.Node().PortLink()
					if !ok {
						log.Printf("ExpandedNodeNew error: output node %s not linked\n", c.Node().Name())
						continue
					}
					ownPort, ok := ret.OutPortByName(portname)
					if !ok {
						log.Printf("ExpandedNodeNew error: linked port %s of output node %s not found\n", portname, c.Node().Name())
						continue
					}
					nodePort, ok := from.OutPortByName(p.Name())
					if !ok {
						log.Printf("ExpandedNodeNew error: port %s of output node %s not found\n", p.Name(), from.Name())
						continue
					}
					ret.portconn = append(ret.portconn, PortConnectorNew(nodePort, ownPort))
				}
			}
		}
	}
	for _, n := range g.InputNodes() {
		for _, p := range n.OutPorts() {
			fromlink, ok := p.Node().PortLink()
			if !ok {
				log.Printf("ExpandedNodeNew error: input node %s not linked\n", p.Node().Name())
				continue
			}
			fromPort, ok := ret.InPortByName(fromlink)
			if !ok {
				log.Printf("ExpandedNodeNew error: linked port %s of input node %s not found\n", fromlink, n.Name())
				continue
			}
			for _, c := range p.Connections() {
				to, ok := ret.ChildByName(c.Node().Name())
				if ok {
					// TODO: connect with node
					toPort, ok := to.InPortByName(c.Name())
					if !ok {
						log.Printf("ExpandedNodeNew error: port %s of node %s not found\n", c.Name(), to.Name())
						continue
					}
					ret.portconn = append(ret.portconn, PortConnectorNew(fromPort, toPort))
				} else {
					tolink, ok := c.Node().PortLink()
					if !ok {
						log.Printf("ExpandedNodeNew error: output node %s not linked\n", c.Node().Name())
						continue
					}
					toPort, ok := ret.OutPortByName(tolink)
					if !ok {
						log.Printf("ExpandedNodeNew error: linked port %s of output node %s not found\n", tolink, c.Node().Name())
						continue
					}
					ret.portconn = append(ret.portconn, PortConnectorNew(fromPort, toPort))
				}
			}
		}
	}
	ret.RegisterOnDraw(func(ctxt interface{}) {
		expandedNodeOnDraw(ret, ctxt)
	})
	return
}

func (n ExpandedNode) CalcInPortPos(index int) (pos image.Point) {
	start := image.Point{global.padX + global.portX0, global.padY + global.portY0}
	shift := image.Point{0, global.portDY}
	pos = n.box.Min.Add(start)
	for i := 0; i < index; i++ {
		pos = pos.Add(shift)
	}
	return
}

func (n ExpandedNode) CalcOutPortPos(index int) (pos image.Point) {
	start := image.Point{n.box.Size().X - global.padX - global.portW - global.portX0, global.padY + global.portY0}
	shift := image.Point{0, global.portDY}
	pos = n.box.Min.Add(start)
	for i := 0; i < index; i++ {
		pos = pos.Add(shift)
	}
	return
}

func (n ExpandedNode) UserObj() bh.NodeIf {
	return n.userObj
}

func (n ExpandedNode) InPorts() []bh.PortIf {
	return n.userObj.InPorts()
}

func (n ExpandedNode) OutPorts() []bh.PortIf {
	return n.userObj.OutPorts()
}

func (n ExpandedNode) NumInPorts() int {
	return len(n.userObj.InPorts())
}

func (n ExpandedNode) InPortIndex(portName string) int {
	return n.userObj.InPortIndex(portName)
}

func (n ExpandedNode) OutPortIndex(portName string) int {
	return n.userObj.OutPortIndex(portName)
}

func (n *ExpandedNode) Expand() {
}

func (n *ExpandedNode) Collapse() {
}

func (n ExpandedNode) InPort(idx int) (p BBoxer) {
	if idx >= len(n.UserObj().InPorts()) {
		log.Panicf("FIXME: Node.InPort(%d): index out of range.\n")
	}
	p = n.ports[idx]
	return
}

func (n ExpandedNode) OutPort(idx int) (p BBoxer) {
	if idx >= len(n.UserObj().InPorts())+len(n.UserObj().OutPorts()) {
		log.Panicf("FIXME: Node.OutPort(%d): index out of range.\n")
	}
	p = n.ports[idx+len(n.UserObj().InPorts())]
	return
}

func (n ExpandedNode) GetSelectedPort(ownId bh.NodeIdIf) (port bh.PortIf, ok bool) {
	if !n.IsHighlighted() {
		return
	}
	for _, ch := range n.Children {
		chId := freesp.NodeIdNew(ownId, ch.(NodeIf).Name())
		port, ok = ch.(NodeIf).GetSelectedPort(chId)
		if ok {
			return
		}
	}
	if n.selectedPort == -1 {
		return
	}
	ok = true
	port = n.ports[n.selectedPort].UserObj.(bh.PortIf)
	return
}

func (n *ExpandedNode) SetPosition(pos image.Point) {
	n.ContainerDefaultSetPosition(pos)
	n.positioner.SetPosition(pos)
}

func (n ExpandedNode) InPortByName(name string) (ret BoxedSelecter, ok bool) {
	for i := 0; i < len(n.UserObj().InPorts()); i++ {
		p := n.ports[i]
		if p.UserObj.(bh.PortIf).Name() == name {
			ok = true
			ret = p
			return
		}
	}
	return
}

func (n ExpandedNode) OutPortByName(name string) (ret BoxedSelecter, ok bool) {
	for i := 0; i < len(n.UserObj().OutPorts()); i++ {
		p := n.ports[i+len(n.UserObj().InPorts())]
		if p.UserObj.(bh.PortIf).Name() == name {
			ok = true
			ret = p
			return
		}
	}
	return
}

func (n ExpandedNode) ChildByName(name string) (chn NodeIf, ok bool) {
	for _, ch := range n.Children {
		if ch.(NodeIf).Name() == name {
			chn = ch.(NodeIf)
			ok = true
			return
		}
	}
	return
}

func (n *ExpandedNode) SelectNode(ownId, selectId bh.NodeIdIf) (modified bool, node NodeIf) {
	if ownId.String() == selectId.String() || ownId.IsAncestor(selectId) {
		n.highlighted = true
		node = n
		modified = n.Select()
	} else {
		n.highlighted = false
		modified = n.Deselect()
	}
	for _, ch := range n.Children {
		nn := ch.(NodeIf)
		chId := freesp.NodeIdNew(ownId, nn.Name())
		m, nd := nn.SelectNode(chId, selectId)
		if nd != nil {
			node = nd
		}
		modified = modified || m
	}
	return
}

func (n ExpandedNode) GetHighlightedNode(ownId bh.NodeIdIf) (id bh.NodeIdIf, ok bool) {
	if !n.IsHighlighted() {
		return
	}
	for _, ch := range n.Children {
		chId := freesp.NodeIdNew(ownId, ch.(NodeIf).Name())
		id, ok = ch.(NodeIf).GetHighlightedNode(chId)
		if ok {
			return
		}
	}
	id, ok = ownId, true
	return
}

func (n ExpandedNode) GetSelectedNode(ownId bh.NodeIdIf) (id bh.NodeIdIf, ok bool) {
	if !n.IsSelected() {
		return
	}
	for _, ch := range n.Children {
		chId := freesp.NodeIdNew(ownId, ch.(NodeIf).Name())
		id, ok = ch.(NodeIf).GetSelectedNode(chId)
		if ok {
			return
		}
	}
	id, ok = ownId, true
	return
}

func (n ExpandedNode) ChildNodes() (nodelist []NodeIf) {
	for _, c := range n.Children {
		nodelist = append(nodelist, c.(NodeIf))
	}
	return
}

func (n *ExpandedNode) SelectPort(port bh.PortIf) {
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

//
//	Drawer interface
//

func expandedNodeOnDraw(n *ExpandedNode, ctxt interface{}) {
	for _, conn := range n.connections {
		conn.Draw(ctxt)
	}
	for _, conn := range n.portconn {
		conn.Draw(ctxt)
	}
}
