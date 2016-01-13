package main

import (
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type GoAppWindow struct {
	window         *gtk.Window
	layout_box     *gtk.Box   // child of window, holds paned and menu
	paned_box      *gtk.Paned // holds navigation and views
	navigation_box *gtk.Box
	graphViews     views.GraphViewCollection
}

func (w *GoAppWindow) Window() *gtk.Window {
	return w.window
}

func (w *GoAppWindow) Init(width, height int) (err error) {
	w.window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Println("Unable to create window:", err)
		return
	}
	w.window.Connect("destroy", func() {
		gtk.MainQuit()
	})
	w.window.SetTitle("Go Application")
	w.window.SetDefaultSize(width, height)
	w.layout_box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Println("Unable to create box:", err)
		return
	}
	w.paned_box, err = gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Println("Unable to create layout box:", err)
		return
	}
	w.navigation_box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Println("Unable to create box:", err)
		return
	}
	w.graphViews, err = views.GraphViewCollectionNew(width, height)
	if err != nil {
		log.Println("Unable to create graphViews:", err)
		return
	}

	w.paned_box.Add1(w.navigation_box)
	w.paned_box.Add2(w.graphViews.Widget())
	w.paned_box.SetPosition(200)
	w.layout_box.PackEnd(w.paned_box, false, true, 0)
	w.window.Add(w.layout_box)

	return
}

func GoAppWindowNew(width, height int) (win *GoAppWindow, err error) {
	win = &GoAppWindow{}
	err = win.Init(width, height)
	return
}
