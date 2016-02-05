package graph

import (
	//	"github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/gotk3/gotk3/cairo"
	"image"
	"log"
	"strings"
)

type Arch struct {
	Container
	userObj    pf.ArchIf
	channelMap map[pf.ChannelIf]*ContainerPort
	processes  []ProcessIf
	mode       gr.PositionMode
	mapping    mp.MappingIf
}

var _ ArchIf = (*Arch)(nil)

func ArchNew(userObj pf.ArchIf) *Arch {
	var processes []ProcessIf
	var Children []ContainerChild
	for _, up := range userObj.Processes() {
		p := ProcessNew(up.ModePosition(gr.PositionModeNormal), up)
		processes = append(processes, p)
		Children = append(Children, p)
	}
	config := DrawConfig{ColorInit(ColorOption(ArchNormal)),
		ColorInit(ColorOption(ArchHighlight)),
		ColorInit(ColorOption(ArchSelected)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{archPortOutBorder, archPortOutBorder}}
	cconfig := ContainerConfig{archPortWidth, archPortHeight, archMinWidth, archMinHeight}
	a := &Arch{ContainerInit(Children, config, userObj, cconfig), userObj,
		make(map[pf.ChannelIf]*ContainerPort), processes, gr.PositionModeNormal, nil}
	a.init()
	a.initPorts()
	return a
}

func findNodeByName(nodes []NodeIf, nId string) (n NodeIf, ok bool) {
	for _, n = range nodes {
		if n.Name() == nId {
			ok = true
			return
		}
	}
	return
}

func findNodeInTree(nodes []NodeIf, nId string) (n NodeIf, ok bool) {
	if strings.Contains(nId, "/") {
		id := strings.Split(nId, "/")
		var n2 NodeIf
		n2, ok = findNodeByName(nodes, id[0])
		if !ok {
			return
		}
		n, ok = findNodeInTree(n2.ChildNodes(), strings.Join(id[1:], "/"))
	} else {
		n, ok = findNodeByName(nodes, nId)
	}
	return
}

func ArchMappingNew(userObj pf.ArchIf, nodes []NodeIf, mapping mp.MappingIf) *Arch {
	var processes []ProcessIf
	var Children []ContainerChild
	for _, up := range userObj.Processes() {
		var mappedNodes []NodeIf
		var mappedIds []bh.NodeIdIf
		for _, nId := range mapping.MappedIds() {
			m, ok := mapping.Mapped(nId.String())
			if ok && m == up {
				log.Printf("ArchMappingNew(p=%s): nId=%s", up.Name(), nId.String())
				n, ok := findNodeInTree(nodes, nId.String())
				if ok {
					mappedNodes = append(mappedNodes, n)
					mappedIds = append(mappedIds, nId)
				}
			}
		}
		p := ProcessMappingNew(mappedNodes, mappedIds, up)
		processes = append(processes, p)
		Children = append(Children, p)
	}
	config := DrawConfig{ColorInit(ColorOption(ArchNormal)),
		ColorInit(ColorOption(ArchHighlight)),
		ColorInit(ColorOption(ArchSelected)),
		ColorInit(ColorOption(BoxFrame)),
		ColorInit(ColorOption(Text)),
		image.Point{archPortOutBorder, archPortOutBorder}}
	cconfig := ContainerConfig{archPortWidth, archPortHeight, archMinWidth, archMinHeight}
	a := &Arch{ContainerInit(Children, config, userObj, cconfig), userObj,
		make(map[pf.ChannelIf]*ContainerPort), processes, gr.PositionModeMapping, mapping}
	a.init()
	/*
		for _, n := range nodes {
			melem, ok := mapping.MappedElement(n.UserObj())
			if !ok {
				log.Printf("ArchMappingNew warning: node %s not mapped.\n", n.Name())
				continue
			}
			n.SetPosition(melem.PathModePosition("", gr.PositionModeMapping))
		}
	*/
	a.initMappingPorts()
	return a
}

func (a *Arch) init() {
	a.ContainerInit()
	a.RegisterOnDraw(func(ctxt interface{}) {
		archOnDraw(a, ctxt)
	})
	for _, pr := range a.Children {
		pr.(ProcessIf).SetArchObject(a)
	}
}

func (a *Arch) initMappingPorts() {
	idx := 0
	for _, up := range a.userObj.Processes() {
		for _, c := range up.InChannels() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, gr.PositionModeMapping, idx)
				idx++
			}
		}
		for _, c := range up.OutChannels() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, gr.PositionModeMapping, idx)
				idx++
			}
		}
	}
	a.Layout()
}

