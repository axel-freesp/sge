package graph

import (
	"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/tool"
)

type ArchPort struct {
	SelectableBox
	userObj interfaces.ArchPortObject
}

func ArchPortNew(pos image.Point, userObj interfaces.ArchPortObject) *ArchPort {
	size := image.Point{archPortWidth, archPortHeight}
	box := image.Rectangle{pos, pos.Add(size)}
	return &ArchPort{SelectableBoxInit(box,
			ColorInit(ColorOption(NormalArchPort)),
			ColorInit(ColorOption(HighlightArchPort)),
			ColorInit(ColorOption(SelectArchPort)),
			ColorInit(ColorOption(BoxFrame)),
			image.Point{}),
			userObj}
}

var _ ArchPortIf = (*ArchPort)(nil)

func (p *ArchPort) SetPosition(pos image.Point) {
	p.userObj.SetPosition(pos)
	p.DefaultSetPosition(pos)
}

//////////////////////

type Arch struct {
	NamedBoxObject
	userObj interfaces.ArchObject
	processes []ProcessIf
	selectedProcess int
	ports []*ArchPort
	selectedPort int
	channelMap map[interfaces.ChannelObject]*ProcessPort
}

var _ ArchIf = (*Arch)(nil)

func ArchNew(box image.Rectangle, userObj interfaces.ArchObject) *Arch {
	a := &Arch{NamedBoxObjectInit(box,
			ColorInit(ColorOption(ArchNormal)),
			ColorInit(ColorOption(ArchHighlight)),
			ColorInit(ColorOption(ArchSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{archPortOutBorder, archPortOutBorder}, userObj), userObj, nil, -1, nil, -1,
		make(map[interfaces.ChannelObject]*ProcessPort)}
	a.RegisterOnHighlight(func(hit bool, pos image.Point) bool{
		return a.onHighlight(hit, pos)
	})
	a.RegisterOnSelect(func(){
		a.onSelect()
	}, func(){
		a.onDeselect()
	})
	for _, up := range userObj.ProcessObjects() {
		p := ProcessNew(up.Position(), up, a)
		a.processes = append(a.processes, p)
	}
	extCnt := a.numExtChannel()
	idx := 0
	empty := image.Point{}
	for _, up := range userObj.ProcessObjects() {
		for _, c := range up.InChannelObjects() {
			if a.channelIsExtern(c) {
				var pos image.Point
				ap := c.ArchPortObject()
				if ap == nil {
					pos = a.calcPortPos(idx, extCnt)
				} else if ap.Position() == empty {
					pos = a.calcPortPos(idx, extCnt)
				} else {
					pos = ap.Position()
				}
				a.ports = append(a.ports, ArchPortNew(pos, ap))
				idx++
			}
		}
		for _, c := range up.OutChannelObjects() {
			if a.channelIsExtern(c) {
				var pos image.Point
				ap := c.ArchPortObject()
				if ap == nil {
					pos = a.calcPortPos(idx, extCnt)
				} else if ap.Position() == empty {
					pos = a.calcPortPos(idx, extCnt)
				} else {
					pos = ap.Position()
				}
				a.ports = append(a.ports, ArchPortNew(pos, ap))
				idx++
			}
		}
	}
	return a
}

func ArchFit(outer, inner image.Rectangle) image.Rectangle {
	borderTop := image.Point{-18, -30}
	borderBottom := image.Point{18, 18}
	test := image.Rectangle{inner.Min.Add(borderTop), inner.Max.Add(borderBottom)}
	if outer.Size().X == 0 {
		return test
	}
	return outer.Union(test)
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
	for _, pr := range a.processes {
		for _, c := range pr.UserObj().InChannelObjects() {
			if c.LinkObject().ProcessObject().ArchObject().Name() == name {
				return true
			}
		}
		for _, c := range pr.UserObj().OutChannelObjects() {
			if c.LinkObject().ProcessObject().ArchObject().Name() == name {
				return true
			}
		}
	}
	return false
}

func (a *Arch) ChannelPort(ch interfaces.ChannelObject) ArchPortIf {
	for _, p := range a.ports {
		if p.userObj.Channel() == ch {
			return p
		}
	}
	return nil
}


//
//	Drawer interface
//

func (a Arch) Draw(ctxt interface{}){
	a.DrawDefaultText(ctxt)
	for _, p := range a.ports {
		p.Draw(ctxt)
	}
	for _, pr := range a.processes {
		pr.Draw(ctxt)
		for _, c := range pr.UserObj().InChannelObjects() {
			a.drawLocalChannel(ctxt, c)
		}
		for _, c := range pr.UserObj().OutChannelObjects() {
			a.drawLocalChannel(ctxt, c)
		}
	}
}


//
//	freesp.Positioner interface
//

func (a *Arch) SetPosition(pos image.Point) {
	if a.selectedProcess != -1 {
		child := a.processes[a.selectedProcess]
		childpos := child.Position()
		offset := childpos.Sub(a.Position())
		child.SetPosition(pos.Add(offset))
		var box image.Rectangle
		for _, p := range a.processes {
			box = ArchFit(box, p.BBox())
		}
		a.box = box
		empty := image.Point{}
		for i, p := range a.ports {
			var pos image.Point
			ap := p.userObj
			if ap == nil {
				pos = a.calcPortPos(i, len(a.ports))
			} else if ap.Position() == empty {
				pos = a.calcPortPos(i, len(a.ports))
			} else {
				pos = ap.Position()
			}
			p.SetPosition(a.portClipPos(pos))
		}
	} else if a.selectedPort != -1 {
		child := a.ports[a.selectedPort]
		childpos := child.Position()
		offset := childpos.Sub(a.Position())
		newPos := pos.Add(offset)
		if a.portIsCorner(childpos) {
			if tool.AbsInt(newPos.X - childpos.X) > tool.AbsInt(newPos.Y - childpos.Y) {
				newPos.Y = childpos.Y
			} else {
				newPos.X = childpos.X
			}
		} else if a.portIsSideBorder(childpos) {
			newPos.X = childpos.X
		} else {
			newPos.Y = childpos.Y
		}
		child.SetPosition(a.portClipPos(newPos))
	} else {
		shift := pos.Sub(a.Position())
		a.userObj.SetPosition(pos)
		a.box = a.box.Add(shift)
		for _, p := range a.processes {
			p.SetPosition(p.Position().Add(shift))
		}
		for _, p := range a.ports {
			p.SetPosition(p.Position().Add(shift))
		}
	}
}


//
//  freesp.Shaper API
//

func (a *Arch) SetShape(shape image.Point) {
	a.box.Max = a.box.Min.Add(shape)
	a.userObj.SetShape(shape)
}


//
//	Selecter interface
//

var _ Selecter  = (*Arch)(nil)

func (a *Arch) onSelect() (selected bool) {
	for _, pr := range a.processes {
		hit, _ := pr.CheckHit(a.Pos)
		if hit {
			selected = pr.Select() && selected
		} else {
			selected = !pr.Deselect() && selected
		}
	}
	for _, p := range a.ports {
		hit, _ := p.CheckHit(a.Pos)
		if hit {
			selected = p.Select() && selected
		} else {
			selected = !p.Deselect() && selected
		}
	}
	return
}

func (a *Arch) onDeselect() (selected bool) {
	for _, pr := range a.processes {
		selected = pr.Deselect() || selected
	}
	for _, p := range a.ports {
		selected = p.Deselect() || selected
	}
	return
}

func (a *Arch) SelectProcess(pr interfaces.ProcessObject) {
	a.selectedProcess = -1
	for i, p := range a.processes {
		if pr == p.UserObj() {
			p.Select()
			a.selectedProcess = i
		} else {
			p.Deselect()
		}
	}
}

func (a Arch) GetSelectedProcess() (ok bool, pr interfaces.ProcessObject) {
	if a.selectedProcess == -1 {
		return
	}
	ok = true
	pr = a.processes[a.selectedProcess].UserObj()
	return
}

func (a *Arch) SelectChannel(ch interfaces.ChannelObject) {
	for _, p := range a.processes {
		p.SelectChannel(ch)
	}
}

func (a Arch) GetSelectedChannel() (ok bool, ch interfaces.ChannelObject) {
	for _, p := range a.processes {
		ok, ch = p.GetSelectedChannel()
		if ok {
			return
		}
	}
	return
}


//
//	Highlighter interface
//

var _ Highlighter  = (*Arch)(nil)

func (a *Arch) onHighlight(hit bool, pos image.Point) (modified bool) {
	a.selectedProcess = -1
	a.selectedPort = -1
	if hit {
		for i, pr := range a.processes {
			phit, mod := pr.CheckHit(pos)
			if phit {
				a.selectedProcess = i
			}
			modified = modified || mod
		}
		for i, p := range a.ports {
			phit, mod := p.CheckHit(pos)
			if phit {
				a.selectedPort = i
			}
			modified = modified || mod
		}
	} else {
		for _, pr := range a.processes {
			pr.CheckHit(pos)
		}
		for _, p := range a.ports {
			p.CheckHit(pos)
		}
	}
	return
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
				pos1 = link.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
				pos2 = ch.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
			} else {
				pos1 = ch.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
				pos2 = link.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
			}
			DrawArrow(context, pos1, pos2)
		}
	}
}

