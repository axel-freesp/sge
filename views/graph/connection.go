package graph

import (
    //"fmt"
    "math"
    "image"
    "github.com/gotk3/gotk3/cairo"
)

const arrowsize = 12.0

var (
    conOut = image.Point{global.nodeWidth - global.padX - global.portX0, global.portY0 + global.padY + global.portH/2}
    conIn  = image.Point{global.padX + global.portX0,                    global.portY0 + global.padY + global.portH/2}
    portDY = image.Point{0, global.portDY}
)

type Connection struct {
    SelectableObject
    from, to NodeIf
    fromPort, toPort int
}

func ConnectionNew(from, to NodeIf, fromPort, toPort int) (ret *Connection) {
    ret = &Connection{SelectableObject{image.Rectangle{}, false, false}, from, to, fromPort, toPort}
    ret.box.Min, ret.box.Max = ret.connectionPoints()
    return ret
}

func (c *Connection) IsLinked(nodeName string) bool {
    return c.from.Name() == nodeName || c.to.Name() == nodeName
}

func (c *Connection) From() NodeIf {
    return c.from
}

func (c *Connection) To() NodeIf {
    return c.to
}

func (c *Connection) FromId() int {
    return c.fromPort
}

func (c *Connection) ToId() int {
    return c.toPort
}

func (c *Connection) connectionPoints() (p1, p2 image.Point) {
    p1 = c.from.BBox().Min.Add(conOut.Add(portDY.Mul(c.fromPort)))
    p2 = c.to.BBox().Min.Add(conIn.Add(portDY.Mul(c.toPort)))
    return
}


var _ SelectableBox = (*Connection)(nil)

/*
 *	Drawer interface
 */

func (c *Connection) Draw(context *cairo.Context){
    var r, g, b float64
    if c.selected {
	r, g, b = ColorOption(SelectLine)
    } else if c.highlighted {
	r, g, b = ColorOption(HighlightLine)
    } else {
	r, g, b = ColorOption(NormalLine)
    }
    context.SetLineWidth(2)
    context.SetSourceRGB(r, g, b)
    p1, p2 := c.connectionPoints()
    DrawArrow(context, p1, p2)
}

/*
 *	Overwrite default implementation
 */

func (c *Connection) CheckHit(pos image.Point) (ok bool) {
    a, b := c.connectionPoints()
    m := a.Add(b).Div(2)
    am := m.Sub(a)
    r := math.Sqrt(float64(am.X*am.X) + float64(am.Y*am.Y))
    cos := float64(m.X - a.X) / r
    sin := float64(a.Y - m.Y) / r
    p := trans(cos, sin, m, pos)
    frame := image.Rect(int(-r), -4, int(r), 4)
    test := image.Rectangle{p, p}
    ok = frame.Overlaps(test)
    c.highlighted = ok
    return
}

func (c *Connection) BBox() image.Rectangle {
    c.box.Min, c.box.Max = c.connectionPoints()
    r := c.box.Canon()
    k := int(math.Ceil(arrowsize/math.Sqrt(1.5)))
    r.Min.X -= k
    r.Max.X += k
    r.Min.Y -= k
    r.Max.Y += k
    return r
}

/*
 *	Private functions
 */

func trans(cos, sin float64, m, x image.Point) (y image.Point) {
    p1 := x.Sub(m)
    a := int(cos*float64(p1.X) - sin*float64(p1.Y))
    b := int(sin*float64(p1.X) + cos*float64(p1.Y))
    y = image.Point{a, b}
    return
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
        h0x, h0y := x1 - arrowsize*dx / r, y1 - arrowsize*dy / r
        h1x, h1y := x1 - arrowsize*dy / r, y1 + arrowsize*dx / r
        h2x, h2y := x1 + arrowsize*dy / r, y1 - arrowsize*dx / r
        a1x, a1y := (2.0*h0x + h1x) / 3.0, (2.0*h0y + h1y) / 3.0
        a2x, a2y := (2.0*h0x + h2x) / 3.0, (2.0*h0y + h2y) / 3.0
        drawLine(gc, p1, p2)
        drawAccLine(gc, a1x, a1y, p2)
        drawAccLine(gc, a2x, a2y, p2)
    }
}

