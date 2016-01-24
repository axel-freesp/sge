package graph

import (
	"github.com/gotk3/gotk3/cairo"
	"image"
)

type SelectObject struct {
	selected             bool
	onSelect, onDeselect OnSelectionFunc
}

type HighlightObject struct {
	highlighted bool
	Pos         image.Point
	onHighlight OnHighlightFunc
}

/*
 *      Default Selecter implementation
 */

var _ Selecter = (*SelectObject)(nil)

func (s *SelectObject) Select() (selected bool) {
	selected = s.selected
	s.selected = true
	s.onSelect()
	return
}

func (s *SelectObject) Deselect() (selected bool) {
	selected = s.selected
	s.selected = false
	s.onDeselect()
	return
}

func (s *SelectObject) IsSelected() bool {
	return s.selected
}

func (s *SelectObject) RegisterOnSelect(onSelect, onDeselect OnSelectionFunc) {
	s.onSelect, s.onDeselect = onSelect, onDeselect
}

var defaultSelection = func() {}

/*
 *      Default Highlighter implementation
 */

var _ Highlighter = (*HighlightObject)(nil)

func (h *HighlightObject) DoHighlight(hit bool, pos image.Point) (modified bool) {
	modified = (h.highlighted != hit)
	h.highlighted = hit
	h.Pos = pos
	if modified || hit {
		modified = h.onHighlight(hit, pos) || modified
	}
	return
}

func (h *HighlightObject) IsHighlighted() bool {
	return h.highlighted
}

func (h *HighlightObject) RegisterOnHighlight(onHighlight OnHighlightFunc) {
	h.onHighlight = onHighlight
}

var defaultHighlight = func(bool, image.Point) bool { return false }

/*
 *      Default SelectableBox implementation
 */

type SelectableBox struct {
	BBoxObject
	SelectObject
	HighlightObject
	config DrawConfig
	onDraw func(ctxt interface{})
}

type Color struct {
	r, g, b float64
}

func ColorInit(r, g, b float64) Color {
	return Color{r, g, b}
}

var _ BoxedSelecter = (*SelectableBox)(nil)

type DrawConfig struct {
	nCol, hCol, sCol, fCol, tCol Color
	pad                          image.Point
}

func SelectableBoxInit(box image.Rectangle, config DrawConfig) SelectableBox {
	return SelectableBox{BBoxInit(box),
		SelectObject{false, defaultSelection, defaultSelection},
		HighlightObject{false, image.Point{}, defaultHighlight},
		config, func(ctxt interface{}) {}}
}

func (b *SelectableBox) CheckHit(pos image.Point) (hit, modified bool) {
	test := image.Rectangle{pos, pos}
	hit = b.BBox().Overlaps(test)
	modified = b.DoHighlight(hit, pos)
	return
}

//
//	Drawer interface
//

func (b SelectableBox) Draw(ctxt interface{}) {
	b.SelectorDefaultDraw(ctxt)
	b.onDraw(ctxt)
}

func (b SelectableBox) SelectorDefaultDraw(ctxt interface{}) {
	switch ctxt.(type) {
	case *cairo.Context:
		context := ctxt.(*cairo.Context)
		var color Color
		if b.IsSelected() {
			color = b.config.sCol
		} else if b.IsHighlighted() {
			color = b.config.hCol
		} else {
			color = b.config.nCol
		}
		x, y, w, h := boxToDraw(&b, b.config.pad)
		context.SetSourceRGB(color.r, color.g, color.b)
		context.Rectangle(x, y, w, h)
		context.FillPreserve()
		context.SetSourceRGB(b.config.fCol.r, b.config.fCol.g, b.config.fCol.b)
		context.SetLineWidth(2)
		context.Stroke()
	}
}

func (b *SelectableBox) RegisterOnDraw(onDraw func(ctxt interface{})) {
	b.onDraw = onDraw
}
