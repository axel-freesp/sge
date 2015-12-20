package freesp

import (
	"log"
)

type implementation struct {
	implementationType ImplementationType
	elementName        string
	graph              SignalGraphType
}

var _ Implementation = (*implementation)(nil)

func ImplementationNew(iName string, iType ImplementationType) *implementation {
	ret := &implementation{iType, iName, nil}
	if iType == NodeTypeGraph {
		ret.graph = SignalGraphTypeNew()
	}
	return ret
}

func (n *implementation) ImplementationType() ImplementationType {
	return n.implementationType
}

func (n *implementation) ElementName() string {
	return n.elementName
}

func (n *implementation) Graph() SignalGraphType {
	return n.graph
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*implementation)(nil)

func (impl *implementation) AddToTree(tree Tree, cursor Cursor) {
	var image Symbol
	var text string
	var prop property
	parentId := tree.Parent(tree.Parent(cursor))
	parent := tree.Object(parentId)
	switch parent.(type) {
	case Library:
		if impl.ImplementationType() == NodeTypeGraph {
			prop = mayAddObject | mayRemove
		} else {
			prop = mayEdit | mayRemove
		}
	case Node:
		prop = 0
	default:
		log.Fatalf("implementation.AddToTree error: invalid parent type: %T\n", parent)
	}
	if impl.ImplementationType() == NodeTypeGraph {
		image = SymbolImplGraph
		text = "Graph"
	} else {
		image = SymbolImplElement
		text = impl.ElementName()
	}
	err := tree.AddEntry(cursor, image, text, impl, prop)
	if err != nil {
		log.Fatalf("implementation.AddToTree error: AddEntry failed: %s\n", err)
	}
	if impl.ImplementationType() == NodeTypeGraph {
		impl.Graph().AddToTree(tree, cursor)
	}
}

func (impl *implementation) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Node:
		log.Println("implementation) AddNewObject")
		err := impl.Graph().AddNode(obj.(Node))
		if err != nil {
			log.Fatal("Implementation.AddNewObject error: AddNode failed: ", err)
		}
		newCursor = tree.Insert(cursor)
		obj.(Node).AddToTree(tree, cursor)

	default:
		log.Fatal("Implementation.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (impl *implementation) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if impl != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Node:
		n := obj.(Node)
		if !IsProcessingNode(n) {
			// Removed Input- and Output nodes are NOT stored (they are
			// created automatically when adding the implementation graph).
			// Within an implementation, it is therefore not allowed to
			// remove such IO nodes.
			return
		}
		// Remove all connections first
		for _, p := range n.OutPorts() {
			pCursor := tree.CursorAt(cursor, p)
			for index, c := range p.Connections() {
				conn := Connection{p, c}
				removed = append(removed, IdWithObject{pCursor.Path, index, conn})
			}
		}
		for _, p := range n.InPorts() {
			pCursor := tree.CursorAt(cursor, p)
			for index, c := range p.Connections() {
				conn := Connection{c, p}
				removed = append(removed, IdWithObject{pCursor.Path, index, conn})
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		impl.Graph().RemoveNode(n)

	default:
		log.Fatal("Implementation.RemoveObject error: invalid type %T", obj)
	}
	return
}

/*
 *      implementationList
 *
 */

type implementationList struct {
	implementations []Implementation
}

func implementationListInit() implementationList {
	return implementationList{nil}
}

func (l *implementationList) Append(nt Implementation) {
	l.implementations = append(l.implementations, nt)
}

func (l *implementationList) Remove(nt Implementation) {
	var i int
	for i = range l.implementations {
		if nt == l.implementations[i] {
			break
		}
	}
	if i >= len(l.implementations) {
		for _, v := range l.implementations {
			log.Printf("implementationList.RemoveImplementation have Implementation %v\n", v)
		}
		log.Fatalf("implementationList.RemoveImplementation error: Implementation %v not in this list\n", nt)
	}
	for i++; i < len(l.implementations); i++ {
		l.implementations[i-1] = l.implementations[i]
	}
	l.implementations = l.implementations[:len(l.implementations)-1]
}

func (l *implementationList) Implementations() []Implementation {
	return l.implementations
}
