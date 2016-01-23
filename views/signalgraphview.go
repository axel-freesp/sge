package views

import (
	//	"fmt"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image"
	"log"
)

type signalGraphView struct {
	parent      *ScaledView
	area        DrawArea
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	sgType      interfaces.GraphObject
	context     interfaces.Context

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*signalGraphView)(nil)
var _ GraphView = (*signalGraphView)(nil)

func SignalGraphViewNew(g interfaces.GraphObject, context interfaces.Context) (viewer *signalGraphView, err error) {
	viewer = &signalGraphView{nil, DrawArea{}, nil, nil, g, context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func SignalGraphViewDestroy(viewer GraphView) {
	DrawAreaDestroy(viewer.(*signalGraphView).area)
}

func (v *signalGraphView) Widget() *gtk.Widget {
	return v.parent.Widget()
}

func (v *signalGraphView) init() (err error) {
	v.parent, err = ScaledViewNew(v)
	if err != nil {
		return
	}
	v.area, err = DrawAreaInit(v)
	if err != nil {
		return
	}
	v.parent.Container().Add(v.area)
	return
}

func (v *signalGraphView) Sync() {
	g := v.sgType
	v.nodes = make([]graph.NodeIf, len(g.NodeObjects()))

	var numberOfConnections = 0
	for _, n := range g.NodeObjects() {
		for _, p := range n.GetOutPorts() {
			numberOfConnections += len(p.ConnectionObjects())
		}
	}
	v.connections = make([]graph.ConnectIf, numberOfConnections)

	for i, n := range g.NodeObjects() {
		v.nodes[i] = graph.NodeNew(n.Position(), n)
	}
	var index = 0
	for _, n := range g.NodeObjects() {
		from := v.findNode(n.Name())
		for _, p := range n.GetOutPorts() {
			fromId := from.OutPortIndex(p.Name())
			for _, c := range p.ConnectionObjects() {
				to := v.findNode(c.NodeObject().Name())
				toId := to.InPortIndex(c.Name())
				v.connections[index] = graph.ConnectionNew(from, to, fromId, toId)
				index++
			}
		}
	}
	v.area.SetSizeRequest(v.calcSceneWidth(), v.calcSceneHeight())
	v.drawAll()
}

func (v signalGraphView) IdentifyGraph(g interfaces.GraphObject) bool {
	return g == v.sgType
}

func (v signalGraphView) IdentifyPlatform(g interfaces.PlatformObject) bool {
	return false
}

func (v signalGraphView) IdentifyMapping(g interfaces.MappingObject) bool {
	return false
}

//
//		Handle selection in treeview
//

func (v *signalGraphView) Select(obj interfaces.GraphElement) {
	switch obj.(type) {
	case interfaces.NodeObject:
		v.deselectConnects()
		v.selectNode(obj.(interfaces.NodeObject))
	case interfaces.PortObject:
		v.deselectConnects()
		n, ok := v.selectNode(obj.(interfaces.PortObject).NodeObject())
		if ok {
			n.(*graph.Node).SelectPort(obj.(interfaces.PortObject))
			v.repaintNode(n)
		}
	case interfaces.ConnectionObject:
		v.deselectNodes()
		v.selectConnect(obj.(interfaces.ConnectionObject))
	}
}

func (v *signalGraphView) selectNode(obj interfaces.NodeObject) (ret graph.NodeIf, ok bool) {
	var n graph.NodeIf
	for _, n = range v.nodes {
		if obj == n.UserObj() {
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

func (v *signalGraphView) deselectNodes() {
	var n graph.NodeIf
	for _, n = range v.nodes {
		if n.Deselect() {
			v.repaintNode(n)
		}
	}
}

func (v *signalGraphView) selectConnect(conn interfaces.ConnectionObject) {
	n1 := v.findNode(conn.FromObject().NodeObject().Name())
	n2 := v.findNode(conn.ToObject().NodeObject().Name())
	for _, c := range v.connections {
		if n1 == c.From() && n2 == c.To() &&
			n1.GetOutPorts()[c.FromId()] == conn.FromObject() &&
			n2.GetInPorts()[c.ToId()] == conn.ToObject() {
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

func (v *signalGraphView) deselectConnects() {
	for _, c := range v.connections {
		if c.Deselect() {
			v.repaintConnection(c)
		}
	}
}

func (v *signalGraphView) Expand(obj interfaces.GraphElement) {
	switch obj.(type) {
	case interfaces.NodeObject:
		n := v.findNode(obj.(interfaces.NodeObject).Name())
		if n != nil {
			n.Expand()
		}
	default:
	}
}

func (v *signalGraphView) Collapse(obj interfaces.GraphElement) {
	switch obj.(type) {
	case interfaces.NodeObject:
		n := v.findNode(obj.(interfaces.NodeObject).Name())
		if n != nil {
			n.Collapse()
		}
	default:
	}
}

//
//		ScaledScene interface
//

func (v *signalGraphView) Update() (width, height int) {
	width = v.calcSceneWidth()
	height = v.calcSceneHeight()
	v.area.SetSizeRequest(width, height)
	v.area.QueueDrawArea(0, 0, width, height)
	return
}

func (v *signalGraphView) calcSceneWidth() int {
	return int(float64(v.bBox().Max.X+50) * v.parent.Scale())
}

func (v *signalGraphView) calcSceneHeight() int {
	return int(float64(v.bBox().Max.Y+50) * v.parent.Scale())
}

//
//		areaButtonCallback
//

func (v *signalGraphView) ButtonCallback(area DrawArea, evType gdk.EventType, position image.Point) {
	pos := v.parent.Position(position)
	switch evType {
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

func (v *signalGraphView) handleNodeSelect(pos image.Point) {
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

func (v *signalGraphView) handleConnectSelect(pos image.Point) {
	for _, c := range v.connections {
		hit, _ := c.CheckHit(pos)
		if hit {
			c.Select()
			n1 := c.From()
			p1 := n1.GetOutPorts()[c.FromId()]
			n2 := c.To()
			p2 := n2.GetInPorts()[c.ToId()]
			v.context.SelectConnect(p1.ConnectionObject(p2))
			v.repaintConnection(c)
		} else {
			if c.Deselect() {
				v.repaintConnection(c)
			}
		}
	}
}

//
//		areaMotionCallback
//

func (v *signalGraphView) MotionCallback(area DrawArea, position image.Point) {
	pos := v.parent.Position(position)
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

func (v *signalGraphView) handleDrag(pos image.Point) {
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

func (v *signalGraphView) handleMouseover(pos image.Point) {
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

//
//		areaDrawCallback
//

func (v *signalGraphView) DrawCallback(area DrawArea, context *cairo.Context) {
	context.Scale(v.parent.Scale(), v.parent.Scale())
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	v.drawNodes(context, r)
	v.drawConnections(context, r)
}

func (v *signalGraphView) drawConnections(context *cairo.Context, r image.Rectangle) {
	for _, c := range v.connections {
		box := c.BBox()
		if r.Overlaps(box) {
			c.Draw(context)
		}
	}
}

func (v *signalGraphView) drawNodes(context *cairo.Context, r image.Rectangle) {
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

//
//		Private functions
//

func (v *signalGraphView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *signalGraphView) drawAll() {
	v.drawScene(v.bBox())
}

func (v *signalGraphView) bBox() (box image.Rectangle) {
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

func (v *signalGraphView) drawScene(r image.Rectangle) {
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.parent.ScaleCoord(x, false), v.parent.ScaleCoord(y, false),
		v.parent.ScaleCoord(w, true)+1, v.parent.ScaleCoord(h, true)+1)
	return
}

func (v *signalGraphView) nodeDrawBox(n graph.NodeIf) image.Rectangle {
	ret := n.BBox()
	for _, c := range v.connections {
		if c.IsLinked(n.Name()) {
			ret = ret.Union(c.BBox())
		}
	}
	return ret
}

func (v *signalGraphView) repaintNode(n graph.NodeIf) {
	v.drawScene(v.nodeDrawBox(n))
}

func (v *signalGraphView) repaintConnection(c graph.ConnectIf) {
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
