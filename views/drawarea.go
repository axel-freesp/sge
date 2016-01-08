package views

import (
    "image"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type DrawAreaClient interface {
    DrawCallback(area DrawArea, context *cairo.Context)
    MotionCallback(area DrawArea, pos image.Point)
    ButtonCallback(area DrawArea, evType gdk.EventType, pos image.Point)
}

type DrawArea struct {
	*gtk.DrawingArea
}

func DrawAreaInit(client DrawAreaClient) (v DrawArea, err error) {
    var area *gtk.DrawingArea
	area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
    v = DrawArea{area}
	v.Connect("draw", drawCallback, client)
	v.Connect("motion-notify-event", motionCallback, client)
	v.Connect("button-press-event", buttonCallback, client)
	v.Connect("button-release-event", buttonCallback, client)
	v.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))
    return
}

func drawCallback(area *gtk.DrawingArea, context *cairo.Context, v DrawAreaClient) bool {
    v.DrawCallback(DrawArea{area}, context)
	return true
}

func motionCallback(area *gtk.DrawingArea, event *gdk.Event, v DrawAreaClient) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x), int(y)}
    v.MotionCallback(DrawArea{area}, pos)
}

func buttonCallback(area *gtk.DrawingArea, event *gdk.Event, v DrawAreaClient) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X()), int(ev.Y())}
    v.ButtonCallback(DrawArea{area}, ev.Type(), pos)
}

