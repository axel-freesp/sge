package graph

import (
    "image"
)

type SelectObject struct {
	selected bool
}

type HighlightObject struct {
	highlighted bool
}

type SelectableBox struct {
	BBoxObject
    SelectObject
    HighlightObject
}

var _ BoxedSelecter = (*SelectableBox)(nil)

func SelectableBoxInit(box image.Rectangle) SelectableBox {
    return SelectableBox{BBoxInit(box), SelectObject{false}, HighlightObject{false}}
}

/*
 *      Default Selecter implementation
 */

var _ Selecter  = (*SelectObject)(nil)

func (s *SelectObject) Select() (selected bool) {
	selected = s.selected
	s.selected = true
	return
}

func (s *SelectObject) Deselect() (selected bool) {
	selected = s.selected
	s.selected = false
	return
}

func (s *SelectObject) IsSelected() bool {
	return s.selected
}

/*
 *      Default Highlighter implementation
 */

var _ Highlighter  = (*HighlightObject)(nil)

func (h *HighlightObject) IsHighlighted() bool {
	return h.highlighted
}


/*
 *      Default SelectableBox implementation
 */

func (b *SelectableBox) CheckHit(p image.Point) (hit, modified bool) {
	test := image.Rectangle{p, p}
	hit = b.BBox().Overlaps(test)
	modified = b.highlighted != hit
	b.highlighted = hit
	return
}

