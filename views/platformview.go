package views

import (
	//	"fmt"
	"github.com/axel-freesp/sge/freesp"
	//"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"image"
	"log"
	"math"
)

type platformView struct {
	parent *ScaledView
	area   DrawArea

	p       freesp.Platform
	context Context

	dragOffs       image.Point
	button1Pressed bool
}

var _ ScaledScene = (*platformView)(nil)
var _ GraphView = (*platformView)(nil)

func PlatformViewNew(p freesp.Platform, context Context) (viewer *platformView, err error) {
	viewer = &platformView{nil, DrawArea{}, p, context, image.Point{}, false}
	err = viewer.init()
	if err != nil {
		return
	}
	viewer.Sync()
	return
}

func (v *platformView) Widget() *gtk.Widget {
	return v.parent.Widget()
}

func (v *platformView) init() (err error) {
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

func (v *platformView) Sync() {

	v.area.SetSizeRequest(v.calcSceneWidth(), v.calcSceneHeight())
	//v.drawAll()
}

/*
 *		Handle selection in treeview
 */

func (v *platformView) Select(obj freesp.TreeElement) {
	switch obj.(type) {
	default:
	}
}

/*
 *		ScaledScene interface
 */

func (v *platformView) Update() (width, height int) {
	width = v.calcSceneWidth()
	height = v.calcSceneHeight()
	v.area.SetSizeRequest(width, height)
	v.area.QueueDrawArea(0, 0, width, height)
	return
}

func (v *platformView) calcSceneWidth() int {
	return int(float64(v.bBox().Max.X+50) * v.parent.Scale())
}

func (v *platformView) calcSceneHeight() int {
	return int(float64(v.bBox().Max.Y+50) * v.parent.Scale())
}

/*
 *		platformButtonCallback
 */

func (v *platformView) ButtonCallback(area DrawArea, evType gdk.EventType, position image.Point) {
	pos := v.parent.Position(position)
	switch evType {
	case gdk.EVENT_BUTTON_PRESS:
		v.button1Pressed = true
		v.dragOffs = pos
	case gdk.EVENT_2BUTTON_PRESS:
		log.Println("areaButtonCallback 2BUTTON_PRESS")

	case gdk.EVENT_BUTTON_RELEASE:
		v.button1Pressed = false
	default:
	}
}

/*
 *		platformMotionCallback
 */

func (v *platformView) MotionCallback(area DrawArea, position image.Point) {
	pos := v.parent.Position(position)
	if v.button1Pressed {
		v.handleDrag(pos)
	} else {
		v.handleMouseover(pos)
	}
}

func (v *platformView) handleDrag(pos image.Point) {
	/*
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
	*/
}

func (v *platformView) handleMouseover(pos image.Point) {
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

/*
 *		platformDrawCallback
 */

func (v *platformView) DrawCallback(area DrawArea, context *cairo.Context) {
	context.Scale(v.parent.Scale(), v.parent.Scale())
	x1, y1, x2, y2 := context.ClipExtents()
	r := image.Rect(int(x1), int(y1), int(x2), int(y2))
	//v.drawNodes(context, r)
	//v.drawConnections(context, r)
	_ = r
}

/*
 *		Private functions
 */

func (v *platformView) drawAll() {
	v.drawScene(v.bBox())
}

func (v *platformView) bBox() (box image.Rectangle) {
	emptyRect := image.Rectangle{}
	box = emptyRect
	/*
		for _, o := range v.nodes {
			if box == emptyRect {
				box = o.BBox()
			} else {
				box = box.Union(o.BBox())
			}
		}
	*/
	return
}

func (v *platformView) drawScene(r image.Rectangle) {
	x, y, w, h := r.Min.X, r.Min.Y, r.Size().X, r.Size().Y
	v.area.QueueDrawArea(v.scaleCoord(x, false), v.scaleCoord(y, false), v.scaleCoord(w, true)+1, v.scaleCoord(h, true)+1)
	return
}

func (v *platformView) scaleCoord(x int, roundUp bool) int {
	if roundUp {
		return int(math.Ceil(float64(x) * v.parent.Scale()))
	} else {
		return int(float64(x) * v.parent.Scale())
	}
}
