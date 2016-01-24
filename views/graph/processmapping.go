package graph

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/gotk3/gotk3/cairo"
	"image"
	"log"
)

type ProcessMapping struct {
	Container
	userObj      pf.ProcessIf
	nodes        []*Node
	selectedPort int
	arch         ArchIf
}

var _ ProcessIf = (*ProcessMapping)(nil)
var _ ContainerChild = (*Node)(nil)

func ProcessMappingNew(nodes []*Node, userObj pf.ProcessIf) (ret *ProcessMapping) {
	config := DrawConfig{ColorInit(ColorOption(ProcessNormal)),
		ColorInit(ColorOption(ProcessHighlight)),
		ColorInit(ColorOption(ProcessSelected)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{procPortOutBorder, procPortOutBorder}}
	cconfig := ContainerConfig{procPortWidth, procPortHeight, procMinWidth, procMinHeight}
	var children []ContainerChild
	for _, n := range nodes {
		children = append(children, n)
	}
	ret = &ProcessMapping{ContainerInit(children, config, userObj, cconfig), userObj, nodes, -1, nil}
	ret.ContainerInit()
	return
}

func (p ProcessMapping) UserObj() pf.ProcessIf {
	return p.userObj
}

func (pr ProcessMapping) ArchObject() ArchIf {
	return pr.arch
}

func processDrawMappedChannel(p *ContainerPort, ctxt interface{}) {
	switch ctxt.(type) {
	case *cairo.Context:
		empty := image.Point{}
		context := ctxt.(*cairo.Context)
		shift1 := image.Point{archPortWidth, archPortHeight}.Div(2)
		shift2 := image.Point{procPortWidth, procPortHeight}.Div(2)
		extPort := p.UserObj2.(pf.ChannelIf).ArchPort()
		if extPort != nil {
			extPPos := extPort.ModePosition(gr.PositionModeMapping)
			if extPPos != empty {
				var pos1, pos2 image.Point
				if p.UserObj2.(pf.ChannelIf).Direction() == gr.InPort {
					pos1 = extPPos.Add(shift1)
					pos2 = p.Position().Add(shift2)
				} else {
					pos1 = p.Position().Add(shift2)
					pos2 = extPPos.Add(shift1)
				}
				var r, g, b float64
				if p.IsSelected() {
					r, g, b = ColorOption(SelectChannelLine)
				} else if p.IsHighlighted() {
					r, g, b = ColorOption(HighlightChannelLine)
				} else {
					r, g, b = ColorOption(NormalChannelLine)
				}
				context.SetLineWidth(2)
				context.SetSourceRGB(r, g, b)
				DrawArrow(context, pos1, pos2)
			}
		}
	}
}

func (pr *ProcessMapping) SetArchObject(a ArchIf) {
	pr.arch = a
	idx := 0
	for _, c := range pr.userObj.InChannels() {
		pr.addPort(c, a, idx)
		idx++
	}
	for _, c := range pr.userObj.OutChannels() {
		pr.addPort(c, a, idx)
		idx++
	}
}

func (pr *ProcessMapping) addPort(c pf.ChannelIf, a ArchIf, idx int) {
	cnt := len(pr.userObj.InChannels()) + len(pr.userObj.OutChannels())
	empty := image.Point{}
	config := DrawConfig{ColorInit(ColorOption(NormalArchPort)),
		ColorInit(ColorOption(HighlightArchPort)),
		ColorInit(ColorOption(SelectArchPort)),
		ColorInit(ColorOption(BoxFrame)),
		Color{},
		image.Point{}}
	pos := c.ModePosition(gr.PositionModeMapping)
	if pos == empty {
		pos = pr.CalcPortPos(idx, cnt)
	}
	p := pr.AddModePort(pos, config, c, gr.PositionModeMapping)
	p.RegisterOnDraw(func(ctxt interface{}) {
		processDrawMappedChannel(p, ctxt)
	})
	a.(*Arch).channelMap[c] = p
	//log.Printf("ProcessMapping.addPort: pos=%v\n", p.Position())
}

func (pr *ProcessMapping) SelectChannel(ch pf.ChannelIf) {
	pr.SelectModePort(ch)
}

func (pr ProcessMapping) GetSelectedChannel() (ok bool, ch pf.ChannelIf) {
	var c gr.ModePositioner
	ok, c = pr.GetSelectedModePort()
	if !ok {
		return
	}
	ch = c.(pf.ChannelIf)
	return
}

func (pr *ProcessMapping) SelectNode(n bh.NodeIf) {
	for _, ch := range pr.Children {
		if n == ch.(NodeIf).UserObj() {
			pr.SelectChild(ch)
			return
		}
	}
}

func (pr ProcessMapping) GetSelectedNode() (ok bool, n NodeIf) {
	var ch ContainerChild
	ok, ch = pr.GetSelectedChild()
	if !ok {
		return
	}
	n = ch.(NodeIf)
	return
}

func (pr *ProcessMapping) SetPosition(pos image.Point) {
	pr.ContainerDefaultSetPosition(pos)
	pr.userObj.SetModePosition(gr.PositionModeMapping, pos)
	if pr.arch != nil { // else is unmapped process
		a := pr.arch.(*Arch)
		if a.mapping != nil {
			for _, ch := range pr.Children {
				n := ch.(*Node)
				melem, ok := a.mapping.MappedElement(n.UserObj())
				if !ok {
					log.Printf("Arch.SetPosition Warning: node %s not mapped\n", n.Name())
					return
				}
				melem.SetPosition(n.Position())
			}
		}
	}
}
