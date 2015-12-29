package views

import (
	//	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image"
	"log"
)

var ( // TODO: move to connection.go
	nodeWidth  = graph.NumericOption(graph.NodeWidth)
	nodeHeight = graph.NumericOption(graph.NodeHeight)
)

type Context interface {
	SelectNode(node graph.NodeObject)     // single click selection
	EditNode(node graph.NodeObject)       // double click selection
	SelectPort(port freesp.Port)          // single click selection
	SelectConnect(conn freesp.Connection) // single click selection
}

type GraphView struct {
	ScrolledView
	area        *gtk.DrawingArea
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	signalGraph freesp.SignalGraph
	context     Context

	dragOffs              image.Point
	selected, highlighted int
	button1Pressed        bool
	conOut, conIn, portDY image.Point
	width, height         int
}

func GraphViewNew(graph freesp.SignalGraph, width, height int, context Context) (viewer *GraphView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &GraphView{*v, nil, nil, nil, graph, context, image.Point{}, -1, -1, false,
		image.Point{}, image.Point{}, image.Point{}, width, height}
	err = viewer.init()
	return
}

func (v *GraphView) init() (err error) {
	v.area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
	v.area.SetSizeRequest(v.width, v.height)

	v.Sync()

	//portW := graph.NumericOption(graph.PortW)
	portH := graph.NumericOption(graph.PortH)
	portX0 := graph.NumericOption(graph.PortX0) + graph.NumericOption(graph.NodePadX)
	portY0 := graph.NumericOption(graph.PortY0) + graph.NumericOption(graph.NodePadY)
	portDY := graph.NumericOption(graph.PortDY)
	XIn := portX0
	XOut := nodeWidth - portX0
	v.conOut = image.Point{XOut, portY0 + portH/2}
	v.conIn = image.Point{XIn, portY0 + portH/2}
	v.portDY = image.Point{0, portDY}

	v.area.Connect("draw", areaDrawCallback, v)
	v.area.Connect("motion-notify-event", areaMotionCallback, v)
	v.area.Connect("button-press-event", areaButtonCallback, v)
	v.area.Connect("button-release-event", areaButtonCallback, v)
	v.area.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))
	v.scrolled.Add(v.area)
	return
}

func (v *GraphView) Sync() {
	g := v.signalGraph
	v.nodes = make([]graph.NodeIf, len(g.ItsType().Nodes()))

	var numberOfConnections = 0
	for _, n := range g.ItsType().Nodes() {
		for _, p := range n.OutPorts() {
			numberOfConnections += len(p.Connections())
		}
	}
	v.connections = make([]graph.ConnectIf, numberOfConnections)

	dy := graph.NumericOption(graph.PortDY)
	for i, n := range g.ItsType().Nodes() {
		pos := n.Position()
		box := image.Rect(pos.X, pos.Y, pos.X+nodeWidth, pos.Y+nodeHeight+numPorts(n)*dy)
		v.nodes[i] = graph.NodeNew(box, n)
	}
	var index = 0
	for _, n := range g.ItsType().Nodes() {
		from := v.findNode(n.Name())
		for _, p := range n.OutPorts() {
			fromId := from.OutPortIndex(p.Name())
			for _, c := range p.Connections() {
				to := v.findNode(c.Node().Name())
				toId := to.InPortIndex(c.Name())
				v.connections[index] = graph.ConnectionNew(from, to, fromId, toId)
				index++
			}
		}
	}
	v.DrawAll()
}

func (v *GraphView) selectNode(obj freesp.Node) (ret graph.NodeIf, ok bool) {
	var i int
	var n graph.NodeIf
	v.selected = -1
	for i, n = range v.nodes {
		if obj == n.UserObj().(freesp.Node) {
			if !n.Select() {
				v.repaintNode(n)
			}
			v.selected = i
			ret = n
			ok = true
		} else {
			if n.Deselect() {
				v.repaintNode(n)
			}
		}
	}
	return
}

// Is this cool?
func (v *GraphView) Select(obj freesp.TreeElement) {
	switch obj.(type) {
	case freesp.Node:
		v.selectNode(obj.(freesp.Node))
	case freesp.Port:
		n, ok := v.selectNode(obj.(freesp.Port).Node())
		if ok {
			n.(*graph.Node).SelectPort(obj.(freesp.Port))
		}
	}
}

func numPorts(n freesp.Porter) int {
	npi := n.NumInPorts()
	npo := n.NumOutPorts()
	return tool.MaxInt(npi, npo)
}

func (v *GraphView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *GraphView) BBox() (box image.Rectangle) {
	emptyRect := image.Rectangle{}
	for _, o := range v.nodes {
		if box == emptyRect {
			box = o.BBox()
		} else {
			box = box.Union(o.BBox())
		}
	}
	return
}

func (v *GraphView) DrawAll() {
	v.DrawScene(v.BBox())
}

func (v *GraphView) DrawScene(r image.Rectangle) {
	v.area.QueueDrawArea(r.Min.X, r.Min.Y, r.Size().X, r.Size().Y)
	return
}

