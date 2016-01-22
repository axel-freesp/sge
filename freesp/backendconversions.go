package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	interfaces "github.com/axel-freesp/sge/interface"
	"image"
	"strings"
)

var validModes = []interfaces.PositionMode{interfaces.PositionModeNormal, interfaces.PositionModeMapping}

var modeFromString = map[string]interfaces.PositionMode{
	"normal":  interfaces.PositionModeNormal,
	"mapping": interfaces.PositionModeMapping,
}

var stringFromMode = map[interfaces.PositionMode]string{
	interfaces.PositionModeNormal:  "normal",
	interfaces.PositionModeMapping: "mapping",
}

func CreateXML(object interface{}) (buf []byte, err error) {
	if object != nil {
		switch object.(type) {
		case SignalGraph:
			s := object.(SignalGraph)
			xmlsignalgraph := CreateXmlSignalGraph(s)
			buf, err = xmlsignalgraph.Write()
		case SignalGraphType:
			s := object.(SignalGraphType)
			xmlsignalgraph := CreateXmlSignalGraphType(s)
			buf, err = xmlsignalgraph.Write()
		case Node:
			n := object.(Node)
			if len(n.InPorts()) == 0 {
				xmlnode := CreateXmlInputNode(n)
				buf, err = xmlnode.Write()
			} else if len(n.OutPorts()) == 0 {
				xmlnode := CreateXmlOutputNode(n)
				buf, err = xmlnode.Write()
			} else {
				xmlnode := CreateXmlProcessingNode(n)
				buf, err = xmlnode.Write()
			}
		case NodeType:
			t := object.(NodeType)
			xmlnodetype := CreateXmlNodeType(t)
			buf, err = xmlnodetype.Write()
		case Port:
			p := object.(Port)
			if p.Direction() == interfaces.OutPort {
				xmlport := CreateXmlOutPort(p)
				buf, err = xmlport.Write()
			} else {
				xmlport := CreateXmlInPort(p)
				buf, err = xmlport.Write()
			}
		case PortType:
			t := object.(PortType)
			if t.Direction() == interfaces.InPort {
				xmlporttype := CreateXmlNamedInPort(t)
				buf, err = xmlporttype.Write()
			} else {
				xmlporttype := CreateXmlNamedOutPort(t)
				buf, err = xmlporttype.Write()
			}
		case Connection:
			xmlconn := CreateXmlConnection(object.(Connection))
			buf, err = xmlconn.Write()
		case SignalType:
			s := object.(SignalType)
			if s != nil {
				xmlsignaltype := CreateXmlSignalType(s)
				buf, err = xmlsignaltype.Write()
			}
		case Library:
			l := object.(Library)
			xmllib := CreateXmlLibrary(l)
			buf, err = xmllib.Write()
		case Implementation:
			impl := object.(Implementation)
			switch impl.ImplementationType() {
			case NodeTypeElement:
				// TODO
			default:
				xmlImpl := CreateXmlSignalGraphType(impl.Graph())
				buf, err = xmlImpl.Write()
			}
		case Platform:
			p := object.(Platform)
			xmlp := CreateXmlPlatform(p)
			buf, err = xmlp.Write()
		case Arch:
			a := object.(Arch)
			xmla := CreateXmlArch(a)
			buf, err = xmla.Write()
		case IOType:
			t := object.(IOType)
			xmlt := CreateXmlIOType(t)
			buf, err = xmlt.Write()
		case Process:
			p := object.(Process)
			xmlp := CreateXmlProcess(p)
			buf, err = xmlp.Write()
		case Channel:
			ch := object.(Channel)
			if ch.Direction() == interfaces.InPort {
				xmlc := CreateXmlInChannel(ch)
				buf, err = xmlc.Write()
			} else {
				xmlc := CreateXmlOutChannel(ch)
				buf, err = xmlc.Write()
			}
		case Mapping:
			m := object.(Mapping)
			xmlm := CreateXmlMapping(m)
			buf, err = xmlm.Write()
		case MappedElement:
			m := object.(*mapelem)
			var pname string
			if m.process != nil {
				pname = fmt.Sprintf("%s/%s", m.process.Arch().Name(), m.process.Name())
			}
			if len(m.node.InPorts()) > 0 && len(m.node.OutPorts()) > 0 {
				xmlm := CreateXmlNodeMap(m.node.Name(), pname, m.Position())
				buf, err = xmlm.Write()
			} else {
				xmlm := CreateXmlIOMap(m.node.Name(), pname, m.Position())
				buf, err = xmlm.Write()
			}
		default:
			err = fmt.Errorf("CreateXML: invalid data type %T (%v)\n", object, object)
		}
	}
	return
}

