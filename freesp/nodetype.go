package freesp

import (
	"github.com/axel-freesp/sge/backend"
	"log"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts portTypeList
	implementation    implementationList
	instances         nodeList
}

var _ NodeType = (*nodeType)(nil)

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, portTypeListInit(),
		portTypeListInit(), implementationListInit(), nodeListInit()}
}

func (t *nodeType) AddNamedPortType(p PortType) {
	if p.Direction() == InPort {
		t.inPorts.Append(p)
	} else {
		t.outPorts.Append(p)
	}
	for _, n := range t.instances.Nodes() {
		if p.Direction() == InPort {
			n.(*node).addInPort(p)
		} else {
			n.(*node).addOutPort(p)
		}
	}
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == NodeTypeGraph {
			if p.Direction() == InPort {
				impl.Graph().(*signalGraphType).addInputNodeFromPortType(p)
			} else {
				impl.Graph().(*signalGraphType).addOutputNodeFromPortType(p)
			}
		}
	}
}

func (t *nodeType) RemoveNamedPortType(p PortType) {
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == NodeTypeGraph {
			if p.Direction() == InPort {
				impl.Graph().(*signalGraphType).removeInputNodeFromPortType(p)
			} else {
				impl.Graph().(*signalGraphType).removeOutputNodeFromPortType(p)
			}
		}
	}
	for _, n := range t.instances.Nodes() {
		n.(*node).removePort(p.(*portType))
	}
	var list *portTypeList
	if p.Direction() == InPort {
		list = &t.inPorts
	} else {
		list = &t.outPorts
	}
	list.Remove(p)
}

func (t *nodeType) Instances() []Node {
	return t.instances.Nodes()
}

func (t *nodeType) addInstance(n *node) {
	t.instances.Append(n)
}

func (t *nodeType) removeInstance(n *node) {
	t.instances.Remove(n)
}

func (t *nodeType) RemoveImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		gt := imp.Graph()
		for len(gt.Nodes()) > 0 {
			gt.RemoveNode(gt.Nodes()[0].(*node))
		}
	}
	t.implementation.Remove(imp)
}

func (t *nodeType) AddImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		if imp.Graph() == nil {
			log.Fatal("nodeType.AddImplementation: missing graph")
		}
		g := imp.Graph().(*signalGraphType)
		for _, p := range t.inPorts.PortTypes() {
			g.addInputNodeFromPortType(p)
		}
		for _, p := range t.outPorts.PortTypes() {
			g.addOutputNodeFromPortType(p)
		}
	}
	t.implementation.Append(imp)
}

func (t *nodeType) TypeName() string {
	return t.name
}

func (t *nodeType) DefinedAt() string {
	return t.definedAt
}

func (t *nodeType) InPorts() []PortType {
	return t.inPorts.PortTypes()
}

func (t *nodeType) OutPorts() []PortType {
	return t.outPorts.PortTypes()
}

func (t *nodeType) Implementation() []Implementation {
	return t.implementation.Implementations()
}

func createNodeTypeFromXmlNode(n backend.XmlNode, ntName string) *nodeType {
	nt := NodeTypeNew(ntName, "")
	for _, p := range n.InPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXmlNode error: signal type '%s' not found\n", p.PType)
		}
		nt.addInPort(p.PName, pType)
	}
	for _, p := range n.OutPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXmlNode error: signal type '%s' not found\n", p.PType)
		}
		nt.addOutPort(p.PName, pType)
	}
	nodeTypes[ntName] = nt
	return nt
}

// TODO: These are possibly redundant..
func (t *nodeType) addInPort(name string, pType SignalType) {
	t.inPorts.Append(PortTypeNew(name, pType.TypeName(), InPort))
}

func (t *nodeType) addOutPort(name string, pType SignalType) {
	t.outPorts.Append(PortTypeNew(name, pType.TypeName(), OutPort))
}

func (t *nodeType) doResolvePort(name string, dir PortDirection) *portType {
	var list portTypeList
	switch dir {
	case InPort:
		list = t.inPorts
	default:
		list = t.outPorts
	}
	for _, p := range list.PortTypes() {
		if name == p.Name() {
			return p.(*portType)
		}
	}
	return nil
}