func (a Arch) calcPortPos(idx, cnt int) (pos image.Point) {
	lX, rX, _, bY := a.portCorners()
	k := float64(idx + 1) / float64(cnt + 1)
	x := int(k * float64(rX - lX))
	pos = image.Point{lX + x, bY}
	return
}

func (a Arch) portIsSideBorder(pos image.Point) bool {
	lX, rX, _, _ := a.portCorners()
	return (pos.X == lX || pos.X == rX)
}

func (a Arch) portIsCorner(pos image.Point) bool {
	lX, rX, tY, bY := a.portCorners()
	return (pos.X == lX || pos.X == rX) && (pos.Y == tY || pos.Y == bY)
}

func (a Arch) portCorners() (lX, rX, tY, bY int) {
	lX = a.BBox().Min.X + 1
	rX = a.BBox().Max.X - archPortWidth - 1
	tY = a.BBox().Min.Y + 1
	bY = a.BBox().Max.Y - archPortHeight - 1
	return
}

func (a Arch) portClipPos(pos image.Point) (ret image.Point) {
	ret = pos
	lX, rX, tY, bY := a.portCorners()
	if ret.X < lX {
		ret.X = lX
	} else if ret.X > rX {
		ret.X = rX
	}
	if ret.Y < tY {
		ret.Y = tY
	} else if ret.Y > bY {
		ret.Y = bY
	}
	dX1 := tool.AbsInt(ret.X - lX)
	dX2 := tool.AbsInt(ret.X - rX)
	dY1 := tool.AbsInt(ret.Y - tY)
	dY2 := tool.AbsInt(ret.Y - bY)
	minDX := tool.MinInt(dX1, dX2)
	minDY := tool.MinInt(dY1, dY2)
	if tool.MinInt(minDX, minDY) > 0 { //pos inside, not on border
		if minDX < minDY {
			if dX1 < dX2 {
				ret.X = lX
			} else {
				ret.X = rX
			}
		} else {
			if dY1 < dY2 {
				ret.Y = tY
			} else {
				ret.Y = bY
			}
		}
	}
	return
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


