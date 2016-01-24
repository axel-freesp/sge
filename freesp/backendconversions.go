package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	gr "github.com/axel-freesp/sge/interface/graph"
	"image"
	"strings"
)

var validModes = []gr.PositionMode{gr.PositionModeNormal, gr.PositionModeMapping}

var ModeFromString = map[string]gr.PositionMode{
	"normal": gr.PositionModeNormal,
	"mp":     gr.PositionModeMapping,
}

var StringFromMode = map[gr.PositionMode]string{
	gr.PositionModeNormal:  "normal",
	gr.PositionModeMapping: "mp",
}
/*
func CreateXML(object interface{}) (buf []byte, err error) {
	if object != nil {
		switch object.(type) {
		case bh.SignalGraphIf:
			s := object.(bh.SignalGraphIf)
			xmlsignalgraph := CreateXmlSignalGraph(s)
			buf, err = xmlsignalgraph.Write()
		case bh.SignalGraphTypeIf:
			s := object.(bh.SignalGraphTypeIf)
			xmlsignalgraph := CreateXmlSignalGraphType(s)
			buf, err = xmlsignalgraph.Write()
		case bh.NodeIf:
			n := object.(bh.NodeIf)
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
		case bh.NodeTypeIf:
			t := object.(bh.NodeTypeIf)
			xmlnodetype := CreateXmlNodeType(t)
			buf, err = xmlnodetype.Write()
		case bh.PortIf:
			p := object.(bh.PortIf)
			if p.Direction() == interfaces.OutPort {
				xmlport := CreateXmlOutPort(p)
				buf, err = xmlport.Write()
			} else {
				xmlport := CreateXmlInPort(p)
				buf, err = xmlport.Write()
			}
		case bh.PortTypeIf:
			t := object.(bh.PortTypeIf)
			if t.Direction() == interfaces.InPort {
				xmlporttype := CreateXmlNamedInPort(t)
				buf, err = xmlporttype.Write()
			} else {
				xmlporttype := CreateXmlNamedOutPort(t)
				buf, err = xmlporttype.Write()
			}
		case bh.ConnectionIf:
			xmlconn := CreateXmlConnection(object.(bh.ConnectionIf))
			buf, err = xmlconn.Write()
		case bh.SignalTypeIf:
			s := object.(bh.SignalTypeIf)
			if s != nil {
				xmlsignaltype := CreateXmlSignalType(s)
				buf, err = xmlsignaltype.Write()
			}
		case bh.LibraryIf:
			l := object.(bh.LibraryIf)
			xmllib := CreateXmlLibrary(l)
			buf, err = xmllib.Write()
		case bh.ImplementationIf:
			impl := object.(bh.ImplementationIf)
			switch impl.ImplementationType() {
			case bh.NodeTypeElement:
				// TODO
			default:
				xmlImpl := CreateXmlSignalGraphType(impl.Graph())
				buf, err = xmlImpl.Write()
			}
		case pf.PlatformIf:
			p := object.(pf.PlatformIf)
			xmlp := CreateXmlPlatform(p)
			buf, err = xmlp.Write()
		case pf.ArchIf:
			a := object.(pf.ArchIf)
			xmla := CreateXmlArch(a)
			buf, err = xmla.Write()
		case pf.IOTypeIf:
			t := object.(pf.IOTypeIf)
			xmlt := CreateXmlIOType(t)
			buf, err = xmlt.Write()
		case pf.ProcessIf:
			p := object.(pf.ProcessIf)
			xmlp := CreateXmlProcess(p)
			buf, err = xmlp.Write()
		case pf.ChannelIf:
			ch := object.(pf.ChannelIf)
			if ch.Direction() == interfaces.InPort {
				xmlc := CreateXmlInChannel(ch)
				buf, err = xmlc.Write()
			} else {
				xmlc := CreateXmlOutChannel(ch)
				buf, err = xmlc.Write()
			}
		case mp.MappingIf:
			m := object.(mp.MappingIf)
			xmlm := CreateXmlMapping(m)
			buf, err = xmlm.Write()
		case mp.MappedElementIf:
			m := object.(mp.MappedElementIf)
			var pname string
			if m.Process() != nil {
				pname = fmt.Sprintf("%s/%s", m.Process().Arch().Name(), m.Process().Name())
			}
			if len(m.Node().InPorts()) > 0 && len(m.Node().OutPorts()) > 0 {
				xmlm := CreateXmlNodeMap(m.Node().Name(), pname, m.Position())
				buf, err = xmlm.Write()
			} else {
				xmlm := CreateXmlIOMap(m.Node().Name(), pname, m.Position())
				buf, err = xmlm.Write()
			}
		default:
			err = fmt.Errorf("CreateXML: invalid data type %T (%v)\n", object, object)
		}
	}
	return
}
*/

