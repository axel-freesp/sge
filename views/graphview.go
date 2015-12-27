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

var (
	nodeWidth  = graph.NumericOption(graph.NodeWidth)
	nodeHeight = graph.NumericOption(graph.NodeHeight)
)

type GraphView struct {
	ScrolledView
	area        *gtk.DrawingArea
	nodes       []graph.Dragable
	connections []graph.Connection
	signalGraph freesp.SignalGraph

	dragOffs              image.Point
	selected, highlighted int
	button1Pressed        bool
	conOut, conIn, portDY image.Point
	width, height         int
}

func GraphViewNew(graph freesp.SignalGraph, width, height int) (viewer *GraphView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &GraphView{*v, nil, nil, nil, graph, image.Point{}, -1, -1, false,
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
	v.nodes = make([]graph.Dragable, len(g.ItsType().Nodes()))

	var numberOfConnections = 0
	for _, n := range g.ItsType().Nodes() {
		for _, p := range n.OutPorts() {
			numberOfConnections += len(p.Connections())
		}
	}
	v.connections = make([]graph.Connection, numberOfConnections)

	dy := graph.NumericOption(graph.PortDY)
	for i, n := range g.ItsType().Nodes() {
		pos := n.Position()
		box := image.Rect(pos.X, pos.Y, pos.X+nodeWidth, pos.Y+nodeHeight+numPorts(n)*dy)
		v.nodes[i] = graph.NodeNew(box, n)
	}
	log.Printf("GraphView.init: initialized %d nodes\n", len(g.ItsType().Nodes()))

	var index = 0
	for _, n := range g.ItsType().Nodes() {
		from := v.findNode(n.Name())
		for _, p := range n.OutPorts() {
			fromId := v.outPortIndex(n, p.PortName())
			for _, c := range p.Connections() {
				to := v.findNode(c.Node().Name())
				toId := v.inPortIndex(c.Node(), c.PortName())
				v.connections[index] = graph.ConnectionNew(from, to, fromId, toId)
				index++
			}
		}
	}
	log.Printf("GraphView.init: initialized %d edges\n", index)
	v.DrawAll()
}

func numPorts(n graph.GObject) int {
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

func (v *GraphView) inPortIndex(n freesp.Node, portName string) int {
	for i, p := range n.InPorts() {
		if p.PortName() == portName {
			return i
		}
	}
	return -1
}

func (v *GraphView) outPortIndex(n freesp.Node, portName string) int {
	for i, p := range n.OutPorts() {
		if p.PortName() == portName {
			return i
		}
	}
	return -1
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

func (v *GraphView) connectionsBox(drg graph.Dragable) image.Rectangle {
	ret := drg.BBox()
	for _, c := range v.connections {
		if drg.Name() == c.From.Name() || drg.Name() == c.To.Name() {
			p1, p2 := v.connectionPoints(c)
			r := c.BBox(p1, p2)
			ret = ret.Union(r)
		}
	}
	return ret
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
	v.area.QueueDrawArea(0, 0, v.width, v.height)
}

func (v *GraphView) DrawScene(r image.Rectangle) {
	v.area.QueueDrawArea(r.Min.X, r.Min.Y, r.Size().X, r.Size().Y)
	return
}

func (v *GraphView) ClearScene(r image.Rectangle) {
	v.area.QueueDrawArea(r.Min.X, r.Min.Y, r.Size().X, r.Size().Y)
	return
}

func (v *GraphView) repaintObj(drg graph.Dragable) {
	v.DrawScene(drg.BBox())
}

func areaButtonCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X()), int(ev.Y())}
	idx, ok := graph.HitObj(v.nodes, pos)
	switch ev.Type() {
	case gdk.EVENT_BUTTON_PRESS:
		v.dragOffs = pos
		oldSelected := v.selected
		v.selected = idx
		if ok {
			v.repaintObj(v.nodes[v.selected])
			if oldSelected >= 0 && v.selected != oldSelected {
				v.repaintObj(v.nodes[oldSelected])
			}
		} else if oldSelected >= 0 {
			v.repaintObj(v.nodes[oldSelected])
		}
		v.button1Pressed = true
	case gdk.EVENT_2BUTTON_PRESS:
		log.Println("areaButtonCallback 2BUTTON_PRESS")

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

func (v *GraphView) check(pos image.Point) (r image.Rectangle, ok bool) {
	return v.nodes[v.selected].Check(pos)
}

func areaMotionCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x), int(y)}
	idx, ok := graph.HitObj(v.nodes, pos)

	if v.button1Pressed {
		if v.selected != -1 {
			s := v.nodes[v.selected]
			box := s.BBox()
			if !graph.Overlaps(v.nodes, v.selected, box.Min.Add(pos.Sub(v.dragOffs))) {
				clearbox := v.connectionsBox(s)
				v.ClearScene(clearbox)
				box = box.Add(pos.Sub(v.dragOffs))
				v.dragOffs = pos
				s.SetPosition(box.Min)
				v.DrawScene(clearbox.Union(v.connectionsBox(s)))
			}
		}
	} else {
		oldHighlighted := v.highlighted
		if idx == v.selected {
			v.highlighted = -1
			if ok {
				r, ok := v.check(pos)
				if ok {
					v.DrawScene(r)
				}
			}
			ok = false
		} else {
			v.highlighted = idx
		}
		if ok {
			v.repaintObj(v.nodes[v.highlighted])
			if oldHighlighted >= 0 && v.highlighted != oldHighlighted {
				v.repaintObj(v.nodes[oldHighlighted])
			}
		} else if oldHighlighted >= 0 {
			v.repaintObj(v.nodes[oldHighlighted])
		}
	}

}

// Draw all connections
func (v *GraphView) drawConnections(context *cairo.Context, r image.Rectangle) {
	context.SetSourceRGB(0.0, 0.0, 0.0)
	context.SetLineWidth(2)
	for _, c := range v.connections {
		p1, p2 := v.connectionPoints(c)
		box := c.BBox(p1, p2)
		if lineRectOverlap(r, box.Min, box.Max) {
			graph.DrawArrow(context, p1, p2)
		}
	}
}

func (v *GraphView) connectionPoints(c graph.Connection) (p1, p2 image.Point) {
	p1 = c.From.BBox().Min.Add(v.conOut.Add(v.portDY.Mul(c.FromPort)))
	p2 = c.To.BBox().Min.Add(v.conIn.Add(v.portDY.Mul(c.ToPort)))
	return
}

// Draw all nodes
func (v *GraphView) drawNodes(context *cairo.Context, r image.Rectangle) {
	for i, o := range v.nodes {
		if o.BBox().Overlaps(r) {
			switch i {
			case v.selected:
				o.Draw(context, graph.SelectedMode)
			case v.highlighted:
				o.Draw(context, graph.HighlightMode)
			default:
				o.Draw(context, graph.NormalMode)
			}
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
