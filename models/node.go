package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type Node struct {
	freesp.Node
}

var _ TreeElement = Node{}

func (n Node) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	var image Symbol
	if len(n.InPorts()) == 0 {
		image = SymbolOutputNode
	} else if len(n.OutPorts()) == 0 {
		image = SymbolInputNode
	} else {
		image = SymbolProcessingNode
	}
	err := tree.AddEntry(cursor, image, n.NodeName(), n.Node)
	if err != nil {
		log.Fatal("Node.AddToTree error: AddEntry failed: %s", err)
	}
	child := tree.Append(cursor)
	NodeType{n.ItsType()}.AddToTree(tree, child)
	for _, p := range n.InPorts() {
		child := tree.Append(cursor)
		Port{p}.AddToTree(tree, child)
	}
	for _, p := range n.OutPorts() {
		child := tree.Append(cursor)
		Port{p}.AddToTree(tree, child)
	}
}

func (n Node) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	switch obj.(type) {
	case freesp.NamedPortType:
		// If this function is called, the freesp model if different from the tree model
		// cursor points to n, cursor.Position may indicate where to insert.
		// Add missing port entry to the tree, matching name and direction with obj.
		newCursor = tree.Insert(cursor)
		pt := obj.(freesp.NamedPortType)
		var list []freesp.Port
		if pt.Direction() == freesp.InPort {
			list = n.InPorts()
		} else {
			list = n.OutPorts()
		}
		for _, p := range list {
			if p.PortName() == pt.Name() {
				Port{p}.AddToTree(tree, newCursor)
				break
			}
		}
	default:
		log.Fatal("Node.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (n Node) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	return
}

func IndexOfNodeInGraph(tree *FilesTreeStore, n freesp.Node) (index int) {
	nCursor := tree.Cursor(n)
	gCursor := tree.Parent(nCursor)
	index = gCursor.Position
	return
}

func IsProcessingNode(n freesp.Node) bool {
	if len(n.InPorts()) == 0 {
		return false
	}
	if len(n.OutPorts()) == 0 {
		return false
	}
	return true
}
