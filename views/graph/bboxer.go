package graph

import (
    //"log"
    "image"
)

//
//	Default BBox implementation
//

type BBoxObject struct {
    box image.Rectangle
}

func BBoxInit(box image.Rectangle) BBoxObject{
    return BBoxObject{box}
}

func (b BBoxObject) BBox() image.Rectangle {
    return b.box
}


//
//	Positioner implementation
//

func (b BBoxObject) Position() image.Point {
    return b.box.Min
}

func (b *BBoxObject) BBoxDefaultSetPosition(pos image.Point) {
    shift := pos.Sub(b.Position())
    b.box = b.box.Add(shift)
}

func (b *BBoxObject) SetPosition(pos image.Point) {
    b.BBoxDefaultSetPosition(pos)
}


//
//	Shaper implementation
//

func (b BBoxObject) Shape() image.Point {
    return b.box.Max.Sub(b.box.Min)
}

func (b *BBoxObject) SetShape(newSize image.Point) {
    b.box.Max = b.box.Min.Add(newSize)
}


//
//      Helper functions
//

func size(b BBoxer) (w, h float64) {
	box := b.BBox()
	r := box.Size()
	w = float64(r.X)
	h = float64(r.Y)
	return
}

func boxToDraw(b BBoxer, pad image.Point) (x, y, w, h float64) {
	dx, dy := pad.X, pad.Y
	w0, h0 := size(b)
	a := b.BBox().Min
	x = float64(a.X) + float64(dx)
	y = float64(a.Y) + float64(dy)
	w = w0 - float64(2*dx)
	h = h0 - float64(2*dy)
	return
}



