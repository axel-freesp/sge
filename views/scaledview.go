package views

import (
	"github.com/gotk3/gotk3/gtk"
	"image"
	"math"
)

const (
	scaleMin   = -4.0
	scaleMax   = 3.0
	scaleStep  = 0.01
	scalePStep = 1.0
	scalePage  = 1.0
)

type ScaledScene interface {
	Update() (width, height int)
}

type ScaledView struct {
	layoutBox *gtk.Box
	scene     *ScrolledView
	scaler    *gtk.Scale
	scale     float64
	guest     ScaledScene
}

func ScaledViewNew(guest ScaledScene) (v *ScaledView, err error) {
	v = &ScaledView{nil, nil, nil, 1.0, guest}
	err = v.init()
	return
}

func (v *ScaledView) Widget() *gtk.Widget {
	return &v.layoutBox.Container.Widget
}

func (v *ScaledView) Scale() float64 {
	return v.scale
}

func (v *ScaledView) Container() *gtk.Container {
	return &v.scene.scrolled.Bin.Container
}

func (v *ScaledView) Position(pos image.Point) image.Point {
	return image.Point{int(float64(pos.X) / v.Scale()), int(float64(pos.Y) / v.Scale())}
}

func (v *ScaledView) ScaleCoord(x int, roundUp bool) int {
	if roundUp {
		return int(math.Ceil(float64(x) * v.scale))
	} else {
		return int(float64(x) * v.scale)
	}
}

func (v *ScaledView) init() (err error) {
	v.layoutBox, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		return
	}
	const (
		initialWidth  = 800
		initialHeight = 600
	)
	v.scene, err = ScrolledViewNew(initialWidth, initialHeight)
	if err != nil {
		return
	}
	var adj *gtk.Adjustment
	adj, err = gtk.AdjustmentNew(0.0, scaleMin, scaleMax, scaleStep, scalePStep, scalePage)
	v.scaler, err = gtk.ScaleNew(gtk.ORIENTATION_HORIZONTAL, adj)
	if err != nil {
		return
	}
	v.scaler.Connect("value-changed", scalerValueCallback, v)
	v.layoutBox.Add(v.scene.Widget())
	v.layoutBox.Add(v.scaler)
	return
}

func scalerValueCallback(r *gtk.Scale, v *ScaledView) {
	oldScale := v.scale
	oldX := v.scene.hadj.GetValue()
	oldY := v.scene.vadj.GetValue()
	visMx := float64(v.layoutBox.GetAllocatedWidth()) / 2.0
	visMy := float64(v.layoutBox.GetAllocatedHeight()) / 2.0
	v.scale = math.Exp2(r.GetValue())
	newW, newH := v.guest.Update()
	v.scene.hadj.SetUpper(float64(newW))
	v.scene.vadj.SetUpper(float64(newH))
	v.scene.hadj.SetValue((oldX+visMx)*v.scale/oldScale - visMx)
	v.scene.vadj.SetValue((oldY+visMy)*v.scale/oldScale - visMy)
}