func createNodeTypeFromXml(n backend.XmlNodeType, filename string) *nodeType {
	nt := NodeTypeNew(n.TypeName, filename)
	for _, p := range n.InPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXmlNode error: signal type '%s' not found\n", p.PType)
		}
		nt.addInPort(p.PName, pType)
	}
	for _, p := range n.OutPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXmlNode error: signal type '%s' not found\n", p.PType)
		}
		nt.addOutPort(p.PName, pType)
	}
	if len(n.Implementation) > 0 {
		for _, i := range n.Implementation {
			var iType ImplementationType
			if len(i.SignalGraph) == 1 {
				iType = NodeTypeGraph
			} else {
				iType = NodeTypeElement
			}
			impl := ImplementationNew(i.Name, iType)
			nt.implementation.Append(impl)
			switch iType {
			case NodeTypeElement:
				impl.elementName = i.Name
			default:
				var err error
				var resolvePort = func(name string, dir PortDirection) *portType {
					return nt.doResolvePort(name, dir)
				}
				impl.graph, err = createSignalGraphTypeFromXml(&i.SignalGraph[0], n.TypeName, resolvePort)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nt
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*nodeType)(nil)

func (t *nodeType) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	parent := tree.Object(parentId)
	switch parent.(type) {
	case Library:
		prop = mayAddObject | mayEdit | mayRemove
	case Node:
		prop = 0
	default:
		log.Fatalf("nodeType.AddToTree error: invalid parent type %T\n", parent)
	}
	err := tree.AddEntry(cursor, SymbolNodeType, t.TypeName(), t, prop)
	if err != nil {
		log.Fatalf("nodeType.AddToTree error: AddEntry failed: %s\n", err)
	}
	for _, impl := range t.Implementation() {
		child := tree.Append(cursor)
		impl.AddToTree(tree, child)
	}
	for _, pt := range t.InPorts() {
		child := tree.Append(cursor)
		pt.AddToTree(tree, child)
	}
	for _, pt := range t.OutPorts() {
		child := tree.Append(cursor)
		pt.AddToTree(tree, child)
	}
}

func (t *nodeType) treeNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Implementation:
		cursor.Position = len(t.Implementation()) - 1
		newCursor = tree.Insert(cursor)
		obj.(Implementation).AddToTree(tree, newCursor)

	case PortType:
		pt := obj.(PortType)
		newCursor = tree.Insert(cursor)
		pt.AddToTree(tree, newCursor)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				// Node linked to outer port
				g := impl.Graph().(*signalGraphType)
				var n Node
				index := -len(t.Implementation())
				if pt.Direction() == InPort {
					n = g.findInputNodeFromPortType(pt)
					if cursor.Position == AppendCursor {
						index += len(g.InputNodes())
					}
				} else {
					n = g.findOutputNodeFromPortType(pt)
					if cursor.Position == AppendCursor {
						index += len(g.InputNodes()) + len(g.OutputNodes())
					}
				}
				if n == nil {
					log.Fatalf("nodeType.AddNewObject error: invalid implementation...\n")
				}
				if cursor.Position != AppendCursor {
					index += cursor.Position
				}
				gCursor := tree.CursorAt(cursor, impl)
				gCursor.Position = index
				n.AddToTree(tree, tree.Insert(gCursor))
			}
		}

	default:
		log.Fatalf("nodeType.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (t *nodeType) treeInstObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Implementation:
		impl := obj.(Implementation)
		// update all instance nodes in the tree with new implementation
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			tCursor := tree.CursorAt(nCursor, t)
			tCursor.Position = len(t.Implementation()) - 1
			newICursor := tree.Insert(tCursor)
			impl.AddToTree(tree, newICursor)
		}

	case PortType:
		pt := obj.(PortType)
		// update all instance nodes in the tree with new port
		for _, n := range t.Instances() {
			var p Port
			var ok bool
			if pt.Direction() == InPort {
				p, ok, _ = n.(*node).inPort.Find(n.Name(), pt.Name())
			} else {
				p, ok, _ = n.(*node).outPort.Find(n.Name(), pt.Name())
			}
			if !ok {
				log.Fatalf("nodeType.treeInstObject error: port %s not found.\n", pt.Name())
			}
			nCursor := tree.Cursor(n)
			// Insert new port at the same position as in the type:
			// Need to deal with implementations in type VS type in node
			if cursor.Position >= 0 {
				nCursor.Position = cursor.Position - len(t.Implementation()) + 1
			}
			newNCursor := tree.Insert(nCursor)
			p.AddToTree(tree, newNCursor)
			// Update mirrored type in node:
			tCursor := tree.CursorAt(nCursor, t)
			tCursor.Position = cursor.Position
			t.treeNewObject(tree, tCursor, obj)
		}

	default:
		log.Fatalf("nodeType.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (t *nodeType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	switch obj.(type) {
	case Implementation:
		t.AddImplementation(obj.(Implementation))
		cursor.Position = len(t.Implementation()) - 1

	case PortType:
		pt := obj.(PortType)
		t.AddNamedPortType(pt) // adds ports of all instances, linked io-nodes in graph implementation
		// Adjust offset to insert:
		if pt.Direction() == InPort {
			cursor.Position = len(t.Implementation()) + len(t.InPorts()) - 1
		} else {
			cursor.Position = -1
		}

	default:
		log.Fatalf("nodeType.AddNewObject error: invalid type %T\n", obj)
	}

	newCursor = t.treeNewObject(tree, cursor, obj)
	t.treeInstObject(tree, cursor, obj)
	return
}

