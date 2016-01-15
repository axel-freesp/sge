package views

import (
	interfaces "github.com/axel-freesp/sge/interface"
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
	mapping     interfaces.MappingObject
	nodes       []graph.NodeIf
	connections []graph.ConnectIf
	arch        []graph.ArchIf
	context     Context

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*mappingView)(nil)
var _ GraphView = (*mappingView)(nil)

func MappingViewNew(m interfaces.MappingObject, context Context) (viewer *mappingView, err error) {
	viewer = &mappingView{nil, DrawArea{}, m, nil, nil, nil, context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func MappingViewDestroy(viewer GraphView) {
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
	g := v.mapping.GraphObject()
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
	p := v.mapping.PlatformObject()
	v.arch = make([]graph.ArchIf, len(p.ArchObjects()))
	for i, a := range p.ArchObjects() {
		v.arch[i] = graph.ArchMappingNew(a, v.nodes, v.mapping)
	}
	v.area.SetSizeRequest(v.calcSceneWidth(), v.calcSceneHeight())
	v.drawAll()
}

func (v mappingView) IdentifyGraph(g interfaces.GraphObject) bool {
	return false
}

func (v mappingView) IdentifyPlatform(p interfaces.PlatformObject) bool {
	return false
}

func (v mappingView) IdentifyMapping(p interfaces.MappingObject) bool {
	return false
}

/*
 *		Handle selection in treeview
 */

func (v *mappingView) Select(obj interfaces.GraphElement) {
	switch obj.(type) {
	case interfaces.ArchObject:
		arch := obj.(interfaces.ArchObject)
		a, ok := v.focusArchFromUserObject(arch)
		if ok {
			v.selectArch(a)
		}
	case interfaces.ProcessObject:
		pr := obj.(interfaces.ProcessObject)
		a, ok := v.focusArchFromUserObject(pr.ArchObject())
		if ok {
			a.SelectProcess(pr)
			v.repaintArch(a)
		}
	case interfaces.ChannelObject:
		ch := obj.(interfaces.ChannelObject)
		a, ok := v.focusArchFromUserObject(ch.ProcessObject().ArchObject())
		if ok {
			a.SelectChannel(ch)
			v.repaintArch(a)
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
		}
	}
}

func (v *mappingView) focusArchFromUserObject(obj interfaces.ArchObject) (ret graph.ArchIf, ok bool) {
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

/*
 *		ScaledScene interface
 */

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

/*
 *		platformButtonCallback
 */

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
		if hit {
			if !a.Select() {
				v.repaintArch(a)
			}
			v.context.SelectArch(a.UserObj())
			var ok bool
			var pr interfaces.ProcessObject
			var ch interfaces.ChannelObject
			ok, pr = a.GetSelectedProcess()
			if ok {
				v.context.SelectProcess(pr)
			}
			ok, ch = a.GetSelectedChannel()
			if ok {
				v.context.SelectChannel(ch)
			}
		} else {
			if a.Deselect() {
				v.repaintArch(a)
			}
		}
	}
}

/*
 *		platformMotionCallback
 */

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
}

func (v *mappingView) handleMouseover(pos image.Point) {
	for _, a := range v.arch {
		_, mod := a.CheckHit(pos)
		if mod {
			v.repaintArch(a)
		}
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

func (v *mappingView) archDrawBox(a graph.ArchIf) image.Rectangle {
	ret := a.BBox()
	for _, aa := range v.arch {
		if aa.IsLinked(a.Name()) {
			ret = ret.Union(aa.BBox())
		}
	}
	return ret
}

/*
 *		platformDrawCallback
 */

func (v *mappingView) DrawCallback(area DrawArea, context *cairo.Context) {
	context.Scale(v.parent.Scale(), v.parent.Scale())
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	v.drawArch(context, r)
	v.drawChannels(context, r)
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
			for _, c := range pr.UserObj().OutChannelObjects() {
				link := c.LinkObject()
				lpr := link.ProcessObject()
				la := lpr.ArchObject()
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

/*
 *		Private functions
 */

func (v *mappingView) findNode(name string) *graph.Node {
	for _, d := range v.nodes {
		if d.Name() == name {
			return d.(*graph.Node)
		}
	}
	return nil
}

func (v *mappingView) drawAll() {
	v.drawScene(v.bBox())
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
	return
}

func (v *mappingView) drawScene(r image.Rectangle) {
	//log.Printf("mappingView.drawScene %v\n", r)
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.parent.ScaleCoord(x, false), v.parent.ScaleCoord(y, false),
		v.parent.ScaleCoord(w, true)+1, v.parent.ScaleCoord(h, true)+1)
	return
}
