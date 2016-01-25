package graph

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	//"github.com/gotk3/gotk3/cairo"
	"image"
	"log"
)

type ExpandedNode struct {
	Container
	userObj bh.NodeIf
}

var _ NodeIf = (*ExpandedNode)(nil)

func ExpandedNodeNew(pos image.Point, userObj bh.NodeIf) (ret *ExpandedNode) {
	config := DrawConfig{ColorInit(ColorOption(ProcessNormal)),
		ColorInit(ColorOption(ProcessHighlight)),
		ColorInit(ColorOption(ProcessSelected)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{procPortWidth, procPortHeight}}
	cconfig := ContainerConfig{global.portW, global.portH, 120, 80}
	ret = &ExpandedNode{ContainerInit(nil, config, userObj, cconfig), userObj}
	dy := NumericOption(PortDY)
	shape := image.Point{global.nodeWidth, global.nodeHeight + numPorts(userObj)*dy}
	ret.box = image.Rectangle{pos, pos.Add(shape)}
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
		ret.AddModePort(ret.portClipPos(pos), config, p, gr.PositionModeExpanded)
	}
	for i, p := range userObj.OutPorts() {
		pos := p.ModePosition(gr.PositionModeExpanded)
		if pos == empty {
			pos = ret.CalcOutPortPos(i)
		}
		ret.AddModePort(ret.portClipPos(pos), config, p, gr.PositionModeExpanded)
	}
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

func (n ExpandedNode) SelectPort(port bh.PortIf) {
	n.SelectModePort(port)
}

func (n ExpandedNode) GetSelectedPort() (ok bool, port bh.PortIf) {
	if n.selectedPort == -1 {
		return
	}
	ok = true
	port = n.ports[n.selectedPort].UserObj2.(bh.PortIf)
	return
}

func (n *ExpandedNode) SetPosition(pos image.Point) {
	n.ContainerDefaultSetPosition(pos)
	n.userObj.SetModePosition(gr.PositionModeExpanded, pos)
}