func (a *Arch) initPorts() {
	idx := 0
	for _, up := range a.userObj.Processes() {
		for _, c := range up.InChannels() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, gr.PositionModeNormal, idx)
				idx++
			}
		}
		for _, c := range up.OutChannels() {
			if a.channelIsExtern(c) {
				a.addExternalPort(c, gr.PositionModeNormal, idx)
				idx++
			}
		}
	}
	a.Layout()
}

func (a *Arch) addExternalPort(c pf.ChannelIf, mode gr.PositionMode, idx int) {
	config := DrawConfig{ColorInit(ColorOption(NormalArchPort)),
		ColorInit(ColorOption(HighlightArchPort)),
		ColorInit(ColorOption(SelectArchPort)),
		ColorInit(ColorOption(BoxFrame)),
		Color{},
		image.Point{}}
	ap := c.ArchPort()
	if ap == nil {
		log.Printf("Arch.addExternalPort error: channel %v has no arch port\n", c)
		return
	}
	positioner := gr.ModePositionerProxyNew(ap, mode)
	a.AddPort(config, ap, positioner)
}

func (a Arch) Processes() []ProcessIf {
	return a.processes
}

func (a Arch) UserObj() pf.ArchIf {
	return a.userObj
}

func (a Arch) IsLinked(name string) bool {
	if name == a.Name() {
		return true
	}
	for _, pr := range a.Children {
		for _, c := range pr.(ProcessIf).UserObj().InChannels() {
			if c.Link().Process().Arch().Name() == name {
				return true
			}
		}
		for _, c := range pr.(ProcessIf).UserObj().OutChannels() {
			if c.Link().Process().Arch().Name() == name {
				return true
			}
		}
	}
	return false
}

func (a *Arch) ChannelPort(ch pf.ChannelIf) ArchPortIf {
	for _, p := range a.ports {
		if p.UserObj.(pf.ArchPortIf).Channel() == ch {
			return p
		}
	}
	return nil
}

func (a *Arch) SelectProcess(pr pf.ProcessIf) (p ProcessIf) {
	for _, ch := range a.Children {
		if pr == ch.(ProcessIf).UserObj() {
			a.SelectChild(ch)
			p = ch.(ProcessIf)
			return
		}
	}
	return
}

func (a Arch) GetSelectedProcess() (ok bool, pr pf.ProcessIf, p ProcessIf) {
	var ch ContainerChild
	ok, ch = a.GetSelectedChild()
	if !ok {
		return
	}
	p = ch.(ProcessIf)
	pr = p.UserObj()
	return
}

func (a *Arch) SelectChannel(ch pf.ChannelIf) {
	for _, p := range a.Children {
		p.(ProcessIf).SelectChannel(ch)
	}
}

func (a Arch) GetSelectedChannel() (ok bool, ch pf.ChannelIf) {
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
		for _, c := range pr.(ProcessIf).UserObj().InChannels() {
			a.drawLocalChannel(ctxt, c)
		}
		for _, c := range pr.(ProcessIf).UserObj().OutChannels() {
			a.drawLocalChannel(ctxt, c)
		}
	}
}

//
//	Private functions
//

const (
	archPortWidth     = 10
	archPortHeight    = 10
	archPortOutBorder = 8
	archMinWidth      = 50
	archMinHeight     = 30
)

func (a Arch) drawLocalChannel(ctxt interface{}, ch pf.ChannelIf) {
	switch ctxt.(type) {
	case *cairo.Context:
		context := ctxt.(*cairo.Context)
		link := ch.Link()
		if ch.Process().Arch().Name() != a.Name() {
			log.Fatal("Arch.drawLocalChannel: channel not in arch %s\n", a.Name())
		}
		if link.Process().Arch().Name() == a.Name() {
			var r, g, b float64
			p1 := a.channelMap[ch]
			p2 := a.channelMap[link]
			if p1.IsSelected() || p2.IsSelected() {
				r, g, b, _ = ColorOption(SelectChannelLine)
			} else if p1.IsHighlighted() || p2.IsHighlighted() {
				r, g, b, _ = ColorOption(HighlightChannelLine)
			} else {
				r, g, b, _ = ColorOption(NormalChannelLine)
			}
			context.SetLineWidth(2)
			context.SetSourceRGB(r, g, b)
			var pos1, pos2 image.Point
			if ch.Direction() == gr.InPort {
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

func (a Arch) channelIsExtern(c pf.ChannelIf) bool {
	link := c.Link()
	cp := link.Process()
	ca := cp.Arch()
	return ca != a.userObj
}

func (a Arch) numExtChannel() (extCnt int) {
	for _, up := range a.userObj.Processes() {
		for _, c := range up.InChannels() {
			if a.channelIsExtern(c) {
				extCnt++
			}
		}
		for _, c := range up.OutChannels() {
			if a.channelIsExtern(c) {
				extCnt++
			}
		}
	}
	return
}
