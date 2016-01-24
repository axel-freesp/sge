package graph

import (
	"github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/tool"
	"image"
)

//
//		Children must implement:
//
type ContainerChild interface {
	Layout() image.Rectangle
	BoxedSelecter
}

//
//		Container Base
//
type Container struct {
	NamedBoxObject
	Children                    []ContainerChild
	ports                       []*ContainerPort
	selectedChild, selectedPort int
	config                      ContainerConfig
}

type ContainerConfig struct {
	portWidth, portHeight, minWidth, minHeight int
}

//
//		Port Base
//
type ContainerPort struct {
	SelectableBox
	UserObj  graph.Positioner
	UserObj2 graph.ModePositioner
	mode     graph.PositionMode
}

//
//		Construct a container with constructed Children (bottom up)
//
func ContainerInit(Children []ContainerChild, config DrawConfig, namer graph.Namer, cconfig ContainerConfig) (c Container) {
	c = Container{NamedBoxObjectInit(image.Rectangle{}, config, namer), Children, nil, -1, -1, cconfig}
	c.box = c.Layout()
	return
}

func ContainerNew(Children []ContainerChild, config DrawConfig, namer graph.Namer, cconfig ContainerConfig) (c *Container) {
	c = &Container{NamedBoxObjectInit(image.Rectangle{}, config, namer), Children, nil, -1, -1, cconfig}
	c.box = c.Layout()
	return
}

func (c *Container) ContainerInit() {
	c.RegisterOnHighlight(func(hit bool, pos image.Point) bool {
		return c.onHighlight(hit, pos)
	})
	c.RegisterOnSelect(func() {
		c.onSelect()
	}, func() {
		c.onDeselect()
	})
}

//
//		Construct port
//
func ContainerPortNew(box image.Rectangle, config DrawConfig, userObj graph.Positioner, userObj2 graph.ModePositioner, mode graph.PositionMode) (p *ContainerPort) {
	return &ContainerPort{SelectableBoxInit(box, config), userObj, userObj2, mode}
}

func (p *ContainerPort) SetPosition(pos image.Point) {
	if p.UserObj != nil {
		p.UserObj.SetPosition(pos)
	}
	if p.UserObj2 != nil {
		p.UserObj2.SetModePosition(p.mode, pos)
	}
	p.BBoxDefaultSetPosition(pos)
}

//
//		AddPort
//
func (c *Container) AddPort(pos image.Point, config DrawConfig, userObj graph.Positioner) (p *ContainerPort) {
	size := image.Point{c.config.portWidth, c.config.portHeight}
	pos = c.portClipPos(pos)
	box := image.Rectangle{pos, pos.Add(size)}
	p = ContainerPortNew(box, config, userObj, nil, graph.PositionMode(-1))
	c.ports = append(c.ports, p)
	return
}

func (c *Container) AddModePort(pos image.Point, config DrawConfig, userObj graph.ModePositioner, mode graph.PositionMode) (p *ContainerPort) {
	size := image.Point{c.config.portWidth, c.config.portHeight}
	box := image.Rectangle{pos, pos.Add(size)}
	p = ContainerPortNew(box, config, nil, userObj, mode)
	c.ports = append(c.ports, p)
	return
}

//
//		Container may be child of container
//
func (c *Container) Layout() (box image.Rectangle) {
	return c.ContainerDefaultLayout()
}

func (c *Container) ContainerDefaultLayout() (box image.Rectangle) {
	if len(c.Children) > 0 {
		for _, ch := range c.Children {
			box = ContainerFit(box, ch.Layout())
		}
	} else {
		box = c.box
	}
	if box.Max.X-box.Min.X < c.config.minWidth {
		box.Max.X = box.Min.X + c.config.minWidth
	}
	if box.Max.Y-box.Min.Y < c.config.minHeight {
		box.Max.Y = box.Min.Y + c.config.minHeight
	}
	c.box = box
	empty := image.Point{}
	for i, p := range c.ports {
		var pos image.Point
		if p.UserObj != nil {
			ap := p.UserObj
			pos = ap.Position()
			if pos == empty {
				pos = c.CalcPortPos(i, len(c.ports))
			}
		} else if p.UserObj2 != nil {
			ap := p.UserObj2
			pos = ap.ModePosition(p.mode)
			if pos == empty {
				pos = c.CalcPortPos(i, len(c.ports))
			}
		} else {
			pos = c.CalcPortPos(i, len(c.ports))
		}
		pos = c.portClipPos(pos)
		p.SetPosition(pos)
	}
	return
}

//
//		Set position
//
//	If child is selected, move child (and re-calculate layout)
//	If port is selected, move port (and clip to own shape)
//
func (c *Container) SetPosition(pos image.Point) {
	c.ContainerDefaultSetPosition(pos)
}

func (c *Container) ContainerDefaultSetPosition(pos image.Point) {
	if c.selectedChild != -1 {
		child := c.Children[c.selectedChild]
		childpos := child.Position()
		offset := childpos.Sub(c.Position())
		child.SetPosition(pos.Add(offset))
		c.box = c.Layout()
	} else if c.selectedPort != -1 {
		child := c.ports[c.selectedPort]
		childpos := child.Position()
		offset := childpos.Sub(c.Position())
		newPos := pos.Add(offset)
		child.SetPosition(c.portClipPos(newPos))
	} else {
		shift := pos.Sub(c.Position())
		c.BBoxDefaultSetPosition(pos)
		for _, p := range c.Children {
			p.SetPosition(p.Position().Add(shift))
		}
		for _, p := range c.ports {
			p.SetPosition(p.Position().Add(shift))
		}
	}
}

