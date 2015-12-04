package freesp

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/tool"
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

func createNodeTypeName(n *XmlNode) string {
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

func createNodeTypeFromXml(n *XmlNode, ntName string) *nodeType {
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

func (s *signalGraph) createNodeFromXml(n *XmlNode) *node {
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
	g := newXmlSignalGraph()
	err := g.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.Read: %v", err))
	}
	s.itsType = newSignalGraphType()
	sgType := s.itsType.(*signalGraphType)
	for _, st := range g.SignalTypes {
		var scope Scope
		switch st.Scope {
		case "local":
			scope = Local
		default:
			scope = Global
		}
		sType := newSignalType(st.Name, st.Ctype, st.Msgid, scope)
		sgType.signalTypes = append(sgType.signalTypes, sType)
	}
	for _, n := range g.InputNodes {
		nnode := s.createNodeFromXml(&n)
		sgType.inputNodes = append(sgType.inputNodes, nnode)
		sgType.nodes = append(sgType.nodes, nnode)
	}
	for _, n := range g.OutputNodes {
		nnode := s.createNodeFromXml(&n)
		sgType.outputNodes = append(sgType.outputNodes, nnode)
		sgType.nodes = append(sgType.nodes, nnode)
	}
	for _, n := range g.ProcessingNodes {
		nnode := s.createNodeFromXml(&n)
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

///////////////////////////////////////

type XmlLibrary struct {
	Name string `xml:"ref,attr"`
}

type XmlPort struct {
	PName string `xml:"port,attr"`
	PType string `xml:"type,attr"`
}

type XmlNode struct {
	NName   string    `xml:"name,attr"`
	NType   string    `xml:"type,attr"`
	InPort  []XmlPort `xml:"intype"`
	OutPort []XmlPort `xml:"outtype"`
	//gHint   Hint   `xml:"graph-hint"`
}

type XmlConnect struct {
	From     string `xml:"from,attr"`
	To       string `xml:"to,attr"`
	FromPort string `xml:"from-port,attr"`
	ToPort   string `xml:"to-port,attr"`
}

type XmlSignalType struct {
	Name  string `xml:"name,attr"`
	Scope string `xml:"scope,attr"`
	Ctype string `xml:"c-type,attr"`
	Msgid string `xml:"message-id,attr"`
}

type XmlSignalGraph struct {
	XMLName         xml.Name        `xml:"http://www.freesp.de/xml/freeSP signal-graph"`
	Version         string          `xml:"version,attr"`
	Libraries       []XmlLibrary    `xml:"library"`
	SignalTypes     []XmlSignalType `xml:"signal-type"`
	InputNodes      []XmlNode       `xml:"nodes>input"`
	OutputNodes     []XmlNode       `xml:"nodes>output"`
	ProcessingNodes []XmlNode       `xml:"nodes>processing-node"`
	Connections     []XmlConnect    `xml:"connections>connect"`
}

func (g *XmlSignalGraph) Read(data []byte) error {
	err := xml.Unmarshal(data, g)
	if err != nil {
		fmt.Printf("SignalGraph.Read error: %v", err)
	}
	return err
}

func (g *XmlSignalGraph) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		fmt.Println("signalgraph.ReadFile error: Failed to read file", filepath)
		return err
	}
	err = g.Read(data)
	if err != nil {
		fmt.Printf("signalgraph.ReadFile error: %v", err)
	}
	return err
}

func newXmlSignalGraph() *XmlSignalGraph {
	return &XmlSignalGraph{xml.Name{"", ""}, "", nil, nil, nil, nil, nil, nil}
}
