package freesp

import (
	"log"
	"github.com/axel-freesp/sge/backend"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	gr "github.com/axel-freesp/sge/interface/graph"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts portTypeList
	implementation    implementationList
	instances         nodeList
}

var _ bh.NodeTypeIf = (*nodeType)(nil)

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, portTypeListInit(),
		portTypeListInit(), implementationListInit(), nodeListInit()}
}

func (t *nodeType) AddNamedPortType(p bh.PortTypeIf) {
	if p.Direction() == gr.InPort {
		t.inPorts.Append(p)
	} else {
		t.outPorts.Append(p)
	}
	for _, n := range t.instances.Nodes() {
		if p.Direction() == gr.InPort {
			n.(*node).addInPort(p)
		} else {
			n.(*node).addOutPort(p)
		}
	}
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			if p.Direction() == gr.InPort {
				impl.Graph().(*signalGraphType).addInputNodeFromPortType(p)
			} else {
				impl.Graph().(*signalGraphType).addOutputNodeFromPortType(p)
			}
		}
	}
}

func (t *nodeType) RemoveNamedPortType(p bh.PortTypeIf) {
	for _, impl := range t.implementation.Implementations() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			if p.Direction() == gr.InPort {
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
	if p.Direction() == gr.InPort {
		list = &t.inPorts
	} else {
		list = &t.outPorts
	}
	list.Remove(p)
}

func (t *nodeType) Instances() []bh.NodeIf {
	return t.instances.Nodes()
}

func (t *nodeType) addInstance(n *node) {
	t.instances.Append(n)
}

func (t *nodeType) removeInstance(n *node) {
	t.instances.Remove(n)
}

func (t *nodeType) RemoveImplementation(imp bh.ImplementationIf) {
	if imp.ImplementationType() == bh.NodeTypeGraph {
		gt := imp.Graph()
		for len(gt.Nodes()) > 0 {
			gt.RemoveNode(gt.Nodes()[0].(*node))
		}
	}
	t.implementation.Remove(imp)
}

func (t *nodeType) AddImplementation(imp bh.ImplementationIf) {
	if imp.ImplementationType() == bh.NodeTypeGraph {
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

func (t *nodeType) SetTypeName(newTypeName string) {
	t.name = newTypeName
}

func (t *nodeType) DefinedAt() string {
	return t.definedAt
}

func (t *nodeType) InPorts() []bh.PortTypeIf {
	return t.inPorts.PortTypes()
}

func (t *nodeType) OutPorts() []bh.PortTypeIf {
	return t.outPorts.PortTypes()
}

func (t *nodeType) Implementation() []bh.ImplementationIf {
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

func (t *nodeType) CreateXml() (buf []byte, err error) {
	xmlnodetype := CreateXmlNodeType(t)
	buf, err = xmlnodetype.Write()
	return
}

// TODO: These are possibly redundant..
func (t *nodeType) addInPort(name string, pType bh.SignalTypeIf) {
	t.inPorts.Append(PortTypeNew(name, pType.TypeName(), gr.InPort))
}

func (t *nodeType) addOutPort(name string, pType bh.SignalTypeIf) {
	t.outPorts.Append(PortTypeNew(name, pType.TypeName(), gr.OutPort))
}

func (t *nodeType) doResolvePort(name string, dir gr.PortDirection) *portType {
	var list portTypeList
	switch dir {
	case gr.InPort:
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

func createNodeTypeFromXml(n backend.XmlNodeType, filename string, context mod.ModelContextIf) *nodeType {
	nt := NodeTypeNew(n.TypeName, filename)
	for _, p := range n.InPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXml error: signal type '%s' not found\n", p.PType)
		}
		nt.addInPort(p.PName, pType)
	}
	for _, p := range n.OutPort {
		pType, ok := signalTypes[p.PType]
		if !ok {
			log.Fatalf("createNodeTypeFromXml error: signal type '%s' not found\n", p.PType)
		}
		nt.addOutPort(p.PName, pType)
	}
	for _, i := range n.Implementation {
		var iType bh.ImplementationType
		if len(i.SignalGraph) == 1 {
			iType = bh.NodeTypeGraph
		} else {
			iType = bh.NodeTypeElement
		}
		impl := ImplementationNew(i.Name, iType, context)
		nt.implementation.Append(impl)
		switch iType {
		case bh.NodeTypeElement:
			impl.elementName = i.Name
		default:
			var err error
			var resolvePort = func(name string, dir gr.PortDirection) *portType {
				return nt.doResolvePort(name, dir)
			}
			impl.graph, err = createSignalGraphTypeFromXml(&i.SignalGraph[0], n.TypeName, context, resolvePort)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nt
}

/*
 *  tr.TreeElement API
 */

var _ tr.TreeElement = (*nodeType)(nil)

func (t *nodeType) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	parent := tree.Object(parentId)
	switch parent.(type) {
	case bh.LibraryIf:
		prop = MayAddObject | MayEdit | MayRemove
	case bh.NodeIf:
		prop = 0
	default:
		log.Fatalf("nodeType.AddToTree error: invalid parent type %T\n", parent)
	}
	err := tree.AddEntry(cursor, tr.SymbolNodeType, t.TypeName(), t, prop)
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

func (t *nodeType) treeNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor) {
	switch obj.(type) {
	case bh.ImplementationIf:
		cursor.Position = len(t.Implementation()) - 1
		newCursor = tree.Insert(cursor)
		obj.(bh.ImplementationIf).AddToTree(tree, newCursor)

	case bh.PortTypeIf:
		pt := obj.(bh.PortTypeIf)
		newCursor = tree.Insert(cursor)
		pt.AddToTree(tree, newCursor)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				// bh.NodeIf linked to outer port
				g := impl.Graph().(*signalGraphType)
				var n bh.NodeIf
				index := -len(t.Implementation())
				if pt.Direction() == gr.InPort {
					n = g.findInputNodeFromPortType(pt)
					if cursor.Position == tr.AppendCursor {
						index += len(g.InputNodes())
					}
				} else {
					n = g.findOutputNodeFromPortType(pt)
					if cursor.Position == tr.AppendCursor {
						index += len(g.InputNodes()) + len(g.OutputNodes())
					}
				}
				if n == nil {
					log.Fatalf("nodeType.AddNewObject error: invalid implementation...\n")
				}
				if cursor.Position != tr.AppendCursor {
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

func (t *nodeType) treeInstObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor) {
	switch obj.(type) {
	case bh.ImplementationIf:
		impl := obj.(bh.ImplementationIf)
		// update all instance nodes in the tree with new implementation
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			tCursor := tree.CursorAt(nCursor, t)
			tCursor.Position = len(t.Implementation()) - 1
			newICursor := tree.Insert(tCursor)
			impl.AddToTree(tree, newICursor)
		}

	case bh.PortTypeIf:
		pt := obj.(bh.PortTypeIf)
		// update all instance nodes in the tree with new port
		for _, n := range t.Instances() {
			var p bh.PortIf
			var ok bool
			if pt.Direction() == gr.InPort {
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

func (t *nodeType) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	switch obj.(type) {
	case bh.ImplementationIf:
		t.AddImplementation(obj.(bh.ImplementationIf))
		cursor.Position = len(t.Implementation()) - 1

	case bh.PortTypeIf:
		pt := obj.(bh.PortTypeIf)
		t.AddNamedPortType(pt) // adds ports of all instances, linked io-nodes in graph implementation
		// Adjust offset to insert:
		if pt.Direction() == gr.InPort {
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

func (t *nodeType) treeRemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.ImplementationIf:
		impl := obj.(bh.ImplementationIf)
		if impl.ImplementationType() == bh.NodeTypeGraph {
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
						removed = append(removed, tr.IdWithObject{pCursor.Path, index, conn})
					}
				}
			}
			// ... and processing nodes
			for _, n := range impl.Graph().ProcessingNodes() {
				nCursor := tree.Cursor(n)
				gCursor := tree.Parent(nCursor)
				nIndex := gCursor.Position
				removed = append(removed, tr.IdWithObject{nCursor.Path, nIndex, n})
			}
		}

	case bh.PortTypeIf:
		nt := obj.(bh.PortTypeIf)
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				// Remove and store all edges connected to the nodes linked to the outer ports
				g := impl.Graph().(*signalGraphType)
				var n bh.NodeIf
				if nt.Direction() == gr.InPort {
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
						removed = append(removed, tr.IdWithObject{pCursor.Path, -1, conn})
					}
				}
				for _, p := range n.OutPorts() {
					pCursor := tree.CursorAt(nCursor, p)
					for _, c := range p.Connections() {
						conn := p.Connection(c)
						removed = append(removed, tr.IdWithObject{pCursor.Path, -1, conn})
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
func (t *nodeType) treeRemoveInstObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parentId := tree.Parent(cursor)
	if t != tree.Object(parentId) {
		log.Fatal("nodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.ImplementationIf:
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			tCursor := tree.CursorAt(nCursor, t)
			iCursor := tree.CursorAt(tCursor, obj)
			iCursor.Position = cursor.Position
			tree.Remove(iCursor)
		}

	case bh.PortTypeIf:
		nt := obj.(bh.PortTypeIf)
		for _, n := range t.Instances() {
			var p bh.PortIf
			var list []bh.PortIf
			nCursor := tree.Cursor(n)
			if nt.Direction() == gr.InPort {
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
			removed = append(removed, tr.IdWithObject{prefix, index, p})
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

func (t *nodeType) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
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
	removed = append(removed, tr.IdWithObject{prefix, index, obj})

	switch obj.(type) {
	case bh.ImplementationIf:
		impl := obj.(bh.ImplementationIf)
		// Remove obj in freesp model
		t.RemoveImplementation(impl)

	case bh.PortTypeIf:
		nt := obj.(bh.PortTypeIf)
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
	nodeTypes []bh.NodeTypeIf
}

func nodeTypeListInit() nodeTypeList {
	return nodeTypeList{nil}
}

func (l *nodeTypeList) Append(nt bh.NodeTypeIf) {
	l.nodeTypes = append(l.nodeTypes, nt)
}

func (l *nodeTypeList) Remove(nt bh.NodeTypeIf) {
	var i int
	for i = range l.nodeTypes {
		if nt == l.nodeTypes[i] {
			break
		}
	}
	if i >= len(l.nodeTypes) {
		for _, v := range l.nodeTypes {
			log.Printf("nodeTypeList.RemoveNodeType have bh.NodeTypeIf %v\n", v)
		}
		log.Fatalf("nodeTypeList.RemoveNodeType error: bh.NodeTypeIf %v not in this list\n", nt)
	}
	for i++; i < len(l.nodeTypes); i++ {
		l.nodeTypes[i-1] = l.nodeTypes[i]
	}
	l.nodeTypes = l.nodeTypes[:len(l.nodeTypes)-1]
}

func (l *nodeTypeList) NodeTypes() []bh.NodeTypeIf {
	return l.nodeTypes
}
