package views

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gtk"
	"log"
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

func (x *XmlTextView) Set(object interface{}) error {
	var buf []byte
	var err error
	if object != nil {
		switch object.(type) {
		case freesp.SignalGraph:
			s := object.(freesp.SignalGraph)
			xmlsignalgraph := freesp.CreateXmlSignalGraph(s)
			buf, err = xmlsignalgraph.Write()
		case freesp.SignalGraphType:
			s := object.(freesp.SignalGraphType)
			xmlsignalgraph := freesp.CreateXmlSignalGraphType(s)
			buf, err = xmlsignalgraph.Write()
		case freesp.Node:
			n := object.(freesp.Node)
			if len(n.InPorts()) == 0 {
				xmlnode := freesp.CreateXmlInputNode(n)
				buf, err = xmlnode.Write()
			} else if len(n.OutPorts()) == 0 {
				xmlnode := freesp.CreateXmlOutputNode(n)
				buf, err = xmlnode.Write()
			} else {
				xmlnode := freesp.CreateXmlProcessingNode(n)
				buf, err = xmlnode.Write()
			}
		case freesp.NodeType:
			t := object.(freesp.NodeType)
			xmlnodetype := freesp.CreateXmlNodeType(t)
			buf, err = xmlnodetype.Write()
		case freesp.Port:
			p := object.(freesp.Port)
			if p.Direction() == freesp.OutPort {
				xmlport := freesp.CreateXmlOutPort(p)
				buf, err = xmlport.Write()
			} else {
				xmlport := freesp.CreateXmlInPort(p)
				buf, err = xmlport.Write()
			}
		case freesp.NamedPortType:
			t := object.(freesp.NamedPortType)
			if t.Direction() == freesp.InPort {
				xmlporttype := freesp.CreateXmlNamedInPort(t)
				buf, err = xmlporttype.Write()
			} else {
				xmlporttype := freesp.CreateXmlNamedOutPort(t)
				buf, err = xmlporttype.Write()
			}
		case freesp.PortType:
			pt := object.(freesp.PortType)
			s := pt.SignalType()
			if s != nil {
				xmlsignaltype := freesp.CreateXmlSignalType(s)
				buf, err = xmlsignaltype.Write()
			}
		case freesp.Connection:
			xmlconn := freesp.CreateXmlConnection(object.(freesp.Connection))
			buf, err = xmlconn.Write()
		case freesp.SignalType:
			s := object.(freesp.SignalType)
			if s != nil {
				xmlsignaltype := freesp.CreateXmlSignalType(s)
				buf, err = xmlsignaltype.Write()
			}
		case freesp.Library:
			l := object.(freesp.Library)
			xmllib := freesp.CreateXmlLibrary(l)
			buf, err = xmllib.Write()
		case freesp.Implementation:
			impl := object.(freesp.Implementation)
			switch impl.ImplementationType() {
			case freesp.NodeTypeElement:
				// TODO
			default:
				xmlImpl := freesp.CreateXmlSignalGraphType(impl.Graph())
				buf, err = xmlImpl.Write()
			}
		default:
			log.Printf("XmlTextView.Set: invalid data type %T (%v)\n", object, object)
		}
	}
	textbuf, err := x.view.GetBuffer()
	if err != nil {
		return fmt.Errorf("XmlTextView.Set: view.GetBuffer failed: %v", err)
	}
	textbuf.SetText(string(buf))
	return nil
}
