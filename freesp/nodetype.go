package freesp

import (
	"github.com/axel-freesp/sge/backend"
	"log"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts namedPortTypeList
	implementation    implementationList
	instances         nodeList
}

var _ NodeType = (*nodeType)(nil)

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, namedPortTypeListInit(),
		namedPortTypeListInit(), implementationListInit(), nodeListInit()}
}

// TODO:
/*
 *
 * Information about TreeCursor is gone in case of redo...
 *
 * Re-think: insertion of deleted objects via AddNewObject - sucks with
 * automatic update of instances and implementation when creating new objects...
 *
 * TODO: carefully maintain list offsets in type definitions (~ #implementation)
 * and nodes (#typedef = 1)
 *
 * */
func (t *nodeType) AddNamedPortType(p NamedPortType) {
	pt := p.(*namedPortType)
	if p.Direction() == InPort {
		t.inPorts.Append(pt)
	} else {
		t.outPorts.Append(pt)
	}
	for _, n := range t.instances.Nodes() {
		if p.Direction() == InPort {
			n.(*node).addInPort(p.(*namedPortType))
		} else {
			n.(*node).addOutPort(p.(*namedPortType))
		}
	}
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == NodeTypeGraph {
			if p.Direction() == InPort {
				impl.Graph().(*signalGraphType).addInputNodeFromNamedPortType(p)
			} else {
				impl.Graph().(*signalGraphType).addOutputNodeFromNamedPortType(p)
			}
		}
	}
}

func (t *nodeType) RemoveNamedPortType(p NamedPortType) {
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == NodeTypeGraph {
			if p.Direction() == InPort {
				impl.Graph().(*signalGraphType).removeInputNodeFromNamedPortType(p)
			} else {
				impl.Graph().(*signalGraphType).removeOutputNodeFromNamedPortType(p)
			}
		}
	}
	for _, n := range t.instances.Nodes() {
		n.(*node).removePort(p.(*namedPortType))
	}
	var list *namedPortTypeList
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
		for _, p := range t.inPorts.NamedPortTypes() {
			g.addInputNodeFromNamedPortType(p)
		}
		for _, p := range t.outPorts.NamedPortTypes() {
			g.addOutputNodeFromNamedPortType(p)
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

func (t *nodeType) InPorts() []NamedPortType {
	return t.inPorts.NamedPortTypes()
}

func (t *nodeType) OutPorts() []NamedPortType {
	return t.outPorts.NamedPortTypes()
}

func (t *nodeType) Implementation() []Implementation {
	return t.implementation.Implementations()
}

func createNodeTypeFromXmlNode(n backend.XmlNode, ntName string) *nodeType {
	nt := NodeTypeNew(ntName, "")
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	nodeTypes[ntName] = nt
	return nt
}

// TODO: These are possibly redundant..
func (t *nodeType) addInPort(name string, pType PortType) {
	t.inPorts.Append(&namedPortType{name, pType.(*portType), InPort})
}

func (t *nodeType) addOutPort(name string, pType PortType) {
	t.outPorts.Append(&namedPortType{name, pType.(*portType), OutPort})
}

func (t *nodeType) doResolvePort(name string, dir PortDirection) *namedPortType {
	var list namedPortTypeList
	switch dir {
	case InPort:
		list = t.inPorts
	default:
		list = t.outPorts
	}
	for _, p := range list.NamedPortTypes() {
		if name == p.Name() {
			return p.(*namedPortType)
		}
	}
	return nil
}

func createNodeTypeFromXml(n backend.XmlNodeType, filename string) *nodeType {
	nt := NodeTypeNew(n.TypeName, filename)
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
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
				var resolvePort = func(name string, dir PortDirection) *namedPortType {
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
		newCursor = tree.Insert(cursor)
		obj.(Implementation).AddToTree(tree, newCursor)

	case NamedPortType:
		pt := obj.(NamedPortType)
		if pt.Direction() == InPort {
			cursor.Position = len(t.InPorts())
		}
		newCursor = tree.Insert(cursor)
		pt.AddToTree(tree, newCursor)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				g := impl.Graph().(*signalGraphType)
				var n Node
				index := -len(t.Implementation())
				if pt.Direction() == InPort {
					n = g.findInputNodeFromNamedPortType(pt)
					if cursor.Position == AppendCursor {
						index += len(g.InputNodes())
					}
				} else {
					n = g.findOutputNodeFromNamedPortType(pt)
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
	case NamedPortType:
		pt := obj.(NamedPortType)
		// update all instance nodes in the tree
		for _, n := range t.Instances() {
			var p Port
			var ok bool
			if pt.Direction() == InPort {
				p, ok = n.(*node).inPort.Find(pt.Name())
			} else {
				p, ok = n.(*node).outPort.Find(pt.Name())
			}
			if !ok {
				log.Fatalf("nodeType.AddNewObject error: port %s not found.\n", pt.Name())
			}
			nCursor := tree.Cursor(n)
			// Insert new port at the same position as in the type:
			nCursor.Position = cursor.Position
			newNCursor := tree.Insert(nCursor)
			p.AddToTree(tree, newNCursor)
			tCursor := tree.CursorAt(nCursor, t)
			tCursor.Position = cursor.Position
			t.treeNewObject(tree, tCursor, obj)
		}

	default:
		log.Fatalf("nodeType.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (t *nodeType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Implementation:
		t.AddImplementation(obj.(Implementation))

	case NamedPortType:
		pt := obj.(NamedPortType)
		t.AddNamedPortType(pt) // adds ports of all instances, linked io-nodes in graph implementation

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
	log.Printf("nodeType) treeRemoveObject cursor=%v, obj=%v\n", cursor, obj)
	switch obj.(type) {
	case Implementation:
	case NamedPortType:
		nt := obj.(NamedPortType)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				g := impl.Graph().(*signalGraphType)
				var n Node
				if nt.Direction() == InPort {
					n = g.findInputNodeFromNamedPortType(nt)
				} else {
					n = g.findOutputNodeFromNamedPortType(nt)
				}
				if n == nil {
					log.Fatalf("nodeType.RemoveObject error: invalid implementation...\n")
				}
				nCursor := tree.CursorAt(parentId, n)
				for _, p := range n.InPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for _, c := range p.Connections() {
						conn := p.(*port).Connection(c.(*port))
						log.Println("nodeType.treeRemoveObject: saving connection", conn)
						removed = append(removed, IdWithObject{pCursor.Path, -1, conn})
					}
				}
				for _, p := range n.OutPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for _, c := range p.Connections() {
						conn := p.(*port).Connection(c.(*port))
						log.Println("nodeType.treeRemoveObject: saving connection", conn)
						removed = append(removed, IdWithObject{pCursor.Path, -1, conn})
					}
				}
				tree.Remove(nCursor)
				//removed = append(removed, IdWithObject{prefix, index, n})
			}
		}

	default:
		log.Fatalf("nodeType.RemoveObject error: invalid type %T\n", obj)
	}
	return
}

