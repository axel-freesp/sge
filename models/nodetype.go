package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type NodeType struct {
	freesp.NodeType
}

var _ TreeElement = NodeType{}

func (t NodeType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, imageNodeType, t.TypeName(), t.NodeType)
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
	for _, impl := range t.Implementation() {
		child := tree.Append(cursor)
		Implementation{impl}.AddToTree(tree, child)
	}
	for _, pt := range t.InPorts() {
		child := tree.Append(cursor)
		NamedPortType{pt}.AddToTree(tree, child)
	}
	for _, pt := range t.OutPorts() {
		child := tree.Append(cursor)
		NamedPortType{pt}.AddToTree(tree, child)
	}
}

func (t NodeType) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	switch obj.(type) {
	case freesp.Implementation:
		t.AddImplementation(obj.(freesp.Implementation))
		newCursor = tree.Insert(cursor)
		Implementation{obj.(freesp.Implementation)}.AddToTree(tree, newCursor)

	case freesp.NamedPortType:
		pt := obj.(freesp.NamedPortType)
		t.AddNamedPortType(pt)
		newCursor = tree.Insert(cursor)
		NamedPortType{pt}.AddToTree(tree, newCursor)
		// update all instance nodes in the tree
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			// Insert new port at the same position as in the type:
			nCursor.Position = cursor.Position
			Node{n}.AddNewObject(tree, nCursor, obj)
		}

	default:
		log.Fatal("NodeType.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (t NodeType) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if t.NodeType != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case freesp.Implementation:
		// TODO: This is redundant with implementation.go
		// Simply remove all nodes?
		impl := obj.(freesp.Implementation)
		// Return all removed edges and nodes
		for _, n := range impl.Graph().InputNodes() {
			nCursor := tree.Cursor(n)
			for _, p := range n.OutPorts() {
				pCursor := tree.CursorAt(nCursor, p)
				for index, c := range p.Connections() {
					conn := freesp.Connection{p, c}
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
					conn := freesp.Connection{p, c}
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

	case freesp.NamedPortType:
		nt := obj.(NamedPortType)
		// TODO: This is redundant with node.go
		// Simply remove port of all nodes?
		for _, n := range t.Instances() {
			nCursor := tree.Cursor(n)
			if nt.Direction() == freesp.InPort {
				for _, p := range n.InPorts() {
					if p.PortName() == nt.Name() {
						pCursor := tree.CursorAt(nCursor, p)
						for index, c := range p.Connections() {
							conn := freesp.Connection{c, p}
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
							conn := freesp.Connection{p, c}
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

// Images:

var (
	imageNodeType *gdk.Pixbuf = nil
	imageFilename string      = "node-type.png"
)

func init_nodetype(iconPath string) (err error) {
	imageNodeType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/%s", iconPath, imageFilename))
	if err != nil {
		err = fmt.Errorf("init_signaltype error loading %s: %s", imageFilename, err)
	}
	return
}
