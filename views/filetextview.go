package views

import (
	"github.com/axel-freesp/sge/tool"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type FileViewer struct {
	ScrolledView
	view *gtk.TextView
}

func (f *FileViewer) Init(filename string) (err error) {
	// Read file content to buf
	buf, err := tool.ReadFile(filename)
	if err != nil {
		log.Println("Warning: could not read file", filename, err)
		return
	}

	// Create a new text view to show the file content
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

	f.scrolled.Add(f.view)
	textbuf.SetText(string(buf))
	return
}

func FileViewerNew(filename string, width, height int) (viewer *FileViewer, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &FileViewer{*v, nil}
	err = viewer.Init(filename)
	return
}
