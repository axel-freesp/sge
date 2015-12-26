package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	//"github.com/gotk3/gotk3/glib"
)

type GoAppWindow struct {
	//gtk.ApplicationWindow
	window         *gtk.Window
	layout_box     *gtk.Box
	navigation_box *gtk.Box
	content_box    *gtk.Box
	header         *gtk.HeaderBar
	tabs           *gtk.StackSwitcher
	stack          *gtk.Stack
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
	w.layout_box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Println("Unable to create layout box:", err)
		return
	}
	w.navigation_box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Println("Unable to create box:", err)
		return
	}
	w.content_box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Println("Unable to create box:", err)
		return
	}
	w.header, err = gtk.HeaderBarNew()
	if err != nil {
		log.Println("Unable to create bar:", err)
		return
	}
	w.tabs, err = gtk.StackSwitcherNew()
	if err != nil {
		log.Println("Unable to create stackswitcher:", err)
		return
	}
	w.stack, err = gtk.StackNew()
	if err != nil {
		log.Println("Unable to create Stack:", err)
		return
	}
	w.layout_box.PackStart(w.navigation_box, true, true, 0)
	w.content_box.PackStart(w.header, false, true, 0)
	w.header.Add(w.tabs)
	w.tabs.SetStack(w.stack)
	w.content_box.Add(w.stack)
	w.layout_box.PackEnd(w.content_box, true, true, 0)
	w.window.Add(w.layout_box)

	return
}

func GoAppWindowNew(width, height int) (win *GoAppWindow, err error) {
	win = &GoAppWindow{nil, nil, nil, nil, nil, nil, nil}
	err = win.Init(width, height)
	return
}
