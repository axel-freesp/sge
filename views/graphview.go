package views

import (
	//	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image"
	"log"
	"math"
)

type Context interface {
	SelectNode(node graph.NodeObject)     // single click selection
	EditNode(node graph.NodeObject)       // double click selection
	SelectPort(port freesp.Port)          // single click selection
	SelectConnect(conn freesp.Connection) // single click selection
}

type GraphView interface {
	Widget() *gtk.Widget
	Sync()
	Select(obj freesp.TreeElement)
}

type graphView struct {
	parent      *ScaledView
	area        *gtk.DrawingArea
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	sgType      freesp.SignalGraphType
	context     Context

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*graphView)(nil)
var _ GraphView = (*graphView)(nil)

func GraphViewNew(g freesp.SignalGraphType, context Context) (viewer *graphView, err error) {
	viewer = &graphView{nil, nil, nil, nil, g, context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func (v *graphView) Widget() *gtk.Widget {
	return v.parent.Widget()
}

func (v *graphView) init() (err error) {
	v.parent, err = ScaledViewNew(v)
	if err != nil {
		return
	}
	v.area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
	v.area.Connect("draw", areaDrawCallback, v)
	v.area.Connect("motion-notify-event", areaMotionCallback, v)
	v.area.Connect("button-press-event", areaButtonCallback, v)
	v.area.Connect("button-release-event", areaButtonCallback, v)
	v.area.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))
	v.parent.Container().Add(v.area)
	return
}

