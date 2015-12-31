package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"image"
	"log"
)

type signalGraphType struct {
	context                                         Context
	libraries                                       []Library
	nodes, inputNodes, outputNodes, processingNodes []Node
}

/*
 *  freesp.SignalGraphType API
 */

var _ SignalGraphType = (*signalGraphType)(nil)

func SignalGraphTypeNew(context Context) *signalGraphType {
	return &signalGraphType{context, nil, nil, nil, nil, nil}
}

func (t *signalGraphType) Nodes() []Node {
	return t.nodes
}

func (t *signalGraphType) NodeByName(name string) Node {
	for _, n := range t.nodes {
		if n.Name() == name {
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
	for _, p := range n.(*node).inPort.Ports() {
		for _, c := range p.Connections() {
			c.RemoveConnection(p)
		}
	}
	RemNode(&t.nodes, n.(*node))
	RemNode(&t.inputNodes, n.(*node))
	RemNode(&t.outputNodes, n.(*node))
	RemNode(&t.processingNodes, n.(*node))
	n.ItsType().(*nodeType).removeInstance(n.(*node))
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

func createSignalGraphTypeFromXml(g *backend.XmlSignalGraph, name string, context Context,
	resolvePort func(portname string, dir PortDirection) *portType) (t *signalGraphType, err error) {
	t = SignalGraphTypeNew(context)
	for _, ref := range g.Libraries {
		l := libraries[ref.Name]
		if l == nil {
			var lib Library
			lib, err = t.context.GetLibrary(ref.Name)
			l = lib.(*library)
			/*
				l = LibraryNew(ref.Name)
				fmt.Println("createSignalGraphTypeFromXml: loading library", ref.Name)
				for _, try := range backend.XmlSearchPaths() {
					fmt.Printf("createSignalGraphTypeFromXml: try %s/%s\n", try, ref.Name)
					err = l.ReadFile(fmt.Sprintf("%s/%s", try, ref.Name))
					if err == nil {
						break
					}
				}
			*/
			if err != nil {
				err = newSignalGraphError(fmt.Sprintf("signalGraph.Read: referenced library file %s not found", ref.Name))
				return
			}
			libraries[ref.Name] = l
			fmt.Println("createSignalGraphTypeFromXml: library", ref.Name, "successfully loaded")
		}
		t.libraries = append(t.libraries, l)
	}
	for _, n := range g.InputNodes {
		var nnode *node
		nnode, err = t.createInputNodeFromXml(n, resolvePort)
		if err != nil {
			return
		}
		t.inputNodes = append(t.inputNodes, nnode)
		t.nodes = append(t.nodes, nnode)
	}
	for _, n := range g.OutputNodes {
		var nnode *node
		nnode, err = t.createOutputNodeFromXml(n, resolvePort)
		if err != nil {
			return
		}
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
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.From, dump))
		}
		n2 := t.NodeByName(c.To)
		if n2 == nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.To, dump))
		}
		p1, err := n1.(*node).outPortFromName(c.FromPort)
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Printf("edge = %v\n", c)
			log.Printf("node = %v, missing port = %s\n", n1, c.FromPort)
			log.Fatal(fmt.Sprintf("invalid edge %d outPortFromName failed: %s\n%s", i, err, dump))
		}
		p2, err := n2.(*node).inPortFromName(c.ToPort)
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d inPortFromName failed: %s\n%s", i, err, dump))
		}
		err = p1.AddConnection(ConnectionNew(p1, p2))
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d PortConnect failed: %s\n%s", i, err, dump))
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
	var err error
	nd, err = NodeNew(nName, nt, t)
	if err != nil {
		log.Fatal("signalGraphType.createNodeFromXml: TODO: error handling")
	}
	nd.position = image.Point{n.Hint.X, n.Hint.Y}
	return
}

func (t *signalGraphType) createInputNodeFromXml(n backend.XmlInputNode,
	resolvePort func(portname string, dir PortDirection) *portType) (ret *node, err error) {
	nName := n.NName
	ntName := createInputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	pt := resolvePort(n.NPort, InPort)
	if pt != nil {
		if len(nt.OutPorts()) == 0 {
			ptCopy := &portType{pt.signalType, "", OutPort}
			nt.AddNamedPortType(ptCopy)
		}
	}
	ret, err = NodeNew(nName, nt, t)
	if err != nil {
		err = fmt.Errorf("signalGraphType.createInputNodeFromXml: %s", err)
		return
	}
	if pt != nil {
		// add link to pt
	}
	ret.position = image.Point{n.Hint.X, n.Hint.Y}
	return
}

func (t *signalGraphType) createOutputNodeFromXml(n backend.XmlOutputNode,
	resolvePort func(portname string, dir PortDirection) *portType) (ret *node, err error) {
	nName := n.NName
	ntName := createOutputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	pt := resolvePort(n.NPort, OutPort) // matches also empty names
	if pt != nil {
		if len(nt.InPorts()) == 0 {
			ptCopy := &portType{pt.signalType, "", InPort}
			nt.AddNamedPortType(ptCopy)
		}
	}
	ret, err = NodeNew(nName, nt, t)
	if err != nil {
		err = fmt.Errorf("signalGraphType.createOutputNodeFromXml: %s", err)
		return
	}
	if pt != nil {
		// add link to pt
	}
	ret.position = image.Point{n.Hint.X, n.Hint.Y}
	return
}

