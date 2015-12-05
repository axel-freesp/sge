package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/backend"
)

var nodeTypes map[string]*nodeType
var portTypes map[string]*portType

func SignalGraphInit() {
	nodeTypes = make(map[string]*nodeType)
	portTypes = make(map[string]*portType)
}

func SignalGraphNew() *signalGraph {
	return &signalGraph{nil}
}

type signalGraph struct {
	itsType SignalGraphType
}

func (s *signalGraph) ItsType() SignalGraphType {
	return s.itsType
}

func createNodeTypeName(n backend.XmlNode) string {
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = fmt.Sprintf("autoTypeOfNode-%s", n.NName)
	}
	return ntName
}

func getPortType(name string) *portType {
	pt := portTypes[name]
	if pt == nil {
		pt = newPortType(name)
		portTypes[name] = pt
	}
	return pt
}

func createNodeTypeFromXml(n backend.XmlNode, ntName string) *nodeType {
	nt := newNodeType(ntName)
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	// TODO: evaluate <implementation>
	return nt
}

func (s *signalGraph) createNodeFromXml(n backend.XmlNode) *node {
	nName := n.NName
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = createNodeTypeName(n)
	}
	nt := nodeTypes[ntName]
	if nt == nil {
		nt = createNodeTypeFromXml(n, ntName)
		nodeTypes[ntName] = nt
	}
	ret := newNode(nName, nt, s)
	for _, p := range nt.InPorts() {
		ret.addInPort(p.(*namedPortType))
	}
	for _, p := range nt.OutPorts() {
		ret.addOutPort(p.(*namedPortType))
	}
	return ret
}

func (s *signalGraph) Read(data []byte) error {
	g := backend.NewXmlSignalGraph()
	err := g.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.Read: %v", err))
	}
	s.itsType = newSignalGraphType()
	sgType := s.itsType.(*signalGraphType)
	for _, st := range g.SignalTypes {
		var scope Scope
		var mode Mode
		switch st.Scope {
		case "local":
			scope = Local
		default:
			scope = Global
		}
		switch st.Mode {
		case "sync":
			mode = Synchronous
		default:
			mode = Asynchronous
		}
		sType := newSignalType(st.Name, st.Ctype, st.Msgid, scope, mode)
		sgType.signalTypes = append(sgType.signalTypes, sType)
	}
	for _, n := range g.InputNodes {
		nnode := s.createNodeFromXml(n.XmlNode)
		sgType.inputNodes = append(sgType.inputNodes, nnode)
		sgType.nodes = append(sgType.nodes, nnode)
	}
	for _, n := range g.OutputNodes {
		nnode := s.createNodeFromXml(n.XmlNode)
		sgType.outputNodes = append(sgType.outputNodes, nnode)
		sgType.nodes = append(sgType.nodes, nnode)
	}
	for _, n := range g.ProcessingNodes {
		nnode := s.createNodeFromXml(n.XmlNode)
		sgType.processingNodes = append(sgType.processingNodes, nnode)
		sgType.nodes = append(sgType.nodes, nnode)
	}
	for i, c := range g.Connections {
		n1 := s.itsType.NodeByName(c.From)
		if n1 == nil {
			return newSignalGraphError(fmt.Sprintf("invalid edge %d: node %s not found", i, c.From))
		}
		n2 := s.itsType.NodeByName(c.To)
		if n2 == nil {
			return newSignalGraphError(fmt.Sprintf("invalid edge %d: node %s not found", i, c.To))
		}
		p1, err := n1.(*node).outPortFromName(c.FromPort)
		if err != nil {
			return newSignalGraphError(fmt.Sprintf("invalid edge %d from: %s", i, err))
		}
		p2, err := n2.(*node).inPortFromName(c.ToPort)
		if err != nil {
			return newSignalGraphError(fmt.Sprintf("invalid edge %d to: %s", i, err))
		}
		err = PortConnect(p1, p2)
		if err != nil {
			return newSignalGraphError(fmt.Sprintf("invalid edge %d: %s", i, err))
		}
	}
	return nil
}

func (s *signalGraph) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.ReadFile: %v", err))
	}
	err = s.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalgraph.ReadFile: %v", err))
	}
	return err
}

func (s *signalGraph) Write() (data []byte, err error) {
	// TODO
	data = nil
	err = newSignalGraphError("Write() interface not implemented")
	return
}

func (s *signalGraph) WriteFile(filepath string) error {
	// TODO
	return newSignalGraphError("WriteFile() interface not implemented")
}

//------------------------------

type signalGraphError struct {
	reason string
}

func (e *signalGraphError) Error() string {
	return fmt.Sprintf("signal graph error: %s", e.reason)
}

func newSignalGraphError(reason string) *signalGraphError {
	return &signalGraphError{reason}
}

