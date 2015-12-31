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
	"unsafe"
)

const (
	scaleMin   = -4.0
	scaleMax   = 3.0
	scaleStep  = 0.01
	scalePStep = 1.0
	scalePage  = 1.0
)

type Context interface {
	SelectNode(node graph.NodeObject)     // single click selection
	EditNode(node graph.NodeObject)       // double click selection
	SelectPort(port freesp.Port)          // single click selection
	SelectConnect(conn freesp.Connection) // single click selection
}

type GraphView struct {
	layoutBox   *gtk.Box
	scene       *ScrolledView
	area        *gtk.DrawingArea
	scaler      *gtk.Scale
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	sgType      freesp.SignalGraphType
	context     Context

	dragOffs       image.Point
	button1Pressed bool
	scale          float64
}

func GraphViewNew(g freesp.SignalGraphType, width, height int, context Context) (viewer *GraphView, err error) {
	viewer = &GraphView{nil, nil, nil, nil, nil, nil, g, context, image.Point{}, false,
		1.0}
	err = viewer.init(width, height)
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func (v *GraphView) Widget() *gtk.Widget {
	return (*gtk.Widget)(unsafe.Pointer(v.layoutBox))
}

func (v *GraphView) init(width, height int) (err error) {
	v.layoutBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}
	v.scene, err = ScrolledViewNew(width, height)
	if err != nil {
		return
	}
	v.area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
	var adj *gtk.Adjustment
	adj, err = gtk.AdjustmentNew(0.0, scaleMin, scaleMax, scaleStep, scalePStep, scalePage)
	v.scaler, err = gtk.ScaleNew(gtk.ORIENTATION_HORIZONTAL, adj)
	if err != nil {
		return
	}
	v.area.Connect("draw", areaDrawCallback, v)
	v.area.Connect("motion-notify-event", areaMotionCallback, v)
	v.area.Connect("button-press-event", areaButtonCallback, v)
	v.area.Connect("button-release-event", areaButtonCallback, v)
	v.scaler.Connect("value-changed", scaleValueCallback, v)
	v.area.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))
	v.layoutBox.Add(v.scene.Widget())
	v.layoutBox.Add(v.scaler)
	v.scene.scrolled.Add(v.area)
	return
}

func (v *GraphView) Sync() {
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

func (v *GraphView) Select(obj freesp.TreeElement) {
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

func (v *GraphView) selectNode(obj freesp.Node) (ret graph.NodeIf, ok bool) {
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

func (v *GraphView) deselectNodes() {
	var n graph.NodeIf
	for _, n = range v.nodes {
		if n.Deselect() {
			v.repaintNode(n)
		}
	}
}

func (v *GraphView) selectConnect(obj freesp.Connection) {
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

func (v *GraphView) deselectConnects() {
	for _, c := range v.connections {
		if c.Deselect() {
			v.repaintConnection(c)
		}
	}
}

/*
 *		scaleValueCallback
 */

func scaleValueCallback(r *gtk.Scale, v *GraphView) {
	oldScale := v.scale
	oldX := v.scene.hadj.GetValue()
	oldY := v.scene.vadj.GetValue()
	visMx := float64(v.layoutBox.GetAllocatedWidth()) / 2.0
	visMy := float64(v.layoutBox.GetAllocatedHeight()) / 2.0
	v.scale = math.Exp2(r.GetValue())
	newW := v.calcSceneWidth()
	newH := v.calcSceneHeight()
	v.area.SetSizeRequest(newW, newH)
	v.area.QueueDrawArea(0, 0, newW, newH)
	v.scene.hadj.SetUpper(float64(newW))
	v.scene.vadj.SetUpper(float64(newH))
	v.scene.hadj.SetValue((oldX+visMx)*v.scale/oldScale - visMx)
	v.scene.vadj.SetValue((oldY+visMy)*v.scale/oldScale - visMy)
}

func (v *GraphView) calcSceneWidth() int {
	return int(float64(v.bBox().Max.X+50) * v.scale)
}

func (v *GraphView) calcSceneHeight() int {
	return int(float64(v.bBox().Max.Y+50) * v.scale)
}

/*
 *		areaButtonCallback
 */

func areaButtonCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X() / v.scale), int(ev.Y() / v.scale)}
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

func (v *GraphView) handleNodeSelect(pos image.Point) {
	for _, n := range v.nodes {
		hit, _ := n.CheckHit(pos)
		if hit {
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
}

func (v *GraphView) handleConnectSelect(pos image.Point) {
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

func areaMotionCallback(area *gtk.DrawingArea, event *gdk.Event, v *GraphView) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x / v.scale), int(y / v.scale)}
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

func (v *GraphView) handleDrag(pos image.Point) {
	for _, n := range v.nodes {
		if n.IsSelected() {
			box := n.BBox()
			if !graph.Overlaps(v.nodes, box.Min.Add(pos.Sub(v.dragOffs))) {
				v.repaintNode(n)
				box = box.Add(pos.Sub(v.dragOffs))
				v.dragOffs = pos
				n.SetPosition(box.Min)
				v.repaintNode(n)
			}
		}
	}
}

func (v *GraphView) handleMouseover(pos image.Point) {
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

func areaDrawCallback(area *gtk.DrawingArea, context *cairo.Context, v *GraphView) bool {
	context.Scale(v.scale, v.scale)
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	v.drawNodes(context, r)
	v.drawConnections(context, r)
	return true
}

func (v *GraphView) drawConnections(context *cairo.Context, r image.Rectangle) {
	for _, c := range v.connections {
		box := c.BBox()
		if r.Overlaps(box) {
			c.Draw(context)
		}
	}
}

func (v *GraphView) drawNodes(context *cairo.Context, r image.Rectangle) {
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

func (v *GraphView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *GraphView) scaleCoord(x int, roundUp bool) int {
	if roundUp {
		return int(math.Ceil(float64(x) * v.scale))
	} else {
		return int(float64(x) * v.scale)
	}
}

func (v *GraphView) drawAll() {
	v.drawScene(v.bBox())
}

func (v *GraphView) bBox() (box image.Rectangle) {
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

func (v *GraphView) drawScene(r image.Rectangle) {
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.scaleCoord(x, false), v.scaleCoord(y, false), v.scaleCoord(w, true)+1, v.scaleCoord(h, true)+1)
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
	v.drawScene(v.nodeDrawBox(n))
}

func (v *GraphView) repaintConnection(c graph.ConnectIf) {
	v.drawScene(c.BBox())
}
