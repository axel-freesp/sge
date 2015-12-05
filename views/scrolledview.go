package views

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"unsafe"
)

type ScrolledView struct {
	hadj, vadj *gtk.Adjustment
	scrolled   *gtk.ScrolledWindow
}

func ScrolledViewNew(width, height int) (viewer *ScrolledView, err error) {
	viewer = &ScrolledView{nil, nil, nil}
	err = viewer.scrolledViewInit(width, height)
	return
}

func (v *ScrolledView) Widget() *gtk.Widget {
	return (*gtk.Widget)(unsafe.Pointer(v.scrolled))
}

func (v *ScrolledView) scrolledViewInit(width, height int) (err error) {
	v.hadj, err = gtk.AdjustmentNew(0, 0, float64(width), 1, 10, float64(width))
	if err != nil {
		log.Println("Unable to create hadj:", err)
		return
	}
	v.vadj, err = gtk.AdjustmentNew(0, 0, float64(height), 1, 10, float64(height))
	if err != nil {
		log.Println("Unable to create vadj:", err)
		return
	}
	v.scrolled, err = gtk.ScrolledWindowNew(v.hadj, v.vadj)
	if err != nil {
		log.Println("Unable to create scrolled:", err)
		return
	}
	v.scrolled.SetHExpand(true)
	v.scrolled.SetVExpand(true)
	return
}
