package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts []NamedPortType
	implementation    []Implementation
	instances         []Node
}

var _ NodeType = (*nodeType)(nil)

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, nil, nil, nil, nil}
}

func (t *nodeType) AddNamedPortType(p NamedPortType) {
	// TODO: Check if type of node will change
	pt := p.(*namedPortType)
	if p.Direction() == InPort {
		t.inPorts = append(t.inPorts, pt)
	} else {
		t.outPorts = append(t.outPorts, pt)
	}
	for _, n := range t.instances {
		if p.Direction() == InPort {
			n.(*node).addInPort(p.(*namedPortType))
		} else {
			n.(*node).addOutPort(p.(*namedPortType))
		}
	}
}

func FindNamedPortType(list []NamedPortType, elem NamedPortType) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemNamedPortType(list *[]NamedPortType, elem NamedPortType) {
	index, ok := FindNamedPortType(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *nodeType) RemoveNamedPortType(p NamedPortType) {
	for _, n := range t.instances {
		n.(*node).removePort(p.(*namedPortType))
	}
	var list *[]NamedPortType
	if p.Direction() == InPort {
		list = &t.inPorts
	} else {
		list = &t.outPorts
	}
	RemNamedPortType(list, p)
}

func (t *nodeType) Instances() []Node {
	return t.instances
}

func (t *nodeType) addInstance(n *node) {
	t.instances = append(t.instances, n)
}

func (t *nodeType) addInPort(name string, pType PortType) {
	t.inPorts = append(t.inPorts, &namedPortType{name, pType.(*portType), InPort})
}

func (t *nodeType) addOutPort(name string, pType PortType) {
	t.outPorts = append(t.outPorts, &namedPortType{name, pType.(*portType), OutPort})
}

func FindImplementation(list []Implementation, elem Implementation) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemImplementation(list *[]Implementation, elem Implementation) {
	index, ok := FindImplementation(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *nodeType) RemoveImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		gt := imp.Graph()
		for len(gt.Nodes()) > 0 {
			gt.RemoveNode(gt.Nodes()[0].(*node))
		}
	}
	RemImplementation(&t.implementation, imp)
}

func (t *nodeType) AddImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		if imp.Graph() == nil {
			log.Fatal("nodeType.AddImplementation: missing graph")
		}
		g := imp.Graph().(*signalGraphType)
		for _, p := range t.inPorts {
			st := p.SignalType()
			ntName := createInputNodeTypeName(st.TypeName())
			nt, ok := nodeTypes[ntName]
			if !ok {
				nt = NodeTypeNew(ntName, "")
				nt.addOutPort("", getPortType(st.TypeName()))
				nodeTypes[ntName] = nt
			}
			if len(nt.outPorts) == 0 {
				log.Fatal("nodeType.AddImplementation: invalid input node type")
			}
			n := NodeNew(fmt.Sprintf("in-%s", p.Name()), nt, imp.Graph())
			n.portlink = p
			err := g.addNode(n)
			if err != nil {
				log.Fatal("nodeType.AddImplementation: AddNode failed:", err)
			}
		}
		for _, p := range t.outPorts {
			st := p.SignalType()
			ntName := createOutputNodeTypeName(st.TypeName())
			nt, ok := nodeTypes[ntName]
			if !ok {
				nt = NodeTypeNew(ntName, "")
				nt.addInPort("", getPortType(st.TypeName()))
				nodeTypes[ntName] = nt
			}
			if len(nt.inPorts) == 0 {
				log.Fatal("nodeType.AddImplementation: invalid input node type")
			}
			n := NodeNew(fmt.Sprintf("out-%s", p.Name()), nt, imp.Graph())
			n.portlink = p
			err := g.addNode(n)
			if err != nil {
				log.Fatal("nodeType.AddImplementation: AddNode failed:", err)
			}
		}
	}
	t.implementation = append(t.implementation, imp)
}

func (t *nodeType) TypeName() string {
	return t.name
}

func (t *nodeType) DefinedAt() string {
	return t.definedAt
}

func (t *nodeType) InPorts() []NamedPortType {
	return t.inPorts
}

func (t *nodeType) OutPorts() []NamedPortType {
	return t.outPorts
}

func (t *nodeType) Implementation() []Implementation {
	return t.implementation
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

func (t *nodeType) doResolvePort(name string, dir PortDirection) *namedPortType {
	var ports []NamedPortType
	switch dir {
	case InPort:
		ports = t.inPorts
	default:
		ports = t.outPorts
	}
	for _, p := range ports {
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
			nt.implementation = append(nt.implementation, impl)
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

var _ TreeElement = (*nodeType)(nil)

func (t *nodeType) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolNodeType, t.TypeName(), t)
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
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

func (t *nodeType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Implementation:
		t.AddImplementation(obj.(Implementation))
		newCursor = tree.Insert(cursor)
		obj.(Implementation).AddToTree(tree, newCursor)

	case NamedPortType:
		pt := obj.(NamedPortType)
		t.AddNamedPortType(pt)
		newCursor = tree.Insert(cursor)
		pt.AddToTree(tree, newCursor)
		// update all instance nodes in the tree
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			// Insert new port at the same position as in the type:
			nCursor.Position = cursor.Position
			n.AddNewObject(tree, nCursor, obj)
		}

	default:
		log.Fatal("NodeType.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (t *nodeType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if t != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
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
		// TODO: This is redundant with node.go
		// Simply remove port of all nodes?
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			if nt.Direction() == InPort {
				for _, p := range n.InPorts() {
					if p.PortName() == nt.Name() {
						pCursor := tree.CursorAt(nCursor, p)
						for index, c := range p.Connections() {
							conn := Connection{c, p}
							removed = append(removed, IdWithObject{pCursor.Path, index, conn})
						}
						break
					}
				}
			} else {
				for _, p := range n.OutPorts() {
					if p.PortName() == nt.Name() {
						pCursor := tree.CursorAt(nCursor, p)
						for index, c := range p.Connections() {
							conn := Connection{p, c}
							removed = append(removed, IdWithObject{pCursor.Path, index, conn})
						}
						break
					}
				}
			}
			nIndex := IndexOfNodeInGraph(tree, n)
			removed = append(removed, IdWithObject{nCursor.Path, nIndex, n})
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		t.RemoveNamedPortType(nt)

	default:
		log.Fatal("NodeType.RemoveObject error: invalid type %T", obj)
	}
	return
}
