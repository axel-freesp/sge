package main

import (
	"github.com/gotk3/gotk3/gtk"
)

func createLabeledRow(labelText string, widget *gtk.Widget) (box *gtk.Box, err error) {
	box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		return
	}
	label, err := gtk.LabelNew(labelText)
	if err != nil {
		return
	}
	box.PackStart(label, false, false, 6)
	if widget != nil {
		box.PackEnd(widget, false, false, 6)
	}
	return
}
