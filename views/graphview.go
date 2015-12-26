package views

import (
//	"fmt"
	"image"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/cairo"
)

const (
	nodeWidth = 60
	nodeHeight = 80
)

type GraphView struct {
	ScrolledView
	area *gtk.DrawingArea
	nodes []graph.DragableObject
}

func GraphViewNew(graph freesp.SignalGraph, width, height int) (viewer *GraphView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &GraphView{*v, nil, nil}
	err = viewer.init(graph, width, height)
	return
}

func (v *GraphView) init(g freesp.SignalGraph, width, height int) (err error) {
	v.area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
	v.area.SetSizeRequest(width, height)

	for i, n := range g.ItsType().Nodes() {
		pos := n.Position()
		box := image.Rect(pos.X, pos.Y, pos.X + nodeWidth, pos.Y + nodeHeight)
		v.nodes[i].Init(box, n)
	}

	v.area.Connect("draw", areaDrawCallback, v)
	v.scrolled.Add(v.area)
	return
}

func areaDrawCallback(area *gtk.DrawingArea, context cairo.Context, v *GraphView) {
}