func (v *graphView) Sync() {
	g := v.sgType
	v.nodes = make([]graph.NodeIf, len(g.Nodes()))

	var numberOfConnections = 0
	for _, n := range g.Nodes() {
		for _, p := range n.OutPorts() {
			numberOfConnections += len(p.Connections())
		}
	}
	v.connections = make([]graph.ConnectIf, numberOfConnections)

	for i, n := range g.Nodes() {
		v.nodes[i] = graph.NodeNew(n.Position(), n)
	}
	var index = 0
	for _, n := range g.Nodes() {
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
	v.area.SetSizeRequest(v.calcSceneWidth(), v.calcSceneHeight())
	v.drawAll()
}

/*
 *		Handle selection in treeview
 */

func (v *graphView) Select(obj freesp.TreeElement) {
	switch obj.(type) {
	case freesp.Node:
		v.deselectConnects()
		v.selectNode(obj.(freesp.Node))
	case freesp.Port:
		v.deselectConnects()
		n, ok := v.selectNode(obj.(freesp.Port).Node())
		if ok {
			n.(*graph.Node).SelectPort(obj.(freesp.Port))
		}
	case freesp.Connection:
		v.deselectNodes()
		v.selectConnect(obj.(freesp.Connection))
	}
}

func (v *graphView) selectNode(obj freesp.Node) (ret graph.NodeIf, ok bool) {
	var n graph.NodeIf
	for _, n = range v.nodes {
		if obj == n.UserObj().(freesp.Node) {
			if !n.Select() {
				v.repaintNode(n)
			}
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

func (v *graphView) deselectNodes() {
	var n graph.NodeIf
	for _, n = range v.nodes {
		if n.Deselect() {
			v.repaintNode(n)
		}
	}
}

func (v *graphView) selectConnect(obj freesp.Connection) {
	n1 := v.findNode(obj.From().Node().Name())
	n2 := v.findNode(obj.To().Node().Name())
	for _, c := range v.connections {
		if n1 == c.From() && n2 == c.To() &&
			n1.GetOutPort(c.FromId()) == obj.From() &&
			n2.GetInPort(c.ToId()) == obj.To() {
			if !c.Select() {
				v.repaintConnection(c)
			}
		} else {
			if c.Deselect() {
				v.repaintConnection(c)
			}
		}
	}
}

func (v *graphView) deselectConnects() {
	for _, c := range v.connections {
		if c.Deselect() {
			v.repaintConnection(c)
		}
	}
}

/*
 *		ScaledScene interface
 */

func (v *graphView) Update() (width, height int) {
	width = v.calcSceneWidth()
	height = v.calcSceneHeight()
	v.area.SetSizeRequest(width, height)
	v.area.QueueDrawArea(0, 0, width, height)
	return
}

func (v *graphView) calcSceneWidth() int {
	return int(float64(v.bBox().Max.X+50) * v.parent.Scale())
}

func (v *graphView) calcSceneHeight() int {
	return int(float64(v.bBox().Max.Y+50) * v.parent.Scale())
}

/*
 *		areaButtonCallback
 */

func areaButtonCallback(area *gtk.DrawingArea, event *gdk.Event, v *graphView) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X() / v.parent.Scale()), int(ev.Y() / v.parent.Scale())}
	switch ev.Type() {
	case gdk.EVENT_BUTTON_PRESS:
		v.button1Pressed = true
		v.dragOffs = pos
		v.handleNodeSelect(pos)
		v.handleConnectSelect(pos)
	case gdk.EVENT_2BUTTON_PRESS:
		log.Println("areaButtonCallback 2BUTTON_PRESS")
		for _, n := range v.nodes {
			if n.IsSelected() {
				v.context.EditNode(n.UserObj())
				break
			}
		}

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

func (v *graphView) handleNodeSelect(pos image.Point) {
	for _, n := range v.nodes {
		hit, _ := n.CheckHit(pos)
		if hit {
			if !n.Select() {
				v.repaintNode(n)
			}
			v.context.SelectNode(n.UserObj())
			ok, port := n.(*graph.Node).GetSelectedPort()
			if ok {
				v.context.SelectPort(port.(freesp.Port))
			}
		} else {
			if n.Deselect() {
				v.repaintNode(n)
			}
		}
	}
}

func (v *graphView) handleConnectSelect(pos image.Point) {
	for _, c := range v.connections {
		hit, _ := c.CheckHit(pos)
		if hit {
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
}

/*
 *		areaMotionCallback
 */

func areaMotionCallback(area *gtk.DrawingArea, event *gdk.Event, v *graphView) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x / v.parent.Scale()), int(y / v.parent.Scale())}
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

func (v *graphView) handleDrag(pos image.Point) {
	for _, n := range v.nodes {
		if n.IsSelected() {
			box := n.BBox()
			if !overlaps(v.nodes, box.Min.Add(pos.Sub(v.dragOffs))) {
				v.repaintNode(n)
				box = box.Add(pos.Sub(v.dragOffs))
				v.dragOffs = pos
				n.SetPosition(box.Min)
				v.repaintNode(n)
			}
		}
	}
}

func (v *graphView) handleMouseover(pos image.Point) {
	for _, n := range v.nodes {
		_, mod := n.CheckHit(pos)
		if mod {
			v.repaintNode(n)
		}
	}
	for _, c := range v.connections {
		_, modified := c.CheckHit(pos)
		if modified {
			v.repaintConnection(c)
		}
	}
}

/*
 *		areaDrawCallback
 */

func areaDrawCallback(area *gtk.DrawingArea, context *cairo.Context, v *graphView) bool {
	context.Scale(v.parent.Scale(), v.parent.Scale())
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	v.drawNodes(context, r)
	v.drawConnections(context, r)
	return true
}

func (v *graphView) drawConnections(context *cairo.Context, r image.Rectangle) {
	for _, c := range v.connections {
		box := c.BBox()
		if r.Overlaps(box) {
			c.Draw(context)
		}
	}
}

func (v *graphView) drawNodes(context *cairo.Context, r image.Rectangle) {
	for _, o := range v.nodes {
		if !o.IsSelected() && o.BBox().Overlaps(r) {
			o.Draw(context)
		}
	}
	for _, o := range v.nodes {
		if o.IsSelected() && o.BBox().Overlaps(r) {
			o.Draw(context)
		}
	}
}

/*
 *		Private functions
 */

func (v *graphView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *graphView) scaleCoord(x int, roundUp bool) int {
	if roundUp {
		return int(math.Ceil(float64(x) * v.parent.Scale()))
	} else {
		return int(float64(x) * v.parent.Scale())
	}
}

func (v *graphView) drawAll() {
	v.drawScene(v.bBox())
}

func (v *graphView) bBox() (box image.Rectangle) {
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

func (v *graphView) drawScene(r image.Rectangle) {
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.scaleCoord(x, false), v.scaleCoord(y, false), v.scaleCoord(w, true)+1, v.scaleCoord(h, true)+1)
	return
}

func (v *graphView) nodeDrawBox(n graph.NodeIf) image.Rectangle {
	ret := n.BBox()
	for _, c := range v.connections {
		if c.IsLinked(n.Name()) {
			ret = ret.Union(c.BBox())
		}
	}
	return ret
}

func (v *graphView) repaintNode(n graph.NodeIf) {
	v.drawScene(v.nodeDrawBox(n))
}

func (v *graphView) repaintConnection(c graph.ConnectIf) {
	v.drawScene(c.BBox())
}

// Check if the selected object would overlap with any other
// if we move it to p coordinates:
func overlaps(list []graph.NodeIf, p image.Point) bool {
	var box image.Rectangle
	for _, n := range list {
		if n.IsSelected() {
			box = n.BBox()
			break
		}
	}
	newBox := box.Add(p.Sub(box.Min))
	for _, n := range list {
		if !n.IsSelected() && newBox.Overlaps(n.BBox()) {
			return true
		}
	}
	return false
}
