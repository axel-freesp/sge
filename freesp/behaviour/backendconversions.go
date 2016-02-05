package behaviour

import (
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/tool"
	//"log"
	"fmt"
	"image"
	"strings"
)

//
//	Hints file
//

func CreateXmlGraphHint(g bh.SignalGraphIf) (xmlg *backend.XmlGraphHint) {
	xmlg = backend.XmlGraphHintNew(g.Filename())
	for _, n := range g.ItsType().InputNodes() {
		xmlg.InputNode = append(xmlg.InputNode, *CreateXmlIONodePosHint(n, ""))
	}
	for _, n := range g.ItsType().OutputNodes() {
		xmlg.OutputNode = append(xmlg.OutputNode, *CreateXmlIONodePosHint(n, ""))
	}
	for _, n := range g.ItsType().ProcessingNodes() {
		hintlist := CreateXmlNodePosHint(n, "")
		for _, h := range hintlist {
			xmlg.ProcessingNode = append(xmlg.ProcessingNode, h)
		}
	}
	return
}

func CreateXmlNodePosHint(nd bh.NodeIf, path string) (xmln []backend.XmlNodePosHint) {
	xmlnd := CreateXmlIONodePosHint(nd, path)
	xmlnd.Expanded = nd.Expanded()
	xmln = append(xmln, *xmlnd)
	nt := nd.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, n := range impl.Graph().ProcessingNodes() {
				var p string
				if len(path) == 0 {
					p = nd.Name()
				} else {
					p = fmt.Sprintf("%s/%s", path, nd.Name())
				}
				hintlist := CreateXmlNodePosHint(n, p)
				for _, h := range hintlist {
					xmln = append(xmln, h)
				}
			}
			break
		}
	}
	return
}

func CreateXmlIONodePosHint(n bh.NodeIf, path string) (xmln *backend.XmlNodePosHint) {
	if len(path) == 0 {
		xmln = backend.XmlNodePosHintNew(n.Name())
	} else {
		xmln = backend.XmlNodePosHintNew(fmt.Sprintf("%s/%s", path, n.Name()))
	}
	empty := image.Point{}
	for _, p := range n.PathList() {
		for _, m := range gr.ValidModes {
			xmlp := string(gr.CreatePathMode(p, m))
			pos := n.PathModePosition(p, m)
			if pos != empty {
				xmln.Entry = append(xmln.Entry, *backend.XmlModeHintEntryNew(xmlp, pos.X, pos.Y))
			}
		}
	}
	for _, p := range n.InPorts() {
		xmlp := backend.XmlPortPosHintNew(p.Name())
		xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
		xmln.InPorts = append(xmln.InPorts, *xmlp)
	}
	for _, p := range n.OutPorts() {
		xmlp := backend.XmlPortPosHintNew(p.Name())
		xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
		xmln.OutPorts = append(xmln.OutPorts, *xmlp)
	}
	return
}

//
//	Behaviour model files:
//

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
	//xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlNamedOutPort(p bh.PortTypeIf) (xmlp *backend.XmlOutPort) {
	xmlp = backend.XmlOutPortNew(p.Name(), p.SignalType().TypeName())
	//xmlp.Entry = freesp.CreateXmlModePosition(p).Entry
	return
}

func CreateXmlInputNode(n bh.NodeIf) *backend.XmlInputNode {
	tName := n.ItsType().TypeName()
	if strings.HasPrefix(tName, "autoInputNodeType-") {
		tName = ""
	}
	ret := backend.XmlInputNodeNew(n.Name(), tName)
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
	if n.(*node).portlink != nil {
		ret.NPort = n.(*node).portlink.Name()
	}
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
	}
	return ret
}

func CreateXmlProcessingNode(n bh.NodeIf) *backend.XmlProcessingNode {
	ret := backend.XmlProcessingNodeNew(n.Name(), n.ItsType().TypeName())
	for _, p := range n.InPorts() {
		ret.InPort = append(ret.InPort, *CreateXmlInPort(p))
	}
	for _, p := range n.OutPorts() {
		ret.OutPort = append(ret.OutPort, *CreateXmlOutPort(p))
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
	fname := l.Filename()
	ret := backend.XmlLibraryNew()
	reflist := tool.StringListInit()
	for _, t := range l.SignalTypes() {
		ref := t.DefinedAt()
		_, ok := reflist.Find(ref)
		if !ok && ref != fname {
			reflist.Append(ref)
		}
	}
	for _, t := range l.NodeTypes() {
		ref := t.DefinedAt()
		_, ok := reflist.Find(ref)
		if !ok && ref != fname {
			reflist.Append(ref)
		}
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				for _, lib := range impl.Graph().Libraries() {
					ref := lib.Filename()
					_, ok := reflist.Find(ref)
					if !ok && ref != fname {
						reflist.Append(ref)
					}
				}
			}
		}
	}
	for _, ref := range reflist.Strings() {
		ret.Libraries = append(ret.Libraries, *CreateXmlLibraryRef(ref))
	}
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
		ret.Libraries = append(ret.Libraries, *CreateXmlLibraryRef(l.Filename()))
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

func CreateXmlLibraryRef(ref string) *backend.XmlLibraryRef {
	return backend.XmlLibraryRefNew(ref)
}