func (t *nodeType) treeRemoveInstObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	log.Printf("nodeType) treeRemoveInstObject cursor=%v, obj=%v\n", cursor, obj)
	switch obj.(type) {
	case Implementation:
	case NamedPortType:
		nt := obj.(NamedPortType)
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
				if p.PortName() == nt.Name() {
					break
				}
			}
			_ = p.(*port)
			log.Printf("nodeType.RemoveObject p=%v - nt.Name()=%s\n", p, nt.Name())
			pCursor := tree.CursorAt(nCursor, p)
			//del := n.RemoveObject(tree, pCursor)
			prefix, index := tree.Remove(pCursor)
			removed = append(removed, IdWithObject{prefix, index, p})
			//for _, d := range del {
			//    removed = append(removed, d)
			//}
			tCursor := tree.CursorAt(nCursor, obj)
			//tCursor.Position = cursor.Position
			del := t.treeRemoveObject(tree, tCursor)
			for _, d := range del {
				removed = append(removed, d)
			}
			tree.Remove(tCursor)
			//removed = append(removed, IdWithObject{prefix, index, obj})
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
	log.Printf("nodeType) RemoveObject cursor=%v, obj=%v\n", cursor, obj)
	switch obj.(type) {
	case Implementation:
		// TODO: This is redundant with implementation.go
		// Simply remove all nodes?
		impl := obj.(Implementation)
		// Return all removed edges and nodes
		for _, n := range impl.Graph().InputNodes() {
			nCursor := tree.Cursor(n)
			for _, p := range n.OutPorts() {
				pCursor := tree.CursorAt(nCursor, p)
				for index, c := range p.Connections() {
					conn := Connection{p, c}
					removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				}
			}
			// Removed Input- and Output nodes are NOT stored (they are
			// created automatically when adding the implementation graph).
		}
		for _, n := range impl.Graph().ProcessingNodes() {
			nCursor := tree.Cursor(n)
			for _, p := range n.OutPorts() {
				pCursor := tree.CursorAt(nCursor, p)
				for index, c := range p.Connections() {
					conn := Connection{p, c}
					removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				}
			}
			nIndex := IndexOfNodeInGraph(tree, n)
			removed = append(removed, IdWithObject{nCursor.Path, nIndex, n})
		}
		// Return removed object
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		// Remove obj in freesp model
		t.RemoveImplementation(impl)

	case NamedPortType:
		nt := obj.(NamedPortType)
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