func CreateXmlInPort(p bh.PortIf) *backend.XmlInPort {
	return backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlOutPort(p bh.PortIf) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlNamedInPort(p bh.PortTypeIf) *backend.XmlInPort {
	return backend.XmlInPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlNamedOutPort(p bh.PortTypeIf) *backend.XmlOutPort {
	return backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
}

func CreateXmlInputNode(n bh.NodeIf) *backend.XmlInputNode {
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

func CreateXmlOutputNode(n bh.NodeIf) *backend.XmlOutputNode {
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

func CreateXmlProcessingNode(n bh.NodeIf) *backend.XmlProcessingNode {
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

func CreateXmlPlatform(p pf.PlatformIf) *backend.XmlPlatform {
	ret := backend.XmlPlatformNew()
	ret.PlatformId = p.PlatformId()
	for _, a := range p.Arch() {
		ret.Arch = append(ret.Arch, *CreateXmlArch(a))
	}
	return ret
}

func CreateXmlArch(a pf.ArchIf) *backend.XmlArch {
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

func CreateXmlIOType(t pf.IOTypeIf) *backend.XmlIOType {
	return backend.XmlIOTypeNew(t.Name(), ioXmlModeMap[t.IOMode()])
}

func CreateXmlProcess(p pf.ProcessIf) *backend.XmlProcess {
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

func CreateXmlInChannel(ch pf.ChannelIf) *backend.XmlInChannel {
	ret := backend.XmlInChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = CreateXmlModePosition(c.archport).Entry
	return ret
}

func CreateXmlOutChannel(ch pf.ChannelIf) *backend.XmlOutChannel {
	ret := backend.XmlOutChannelNew(ch.Name(), ch.IOType().Name(), ch.(*channel).linkText)
	ret.Entry = CreateXmlModePosition(ch).Entry
	c := ch.(*channel)
	ret.ArchPortHints.Entry = CreateXmlModePosition(c.archport).Entry
	return ret
}

func CreateXmlIOMap(node, process string, pos image.Point) *backend.XmlIOMap {
	return backend.XmlIOMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlNodeMap(node, process string, pos image.Point) *backend.XmlNodeMap {
	return backend.XmlNodeMapNew(node, process, pos.X, pos.Y)
}

func CreateXmlMapping(m mp.MappingIf) (xmlm *backend.XmlMapping) {
	xmlm = backend.XmlMappingNew(m.Graph().Filename(), m.Platform().Filename())
	g := m.Graph().ItsType()
	for _, n := range g.InputNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.OutputNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.IOMappings = append(xmlm.IOMappings, *CreateXmlIOMap(n.Name(), pname, melem.Position()))
		}
	}
	for _, n := range g.ProcessingNodes() {
		melem, _ := m.MappedElement(n)
		p, ok := m.Mapped(n)
		if ok {
			pname := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			xmlm.Mappings = append(xmlm.Mappings, *CreateXmlNodeMap(n.Name(), pname, melem.Position()))
		}
	}
	return
}

func CreateXmlModePosition(x gr.ModePositioner) (h *backend.XmlModeHint) {
	h = backend.XmlModeHintNew()
	for _, m := range validModes {
		pos := x.ModePosition(m)
		h.Entry = append(h.Entry, *backend.XmlModeHintEntryNew(StringFromMode[m], pos.X, pos.Y))
	}
	return
}