func (g *signalGraphType) addInputNodeFromPortType(p PortType) {
	st := p.SignalType()
	ntName := createInputNodeTypeName(st.TypeName())
	nt, ok := nodeTypes[ntName]
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.addOutPort("", st)
		nodeTypes[ntName] = nt
	}
	if len(nt.outPorts.PortTypes()) == 0 {
		log.Fatal("signalGraphType.addInputNodeFromNamedPortType: invalid input node type")
	}
	n, err := NodeNew(fmt.Sprintf("in-%s", p.Name()), nt, g)
	if err != nil {
		log.Fatal("signalGraphType.addInputNodeFromPortType: TODO: error handling")
	}
	n.portlink = p
	err = g.addNode(n)
	if err != nil {
		log.Fatal("signalGraphType.addInputNodeFromNamedPortType: AddNode failed:", err)
	}
}

func (g *signalGraphType) addOutputNodeFromPortType(p PortType) {
	st := p.SignalType()
	ntName := createOutputNodeTypeName(st.TypeName())
	nt, ok := nodeTypes[ntName]
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.addInPort("", st)
		nodeTypes[ntName] = nt
	}
	if len(nt.inPorts.PortTypes()) == 0 {
		log.Fatal("signalGraphType.addOutputNodeFromNamedPortType: invalid output node type")
	}
	n, err := NodeNew(fmt.Sprintf("out-%s", p.Name()), nt, g)
	if err != nil {
		log.Fatal("signalGraphType.addOutputNodeFromPortType: TODO: error handling")
	}
	n.portlink = p
	err = g.addNode(n)
	if err != nil {
		log.Fatal("signalGraphType.addOutputNodeFromNamedPortType: AddNode failed:", err)
	}
}

func (g *signalGraphType) removeInputNodeFromPortType(p PortType) {
	for _, n := range g.InputNodes() {
		if n.Name() == fmt.Sprintf("in-%s", p.Name()) {
			g.RemoveNode(n)
			return
		}
	}
}

func (g *signalGraphType) removeOutputNodeFromPortType(p PortType) {
	for _, n := range g.OutputNodes() {
		if n.Name() == fmt.Sprintf("out-%s", p.Name()) {
			g.RemoveNode(n)
			return
		}
	}
}

func (g *signalGraphType) findInputNodeFromPortType(p PortType) Node {
	for _, n := range g.InputNodes() {
		if n.Name() == fmt.Sprintf("in-%s", p.Name()) {
			return n
		}
	}
	return nil
}

func (g *signalGraphType) findOutputNodeFromPortType(p PortType) Node {
	for _, n := range g.OutputNodes() {
		if n.Name() == fmt.Sprintf("out-%s", p.Name()) {
			return n
		}
	}
	return nil
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

func (t *signalGraphType) treeAddNewObject(tree Tree, cursor Cursor, n Node) (newCursor Cursor) {
	newCursor = tree.Insert(cursor)
	n.AddToTree(tree, newCursor)
	return
}

func (t *signalGraphType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	switch obj.(type) {
	case Node:
		n := obj.(Node)
		err = t.AddNode(n)
		if err != nil {
			err = fmt.Errorf("signalGraphType.AddNewObject error: %s", err)
			nt := n.ItsType().(*nodeType)
			if nt != nil {
				ok, _ := nt.instances.Find(n)
				if ok {
					nt.instances.Remove(n)
				}
			}
			return
		}
		newCursor = t.treeAddNewObject(tree, cursor, n)

		parent := tree.Object(cursor)
		switch parent.(type) {
		case SignalGraph:
		case Implementation:
			// propagate new node to all instances of embracing type
			pCursor := tree.Parent(cursor)
			nt := tree.Object(pCursor)
			for _, nn := range nt.(NodeType).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, parent)
				tCursor.Position = cursor.Position
				t.treeAddNewObject(tree, tCursor, n)
			}

		default:
			log.Fatalf("signalGraphType.AddNewObject error: wrong parent type %T: %v\n", parent, parent)
		}

	default:
		log.Fatalf("signalGraphType.AddNewObject error: wrong type %t: %v\n", obj, obj)
	}
	return
}

func (t *signalGraphType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Node:
		n := obj.(Node)
		// Remove all connections first
		for _, p := range n.OutPorts() {
			for _, c := range p.Connections() {
				conn := p.Connection(c)
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		for _, p := range n.InPorts() {
			for _, c := range p.Connections() {
				conn := p.Connection(c)
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		parentCursor := tree.Parent(cursor)
		parent := tree.Object(parentCursor)
		switch parent.(type) {
		case SignalGraph:
		case Implementation:
			// propagate new node to all instances of embracing type
			pCursor := tree.Parent(parentCursor)
			nt := tree.Object(pCursor)
			for _, nn := range nt.(NodeType).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, parent)
				tree.Remove(tree.CursorAt(tCursor, n))
			}

		default:
			log.Fatalf("signalGraphType.RemoveObject error: wrong parent type %t: %v\n", parent, parent)
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		t.RemoveNode(n)

	default:
		log.Fatalf("signalGraphType.RemoveObject error: wrong type %t: %v", obj, obj)
	}
	return
}
