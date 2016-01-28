package graph

import (
	"github.com/gotk3/gotk3/cairo"
	"image"
	"math"
)

//
//	Line based BBox implementation
//

type LineObject struct {
	BBoxObject
	SelectObject
	HighlightObject
	p1, p2 image.Point
}

var _ BBoxer = (*LineObject)(nil)
var _ Selecter = (*LineObject)(nil)
var _ Highlighter = (*LineObject)(nil)

func LineObjectInit(p1, p2 image.Point) LineObject {
	return LineObject{BBoxObject{image.Rectangle{p1, p2}.Canon()},
		SelectObject{false, defaultSelection, defaultSelection},
		HighlightObject{false, image.Point{}, defaultHighlight},
		p1, p2}
}

func (l *LineObject) CheckHit(pos image.Point) (hit, modified bool) {
	return l.LineDefaultCheckHit(pos)
}

func (l *LineObject) LineDefaultCheckHit(pos image.Point) (hit, modified bool) {
	f, r := l.Transformation()
	p := f(pos)
	px := math.Abs(float64(p.X))
	py := math.Abs(float64(p.Y))
	hit = (px <= r && py <= 3.0)
	modified = l.DoHighlight(hit, pos)
	return
}

//
//      Helper functions
//

func (l LineObject) Transformation() (func(image.Point) image.Point, float64) {
	a, b := l.p1, l.p2
	m := a.Add(b).Div(2)
	am := m.Sub(a)
	r := math.Sqrt(float64(am.X*am.X) + float64(am.Y*am.Y))
	cos := float64(m.X-a.X) / r
	sin := float64(a.Y-m.Y) / r
	return func(x image.Point) (y image.Point) {
		p1 := x.Sub(m)
		a := int(cos*float64(p1.X) - sin*float64(p1.Y))
		b := int(sin*float64(p1.X) + cos*float64(p1.Y))
		y = image.Point{a, b}
		return
	}, r
}

func DrawLine(gc *cairo.Context, p1, p2 image.Point) {
	gc.MoveTo(float64(p1.X), float64(p1.Y))
	gc.LineTo(float64(p2.X), float64(p2.Y))
	gc.Stroke()
}

func DrawAccLine(gc *cairo.Context, x1, y1 float64, p2 image.Point) {
	gc.MoveTo(x1, y1)
	gc.LineTo(float64(p2.X), float64(p2.Y))
	gc.Stroke()
}

func DrawArrow(gc *cairo.Context, p1, p2 image.Point) {
	x1 := float64(p2.X)
	y1 := float64(p2.Y)
	dx := x1 - float64(p1.X)
	dy := y1 - float64(p1.Y)
	r := math.Sqrt(dx*dx + dy*dy)
	if r > 0 {
		h0x, h0y := x1-arrowsize*dx/r, y1-arrowsize*dy/r
		h1x, h1y := x1-arrowsize*dy/r, y1+arrowsize*dx/r
		h2x, h2y := x1+arrowsize*dy/r, y1-arrowsize*dx/r
		a1x, a1y := (2.0*h0x+h1x)/3.0, (2.0*h0y+h1y)/3.0
		a2x, a2y := (2.0*h0x+h2x)/3.0, (2.0*h0y+h2y)/3.0
		DrawLine(gc, p1, p2)
		DrawAccLine(gc, a1x, a1y, p2)
		DrawAccLine(gc, a2x, a2y, p2)
	}
}
