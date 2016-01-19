package graph

import (
	"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
)

type Arch struct {
	Container
	userObj interfaces.ArchObject
	channelMap map[interfaces.ChannelObject]*ContainerPort
	processes []ProcessIf
	mode interfaces.PositionMode
	mapping interfaces.MappingObject
}

var _ ArchIf = (*Arch)(nil)

func ArchNew(userObj interfaces.ArchObject) *Arch {
	var processes []ProcessIf
	var Children []ContainerChild
	for _, up := range userObj.ProcessObjects() {
		p := ProcessNew(up.ModePosition(interfaces.PositionModeNormal), up)
		processes = append(processes, p)
		Children = append(Children, p)
	}
	config := DrawConfig{ColorInit(ColorOption(ArchNormal)),
			ColorInit(ColorOption(ArchHighlight)),
			ColorInit(ColorOption(ArchSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{archPortOutBorder, archPortOutBorder}}
	cconfig := ContainerConfig{archPortWidth, archPortHeight}
	a := &Arch{ContainerInit(Children, config, userObj, cconfig), userObj,
		make(map[interfaces.ChannelObject]*ContainerPort), processes, interfaces.PositionModeNormal, nil}
	a.init()
	a.initPorts()
	return a
}

func ArchMappingNew(userObj interfaces.ArchObject, nodes []NodeIf, mapping interfaces.MappingObject) *Arch {
	var processes []ProcessIf
	var Children []ContainerChild
	for _, up := range userObj.ProcessObjects() {
		var mappedNodes []*Node
		for _, n := range nodes {
			m, ok := mapping.MappedObject(n.UserObj())
			if ok && m == up {
				mappedNodes = append(mappedNodes, n.(*Node))
			}
		}
		p := ProcessMappingNew(mappedNodes, up)
		processes = append(processes, p)
		Children = append(Children, p)
	}
	config := DrawConfig{ColorInit(ColorOption(ArchNormal)),
			ColorInit(ColorOption(ArchHighlight)),
			ColorInit(ColorOption(ArchSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{archPortOutBorder, archPortOutBorder}}
	cconfig := ContainerConfig{archPortWidth, archPortHeight}
	a := &Arch{ContainerInit(Children, config, userObj, cconfig), userObj,
		make(map[interfaces.ChannelObject]*ContainerPort), processes, interfaces.PositionModeMapping, mapping}
	a.init()
	for _, n := range nodes {
		melem, ok := mapping.MapElemObject(n.UserObj())
		if !ok {
			log.Printf("ArchMappingNew warning: node %s not mapped.\n", n.Name())
			continue
		}
		n.SetPosition(melem.Position())
	}
	a.initMappingPorts()
	return a
}

func (a *Arch) init() {
	a.ContainerInit()
	a.RegisterOnDraw(func(ctxt interface{}){
		archOnDraw(a, ctxt)
	})
	for _, pr := range a.Children {
		pr.(ProcessIf).SetArchObject(a)
	}
}

func (a *Arch) initMappingPorts() {
	idx := 0
	for _, up := range a.userObj.ProcessObjects() {
		for _, c := range up.InChannelObjects() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, interfaces.PositionModeMapping, idx)
				idx++
			}
		}
		for _, c := range up.OutChannelObjects() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, interfaces.PositionModeMapping, idx)
				idx++
			}
		}
	}
	a.Layout()
}

func (a *Arch) initPorts() {
	idx := 0
	for _, up := range a.userObj.ProcessObjects() {
		for _, c := range up.InChannelObjects() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, interfaces.PositionModeNormal, idx)
				idx++
			}
		}
		for _, c := range up.OutChannelObjects() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, interfaces.PositionModeNormal, idx)
				idx++
			}
		}
	}
	a.Layout()
}

func (a *Arch) addExternalPort(c interfaces.ChannelObject, mode interfaces.PositionMode, idx int) {
	config := DrawConfig{ColorInit(ColorOption(NormalArchPort)),
			ColorInit(ColorOption(HighlightArchPort)),
			ColorInit(ColorOption(SelectArchPort)),
			ColorInit(ColorOption(BoxFrame)),
			Color{},
			image.Point{}}
	empty := image.Point{}
	ap := c.ArchPortObject()
	pos := ap.ModePosition(mode)
	if pos == empty {
		pos = a.CalcPortPos(idx, a.numExtChannel())
	}
	cp := a.AddModePort(pos, config, ap, mode)
	ap.SetModePosition(mode, cp.Position())
}

func (a Arch) Processes() []ProcessIf {
	return a.processes
}

func (a Arch) UserObj() interfaces.ArchObject {
	return a.userObj
}

