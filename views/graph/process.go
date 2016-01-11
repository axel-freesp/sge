package graph

import (
	//"log"
	"image"
	"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/tool"
)

type ProcessPort struct {
	SelectableBox
	userObj interfaces.ChannelObject
}

func ProcessPortNew(pos image.Point, userObj interfaces.ChannelObject) *ProcessPort {
	size := image.Point{procPortWidth, procPortHeight}
	box := image.Rectangle{pos, pos.Add(size)}
	return &ProcessPort{SelectableBoxInit(box,
			ColorInit(ColorOption(NormalArchPort)),
			ColorInit(ColorOption(HighlightArchPort)),
			ColorInit(ColorOption(SelectArchPort)),
			ColorInit(ColorOption(BoxFrame)),
			image.Point{}),
			userObj}
}

func (p ProcessPort) Draw(ctxt interface{}) {
	p.DrawDefault(ctxt)
    switch ctxt.(type) {
    case *cairo.Context:
		empty := image.Point{}
		context := ctxt.(*cairo.Context)
		extPort := p.userObj.ArchPortObject()
		if extPort != nil {
			extPPos := extPort.Position()
			if extPPos != empty {
				var pos1, pos2 image.Point
				if p.userObj.Direction() == interfaces.InPort {
					pos1 = extPPos.Add(image.Point{archPortWidth, archPortHeight}.Div(2))
					pos2 = p.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
				} else {
					pos1 = p.Position().Add(image.Point{procPortWidth, procPortHeight}.Div(2))
					pos2 = extPPos.Add(image.Point{archPortWidth, archPortHeight}.Div(2))
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

func (p *ProcessPort) SetPosition(pos image.Point) {
	p.userObj.SetPosition(pos)
	p.DefaultSetPosition(pos)
}


type Process struct {
	NamedBoxObject
	userObj interfaces.ProcessObject
	ports []*ProcessPort
	selectedPort int
	arch ArchIf
}

var _ ProcessIf = (*Process)(nil)

func ProcessNew(pos image.Point, userObj interfaces.ProcessObject, arch ArchIf) (ret *Process) {
	shape := image.Point{global.processWidth, global.processHeight}
	box := image.Rectangle{pos, pos.Add(shape)}
	ret = &Process{NamedBoxObjectInit(box,
			ColorInit(ColorOption(ProcessNormal)),
			ColorInit(ColorOption(ProcessHighlight)),
			ColorInit(ColorOption(ProcessSelected)),
			ColorInit(ColorOption(BoxFrame)),
			ColorInit(ColorOption(Text)),
			image.Point{procPortOutBorder, procPortOutBorder}, userObj), userObj, nil, -1, arch}
	ret.RegisterOnHighlight(func(hit bool, pos image.Point) bool {
		return ret.onHighlight(hit, pos)
	})
	ret.RegisterOnSelect(func(){
		ret.onSelect()
	}, func(){
		ret.onDeselect()
	})
	idx := 0
	inCnt := len(userObj.InChannelObjects())
	outCnt := len(userObj.OutChannelObjects())
	empty := image.Point{}
	for _, c := range userObj.InChannelObjects() {
		var pos image.Point
		if c.Position() == empty {
			pos = ret.calcPortPos(idx, inCnt + outCnt)
		} else {
			pos = c.Position()
		}
		arch.(*Arch).channelMap[c] = ProcessPortNew(pos, c)
		ret.ports = append(ret.ports, arch.(*Arch).channelMap[c])
		idx++
	}
	for _, c := range userObj.OutChannelObjects() {
		var pos image.Point
		if c.Position() == empty {
			pos = ret.calcPortPos(idx, inCnt + outCnt)
		} else {
			pos = c.Position()
		}
		arch.(*Arch).channelMap[c] = ProcessPortNew(pos, c)
		ret.ports = append(ret.ports, arch.(*Arch).channelMap[c])
		idx++
	}
	return
}

func ProcessBox(pos image.Point) image.Rectangle {
	size := image.Point{global.processWidth, global.processHeight}
	return image.Rectangle{pos, pos.Add(size)}
}

func (p Process) UserObj() interfaces.ProcessObject {
	return p.userObj
}

func (pr Process) ArchObject() ArchIf {
	return pr.arch
}

func (pr *Process) SelectChannel(ch interfaces.ChannelObject) {
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

func (pr Process) GetSelectedChannel() (ok bool, ch interfaces.ChannelObject) {
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


func (pr Process) Draw(ctxt interface{}){
	pr.DrawDefaultText(ctxt)
	for _, port := range pr.ports {
		port.Draw(ctxt)
	}
}


//
//	freesp.Positioner interface
//

func (pr *Process) SetPosition(pos image.Point) {
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

var _ Selecter  = (*Process)(nil)

func (pr *Process) onSelect() (selected bool) {
	for i, p := range pr.ports {
		if i == pr.selectedPort {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	return
}

func (pr *Process) onDeselect() (selected bool) {
	for _, p := range pr.ports {
		p.Deselect()
	}
	return
}



//
//	Highlighter interface
//

var _ Highlighter  = (*Process)(nil)

func (pr *Process) onHighlight(hit bool, pos image.Point) (modified bool) {
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

const (
	procPortWidth = 8
	procPortHeight = 8
	procPortOutBorder = 6
)

func processDrawBody(context *cairo.Context, x, y, w, h float64, name string, mode ColorMode) {
	context.SetSourceRGB(processChooseColor(mode))
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

func processChooseColor(mode ColorMode) (r, g, b float64) {
	switch mode {
	case NormalMode:
		r, g, b = ColorOption(ProcessNormal)
	case HighlightMode:
		r, g, b = ColorOption(ProcessHighlight)
	case SelectedMode:
		r, g, b = ColorOption(ProcessSelected)
	}
	return
}

// obsolete
func (p Process) calcPortBox(idx, cnt int) image.Rectangle {
	lX, rX, _, bY := p.portCorners()
	size := image.Point{procPortWidth, procPortHeight}
	k := float64(idx + 1) / float64(cnt + 1)
	x := int(k * float64(rX - lX))
	pos := image.Point{lX + x, bY}
	return image.Rectangle{pos, pos.Add(size)}
}

func (p Process) calcPortPos(idx, cnt int) (pos image.Point) {
	lX, rX, _, bY := p.portCorners()
	k := float64(idx + 1) / float64(cnt + 1)
	x := int(k * float64(rX - lX))
	pos = image.Point{lX + x, bY}
	return
}

func (p Process) portIsSideBorder(pos image.Point) bool {
	lX, rX, _, _ := p.portCorners()
	return (pos.X == lX || pos.X == rX)
}

func (p Process) portIsCorner(pos image.Point) bool {
	lX, rX, tY, bY := p.portCorners()
	return (pos.X == lX || pos.X == rX) && (pos.Y == tY || pos.Y == bY)
}

func (p Process) portCorners() (lX, rX, tY, bY int) {
	lX = p.BBox().Min.X + 1
	rX = p.BBox().Max.X - procPortWidth - 1
	tY = p.BBox().Min.Y + 1
	bY = p.BBox().Max.Y - procPortHeight - 1
	return
}

func (p Process) portClipPos(pos image.Point) (ret image.Point) {
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