func (v *GraphView) nodeDrawBox(n graph.NodeIf) image.Rectangle {
	ret := n.BBox()
	for _, c := range v.connections {
		if c.IsLinked(n.Name()) {
			ret = ret.Union(c.BBox())
		}
	}
	return ret
}

func (v *GraphView) repaintNode(n graph.NodeIf) {
	v.DrawScene(v.nodeDrawBox(n))
}

func (v *GraphView) repaintConnection(c graph.ConnectIf) {
	v.DrawScene(c.BBox())
}

func areaButtonCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X()), int(ev.Y())}
	switch ev.Type() {
	case gdk.EVENT_BUTTON_PRESS:
		v.button1Pressed = true
		v.dragOffs = pos
		v.selected = -1
		for i, n := range v.nodes {
			if n.CheckHit(pos) {
				v.selected = i
				if !n.Select() {
					v.repaintNode(n)
				}
				v.context.SelectNode(n.UserObj())
				ok, port := n.(*graph.Node).GetSelectedPort()
				if ok {
					v.context.SelectPort(port)
				}
			} else {
				if n.Deselect() {
					v.repaintNode(n)
				}
			}
		}
		for _, c := range v.connections {
			if c.CheckHit(pos) {
				c.Select()
				n1 := c.From()
				p1 := n1.GetOutPort(c.FromId())
				n2 := c.To()
				p2 := n2.GetInPort(c.ToId())
				v.context.SelectConnect(p1.(freesp.Port).Connection(p2.(freesp.Port)))
				v.repaintConnection(c)
			} else {
				if c.Deselect() {
					v.repaintConnection(c)
				}
			}
		}
	case gdk.EVENT_2BUTTON_PRESS:
		log.Println("areaButtonCallback 2BUTTON_PRESS")
		if v.selected != -1 {
			v.context.EditNode(v.nodes[v.selected].UserObj())
		}

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

func (v *GraphView) handleDrag(pos image.Point) {
	if v.selected != -1 {
		s := v.nodes[v.selected]
		box := s.BBox()
		if !graph.Overlaps(v.nodes, v.selected, box.Min.Add(pos.Sub(v.dragOffs))) {
			v.repaintNode(s)
			box = box.Add(pos.Sub(v.dragOffs))
			v.dragOffs = pos
			s.SetPosition(box.Min)
			v.repaintNode(s)
		}
	}
}

func (v *GraphView) handleMouseover(pos image.Point) {
	var idx int
	ok := false
	for i, n := range v.nodes {
		if n.CheckHit(pos) {
			ok = true
			idx = i
			break
		}
	}
	oldHighlighted := v.highlighted
	if idx == v.selected {
		v.highlighted = -1
		if ok {
			ok = v.nodes[v.selected].CheckHit(pos)
			if ok {
				v.repaintNode(v.nodes[v.selected])
			}
		}
		ok = false
	} else {
		v.highlighted = idx
	}
	if ok {
		v.repaintNode(v.nodes[v.highlighted])
		if oldHighlighted >= 0 && v.highlighted != oldHighlighted {
			v.repaintNode(v.nodes[oldHighlighted])
		}
	} else if oldHighlighted >= 0 {
		v.repaintNode(v.nodes[oldHighlighted])
	}
	for _, c := range v.connections {
		if c.CheckHit(pos) {
			v.repaintConnection(c)
		}
	}
}

func areaMotionCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x), int(y)}
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

// Draw all connections
func (v *GraphView) drawConnections(context *cairo.Context, r image.Rectangle) {
	for _, c := range v.connections {
		box := c.BBox()
		if lineRectOverlap(r, box.Min, box.Max) {
			c.Draw(context)
		}
	}
}

// Draw all nodes
func (v *GraphView) drawNodes(context *cairo.Context, r image.Rectangle) {
	for _, o := range v.nodes {
		if o.BBox().Overlaps(r) {
			o.Draw(context)
		}
	}
}

func areaDrawCallback(area *gtk.DrawingArea, context *cairo.Context, v *GraphView) bool {
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	v.drawNodes(context, r)
	v.drawConnections(context, r)
	return true
}

func lineRectOverlap(r image.Rectangle, p1, p2 image.Point) bool {
	if !r.Overlaps(image.Rectangle{p1, p2}.Canon()) {
		return false
	}
	ps := r.Size()
	pw := image.Point{ps.X, 0}
	ph := image.Point{0, ps.Y}
	l := p2.Sub(p1)
	crossProd := func(p1, p2 image.Point) int {
		return p1.X*p2.Y - p1.Y*p2.X
	}
	sign := func(x int) bool {
		return x < 0
	}
	c1 := sign(crossProd(l, r.Min.Sub(p1)))
	c2 := sign(crossProd(l, r.Min.Add(pw).Sub(p1)))
	c3 := sign(crossProd(l, r.Min.Add(ph).Sub(p1)))
	c4 := sign(crossProd(l, r.Max.Sub(p1)))
	return (c1 || c2 || c3 || c4) != (c1 && c2 && c3 && c4)
}