func (t *nodeType) treeRemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Implementation:
		impl := obj.(Implementation)
		if impl.ImplementationType() == NodeTypeGraph {
			// TODO: This is redundant with implementation.go
			// Simply remove all nodes? Do not traverse a modifying list...
			// Removed Input- and Output nodes are NOT stored (they are
			// created automatically when adding the implementation graph).
			// Return all removed edges ...
			for _, n := range impl.Graph().Nodes() {
				nCursor := tree.Cursor(n)
				for _, p := range n.OutPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for index, c := range p.Connections() {
						conn := p.Connection(c)
						removed = append(removed, IdWithObject{pCursor.Path, index, conn})
					}
				}
			}
			// ... and processing nodes
			for _, n := range impl.Graph().ProcessingNodes() {
				nCursor := tree.Cursor(n)
				gCursor := tree.Parent(nCursor)
				nIndex := gCursor.Position
				removed = append(removed, IdWithObject{nCursor.Path, nIndex, n})
			}
		}

	case PortType:
		nt := obj.(PortType)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				// Remove and store all edges connected to the nodes linked to the outer ports
				g := impl.Graph().(*signalGraphType)
				var n Node
				if nt.Direction() == InPort {
					n = g.findInputNodeFromPortType(nt)
				} else {
					n = g.findOutputNodeFromPortType(nt)
				}
				if n == nil {
					log.Fatalf("nodeType.RemoveObject error: invalid implementation...\n")
				}
				nCursor := tree.CursorAt(parentId, n)
				for _, p := range n.InPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for _, c := range p.Connections() {
						conn := p.Connection(c)
						removed = append(removed, IdWithObject{pCursor.Path, -1, conn})
					}
				}
				for _, p := range n.OutPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for _, c := range p.Connections() {
						conn := p.Connection(c)
						removed = append(removed, IdWithObject{pCursor.Path, -1, conn})
					}
				}
				// Remove (but dont store) the nodes linked to the outer ports:
				tree.Remove(nCursor)
			}
		}

	default:
		log.Fatalf("nodeType.RemoveObject error: invalid type %T\n", obj)
	}
	return
}

// Remove object mirrored in all instance node type
func (t *nodeType) treeRemoveInstObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Implementation:
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			tCursor := tree.CursorAt(nCursor, t)
			iCursor := tree.CursorAt(tCursor, obj)
			iCursor.Position = cursor.Position
			tree.Remove(iCursor)
		}

	case PortType:
		nt := obj.(PortType)
		for _, n := range t.Instances() {
			var p Port
			var list []Port
			nCursor := tree.Cursor(n)
			if nt.Direction() == InPort {
				list = n.InPorts()
			} else {
				list = n.OutPorts()
			}
			for _, p = range list {
				if p.Name() == nt.Name() {
					break
				}
			}
			_ = p.(*port)
			pCursor := tree.CursorAt(nCursor, p)
			prefix, index := tree.Remove(pCursor)
			removed = append(removed, IdWithObject{prefix, index, p})
			tCursor := tree.CursorAt(nCursor, obj)
			del := t.treeRemoveObject(tree, tCursor)
			for _, d := range del {
				removed = append(removed, d)
			}
			tree.Remove(tCursor)
		}

	default:
		log.Fatalf("nodeType.RemoveObject error: invalid type %T\n", obj)
	}
	return
}

func (t *nodeType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	del := t.treeRemoveObject(tree, cursor)
	for _, d := range del {
		removed = append(removed, d)
	}
	del = t.treeRemoveInstObject(tree, cursor)
	for _, d := range del {
		removed = append(removed, d)
	}
	prefix, index := tree.Remove(cursor)
	removed = append(removed, IdWithObject{prefix, index, obj})

	switch obj.(type) {
	case Implementation:
		impl := obj.(Implementation)
		// Remove obj in freesp model
		t.RemoveImplementation(impl)

	case PortType:
		nt := obj.(PortType)
		t.RemoveNamedPortType(nt)

	default:
		log.Fatalf("nodeType.RemoveObject error: invalid type %T\n", obj)
	}
	return
}

/*
 *      nodeTypeList
 *
 */

type nodeTypeList struct {
	nodeTypes []NodeType
}

func nodeTypeListInit() nodeTypeList {
	return nodeTypeList{nil}
}

func (l *nodeTypeList) Append(nt NodeType) {
	l.nodeTypes = append(l.nodeTypes, nt)
}

func (l *nodeTypeList) Remove(nt NodeType) {
	var i int
	for i = range l.nodeTypes {
		if nt == l.nodeTypes[i] {
			break
		}
	}
	if i >= len(l.nodeTypes) {
		for _, v := range l.nodeTypes {
			log.Printf("nodeTypeList.RemoveNodeType have NodeType %v\n", v)
		}
		log.Fatalf("nodeTypeList.RemoveNodeType error: NodeType %v not in this list\n", nt)
	}
	for i++; i < len(l.nodeTypes); i++ {
		l.nodeTypes[i-1] = l.nodeTypes[i]
	}
	l.nodeTypes = l.nodeTypes[:len(l.nodeTypes)-1]
}

func (l *nodeTypeList) NodeTypes() []NodeType {
	return l.nodeTypes
}
