package views

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
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
			t := object.(freesp.NodeType)
			xmlnodetype := createXmlNodeType(t)
			buf, err = xmlnodetype.Write()
		case freesp.Port:
			p := object.(freesp.Port)
			if p.Direction() == freesp.OutPort {
				xmlport := createXmlOutPort(p)
				buf, err = xmlport.Write()
			} else {
				xmlport := createXmlInPort(p)
				buf, err = xmlport.Write()
			}
		case freesp.NamedPortType:
			t := object.(freesp.NamedPortType)
			if t.Direction() == freesp.InPort {
				xmlporttype := createXmlNamedInPort(t)
				buf, err = xmlporttype.Write()
			} else {
				xmlporttype := createXmlNamedOutPort(t)
				buf, err = xmlporttype.Write()
			}
		case freesp.PortType:
			log.Println("XmlTextView.Set: freesp.PortType")
		case freesp.Connection:
			xmlconn := createXmlConnection(object.(freesp.Connection))
			buf, err = xmlconn.Write()
		case freesp.SignalType:
			s := object.(freesp.SignalType)
			if s != nil {
				xmlsignaltype := createXmlSignalType(s)
				buf, err = xmlsignaltype.Write()
			}
		case freesp.Library:
			log.Println("XmlTextView.Set: freesp.Library")
			l := object.(freesp.Library)
			xmllib := createXmlLibrary(l)
			buf, err = xmllib.Write()
		case freesp.Implementation:
			log.Println("XmlTextView.Set: freesp.Implementation")
			impl := object.(freesp.Implementation)
			switch impl.ImplementationType() {
			case freesp.NodeTypeElement:
				// TODO
			default:
				xmlImpl := createXmlSignalGraphType(impl.Graph())
				buf, err = xmlImpl.Write()
			}
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

func createXmlInPort(p freesp.Port) *backend.XmlInPort {
	return backend.XmlInPortNew(p.PortName(), p.ItsType().TypeName())
}

func createXmlOutPort(p freesp.Port) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.PortName(), p.ItsType().TypeName())
}

func createXmlNamedInPort(p freesp.NamedPortType) *backend.XmlInPort {
	return backend.XmlInPortNew(p.Name(), p.TypeName())
}

func createXmlNamedOutPort(p freesp.NamedPortType) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.Name(), p.TypeName())
}

func createXmlInputNode(n freesp.Node) *backend.XmlInputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoInputNodeType-") {
		tName = ""
	}
	ret := backend.XmlInputNodeNew(n.NodeName(), tName)
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *createXmlOutPort(p))
	}
	return ret
}

func createXmlOutputNode(n freesp.Node) *backend.XmlOutputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoOutputNodeType-") {
		tName = ""
	}
	ret := backend.XmlOutputNodeNew(n.NodeName(), tName)
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *createXmlInPort(p))
	}
	return ret
}

func createXmlProcessingNode(n freesp.Node) *backend.XmlProcessingNode {
	ret := backend.XmlProcessingNodeNew(n.NodeName(), n.ItsType().TypeName())
	if len(n.ItsType().DefinedAt()) == 0 {
		for _, p := range n.InPorts() {
			ret.InPort = append(ret.InPort, *createXmlInPort(p))
		}
		for _, p := range n.OutPorts() {
			ret.OutPort = append(ret.OutPort, *createXmlOutPort(p))
		}
	}
	return ret
}

func createXmlNodeType(t freesp.NodeType) *backend.XmlNodeType {
	ret := backend.XmlNodeTypeNew(t.TypeName())
	for _, p := range t.InPorts() {
		ret.InPort = append(ret.InPort, *createXmlNamedInPort(p))
	}
	for _, p := range t.OutPorts() {
		ret.OutPort = append(ret.OutPort, *createXmlNamedOutPort(p))
	}
	impl := t.Implementation()
	if impl != nil {
		ret.Implementation = append(ret.Implementation, *createXmlImplementation(impl))
	}
	return ret
}

func createXmlImplementation(impl freesp.Implementation) *backend.XmlImplementation {
	ret := backend.XmlImplementationNew(impl.ElementName())
	if impl.ImplementationType() == freesp.NodeTypeGraph {
		ret.SignalGraph = append(ret.SignalGraph, *createXmlSignalGraphType(impl.Graph()))
	}
	return ret
}

func createXmlConnection(p freesp.Connection) *backend.XmlConnect {
	switch p.From.Direction() {
	case freesp.OutPort:
		return backend.XmlConnectNew(p.From.Node().NodeName(), p.To.Node().NodeName(), p.From.PortName(), p.To.PortName())
	default:
		return backend.XmlConnectNew(p.To.Node().NodeName(), p.From.Node().NodeName(), p.To.PortName(), p.From.PortName())
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

func createXmlLibrary(l freesp.Library) *backend.XmlLibrary {
	ret := backend.XmlLibraryNew()
	for _, t := range l.SignalTypes() {
		ret.SignalTypes = append(ret.SignalTypes, *createXmlSignalType(t))
	}
	for _, t := range l.NodeTypes() {
		ret.NodeTypes = append(ret.NodeTypes, *createXmlNodeType(t))
	}
	return ret
}

func createXmlSignalGraph(g freesp.SignalGraph) *backend.XmlSignalGraph {
	return createXmlSignalGraphType(g.ItsType())
}

func createXmlSignalGraphType(t freesp.SignalGraphType) *backend.XmlSignalGraph {
	ret := backend.XmlSignalGraphNew()
	for _, l := range t.Libraries() {
		ret.Libraries = append(ret.Libraries, *createXmlLibraryRef(l))
	}
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

func createXmlLibraryRef(l freesp.Library) *backend.XmlLibraryRef {
	return backend.XmlLibraryRefNew(l.Filename())
}
