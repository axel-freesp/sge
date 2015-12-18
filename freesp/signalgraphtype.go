package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

var signalTypes map[string]*signalType
var nodeTypes map[string]*nodeType
var portTypes map[string]*portType
var libraries map[string]*library
var registeredNodeTypes []string
var registeredSignalTypes []string

func Init() {
	signalTypes = make(map[string]*signalType)
	nodeTypes = make(map[string]*nodeType)
	portTypes = make(map[string]*portType)
	libraries = make(map[string]*library)
}

type signalGraphType struct {
	libraries                                       []Library
	nodes, inputNodes, outputNodes, processingNodes []Node
}

/*
 *  freesp.SignalGraphType API
 */

var _ SignalGraphType = (*signalGraphType)(nil)

func SignalGraphTypeNew() *signalGraphType {
	return &signalGraphType{nil, nil, nil, nil, nil}
}

func (t *signalGraphType) Nodes() []Node {
	return t.nodes
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

func (t *signalGraphType) AddNode(n Node) error {
	nType := n.ItsType()
	libname := nType.DefinedAt()
	if len(libname) == 0 {
		return fmt.Errorf("signalGraphType.AddNode error: node type %s has no DefinedAt...", nType.TypeName())
	}
	if !t.containsLibRef(libname) {
		lib := libraries[libname]
		if lib == nil {
			return fmt.Errorf("signalGraphType.AddNode error: library %s not registered", libname)
		}
		t.libraries = append(t.libraries, lib)
	}
	return t.addNode(n)
}

func (t *signalGraphType) RemoveNode(n Node) {
	for _, p := range n.(*node).inPort {
		for _, c := range p.(*port).connections {
			c.RemoveConnection(p)
		}
	}
	RemNode(&t.nodes, n.(*node))
	RemNode(&t.inputNodes, n.(*node))
	RemNode(&t.outputNodes, n.(*node))
	RemNode(&t.processingNodes, n.(*node))
}

func (t *signalGraphType) containsLibRef(libname string) bool {
	for _, l := range t.libraries {
		if l.Filename() == libname {
			return true
		}
	}
	return false
}

func FindNode(list []Node, elem Node) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemNode(list *[]Node, elem Node) {
	index, ok := FindNode(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *signalGraphType) addNode(n Node) error {
	if len(n.InPorts()) > 0 {
		if len(n.OutPorts()) > 0 {
			t.processingNodes = append(t.processingNodes, n.(*node))
		} else {
			t.outputNodes = append(t.outputNodes, n.(*node))
		}
	} else {
		if len(n.OutPorts()) > 0 {
			t.inputNodes = append(t.inputNodes, n.(*node))
		} else {
			return fmt.Errorf("signalGraphType.AddNode error: node has no ports")
		}
	}
	t.nodes = append(t.nodes, n.(*node))
	return nil
}

func createSignalGraphTypeFromXml(g *backend.XmlSignalGraph, name string, resolvePort func(portname string, dir PortDirection) *namedPortType) (t *signalGraphType, err error) {
	t = SignalGraphTypeNew()
	for _, ref := range g.Libraries {
		l := libraries[ref.Name]
		if l == nil {
			l = LibraryNew(ref.Name)
			libraries[ref.Name] = l
			fmt.Println("createSignalGraphTypeFromXml: loading library", ref.Name)
			for _, try := range backend.XmlSearchPaths() {
				fmt.Printf("createSignalGraphTypeFromXml: try %s/%s\n", try, ref.Name)
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
		st := signalTypes[name]
		if st == nil {
			log.Fatal("getPortType: signalType '", name, "' is not defined")
		}
		pt = PortTypeNew(name, st)
		portTypes[name] = pt
	}
	return pt
}

func (t *signalGraphType) createNodeFromXml(n backend.XmlNode) (nd *node) {
	nName := n.NName
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = createNodeTypeName(n)
	}
	nt := nodeTypes[ntName]
	if nt == nil {
		nt = createNodeTypeFromXmlNode(n, ntName)
	}
	nd = NodeNew(nName, nt, t)
	return
}

func (t *signalGraphType) createInputNodeFromXml(n backend.XmlInputNode, resolvePort func(portname string, dir PortDirection) *namedPortType) *node {
	nName := n.NName
	ntName := createInputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	ret := NodeNew(nName, nt, t)
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
	ret := NodeNew(nName, nt, t)
	pt := resolvePort(n.NPort, OutPort) // matches also empty names
	if pt != nil {
		ret.addInPort(pt)
		ret.addOutPort(pt)
	}
	return ret
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*signalGraphType)(nil)

func (t *signalGraphType) AddToTree(tree Tree, cursor Cursor) {
	for _, n := range t.InputNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
	for _, n := range t.OutputNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
	for _, n := range t.ProcessingNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
}

func (t *signalGraphType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Printf("SignalGraphType.AddNewObject %T, %v\n", obj, obj)
	switch obj.(type) {
	case Node:
		n := obj.(Node)
		err := t.AddNode(n)
		if err != nil {
			log.Fatal("SignalGraphType.AddNewObject error: %s", err)
		}
		newCursor = tree.Insert(cursor)
		n.AddToTree(tree, newCursor)

	default:
		log.Fatal("SignalGraphType.AddNewObject error: wrong type %t: %v", obj, obj)
	}
	return
}

func (t *signalGraphType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Println("SignalGraphType.RemoveObject ", cursor)
	parent := tree.Parent(cursor)
	if t != tree.Object(parent) {
		log.Fatal("SignalGraphType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Node:
		n := obj.(Node)
		log.Println("SignalGraphType.RemoveObject remove node ", n)
		// Remove all connections first
		for _, p := range n.OutPorts() {
			//pCursor := tree.CursorAt(cursor, p)
			for _, c := range p.Connections() {
				conn := Connection{p, c}
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				//removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		for _, p := range n.InPorts() {
			//pCursor := tree.CursorAt(cursor, p)
			for _, c := range p.Connections() {
				conn := Connection{c, p}
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				//removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		t.RemoveNode(n)

	default:
		log.Fatal("SignalGraphType.RemoveObject error: wrong type %t: %v", obj, obj)
	}
	return
}
