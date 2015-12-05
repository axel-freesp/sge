package views

import (
	"log"
	"fmt"
	"encoding/xml"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gtk"
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
			log.Println("XmlTextView.Set: freesp.SignalGraph")
			s := object.(freesp.SignalGraph)
			xmlsignalgraph := createXmlSignalGraph(s)
			buf, err = xmlsignalgraph.Write()
		case freesp.Node:
			n := object.(freesp.Node)
			if len(n.InPorts()) == 0 {
				xmlnode := createXmlInputNode(n)
				buf, err = xmlnode.Write()
			} else if len(n.OutPorts()) == 0 {
				xmlnode := createXmlOutputNode(n)
				buf, err = xmlnode.Write()
			} else {
				xmlnode := createXmlProcessingNode(n)
				buf, err = xmlnode.Write()
			}
		case freesp.NodeType:
			log.Println("XmlTextView.Set: freesp.NodeType")
		case freesp.Port:
			p := object.(freesp.Port)
			switch p.Direction() {
			case freesp.OutPort:
				xmlport := createXmlOutPort(p)
				buf, err = xmlport.Write()
			default:
				xmlport := createXmlInPort(p)
				buf, err = xmlport.Write()
			}
		case freesp.PortType:
			log.Println("XmlTextView.Set: freesp.PortType")
		case freesp.NamedPortType:
			log.Println("XmlTextView.Set: freesp.NamedPortType")
		case freesp.Connection:
			xmlconn := createXmlConnection(object.(freesp.Connection))
			buf, err = xmlconn.Write()
		case freesp.SignalType:
			s := object.(freesp.SignalType)
			xmlsignaltype := createXmlSignalType(s)
			buf, err = xmlsignaltype.Write()
		default:
			log.Println("XmlTextView.Set: invalid data type")
		}
	}
	textbuf, err := x.view.GetBuffer()
	if err != nil {
		return fmt.Errorf("XmlTextView.Set: view.GetBuffer failed: %v", err)
	}
	textbuf.SetText(string(buf))
	return nil
}

// Conversions: freesp interface type -> backend XML type

func createXmlInPort(p freesp.Port) *backend.XmlInPort{
	return backend.XmlInPortNew(p.PortName(), p.ItsType().TypeName())
}

func createXmlOutPort(p freesp.Port) *backend.XmlOutPort{
	return backend.XmlOutPortNew(p.PortName(), p.ItsType().TypeName())
}

func createXmlInputNode(n freesp.Node) *backend.XmlInputNode{
	ret := backend.XmlInputNodeNew(n.NodeName(), n.ItsType().TypeName())
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *createXmlOutPort(p))
	}
	return ret
}

func createXmlOutputNode(n freesp.Node) *backend.XmlOutputNode{
	ret := backend.XmlOutputNodeNew(n.NodeName(), n.ItsType().TypeName())
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *createXmlInPort(p))
	}
	return ret
}

func createXmlProcessingNode(n freesp.Node) *backend.XmlProcessingNode{
	ret := backend.XmlProcessingNodeNew(n.NodeName(), n.ItsType().TypeName())
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *createXmlInPort(p))
	}
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *createXmlOutPort(p))
	}
	return ret
}

func createXmlConnection(p freesp.Connection) *backend.XmlConnect{
	switch p.From.Direction() {
	case freesp.OutPort:
		return &backend.XmlConnect{xml.Name{"http://www.freesp.de/xml/freeSP", "connect"}, p.From.Node().NodeName(), p.To.Node().NodeName(), p.From.PortName(), p.To.PortName()}
	default:
		return &backend.XmlConnect{xml.Name{"http://www.freesp.de/xml/freeSP", "connect"}, p.To.Node().NodeName(), p.From.Node().NodeName(), p.To.PortName(), p.From.PortName()}
	}
}

func createXmlSignalType(s freesp.SignalType) *backend.XmlSignalType {
	var scope, mode string
	if s.Scope() == freesp.Local {
		scope = "local"
	}
	if s.Mode() == freesp.Synchronous {
		mode = "sync"
	}
	return backend.XmlSignalTypeNew(s.TypeName(), scope, mode, s.CType(), s.ChannelId())
}

func createXmlSignalGraph(g freesp.SignalGraph) *backend.XmlSignalGraph{
	t := g.ItsType()
	ret := backend.XmlSignalGraphNew()
	for _, s := range t.SignalTypes() {
		ret.SignalTypes = append(ret.SignalTypes, *createXmlSignalType(s))
	}
	for _, n := range t.InputNodes() {
		ret.InputNodes = append(ret.InputNodes, *createXmlInputNode(n))
	}
	for _, n := range t.OutputNodes() {
		ret.OutputNodes = append(ret.OutputNodes, *createXmlOutputNode(n))
	}
	for _, n := range t.ProcessingNodes() {
		ret.ProcessingNodes = append(ret.ProcessingNodes, *createXmlProcessingNode(n))
	}
	for _, n := range t.Nodes() {
		for _, p := range n.OutPorts() {
			for _, c := range p.Connections() {
				conn := freesp.Connection{p, c}
				ret.Connections = append(ret.Connections, *createXmlConnection(conn))
			}
		}
	}
	return ret
}



