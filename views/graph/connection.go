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
    LineObject
    from, to NodeIf
    fromPort, toPort int
}

func ConnectionNew(from, to NodeIf, fromPort, toPort int) (ret *Connection) {
    ret = &Connection{LineObjectInit(connectionPoints(from, to, fromPort, toPort)), from, to, fromPort, toPort}
    return ret
}

func (c Connection) IsLinked(nodeName string) bool {
    return c.from.Name() == nodeName || c.to.Name() == nodeName
}

func (c Connection) From() NodeIf {
    return c.from
}

func (c Connection) To() NodeIf {
    return c.to
}

func (c Connection) FromId() int {
    return c.fromPort
}

func (c Connection) ToId() int {
    return c.toPort
}


//
//	Drawer interface
//

func (c Connection) Draw(ctxt interface{}){
    switch ctxt.(type) {
    case *cairo.Context:
        context := ctxt.(*cairo.Context)
        var r, g, b float64
        if c.IsSelected() {
            r, g, b = ColorOption(SelectLine)
        } else if c.IsHighlighted() {
            r, g, b = ColorOption(HighlightLine)
        } else {
            r, g, b = ColorOption(NormalLine)
        }
        context.SetLineWidth(2)
        context.SetSourceRGB(r, g, b)
        p1, p2 := connectionPoints(c.from, c.to, c.fromPort, c.toPort)
        DrawArrow(context, p1, p2)
    }
}


//
//	Overwrite default implementation
//
var k = int(math.Ceil(arrowsize/math.Sqrt(1.5)))

func (c *Connection) BBox() image.Rectangle {
    c.box.Min, c.box.Max = connectionPoints(c.from, c.to, c.fromPort, c.toPort)
    c.box = c.box.Canon()
    c.box.Min.X -= k
    c.box.Max.X += k
    c.box.Min.Y -= k
    c.box.Max.Y += k
    return c.box
}


//
//	Private functions
//

func connectionPoints(from, to NodeIf, fromPort, toPort int) (p1, p2 image.Point) {
    p1 = from.BBox().Min.Add(conOut.Add(portDY.Mul(fromPort)))
    p2 = to.BBox().Min.Add(conIn.Add(portDY.Mul(toPort)))
    return
}

