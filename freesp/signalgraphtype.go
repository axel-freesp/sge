package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
	//"encoding/xml"
)

var nodeTypes map[string]*nodeType
var portTypes map[string]*portType
var libraries map[string]*library

func Init() {
	nodeTypes = make(map[string]*nodeType)
	portTypes = make(map[string]*portType)
	libraries = make(map[string]*library)
}

type signalGraphType struct {
	name                                            string
	libraries                                       []Library
	signalTypes                                     []SignalType
	nodes, inputNodes, outputNodes, processingNodes []Node
}

func newSignalGraphType(name string) *signalGraphType {
	return &signalGraphType{name, nil, nil, nil, nil, nil, nil}
}

func (t *signalGraphType) Name() string {
	return t.name
}

func (t *signalGraphType) Nodes() []Node {
	return t.nodes
}

func (t *signalGraphType) SignalTypes() []SignalType {
	return t.signalTypes
}

func (t *signalGraphType) NodeByName(name string) Node {
	for _, n := range t.nodes {
		if n.NodeName() == name {
			return n
		}
	}
	return nil
}

func (t *signalGraphType) Libraries() []Library {
	return t.libraries
}

func (t *signalGraphType) InputNodes() []Node {
	return t.inputNodes
}

func (t *signalGraphType) OutputNodes() []Node {
	return t.outputNodes
}

func (t *signalGraphType) ProcessingNodes() []Node {
	return t.processingNodes
}

func createSignalGraphTypeFromXml(g *backend.XmlSignalGraph, name string, resolvePort func(portname string, dir PortDirection) *namedPortType) (t *signalGraphType, err error) {
	t = newSignalGraphType(name)
	for _, ref := range g.Libraries {
		l := libraries[ref.Name]
		if l == nil {
			l = LibraryNew(ref.Name)
			libraries[ref.Name] = l
			fmt.Println("createSignalGraphTypeFromXml: loading library", ref.Name)
			for _, try := range backend.XmlSearchPaths() {
				err = l.ReadFile(fmt.Sprintf("%s/%s", try, ref.Name))
				if err == nil {
					break
				}
			}
			if err != nil {
				err = newSignalGraphError(fmt.Sprintf("signalGraph.Read: referenced library file %s not found", ref.Name))
				return
			}
			fmt.Println("createSignalGraphTypeFromXml: library", ref.Name, "successfully loaded")
		}
		t.libraries = append(t.libraries, l)
	}
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
		t.signalTypes = append(t.signalTypes, sType)
	}
	for _, n := range g.InputNodes {
		nnode := t.createInputNodeFromXml(n, resolvePort)
		t.inputNodes = append(t.inputNodes, nnode)
		t.nodes = append(t.nodes, nnode)
	}
	for _, n := range g.OutputNodes {
		nnode := t.createOutputNodeFromXml(n, resolvePort)
		t.outputNodes = append(t.outputNodes, nnode)
		t.nodes = append(t.nodes, nnode)
	}
	for _, n := range g.ProcessingNodes {
		nnode := t.createNodeFromXml(n.XmlNode)
		t.processingNodes = append(t.processingNodes, nnode)
		t.nodes = append(t.nodes, nnode)
	}
	for i, c := range g.Connections {
		n1 := t.NodeByName(c.From)
		if n1 == nil {
			dump, _ := g.Write()
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.From, dump))
			//return nil, newSignalGraphError(fmt.Sprintf("invalid edge %d: node %s not found", i, c.From))
		}
		n2 := t.NodeByName(c.To)
		if n2 == nil {
			dump, _ := g.Write()
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.To, dump))
		}
		p1, err := n1.(*node).outPortFromName(c.FromPort)
		if err != nil {
			dump, _ := g.Write()
			log.Fatal(fmt.Sprintf("invalid edge %d from: %s\n%s", i, err, dump))
		}
		p2, err := n2.(*node).inPortFromName(c.ToPort)
		if err != nil {
			dump, _ := g.Write()
			log.Fatal(fmt.Sprintf("invalid edge %d to: %s\n%s", i, err, dump))
		}
		err = PortConnect(p1, p2)
		if err != nil {
			dump, _ := g.Write()
			log.Fatal(fmt.Sprintf("invalid edge %d: %s\n%s", i, err, dump))
		}
	}
	return
}

func createNodeTypeName(n backend.XmlNode) string {
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = fmt.Sprintf("autoTypeOfNode-%s", n.NName)
	}
	return ntName
}

func createInputNodeTypeName(name string) string {
	return fmt.Sprintf("autoInputNodeType-%s", name)
}

func createOutputNodeTypeName(name string) string {
	return fmt.Sprintf("autoOutputNodeType-%s", name)
}

func getPortType(name string) *portType {
	pt := portTypes[name]
	if pt == nil {
		pt = newPortType(name)
		portTypes[name] = pt
	}
	return pt
}

func (t *signalGraphType) createNodeFromXml(n backend.XmlNode) *node {
	nName := n.NName
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = createNodeTypeName(n)
	}
	nt := nodeTypes[ntName]
	if nt == nil {
		nt = createNodeTypeFromXmlNode(n, ntName)
	}
	ret := newNode(nName, nt, t)
	for _, p := range nt.InPorts() {
		ret.addInPort(p.(*namedPortType))
	}
	for _, p := range nt.OutPorts() {
		ret.addOutPort(p.(*namedPortType))
	}
	return ret
}

func (t *signalGraphType) createInputNodeFromXml(n backend.XmlInputNode, resolvePort func(portname string, dir PortDirection) *namedPortType) *node {
	nName := n.NName
	ntName := createInputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	ret := newNode(nName, nt, t)
	for _, p := range nt.OutPorts() {
		ret.addOutPort(p.(*namedPortType))
	}
	pt := resolvePort(n.NPort, InPort)
	if pt != nil {
		ret.addInPort(pt)
		ret.addOutPort(pt)
	}
	return ret
}

func (t *signalGraphType) createOutputNodeFromXml(n backend.XmlOutputNode, resolvePort func(portname string, dir PortDirection) *namedPortType) *node {
	nName := n.NName
	ntName := createInputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	ret := newNode(nName, nt, t)
	for _, p := range nt.InPorts() {
		ret.addInPort(p.(*namedPortType))
	}
	pt := resolvePort(n.NPort, OutPort) // matches also empty names
	if pt != nil {
		ret.addInPort(pt)
		ret.addOutPort(pt)
	}
	return ret
}