func CreateXmlInPort(p Port) *backend.XmlInPort {
	return backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlOutPort(p Port) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
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
	}
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
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
	}
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
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

func CreateXmlConnection(c Connection) *backend.XmlConnect {
	from := c.From()
	to := c.To()
	fromNode := from.Node()
	toNode := to.Node()
	switch from.Direction() {
	case interfaces.OutPort:
		return backend.XmlConnectNew(fromNode.Name(), toNode.Name(), from.Name(), to.Name())
	default:
		return backend.XmlConnectNew(toNode.Name(), fromNode.Name(), to.Name(), from.Name())
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
				conn := p.Connection(c)
				ret.Connections = append(ret.Connections, *CreateXmlConnection(conn))
			}
		}
	}
	return ret
}

func CreateXmlLibraryRef(l Library) *backend.XmlLibraryRef {
	return backend.XmlLibraryRefNew(l.Filename())
}

func CreateXmlPlatform(p Platform) *backend.XmlPlatform {
	ret := backend.XmlPlatformNew()
	ret.PlatformId = p.PlatformId()
	for _, a := range p.Arch() {
		ret.Arch = append(ret.Arch, *CreateXmlArch(a))
	}
	return ret
}

func CreateXmlArch(a Arch) *backend.XmlArch {
	ret := backend.XmlArchNew(a.Name())
	for _, t := range a.IOTypes() {
		ret.IOType = append(ret.IOType, *CreateXmlIOType(t))
	}
	for _, p := range a.Processes() {
		ret.Processes = append(ret.Processes, *CreateXmlProcess(p))
	}
	ret.Entry = CreateXmlModePosition(a).Entry
	return ret
}

func CreateXmlIOType(t IOType) *backend.XmlIOType {
	return backend.XmlIOTypeNew(t.Name(), ioXmlModeMap[t.IOMode()])
}

func CreateXmlProcess(p Process) *backend.XmlProcess {
	ret := backend.XmlProcessNew(p.Name())
	for _, c := range p.InChannels() {
		ret.InputChannels = append(ret.InputChannels, *CreateXmlInChannel(c))
	}
	for _, c := range p.OutChannels() {
		ret.OutputChannels = append(ret.OutputChannels, *CreateXmlOutChannel(c))
	}
	ret.Entry = CreateXmlModePosition(p).Entry
	return ret
}

func CreateXmlInChannel(ch Channel) *backend.XmlInChannel {
	ret := backend.XmlInChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = CreateXmlModePosition(c.archPort).Entry
	return ret
}

func CreateXmlOutChannel(ch Channel) *backend.XmlOutChannel {
	ret := backend.XmlOutChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = CreateXmlModePosition(c.archPort).Entry
	return ret
}

func CreateXmlIOMap(node, process string, pos image.Point) *backend.XmlIOMap {
	return backend.XmlIOMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlNodeMap(node, process string, pos image.Point) *backend.XmlNodeMap {
	return backend.XmlNodeMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlMapping(m Mapping) (xmlm *backend.XmlMapping) {
	xmlm = backend.XmlMappingNew(m.Graph().Filename(), m.Platform().Filename())
	g := m.Graph().ItsType()
	for _, n := range g.InputNodes() {
		melem := m.(*mapping).maps[n.Name()]
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.OutputNodes() {
		melem := m.(*mapping).maps[n.Name()]
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.ProcessingNodes() {
		melem := m.(*mapping).maps[n.Name()]
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.Mappings = append(xmlm.Mappings, *CreateXmlNodeMap(n.Name(), pname, melem.Position()))
		}
	}
	return
}

func CreateXmlModePosition(x interfaces.ModePositioner) (h *backend.XmlModeHint) {
	h = backend.XmlModeHintNew()
	for _, m := range validModes {
		pos := x.ModePosition(m)
		h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(stringFromMode[m], pos.X, pos.Y))
	}
	return
}
