package graph

import (
    "math"
    "image"
    "github.com/gotk3/gotk3/cairo"
)

const s = 16.0

type Connection struct {
    From, To *Node
    FromPort, ToPort int
}

func ConnectionNew(From, To Dragable, FromPort, ToPort int) Connection{
    return Connection{From.(*Node), To.(*Node), FromPort, ToPort}
}

func (c Connection) BBox(p1, p2 image.Point) image.Rectangle {
    r := image.Rectangle{p1, p2}.Canon()
    k := int(math.Ceil(s/math.Sqrt(1.5)))
    r.Min.X -= k
    r.Max.X += k
    r.Min.Y -= k
    r.Max.Y += k
    return r
}

func drawLine(gc *cairo.Context, p1, p2 image.Point) {
    gc.MoveTo(float64(p1.X), float64(p1.Y))
    gc.LineTo(float64(p2.X), float64(p2.Y))
    gc.Stroke()
}

func drawAccLine(gc *cairo.Context, x1, y1 float64, p2 image.Point) {
    gc.MoveTo(x1, y1)
    gc.LineTo(float64(p2.X), float64(p2.Y))
    gc.Stroke()
}

func DrawArrow(gc *cairo.Context, p1, p2 image.Point){
    x1 := float64(p2.X)
    y1 := float64(p2.Y)
    dx := x1 - float64(p1.X)
    dy := y1 - float64(p1.Y)
    r := math.Sqrt(dx*dx + dy*dy)
    if r > 0 {
        h0x, h0y := x1 - s*dx / r, y1 - s*dy / r
        h1x, h1y := x1 - s*dy / r, y1 + s*dx / r
        h2x, h2y := x1 + s*dy / r, y1 - s*dx / r
        a1x, a1y := (2.0*h0x + h1x) / 3.0, (2.0*h0y + h1y) / 3.0
        a2x, a2y := (2.0*h0x + h2x) / 3.0, (2.0*h0y + h2y) / 3.0
        drawLine(gc, p1, p2)
        drawAccLine(gc, a1x, a1y, p2)
        drawAccLine(gc, a2x, a2y, p2)
    }
}

