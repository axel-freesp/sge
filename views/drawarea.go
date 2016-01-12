package views

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"image"
)

type DrawAreaClient interface {
	DrawCallback(area DrawArea, context *cairo.Context)
	MotionCallback(area DrawArea, pos image.Point)
	ButtonCallback(area DrawArea, evType gdk.EventType, pos image.Point)
}

type DrawArea struct {
	*gtk.DrawingArea
	shDraw, shMotion, shButtonPress, shButtonRelease glib.SignalHandle
}

func DrawAreaInit(client DrawAreaClient) (v DrawArea, err error) {
	var area *gtk.DrawingArea
	area, err = gtk.DrawingAreaNew()
	if err != nil {
		return
	}
	v = DrawArea{area, 0, 0, 0, 0}
	v.shDraw, err = v.Connect("draw", drawCallback, client)
	if err != nil {
		return
	}
	v.shMotion, err = v.Connect("motion-notify-event", motionCallback, client)
	if err != nil {
		return
	}
	v.shButtonPress, err = v.Connect("button-press-event", buttonCallback, client)
	if err != nil {
		return
	}
	v.shButtonRelease, err = v.Connect("button-release-event", buttonCallback, client)
	if err != nil {
		return
	}
	v.SetEvents(int(gdk.POINTER_MOTION_MASK | gdk.BUTTON_PRESS_MASK | gdk.BUTTON_RELEASE_MASK))
	return
}

func DrawAreaDestroy(v DrawArea) {
	//v.Unmap()
	v.SetEvents(0)
	v.HandlerDisconnect(v.shButtonRelease)
	v.HandlerDisconnect(v.shButtonPress)
	v.HandlerDisconnect(v.shMotion)
	v.HandlerDisconnect(v.shDraw)
}

func drawCallback(area *gtk.DrawingArea, context *cairo.Context, v DrawAreaClient) bool {
	v.DrawCallback(DrawArea{area, 0, 0, 0, 0}, context)
	return true
}

func motionCallback(area *gtk.DrawingArea, event *gdk.Event, v DrawAreaClient) {
	ev := gdk.EventMotion{event}
	x, y := ev.MotionVal()
	pos := image.Point{int(x), int(y)}
	v.MotionCallback(DrawArea{area, 0, 0, 0, 0}, pos)
}

func buttonCallback(area *gtk.DrawingArea, event *gdk.Event, v DrawAreaClient) {
	ev := gdk.EventButton{event}
	pos := image.Point{int(ev.X()), int(ev.Y())}
	v.ButtonCallback(DrawArea{area, 0, 0, 0, 0}, ev.Type(), pos)
}