func (c Container) CalcPortPos(idx, cnt int) (pos image.Point) {
	lX, rX, _, bY := c.portCorners()
	k := float64(idx+1) / float64(cnt+1)
	x := int(k * float64(rX-lX))
	pos = image.Point{lX + x, bY}
	return
}

//
//	Selecter interface
//

var _ Selecter = (*Arch)(nil)

func (c *Container) onSelect() {
	for _, ch := range c.Children {
		hit, _ := ch.CheckHit(c.Pos)
		if hit {
			ch.Select()
		} else {
			ch.Deselect()
		}
	}
	for _, p := range c.ports {
		hit, _ := p.CheckHit(c.Pos)
		if hit {
			p.Select()
		} else {
			p.Deselect()
		}
	}
	return
}

func (c *Container) onDeselect() {
	for _, ch := range c.Children {
		ch.Deselect()
	}
	for _, p := range c.ports {
		p.Deselect()
	}
	return
}

func (c *Container) SelectChild(child ContainerChild) {
	c.selectedChild = -1
	for i, ch := range c.Children {
		if ch == child {
			ch.Select()
			c.selectedChild = i
		} else {
			ch.Deselect()
		}
	}
}

func (c Container) GetSelectedChild() (ok bool, child ContainerChild) {
	if c.selectedChild == -1 {
		return
	}
	ok = true
	child = c.Children[c.selectedChild]
	return
}

func (c *Container) SelectPort(userObj graph.Positioner) {
	c.selectedPort = -1
	for i, p := range c.ports {
		if p.UserObj == userObj {
			c.selectedPort = i
			p.Select()
		} else {
			p.Deselect()
		}
	}
}

func (c *Container) SelectModePort(userObj graph.ModePositioner) {
	c.selectedPort = -1
	for i, p := range c.ports {
		if p.UserObj2 == userObj {
			c.selectedPort = i
			p.Select()
		} else {
			p.Deselect()
		}
	}
}

func (c Container) GetSelectedPort() (ok bool, userObj graph.Positioner) {
	if c.selectedPort == -1 {
		return
	}
	ok = true
	userObj = c.ports[c.selectedPort].UserObj
	return
}

func (c Container) GetSelectedModePort() (ok bool, userObj graph.ModePositioner) {
	if c.selectedPort == -1 {
		return
	}
	ok = true
	userObj = c.ports[c.selectedPort].UserObj2
	return
}

//
//	Highlighter interface
//

var _ Highlighter = (*Container)(nil)

func (c *Container) onHighlight(hit bool, pos image.Point) (modified bool) {
	c.selectedChild = -1
	c.selectedPort = -1
	if hit {
		for i, pr := range c.Children {
			phit, mod := pr.CheckHit(pos)
			if phit {
				c.selectedChild = i
			}
			modified = modified || mod
		}
		for i, p := range c.ports {
			phit, mod := p.CheckHit(pos)
			if phit {
				c.selectedPort = i
			}
			modified = modified || mod
		}
	} else {
		for _, pr := range c.Children {
			_, mod := pr.CheckHit(pos)
			modified = modified || mod
		}
		for _, p := range c.ports {
			_, mod := p.CheckHit(pos)
			modified = modified || mod
		}
	}
	return
}

//
//	Drawer interface
//

func (c Container) Draw(ctxt interface{}) {
	c.ContainerDefaultDraw(ctxt)
}

func (c Container) ContainerDefaultDraw(ctxt interface{}) {
	c.NamedBoxDefaultDraw(ctxt)
	for _, p := range c.ports {
		p.Draw(ctxt)
	}
	for _, pr := range c.Children {
		pr.Draw(ctxt)
	}
	c.onDraw(ctxt)
}

//
//		Private functions
//
func ContainerFit(outer, inner image.Rectangle) image.Rectangle {
	borderTop := image.Point{-18, -30}
	borderBottom := image.Point{18, 18}
	test := image.Rectangle{inner.Min.Add(borderTop), inner.Max.Add(borderBottom)}
	if outer.Size().X == 0 {
		return test
	}
	return outer.Union(test)
}

func (c Container) portIsSideBorder(pos image.Point) bool {
	lX, rX, _, _ := c.portCorners()
	return (pos.X == lX || pos.X == rX)
}

func (c Container) portIsCorner(pos image.Point) bool {
	lX, rX, tY, bY := c.portCorners()
	return (pos.X == lX || pos.X == rX) && (pos.Y == tY || pos.Y == bY)
}

func (c Container) portCorners() (lX, rX, tY, bY int) {
	lX = c.BBox().Min.X + 1
	rX = c.BBox().Max.X - c.config.portWidth - 1
	tY = c.BBox().Min.Y + 1
	bY = c.BBox().Max.Y - c.config.portHeight - 1
	return
}

func (c Container) portClipPos(pos image.Point) (ret image.Point) {
	ret = pos
	lX, rX, tY, bY := c.portCorners()
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
