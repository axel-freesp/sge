package views

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	gr "github.com/axel-freesp/sge/interface/graph"
)

type XmlTextView struct {
	ScrolledView
	view *gtk.TextView
}

var _ XmlTextViewIf = (*XmlTextView)(nil)

func XmlTextViewNew(width, height int) (view *XmlTextView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		return
	}
	view = &XmlTextView{*v, nil}
	err = view.xmlTextViewInit()
	return
}

func (x *XmlTextView) xmlTextViewInit() (err error) {
	x.view, err = gtk.TextViewNew()
	if err != nil {
		return fmt.Errorf("Unable to create textview: %v", err)
	}
	x.scrolled.Add(x.view)
	return nil
}

func (x *XmlTextView) Set(object gr.XmlCreator) (err error) {
	switch object.(type) {
	case gr.XmlCreator:
	default:
		return fmt.Errorf("XmlTextView.Set: invalid type %T\n", object)
	}
	var textbuf *gtk.TextBuffer
	textbuf, err = x.view.GetBuffer()
	if err != nil {
		return fmt.Errorf("XmlTextView.Set: view.GetBuffer failed: %v", err)
	}
	if object == nil {
		textbuf.SetText("")
		return
	}
	var buf []byte
	buf, err = object.(gr.XmlCreator).CreateXml()
	if err != nil {
		return err
	}
	textbuf.SetText(string(buf))
	return nil
}
