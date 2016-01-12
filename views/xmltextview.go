package views

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gtk"
	//"log"
)

type XmlTextView struct {
	ScrolledView
	view *gtk.TextView
}

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

func (x *XmlTextView) Set(object interface{}) (err error) {
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
	buf, err = freesp.CreateXML(object)
	if err != nil {
		return err
	}
	textbuf.SetText(string(buf))
	return nil
}
