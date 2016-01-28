package views

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image"
	"log"
)

type mappingView struct {
	parent      *ScaledView
	area        DrawArea
	mapping     mp.MappingIf
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	arch        []graph.ArchIf
	context     ContextIf
	unmapped    graph.ProcessIf
	unmappedObj pf.ProcessIf

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*mappingView)(nil)
var _ GraphViewIf = (*mappingView)(nil)

func MappingViewNew(m mp.MappingIf, context ContextIf) (viewer *mappingView, err error) {
	viewer = &mappingView{nil, DrawArea{}, m, nil, nil, nil, context, nil, &unmappedProcess{}, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func MappingViewDestroy(viewer GraphViewIf) {
	DrawAreaDestroy(viewer.(*mappingView).area)
}

func (v *mappingView) Widget() *gtk.Widget {
	return v.parent.Widget()
}

func (v *mappingView) init() (err error) {
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

func (v *mappingView) Sync() {
	log.Printf("mappingView.Sync()\n")
	g := v.mapping.Graph()
	v.nodes = make([]graph.NodeIf, len(g.Nodes()))
	var numberOfConnections = 0
	for _, n := range g.Nodes() {
		for _, p := range n.OutPorts() {
			numberOfConnections += len(p.Connections())
		}
	}
	v.connections = make([]graph.ConnectIf, numberOfConnections)
	for i, n := range g.Nodes() {
		v.nodes[i] = graph.NodeNew(n.PathModePosition("", gr.PositionModeMapping), n, "")
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
	p := v.mapping.Platform()
	v.arch = make([]graph.ArchIf, len(p.Arch()))
	for i, a := range p.Arch() {
		v.arch[i] = graph.ArchMappingNew(a, v.nodes, v.mapping)
	}
	var unmappedNodes []*graph.Node
	for _, n := range v.nodes {
		_, ok := v.mapping.Mapped(n.UserObj())
		if !ok {
			unmappedNodes = append(unmappedNodes, n.(*graph.Node))
		}
	}
	v.unmapped = graph.ProcessMappingNew(unmappedNodes, v.unmappedObj)
	// TODO: add v.unmapped to scene
	v.area.SetSizeRequest(v.calcSceneWidth(), v.calcSceneHeight())
	v.drawAll()
}

func (v mappingView) IdentifyGraph(g bh.SignalGraphIf) bool {
	return false
}

func (v mappingView) IdentifyPlatform(p pf.PlatformIf) bool {
	return false
}

func (v mappingView) IdentifyMapping(m mp.MappingIf) bool {
	return v.mapping == m
}

//
//		Handle selection in treeview
//

func (v *mappingView) Select(obj interface{}) {
	switch obj.(type) {
	case pf.ArchIf:
		arch := obj.(pf.ArchIf)
		a, ok := v.focusArchFromUserObject(arch)
		if ok {
			v.selectArch(a)
		}
	case pf.ProcessIf:
		pr := obj.(pf.ProcessIf)
		a, ok := v.focusArchFromUserObject(pr.Arch())
		if ok {
			a.SelectProcess(pr)
			v.repaintArch(a)
		}
	case pf.ChannelIf:
		ch := obj.(pf.ChannelIf)
		a, ok := v.focusArchFromUserObject(ch.Process().Arch())
		if ok {
			a.SelectChannel(ch)
			v.repaintArch(a)
		}
	case bh.NodeIf:
		n := obj.(bh.NodeIf)
		melem, ok := v.mapping.MappedElement(n)
		if !ok {
			// todo: unmapped nodes
			return
		}
		p, ok := melem.Process()
		if !ok {
			return
		}
		var a graph.ArchIf
		a, ok = v.focusArchFromUserObject(p.Arch())
		if !ok {
			log.Printf("mappingView.Select error: arch %s not found\n", a.Name())
			return
		}
		pr := a.SelectProcess(p)
		if pr == nil {
			log.Printf("mappingView.Select error: process %v not found\n", p)
			return
		}
		pr.(*graph.ProcessMapping).SelectNode(n)
		v.repaintArch(a)
	case mp.MappedElementIf:
		melem := obj.(mp.MappedElementIf)
		p, isMapped := melem.Process()
		var pr graph.ProcessIf
		var a graph.ArchIf
		if !isMapped {
			pr = v.unmapped
		} else {
			var ok bool
			a, ok = v.focusArchFromUserObject(p.Arch())
			if !ok {
				log.Printf("mappingView.Select error: arch %s not found\n", a.Name())
				return
			}
			pr = a.SelectProcess(p)
			if pr == nil {
				log.Printf("mappingView.Select error: process %v not found\n", p)
				return
			}
		}
		n := melem.Node()
		if n != nil {
			pr.(*graph.ProcessMapping).SelectNode(n)
		}
		if isMapped {
			v.repaintArch(a)
		} else {
			v.repaintUnmapped(v.unmapped)
		}
	default:
	}
}

func (v *mappingView) selectArch(toSelect graph.ArchIf) {
	var a graph.ArchIf
	for _, a = range v.arch {
		if a.Name() == toSelect.Name() {
			if !a.Select() {
				v.repaintArch(a)
			}
			return
		}
	}
}

func (v *mappingView) focusArchFromUserObject(obj pf.ArchIf) (ret graph.ArchIf, ok bool) {
	var a graph.ArchIf
	for _, a = range v.arch {
		if obj.Name() == a.UserObj().Name() {
			ret = a
			ok = true
		} else {
			if a.Deselect() {
				v.repaintArch(a)
			}
		}
	}
	return
}

func (v *mappingView) Expand(obj interface{}) {
}

func (v *mappingView) Collapse(obj interface{}) {
}

//
//		ScaledScene interface
//

func (v *mappingView) Update() (width, height int) {
	width = v.calcSceneWidth()
	height = v.calcSceneHeight()
	v.area.SetSizeRequest(width, height)
	v.area.QueueDrawArea(0, 0, width, height)
	return
}

func (v *mappingView) calcSceneWidth() int {
	return int(float64(v.bBox().Max.X+50) * v.parent.Scale())
}

func (v *mappingView) calcSceneHeight() int {
	return int(float64(v.bBox().Max.Y+50) * v.parent.Scale())
}

//
//		platformButtonCallback
//

func (v *mappingView) ButtonCallback(area DrawArea, evType gdk.EventType, position image.Point) {
	pos := v.parent.Position(position)
	switch evType {
	case gdk.EVENT_BUTTON_PRESS:
		v.button1Pressed = true
		v.dragOffs = pos
		v.handleArchSelect(pos)
	case gdk.EVENT_2BUTTON_PRESS:
		log.Println("areaButtonCallback 2BUTTON_PRESS")

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

func (v *mappingView) handleArchSelect(pos image.Point) {
	for _, a := range v.arch {
		hit, _ := a.CheckHit(pos)
		if !hit {
			if a.Deselect() {
				v.repaintArch(a)
			}
			continue
		}
		if !a.Select() {
			v.repaintArch(a)
		}
		v.context.SelectArch(a.UserObj())
		var ok bool
		var pr pf.ProcessIf
		var ch pf.ChannelIf
		var n graph.NodeIf
		var p graph.ProcessIf
		ok, pr, p = a.GetSelectedProcess()
		if !ok {
			continue
		}
		//log.Printf("mappingView.handleArchSelect: select process %v\n", pr)
		v.context.SelectProcess(pr)
		ok, ch = a.GetSelectedChannel()
		if ok {
			//log.Printf("mappingView.handleArchSelect: select channel %v\n", ch)
			v.context.SelectChannel(ch)
			continue
		}
		ok, n = p.(*graph.ProcessMapping).GetSelectedNode()
		if !ok {
			continue
		}
		var melem mp.MappedElementIf
		melem, ok = v.mapping.MappedElement(n.UserObj())
		if !ok {
			// todo: unmapped node
		} else {
			v.context.SelectMapElement(melem)
		}
	}
	hit, _ := v.unmapped.CheckHit(pos)
	if !hit {
		if v.unmapped.Deselect() {
			v.repaintUnmapped(v.unmapped)
		}
		return
	}
	if !v.unmapped.Select() {
		v.repaintUnmapped(v.unmapped)
	}
	var ok bool
	var n graph.NodeIf
	//log.Printf("mappingView.handleArchSelect: select process %v\n", pr)
	ok, n = v.unmapped.(*graph.ProcessMapping).GetSelectedNode()
	if !ok {
		return
	}
	var melem mp.MappedElementIf
	melem, ok = v.mapping.MappedElement(n.UserObj())
	if !ok {
		// todo: unmapped node
	} else {
		v.context.SelectMapElement(melem)
	}
}

//
//		platformMotionCallback
//

func (v *mappingView) MotionCallback(area DrawArea, position image.Point) {
	pos := v.parent.Position(position)
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

func (v *mappingView) handleDrag(pos image.Point) {
	for _, a := range v.arch {
		if a.IsSelected() {
			box := a.BBox()
			v.repaintArch(a)
			box = box.Add(pos.Sub(v.dragOffs))
			v.dragOffs = pos
			a.SetPosition(box.Min)
			v.repaintArch(a)
		}
	}
	if v.unmapped.IsSelected() {
		box := v.unmapped.BBox()
		v.repaintUnmapped(v.unmapped)
		box = box.Add(pos.Sub(v.dragOffs))
		v.dragOffs = pos
		v.unmapped.SetPosition(box.Min)
		v.repaintUnmapped(v.unmapped)
	}
}

func (v *mappingView) handleMouseover(pos image.Point) {
	for _, a := range v.arch {
		_, mod := a.CheckHit(pos)
		if mod {
			v.repaintArch(a)
		}
	}
	_, mod := v.unmapped.CheckHit(pos)
	if mod {
		v.repaintUnmapped(v.unmapped)
	}
	/*
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
	*/
}

func (v *mappingView) repaintArch(a graph.ArchIf) {
	v.drawScene(v.archDrawBox(a))
}

func (v *mappingView) repaintUnmapped(p graph.ProcessIf) {
	v.drawScene(v.unmappedDrawBox())
}

func (v *mappingView) archDrawBox(a graph.ArchIf) image.Rectangle {
	ret := a.BBox()
	for _, aa := range v.arch {
		if aa.IsLinked(a.Name()) {
			ret = ret.Union(aa.BBox())
		}
	}
	return ret
}

func (v *mappingView) unmappedDrawBox() image.Rectangle {
	ret := v.unmapped.BBox()
	for _, c := range v.connections {
		if v.connectionLinksUnmapped(c) {
			ret = ret.Union(c.BBox())
		}
	}
	return ret
}

func (v *mappingView) connectionLinksUnmapped(c graph.ConnectIf) bool {
	for _, ch := range v.unmapped.(*graph.ProcessMapping).Children {
		name := ch.(graph.NamedBox).Name()
		if name == c.From().Name() || name == c.To().Name() {
			return true
		}
	}
	return false
}

//
//		platformDrawCallback
//

func (v *mappingView) DrawCallback(area DrawArea, context *cairo.Context) {
	context.Scale(v.parent.Scale(), v.parent.Scale())
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	if r.Overlaps(v.unmapped.BBox()) {
		v.unmapped.Draw(context)
	}
	v.drawArch(context, r)
	v.drawChannels(context, r)
	v.drawConnections(context, r)
}

func (v *mappingView) drawConnections(context *cairo.Context, r image.Rectangle) {
	for _, c := range v.connections {
		box := c.BBox()
		if r.Overlaps(box) {
			c.Draw(context)
		}
	}
}

func (v *mappingView) drawArch(context *cairo.Context, r image.Rectangle) {
	for _, a := range v.arch {
		if r.Overlaps(a.BBox()) {
			a.Draw(context)
		}
	}
}

func (v *mappingView) drawChannels(context *cairo.Context, r image.Rectangle) {
	for _, a := range v.arch {
		for _, pr := range a.Processes() {
			for _, c := range pr.UserObj().OutChannels() {
				link := c.Link()
				lpr := link.Process()
				la := lpr.Arch()
				if la.Name() != a.Name() {
					var a2 graph.ArchIf
					for _, a2 = range v.arch {
						if a2.UserObj() == la {
							break
						}
					}
					p1 := a.ChannelPort(c)
					p2 := a2.ChannelPort(link)
					if p1 == nil || p2 == nil {
						log.Printf("mappingView.drawChannels error: invalid nil port (%s - %s).\n", a.Name(), la.Name())
						continue
					}
					r, g, b := graph.ColorOption(graph.NormalLine)
					context.SetLineWidth(2)
					context.SetSourceRGB(r, g, b)
					pos1 := p1.Position().Add(image.Point{5, 5})
					pos2 := p2.Position().Add(image.Point{5, 5})
					graph.DrawLine(context, pos1, pos2)
				}
			}
		}
	}
}

//
//		Private methods
//

func (v *mappingView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *mappingView) drawAll() {
	r := image.Rect(0, 0, v.calcSceneWidth(), v.calcSceneHeight())
	v.drawScene(r)
}

func (v *mappingView) bBox() (box image.Rectangle) {
	emptyRect := image.Rectangle{}
	box = emptyRect
	for _, a := range v.arch {
		if box == emptyRect {
			box = a.BBox()
		} else {
			box = box.Union(a.BBox())
		}
	}
	box = box.Union(v.unmapped.BBox())
	return
}

func (v *mappingView) drawScene(r image.Rectangle) {
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.parent.ScaleCoord(x, false), v.parent.ScaleCoord(y, false),
		v.parent.ScaleCoord(w, true)+1, v.parent.ScaleCoord(h, true)+1)
	return
}

type unmappedProcess struct {
	pos image.Point
}

var _ pf.ProcessIf = (*unmappedProcess)(nil)

func (p *unmappedProcess) Name() string {
	return "unmapped"
}

func (p *unmappedProcess) SetName(string) {
}

func (p *unmappedProcess) ModePosition(gr.PositionMode) image.Point {
	return p.pos
}

func (p *unmappedProcess) SetModePosition(m gr.PositionMode, pos image.Point) {
	p.pos = pos
}

func (p *unmappedProcess) Arch() pf.ArchIf {
	return nil
}

func (p *unmappedProcess) InChannels() []pf.ChannelIf {
	return nil
}

func (p *unmappedProcess) OutChannels() []pf.ChannelIf {
	return nil
}

func (p *unmappedProcess) AddToTree(tr.TreeIf, tr.Cursor) {
	return
}
func (p *unmappedProcess) AddNewObject(tr.TreeIf, tr.Cursor, tr.TreeElement) (c tr.Cursor, e error) {
	return
}
func (p *unmappedProcess) RemoveObject(tr.TreeIf, tr.Cursor) (r []tr.IdWithObject) {
	return
}

func (p *unmappedProcess) CreateXml() (buf []byte, err error) {
	return
}
