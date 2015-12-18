package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type Implementation struct {
	freesp.Implementation
}

var _ TreeElement = Implementation{}

func (impl Implementation) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	var image *gdk.Pixbuf
	var text string
	if impl.ImplementationType() == freesp.NodeTypeGraph {
		image = imageImplGraph
		text = "Graph"
	} else {
		image = imageImplElement
		text = impl.ElementName()
	}
	err := tree.AddEntry(cursor, image, text, impl.Implementation)
	if err != nil {
		log.Fatal("Implementation.AddToTree error: AddEntry failed: %s", err)
	}
	if impl.ImplementationType() == freesp.NodeTypeGraph {
		SignalGraphType{impl.Graph()}.AddToTree(tree, cursor)
	}
}

func (impl Implementation) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	switch obj.(type) {
	case freesp.Node:
		err := impl.Graph().AddNode(obj.(freesp.Node))
		if err != nil {
			log.Fatal("Implementation.AddNewObject error: AddNode failed: ", err)
		}
		newCursor = tree.Insert(cursor)
		Node{obj.(freesp.Node)}.AddToTree(tree, cursor)

	default:
		log.Fatal("Implementation.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (impl Implementation) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if impl.Implementation != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case freesp.Node:
		n := obj.(freesp.Node)
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
				conn := freesp.Connection{p, c}
				removed = append(removed, IdWithObject{pCursor.Path, index, conn})
			}
		}
		for _, p := range n.InPorts() {
			pCursor := tree.CursorAt(cursor, p)
			for index, c := range p.Connections() {
				conn := freesp.Connection{c, p}
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

var (
	imageImplElement *gdk.Pixbuf = nil
	imageImplGraph   *gdk.Pixbuf = nil
)

func init_implementation(iconPath string) (err error) {
	imageImplElement, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test0.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_signaltype error loading test0.png: %s", err)
		return
	}
	imageImplGraph, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test1.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_signaltype error loading test0.png: %s", err)
	}
	return
}
