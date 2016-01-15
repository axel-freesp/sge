package graph

import (
	//"log"
	"image"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/tool"
)

type ProcessMapping struct {
	Container
	userObj interfaces.ProcessObject
	ports []*ProcessPort
	nodes []*Node
	selectedPort int
	arch ArchIf
}

var _ ProcessIf = (*ProcessMapping)(nil)
var _ ContainerChild = (*Node)(nil)

func ProcessMappingNew(nodes []*Node, userObj interfaces.ProcessObject) (ret *ProcessMapping) {
	config := DrawConfig{ColorInit(ColorOption(ProcessNormal)),
			ColorInit(ColorOption(ProcessHighlight)),
			ColorInit(ColorOption(ProcessSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{procPortOutBorder, procPortOutBorder}}
	var children []ContainerChild
	for _, n := range nodes {
		children = append(children, n)
	}
	ret = &ProcessMapping{ContainerInit(children, config, userObj), userObj, nil, nodes, -1, nil}
	ret.RegisterOnHighlight(func(hit bool, pos image.Point) bool {
		return ret.onHighlight(hit, pos)
	})
	ret.RegisterOnSelect(func(){
		ret.onSelect()
	}, func(){
		ret.onDeselect()
	})
	return
}

func ProcessMappingBox(pos image.Point) image.Rectangle {
	size := image.Point{global.processWidth, global.processHeight}
	return image.Rectangle{pos, pos.Add(size)}
}

func ProcessMappingFit(outer, inner image.Rectangle) image.Rectangle {
	borderTop := image.Point{-18, -30}
	borderBottom := image.Point{18, 18}
	test := image.Rectangle{inner.Min.Add(borderTop), inner.Max.Add(borderBottom)}
	if outer.Size().X == 0 {
		return test
	}
	return outer.Union(test)
}

func (p ProcessMapping) UserObj() interfaces.ProcessObject {
	return p.userObj
}

func (pr ProcessMapping) ArchObject() ArchIf {
	return pr.arch
}

func (pr *ProcessMapping) SetArchObject(a ArchIf) {
	pr.arch = a
	idx := 0
	inCnt := len(pr.userObj.InChannelObjects())
	outCnt := len(pr.userObj.OutChannelObjects())
	empty := image.Point{}
	for _, c := range pr.userObj.InChannelObjects() {
		var pos image.Point
		if c.Position() == empty {
			pos = pr.calcPortPos(idx, inCnt + outCnt)
		} else {
			pos = c.Position()
		}
		a.(*Arch).channelMap[c] = ProcessPortNew(pos, c)
		pr.ports = append(pr.ports, a.(*Arch).channelMap[c])
		idx++
	}
	for _, c := range pr.userObj.OutChannelObjects() {
		var pos image.Point
		if c.Position() == empty {
			pos = pr.calcPortPos(idx, inCnt + outCnt)
		} else {
			pos = c.Position()
		}
		a.(*Arch).channelMap[c] = ProcessPortNew(pos, c)
		pr.ports = append(pr.ports, a.(*Arch).channelMap[c])
		idx++
	}
}

func (pr *ProcessMapping) SelectChannel(ch interfaces.ChannelObject) {
	pr.selectedPort = -1
	for i, p := range pr.ports {
		if ch == p.userObj {
			p.Select()
			pr.selectedPort = i
		} else {
			p.Deselect()
		}
	}
}

func (pr ProcessMapping) GetSelectedChannel() (ok bool, ch interfaces.ChannelObject) {
	if pr.selectedPort == -1 {
		return
	}
	ok = true
	ch = pr.ports[pr.selectedPort].userObj
	return
}


//
//	Drawer interface
//


func (pr ProcessMapping) Draw(ctxt interface{}){
	pr.DrawDefaultText(ctxt)
	for _, port := range pr.ports {
		port.Draw(ctxt)
	}
}


//
//	freesp.Positioner interface
//

func (pr *ProcessMapping) SetPosition(pos image.Point) {
	if pr.selectedPort != -1 {
		child := pr.ports[pr.selectedPort]
		childpos := child.Position()
		offset := childpos.Sub(pr.Position())
		newPos := pos.Add(offset)
		if pr.portIsCorner(childpos) {
			if tool.AbsInt(newPos.X - childpos.X) > tool.AbsInt(newPos.Y - childpos.Y) {
				newPos.Y = childpos.Y
			} else {
				newPos.X = childpos.X
			}
		} else if pr.portIsSideBorder(childpos) {
			newPos.X = childpos.X
		} else {
			newPos.Y = childpos.Y
		}
		child.SetPosition(pr.portClipPos(newPos))
	} else {
		shift := pos.Sub(pr.Position())
		pr.userObj.SetPosition(pos)
		pr.box = pr.box.Add(shift)
		for _, p := range pr.ports {
			p.SetPosition(p.Position().Add(shift))
		}
	}
}


//
//	Selecter interface
//

var _ Selecter  = (*ProcessMapping)(nil)

func (pr *ProcessMapping) onSelect() (selected bool) {
	for i, p := range pr.ports {
		if i == pr.selectedPort {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	return
}

func (pr *ProcessMapping) onDeselect() (selected bool) {
	for _, p := range pr.ports {
		p.Deselect()
	}
	return
}



//
//	Highlighter interface
//

var _ Highlighter  = (*ProcessMapping)(nil)

func (pr *ProcessMapping) onHighlight(hit bool, pos image.Point) (modified bool) {
	pr.selectedPort = -1
	if hit {
		for i, p := range pr.ports {
			phit, mod := p.CheckHit(pos)
			if phit {
				pr.selectedPort = i
			}
			modified = modified || mod
		}
	} else {
		for _, p := range pr.ports {
			p.CheckHit(pos)
		}
	}
	return
}

//
//	Private functions
//

func (p ProcessMapping) calcPortPos(idx, cnt int) (pos image.Point) {
	lX, rX, _, bY := p.portCorners()
	k := float64(idx + 1) / float64(cnt + 1)
	x := int(k * float64(rX - lX))
	pos = image.Point{lX + x, bY}
	return
}

func (p ProcessMapping) portIsSideBorder(pos image.Point) bool {
	lX, rX, _, _ := p.portCorners()
	return (pos.X == lX || pos.X == rX)
}

func (p ProcessMapping) portIsCorner(pos image.Point) bool {
	lX, rX, tY, bY := p.portCorners()
	return (pos.X == lX || pos.X == rX) && (pos.Y == tY || pos.Y == bY)
}

func (p ProcessMapping) portCorners() (lX, rX, tY, bY int) {
	lX = p.BBox().Min.X + 1
	rX = p.BBox().Max.X - procPortWidth - 1
	tY = p.BBox().Min.Y + 1
	bY = p.BBox().Max.Y - procPortHeight - 1
	return
}

func (p ProcessMapping) portClipPos(pos image.Point) (ret image.Point) {
	ret = pos
	lX, rX, tY, bY := p.portCorners()
	if ret.X < lX {
		ret.X = lX
	}
	if ret.X > rX {
		ret.X = rX
	}
	if ret.Y < tY {
		ret.Y = tY
	}
	if ret.Y > bY {
		ret.Y = bY
	}
	return
}


