package graph

import (
    "image"
)

type SelectObject struct {
	selected bool
	onSelect, onDeselect OnSelectionFunc
}

type HighlightObject struct {
	highlighted bool
	Pos image.Point
	onHighlight OnHighlightFunc
}

/*
 *      Default Selecter implementation
 */

var _ Selecter  = (*SelectObject)(nil)

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

var defaultSelection = func(){}

/*
 *      Default Highlighter implementation
 */

var _ Highlighter  = (*HighlightObject)(nil)

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

var defaultHighlight = func(bool, image.Point)bool{return false}

/*
 *      Default SelectableBox implementation
 */

type SelectableBox struct {
    BBoxObject
    SelectObject
    HighlightObject
}

var _ BoxedSelecter = (*SelectableBox)(nil)

func SelectableBoxInit(box image.Rectangle) SelectableBox {
    return SelectableBox{BBoxInit(box),
		SelectObject{false, defaultSelection, defaultSelection},
		HighlightObject{false, image.Point{}, defaultHighlight}}
}

func (b *SelectableBox) CheckHit(pos image.Point) (hit, modified bool) {
	test := image.Rectangle{pos, pos}
	hit = b.BBox().Overlaps(test)
	modified = b.DoHighlight(hit, pos)
	return
}

