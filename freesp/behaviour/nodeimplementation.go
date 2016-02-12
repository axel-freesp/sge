package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
)

type implementation struct {
	implementationType bh.ImplementationType
	elementName        string
	graph              bh.SignalGraphTypeIf
}

var _ bh.ImplementationIf = (*implementation)(nil)

func ImplementationNew(iName string, iType bh.ImplementationType, context mod.ModelContextIf) *implementation {
	ret := &implementation{iType, iName, nil}
	if iType == bh.NodeTypeGraph {
		ret.graph = SignalGraphTypeNew(context)
	}
	return ret
}

func (n *implementation) ImplementationType() bh.ImplementationType {
	return n.implementationType
}

func (n *implementation) ElementName() string {
	return n.elementName
}

func (n *implementation) SetElemName(newName string) {
	n.elementName = newName
}

func (n *implementation) Graph() bh.SignalGraphTypeIf {
	return n.graph
}

func (n *implementation) CreateXml() (buf []byte, err error) {
	switch n.ImplementationType() {
	case bh.NodeTypeElement:
		// TODO
	default:
		xmlImpl := CreateXmlSignalGraphType(n.Graph())
		buf, err = xmlImpl.Write()
	}
	return
}

/*
 *  fmt.Stringer API
 */

func (n *implementation) String() string {
	if n.implementationType == bh.NodeTypeGraph {
		return fmt.Sprintf("bh.ImplementationIf graph {\n%v\n}", n.graph)
	} else {
		return fmt.Sprintf("bh.ImplementationIf module %s", n.elementName)
	}

}

/*
 *  tr.TreeElementIf API
 */

var _ tr.TreeElementIf = (*implementation)(nil)

func (impl *implementation) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var image tr.Symbol
	var text string
	var prop tr.Property
	parentId := tree.Parent(tree.Parent(cursor))
	parent := tree.Object(parentId)
	switch parent.(type) {
	case bh.LibraryIf:
		if impl.ImplementationType() == bh.NodeTypeGraph {
			prop = freesp.PropertyNew(true, false, true)
		} else {
			prop = freesp.PropertyNew(true, true, true)
		}
	case bh.NodeIf:
		prop = freesp.PropertyNew(false, false, false)
	default:
		log.Fatalf("implementation.AddToTree error: invalid parent type: %T\n", parent)
	}
	if impl.ImplementationType() == bh.NodeTypeGraph {
		image = tr.SymbolImplGraph
		text = "Graph"
	} else {
		image = tr.SymbolImplElement
		text = impl.ElementName()
	}
	err := tree.AddEntry(cursor, image, text, impl, prop)
	if err != nil {
		log.Fatalf("implementation.AddToTree error: AddEntry failed: %s\n", err)
	}
	if impl.ImplementationType() == bh.NodeTypeGraph {
		impl.Graph().AddToTree(tree, cursor)
	}
}

func (impl *implementation) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	switch obj.(type) {
	case bh.NodeIf:
		if impl.ImplementationType() == bh.NodeTypeGraph {
			return impl.Graph().AddNewObject(tree, cursor, obj)
		} else {
			log.Fatalf("implementation.AddNewObject error: cannot add node to elementary implementation.\n")
		}

	case bh.ConnectionIf:
		if impl.ImplementationType() == bh.NodeTypeGraph {
			return impl.Graph().AddNewObject(tree, cursor, obj)
		} else {
			log.Fatalf("implementation.AddNewObject error: cannot add connection to elementary implementation.\n")
		}

	default:
		log.Fatalf("implementation.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (impl *implementation) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if impl != tree.Object(parent) {
		log.Fatal("implementation.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.NodeIf:
		if impl.ImplementationType() == bh.NodeTypeGraph {
			log.Println("implementation.RemoveObject: delegate to signalGraphType\n")
			return impl.Graph().RemoveObject(tree, cursor)
		} else {
			log.Fatalf("implementation.RemoveObject error: cannot remove node from elementary implementation.\n")
		}
		return

		n := obj.(bh.NodeIf)
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
				conn := p.Connection(c)
				removed = append(removed, tr.IdWithObject{pCursor.Path, index, conn})
			}
		}
		for _, p := range n.InPorts() {
			pCursor := tree.CursorAt(cursor, p)
			for index, c := range p.Connections() {
				conn := p.Connection(c)
				removed = append(removed, tr.IdWithObject{pCursor.Path, index, conn})
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, obj})
		impl.Graph().RemoveNode(n)

	default:
		log.Fatalf("implementation.RemoveObject error: invalid type %T", obj)
	}
	return
}

/*
 *      implementationList
 *
 */

type implementationList struct {
	implementations []bh.ImplementationIf
}

func implementationListInit() implementationList {
	return implementationList{nil}
}

func (l *implementationList) Append(nt bh.ImplementationIf) {
	l.implementations = append(l.implementations, nt)
}

func (l *implementationList) Remove(nt bh.ImplementationIf) {
	var i int
	for i = range l.implementations {
		if nt == l.implementations[i] {
			break
		}
	}
	if i >= len(l.implementations) {
		for _, v := range l.implementations {
			log.Printf("implementationList.RemoveImplementation have bh.ImplementationIf %v\n", v)
		}
		log.Fatalf("implementationList.RemoveImplementation error: bh.ImplementationIf %v not in this list\n", nt)
	}
	for i++; i < len(l.implementations); i++ {
		l.implementations[i-1] = l.implementations[i]
	}
	l.implementations = l.implementations[:len(l.implementations)-1]
}

func (l *implementationList) Implementations() []bh.ImplementationIf {
	return l.implementations
}
