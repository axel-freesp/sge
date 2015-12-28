package freesp

import (
	"github.com/axel-freesp/sge/backend"
	"strings"
)

func CreateXmlInPort(p Port) *backend.XmlInPort {
	return backend.XmlInPortNew(p.PortName(), p.ItsType().TypeName())
}

func CreateXmlOutPort(p Port) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.PortName(), p.ItsType().TypeName())
}

func CreateXmlNamedInPort(p PortType) *backend.XmlInPort {
	return backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlNamedOutPort(p PortType) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlInputNode(n Node) *backend.XmlInputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoInputNodeType-") {
		tName = ""
	}
	pos := n.Position()
	ret := backend.XmlInputNodeNew(n.Name(), tName, pos.X, pos.Y)
	if n.(*node).portlink != nil {
		ret.NPort = n.(*node).portlink.Name()
	} else {
		for _, p := range n.OutPorts() {
			ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
		}
	}
	return ret
}

func CreateXmlOutputNode(n Node) *backend.XmlOutputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoOutputNodeType-") {
		tName = ""
	}
	pos := n.Position()
	ret := backend.XmlOutputNodeNew(n.Name(), tName, pos.X, pos.Y)
	if n.(*node).portlink != nil {
		ret.NPort = n.(*node).portlink.Name()
	} else {
		for _, p := range n.InPorts() {
			ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
		}
	}
	return ret
}

func CreateXmlProcessingNode(n Node) *backend.XmlProcessingNode {
	pos := n.Position()
	ret := backend.XmlProcessingNodeNew(n.Name(), n.ItsType().TypeName(), pos.X, pos.Y)
	if len(n.ItsType().DefinedAt()) == 0 {
		for _, p := range n.InPorts() {
			ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
		}
		for _, p := range n.OutPorts() {
			ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
		}
	}
	return ret
}

func CreateXmlNodeType(t NodeType) *backend.XmlNodeType {
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

func CreateXmlImplementation(impl Implementation) *backend.XmlImplementation {
	ret := backend.XmlImplementationNew(impl.ElementName())
	if impl.ImplementationType() == NodeTypeGraph {
		ret.SignalGraph = append(ret.SignalGraph, *CreateXmlSignalGraphType(impl.Graph()))
	}
	return ret
}

func CreateXmlConnection(p Connection) *backend.XmlConnect {
	switch p.From.Direction() {
	case OutPort:
		return backend.XmlConnectNew(p.From.Node().Name(), p.To.Node().Name(), p.From.PortName(), p.To.PortName())
	default:
		return backend.XmlConnectNew(p.To.Node().Name(), p.From.Node().Name(), p.To.PortName(), p.From.PortName())
	}
}

func CreateXmlSignalType(s SignalType) *backend.XmlSignalType {
	var scope, mode string
	if s.Scope() == Local {
		scope = "local"
	}
	if s.Mode() == Synchronous {
		mode = "sync"
	}
	return backend.XmlSignalTypeNew(s.TypeName(), scope, mode, s.CType(), s.ChannelId())
}

func CreateXmlLibrary(l Library) *backend.XmlLibrary {
	ret := backend.XmlLibraryNew()
	for _, t := range l.SignalTypes() {
		ret.SignalTypes = append(ret.SignalTypes, *CreateXmlSignalType(t))
	}
	for _, t := range l.NodeTypes() {
		ret.NodeTypes = append(ret.NodeTypes, *CreateXmlNodeType(t))
	}
	return ret
}

func CreateXmlSignalGraph(g SignalGraph) *backend.XmlSignalGraph {
	return CreateXmlSignalGraphType(g.ItsType())
}

func CreateXmlSignalGraphType(t SignalGraphType) *backend.XmlSignalGraph {
	ret := backend.XmlSignalGraphNew()
	for _, l := range t.Libraries() {
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
				conn := Connection{p, c}
				ret.Connections = append(ret.Connections, *CreateXmlConnection(conn))
			}
		}
	}
	return ret
}

func CreateXmlLibraryRef(l Library) *backend.XmlLibraryRef {
	return backend.XmlLibraryRefNew(l.Filename())
}
