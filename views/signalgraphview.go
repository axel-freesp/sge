package views

import (
	freesp "github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
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
	sgType      bh.SignalGraphTypeIf
	graphId     string
	context     ContextIf

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*signalGraphView)(nil)
var _ GraphViewIf = (*signalGraphView)(nil)

func SignalGraphViewNew(g bh.SignalGraphIf, context ContextIf) (viewer *signalGraphView, err error) {
	viewer = &signalGraphView{nil, DrawArea{}, nil, nil, g.ItsType(), g.Filename(), context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func SignalGraphViewNewFromType(g bh.SignalGraphTypeIf, context ContextIf) (viewer *signalGraphView, err error) {
	viewer = &signalGraphView{nil, DrawArea{}, nil, nil, g, "", context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func SignalGraphViewDestroy(viewer GraphViewIf) {
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
	var selected bh.NodeIdIf
	var wasSelected bool
	if v.nodes != nil {
		selected, wasSelected = v.getSelectedNode()
	}
	v.nodes = make([]graph.NodeIf, len(g.Nodes()))

	var numberOfConnections = 0
	for _, n := range g.Nodes() {
		for _, p := range n.OutPorts() {
			numberOfConnections += len(p.Connections())
		}
	}
	v.connections = make([]graph.ConnectIf, numberOfConnections)

	getPositioner := func(nId bh.NodeIdIf) gr.ModePositioner {
		n, ok := g.NodeByPath(nId.String())
		if !ok {
			log.Panicf("getPositioner: could not find node %v\n", nId)
		}
		log.Printf("getPositioner(%v): node=%s\n", nId, n.Name())
		proxy := gr.PathModePositionerProxyNew(n)
		proxy.SetActivePath(nId.Parent().String())
		return proxy
	}
	for i, n := range g.Nodes() {
		nId := freesp.NodeIdFromString(n.Name(), v.graphId)
		if n.Expanded() {
			n.SetActiveMode(gr.PositionModeExpanded)
			v.nodes[i] = graph.ExpandedNodeNew(getPositioner, n, nId)
		} else {
			n.SetActiveMode(gr.PositionModeNormal)
			v.nodes[i] = graph.NodeNew(getPositioner, n, nId)
		}
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
	if wasSelected {
		log.Printf("signalGraphView.Sync: select %v after sync\n", selected)
		v.selectNode(selected)
	}
}

func (v signalGraphView) IdentifyGraph(g bh.SignalGraphIf) bool {
	return g.ItsType() == v.sgType
}

func (v signalGraphView) IdentifyPlatform(g pf.PlatformIf) bool {
	return false
}

func (v signalGraphView) IdentifyMapping(g mp.MappingIf) bool {
	return false
}

//
//		Handle selection in treeview
//

func (v *signalGraphView) Select(obj interface{}) {
	switch obj.(type) {
	case bh.ConnectionIf:
		v.deselectNodes()
		v.selectConnect(obj.(bh.ConnectionIf))
	}
}

func (v *signalGraphView) Select2(obj interface{}, id string) {
	switch obj.(type) {
	case bh.NodeIf:
		v.deselectConnects()
		v.selectNode(freesp.NodeIdFromString(id, v.graphId))
	case bh.PortIf:
		v.deselectConnects()
		n := v.selectNode(freesp.NodeIdFromString(id, v.graphId))
		if n != nil {
			n.SelectPort(obj.(bh.PortIf))
			v.repaintNode(n)
		}
	}
}

func (v *signalGraphView) selectNode(id bh.NodeIdIf) (node graph.NodeIf) {
	var n graph.NodeIf
	for _, n = range v.nodes {
		ok, nd := n.SelectNode(freesp.NodeIdFromString(n.Name(), v.graphId), id)
		if ok {
			v.repaintNode(n)
		}
		if nd != nil {
			node = nd
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

func (v *signalGraphView) selectConnect(conn bh.ConnectionIf) {
	n1 := v.findNode(conn.From().Node().Name())
	n2 := v.findNode(conn.To().Node().Name())
	for _, c := range v.connections {
		if n1 == c.From() && n2 == c.To() &&
			n1.OutPorts()[c.FromId()] == conn.From() &&
			n2.InPorts()[c.ToId()] == conn.To() {
			if c.Select() {
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

func (v *signalGraphView) Expand(obj interface{}) {
	switch obj.(type) {
	case bh.NodeIf:
		log.Printf("signalGraphView.Expand\n")
		if obj.(bh.NodeIf).Expanded() {
			log.Printf("signalGraphView.Expand: nothing to do\n")
			return
		}
		v.drawAll()
		obj.(bh.NodeIf).SetExpanded(true)
		v.Sync()
	default:
	}
}

func (v *signalGraphView) Collapse(obj interface{}) {
	switch obj.(type) {
	case bh.NodeIf:
		log.Printf("signalGraphView.Collapse\n")
		if !obj.(bh.NodeIf).Expanded() {
			log.Printf("signalGraphView.Collapse: nothing to do\n")
			return
		}
		v.drawAll()
		obj.(bh.NodeIf).SetExpanded(false)
		v.Sync()
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
			selected, ok := n.GetHighlightedNode(freesp.NodeIdFromString(n.Name(), v.graphId))
			if ok {
				v.context.EditNode(selected)
				break
			}
		}

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

func (v signalGraphView) getSelectedNode() (id bh.NodeIdIf, ok bool) {
	for _, n := range v.nodes {
		id, ok = n.GetSelectedNode(freesp.NodeIdFromString(n.Name(), v.graphId))
		if ok {
			break
		}
	}
	return
}

func (v *signalGraphView) handleNodeSelect(pos image.Point) {
	for _, n := range v.nodes {
		hit, _ := n.CheckHit(pos)
		if hit {
			nodeId := freesp.NodeIdFromString(n.Name(), v.graphId)
			selected, ok := n.GetHighlightedNode(nodeId)
			if !ok {
				log.Printf("signalGraphView.handleNodeSelect: FIXME: hit=true, GetSelectedNode() is not ok\n")
				return
			}
			selected.SetFilename(v.graphId)
			v.context.SelectNode(selected)
			port, ok := n.GetSelectedPort(nodeId)
			if ok {
				v.context.SelectPort(port, selected)
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
			p1 := n1.OutPorts()[c.FromId()]
			n2 := c.To()
			p2 := n2.InPorts()[c.ToId()]
			v.context.SelectConnect(p1.Connection(p2))
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
			//if !overlaps(v.nodes, box.Min.Add(pos.Sub(v.dragOffs))) {
			v.repaintNode(n)
			box = box.Add(pos.Sub(v.dragOffs))
			v.dragOffs = pos
			n.SetPosition(box.Min)
			v.repaintNode(n)
			//}
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

func (v *signalGraphView) findNode(name string) graph.NodeIf {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d
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
