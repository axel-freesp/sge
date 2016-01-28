package behaviour

import (
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	//"log"
	"image"
	"strings"
)

func CreateXmlInPort(p bh.PortIf) (xmlp *backend.XmlInPort) {
	xmlp = backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlOutPort(p bh.PortIf) (xmlp *backend.XmlOutPort) {
	xmlp = backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlNamedInPort(p bh.PortTypeIf) (xmlp *backend.XmlInPort) {
	xmlp = backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlNamedOutPort(p bh.PortTypeIf) (xmlp *backend.XmlOutPort) {
	xmlp = backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
	xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlInputNode(n bh.NodeIf) *backend.XmlInputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoInputNodeType-") {
		tName = ""
	}
	ret := backend.XmlInputNodeNew(n.Name(), tName)
	converter := gr.CreateModePositioner("", n)
	ret.Entry = freesp.CreateXmlModePosition(converter).Entry
	if n.(*node).portlink != nil {
		ret.NPort = n.(*node).portlink.Name()
	}
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
	}
	return ret
}

func CreateXmlOutputNode(n bh.NodeIf) *backend.XmlOutputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoOutputNodeType-") {
		tName = ""
	}
	ret := backend.XmlOutputNodeNew(n.Name(), tName)
	converter := gr.CreateModePositioner("", n)
	ret.Entry = freesp.CreateXmlModePosition(converter).Entry
	if n.(*node).portlink != nil {
		ret.NPort = n.(*node).portlink.Name()
	}
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
	}
	return ret
}

func CreateXmlProcessingNodeHint(n bh.NodeIf) (xmlh *backend.XmlNodeHint) {
	//log.Printf("CreateXmlProcessingNodeHint(%s): pathlist = %v, position = %v\n", n.Name(), n.PathList(), n.(*node).position)
	xmlh = backend.XmlNodeHintNew(n.Expanded())
	empty := image.Point{}
	for _, p := range n.PathList() {
		for _, m := range freesp.ValidModes {
			xmlp := gr.CreatePathMode(p, m)
			pos := n.PathModePosition(p, m)
			if pos != empty {
				xmlh.Entry = append(xmlh.Entry, *backend.XmlModeHintEntryNew(xmlp, pos.X, pos.Y))
			}
		}
	}
	nt := n.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			g := impl.Graph()
			for _, chn := range g.ProcessingNodes() {
				xmlh.Children = append(xmlh.Children, *CreateXmlProcessingNodeHint(chn))
			}
			break
		}
	}
	return
}

func CreateXmlProcessingNode(n bh.NodeIf) *backend.XmlProcessingNode {
	ret := backend.XmlProcessingNodeNew(n.Name(), n.ItsType().TypeName())
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
	}
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
	}
	ret.XmlNodeHint = *CreateXmlProcessingNodeHint(n)
	return ret
}

func CreateXmlNodeType(t bh.NodeTypeIf) *backend.XmlNodeType {
	ret := backend.XmlNodeTypeNew(t.TypeName())
	for _, p := range t.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlNamedInPort(p))
	}
	for _, p := range t.OutPorts() {
		ret.OutPort = append(ret.OutPort, *CreateXmlNamedOutPort(p))
	}
	for _, impl := range t.Implementation() {
		ret.Implementation = append(ret.Implementation, *CreateXmlImplementation(impl))
	}
	return ret
}

func CreateXmlImplementation(impl bh.ImplementationIf) *backend.XmlImplementation {
	ret := backend.XmlImplementationNew(impl.ElementName())
	if impl.ImplementationType() == bh.NodeTypeGraph {
		ret.SignalGraph = append(ret.SignalGraph, *CreateXmlSignalGraphType(impl.Graph()))
	}
	return ret
}

func CreateXmlConnection(c bh.ConnectionIf) *backend.XmlConnect {
	from := c.From()
	to := c.To()
	fromNode := from.Node()
	toNode := to.Node()
	switch from.Direction() {
	case gr.OutPort:
		return backend.XmlConnectNew(fromNode.Name(), toNode.Name(), from.Name(), to.Name())
	default:
		return backend.XmlConnectNew(toNode.Name(), fromNode.Name(), to.Name(), from.Name())
	}
}

func CreateXmlSignalType(s bh.SignalTypeIf) *backend.XmlSignalType {
	var scope, mode string
	if s.Scope() == bh.Local {
		scope = "local"
	}
	if s.Mode() == bh.Synchronous {
		mode = "sync"
	}
	return backend.XmlSignalTypeNew(s.TypeName(), scope, mode, s.CType(), s.ChannelId())
}

func CreateXmlLibrary(l bh.LibraryIf) *backend.XmlLibrary {
	ret := backend.XmlLibraryNew()
	for _, t := range l.SignalTypes() {
		ret.SignalTypes = append(ret.SignalTypes, *CreateXmlSignalType(t))
	}
	for _, t := range l.NodeTypes() {
		ret.NodeTypes = append(ret.NodeTypes, *CreateXmlNodeType(t))
	}
	return ret
}

func CreateXmlSignalGraph(g bh.SignalGraphIf) *backend.XmlSignalGraph {
	return CreateXmlSignalGraphType(g.ItsType())
}

func CreateXmlSignalGraphType(t bh.SignalGraphTypeIf) *backend.XmlSignalGraph {
	ret := backend.XmlSignalGraphNew()
	for _, l := range t.Libraries() {
		//log.Printf("CreateXmlSignalGraphType: l=%v\n", l)
		ret.Libraries = append(ret.Libraries, *CreateXmlLibraryRef(l))
	}
	for _, n := range t.InputNodes() {
		ret.InputNodes = append(ret.InputNodes, *CreateXmlInputNode(n))
	}
	for _, n := range t.OutputNodes() {
		ret.OutputNodes = append(ret.OutputNodes, *CreateXmlOutputNode(n))
	}
	for _, n := range t.ProcessingNodes() {
		ret.ProcessingNodes = append(ret.ProcessingNodes, *CreateXmlProcessingNode(n))
	}
	for _, n := range t.Nodes() {
		for _, p := range n.OutPorts() {
			for _, c := range p.Connections() {
				conn := p.Connection(c)
				ret.Connections = append(ret.Connections, *CreateXmlConnection(conn))
			}
		}
	}
	return ret
}

func CreateXmlLibraryRef(l bh.LibraryIf) *backend.XmlLibraryRef {
	return backend.XmlLibraryRefNew(l.Filename())
}
