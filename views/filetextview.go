package views

import (
	"github.com/axel-freesp/tool"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type FileViewer struct {
	hadj, vadj *gtk.Adjustment
	scrolled   *gtk.ScrolledWindow
	view       *gtk.TextView
}

func (f *FileViewer) Init(filename string, width, height int) (err error) {
	// Read file content to buf
	buf, err := tool.ReadFile(filename)
	if err != nil {
		log.Println("Warning: could not read file", filename, err)
		return
	}

	// Create a new text view to show the file content
	f.hadj, err = gtk.AdjustmentNew(0, 0, float64(width), 1, 10, float64(width))
	if err != nil {
		log.Println("Unable to create hadj:", err)
		return
	}
	f.vadj, err = gtk.AdjustmentNew(0, 0, float64(height), 1, 10, float64(height))
	if err != nil {
		log.Println("Unable to create vadj:", err)
		return
	}
	f.scrolled, err = gtk.ScrolledWindowNew(f.hadj, f.vadj)
	if err != nil {
		log.Println("Unable to create scrolled:", err)
		return
	}
	f.view, err = gtk.TextViewNew()
	if err != nil {
		log.Println("Unable to create textview:", err)
		return
	}
	textbuf, err := f.view.GetBuffer()
	if err != nil {
		log.Println("view.GetBuffer failed:", err)
		return
	}

	f.scrolled.SetHExpand(true)
	f.scrolled.SetVExpand(true)
	f.view.SetEditable(false)
	f.view.SetCursorVisible(false)
	f.scrolled.Add(f.view)
	textbuf.SetText(string(buf))
	return
}

func FileViewerNew(filename string, width, height int) (viewer *FileViewer, err error) {
	viewer = &FileViewer{nil, nil, nil, nil}
	err = viewer.Init(filename, width, height)
	return
}