func (a Arch) IsLinked(name string) bool {
	if name == a.Name() {
		return true
	}
	for _, pr := range a.Children {
		for _, c := range pr.(ProcessIf).UserObj().InChannelObjects() {
			if c.LinkObject().ProcessObject().ArchObject().Name() == name {
				return true
			}
		}
		for _, c := range pr.(ProcessIf).UserObj().OutChannelObjects() {
			if c.LinkObject().ProcessObject().ArchObject().Name() == name {
				return true
			}
		}
	}
	return false
}

func (a *Arch) ChannelPort(ch interfaces.ChannelObject) ArchPortIf {
	for _, p := range a.ports {
		if p.UserObj2.(interfaces.ArchPortObject).Channel() == ch {
			return p
		}
	}
	return nil
}



func (a *Arch) SelectProcess(pr interfaces.ProcessObject) (p ProcessIf) {
	for _, ch := range a.Children {
		if pr == ch.(ProcessIf).UserObj() {
			a.SelectChild(ch)
			p = ch.(ProcessIf)
			return
		}
	}
	return
}

func (a Arch) GetSelectedProcess() (ok bool, pr interfaces.ProcessObject, p ProcessIf) {
	var ch ContainerChild
	ok, ch = a.GetSelectedChild()
	if !ok {
		return
	}
	p = ch.(ProcessIf)
	pr = p.UserObj()
	return
}

func (a *Arch) SelectChannel(ch interfaces.ChannelObject) {
	for _, p := range a.Children {
		p.(ProcessIf).SelectChannel(ch)
	}
}

func (a Arch) GetSelectedChannel() (ok bool, ch interfaces.ChannelObject) {
	for _, p := range a.Children {
		ok, ch = p.(ProcessIf).GetSelectedChannel()
		if ok {
			return
		}
	}
	return
}

func (a *Arch) SetPosition(pos image.Point) {
	a.ContainerDefaultSetPosition(pos)
	a.userObj.SetModePosition(a.mode, pos)
}


//
//	Drawer interface
//

func archOnDraw(a *Arch, ctxt interface{}) {
	for _, pr := range a.Children {
		for _, c := range pr.(ProcessIf).UserObj().InChannelObjects() {
			a.drawLocalChannel(ctxt, c)
		}
		for _, c := range pr.(ProcessIf).UserObj().OutChannelObjects() {
			a.drawLocalChannel(ctxt, c)
		}
	}
}

//
//	Private functions
//

const (
	archPortWidth = 10
	archPortHeight = 10
	archPortOutBorder = 8
)

func (a Arch) drawLocalChannel(ctxt interface{}, ch interfaces.ChannelObject) {
    switch ctxt.(type) {
    case *cairo.Context:
		context := ctxt.(*cairo.Context)
		link := ch.LinkObject()
		if ch.ProcessObject().ArchObject().Name() != a.Name() {
			log.Fatal("Arch.drawLocalChannel: channel not in arch %s\n", a.Name())
		}
		if link.ProcessObject().ArchObject().Name() == a.Name() {
			var r, g, b float64
			p1 := a.channelMap[ch]
			p2 := a.channelMap[link]
			if p1.IsSelected() || p2.IsSelected() {
				r, g, b = ColorOption(SelectChannelLine)
			} else if p1.IsHighlighted() || p2.IsHighlighted() {
				r, g, b = ColorOption(HighlightChannelLine)
			} else {
				r, g, b = ColorOption(NormalChannelLine)
			}
			context.SetLineWidth(2)
			context.SetSourceRGB(r, g, b)
			var pos1, pos2 image.Point
			if ch.Direction() == interfaces.InPort {
				pos1 = link.ModePosition(a.mode).Add(image.Point{procPortWidth, procPortHeight}.Div(2))
				pos2 = ch.ModePosition(a.mode).Add(image.Point{procPortWidth, procPortHeight}.Div(2))
			} else {
				pos1 = ch.ModePosition(a.mode).Add(image.Point{procPortWidth, procPortHeight}.Div(2))
				pos2 = link.ModePosition(a.mode).Add(image.Point{procPortWidth, procPortHeight}.Div(2))
			}
			DrawArrow(context, pos1, pos2)
		}
	}
}

func (a Arch) channelIsExtern(c interfaces.ChannelObject) bool {
	link := c.LinkObject()
	cp := link.ProcessObject()
	ca := cp.ArchObject()
	return ca != a.userObj
}

func (a Arch) numExtChannel() (extCnt int) {
	for _, up := range a.userObj.ProcessObjects() {
		for _, c := range up.InChannelObjects() {
			if a.channelIsExtern(c) {
				extCnt++
			}
		}
		for _, c := range up.OutChannelObjects() {
			if a.channelIsExtern(c) {
				extCnt++
			}
		}
	}
	return
}


