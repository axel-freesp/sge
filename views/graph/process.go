package graph

import (
	"image"
	"github.com/gotk3/gotk3/cairo"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
)

type Process struct {
	Container
	userObj pf.ProcessIf
	arch ArchIf
}

var _ ProcessIf = (*Process)(nil)

func ProcessNew(pos image.Point, userObj pf.ProcessIf) (ret *Process) {
	config := DrawConfig{ColorInit(ColorOption(ProcessNormal)),
			ColorInit(ColorOption(ProcessHighlight)),
			ColorInit(ColorOption(ProcessSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{procPortWidth, procPortHeight}}
	cconfig := ContainerConfig{procPortWidth, procPortHeight, procMinWidth, procMinHeight}
	ret = &Process{ContainerInit(nil, config, userObj, cconfig), userObj, nil}
	shape := image.Point{global.processWidth, global.processHeight}
	ret.box = image.Rectangle{pos, pos.Add(shape)}
	ret.ContainerInit()
	return
}

func (p *Process) Layout() image.Rectangle {
	return p.box
}

func (p Process) UserObj() pf.ProcessIf {
	return p.userObj
}

func (pr Process) ArchObject() ArchIf {
	return pr.arch
}

func processDrawChannel(p *ContainerPort, ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
		empty := image.Point{}
		context := ctxt.(*cairo.Context)
		shift1 := image.Point{archPortWidth, archPortHeight}.Div(2)
		shift2 := image.Point{procPortWidth, procPortHeight}.Div(2)
		extPort := p.UserObj2.(pf.ChannelIf).ArchPort()
		if extPort != nil {
			extPPos := extPort.ModePosition(gr.PositionModeNormal)
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

func (pr *Process) SetArchObject(a ArchIf) {
	pr.arch = a
	idx := 0
	inCnt := len(pr.userObj.InChannels())
	outCnt := len(pr.userObj.OutChannels())
	empty := image.Point{}
	config := DrawConfig{ColorInit(ColorOption(NormalArchPort)),
			ColorInit(ColorOption(HighlightArchPort)),
			ColorInit(ColorOption(SelectArchPort)),
			ColorInit(ColorOption(BoxFrame)),
			Color{},
			image.Point{}}
	for _, c := range pr.userObj.InChannels() {
		pos := c.ModePosition(gr.PositionModeNormal)
		if pos == empty {
			pos = pr.CalcPortPos(idx, inCnt + outCnt)
		}
		p := pr.AddModePort(pos, config, c, gr.PositionModeNormal)
		p.RegisterOnDraw(func(ctxt interface{}){
			processDrawChannel(p, ctxt)
		})
		a.(*Arch).channelMap[c] = p
		idx++
	}
	for _, c := range pr.userObj.OutChannels() {
		pos := c.ModePosition(gr.PositionModeNormal)
		if pos == empty {
			pos = pr.CalcPortPos(idx, inCnt + outCnt)
		}
		p := pr.AddModePort(pos, config, c, gr.PositionModeNormal)
		p.RegisterOnDraw(func(ctxt interface{}){
			processDrawChannel(p, ctxt)
		})
		a.(*Arch).channelMap[c] = p
		idx++
	}
}

func (pr *Process) SelectChannel(ch pf.ChannelIf) {
	pr.SelectModePort(ch)
}

func (pr Process) GetSelectedChannel() (ok bool, ch pf.ChannelIf) {
	var p gr.ModePositioner
	ok, p = pr.GetSelectedModePort()
	if !ok {
		return
	}
	ch = p.(pf.ChannelIf)
	return
}

func (pr *Process) SetPosition(pos image.Point) {
	pr.ContainerDefaultSetPosition(pos)
	pr.userObj.SetModePosition(gr.PositionModeNormal, pos)
}

const (
	procPortWidth = 8
	procPortHeight = 8
	procPortOutBorder = 6
	procMinWidth = 120
	procMinHeight = 80
)

