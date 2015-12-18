package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type SignalGraphType struct {
	freesp.SignalGraphType
}

var _ TreeElement = SignalGraphType{}

func (t SignalGraphType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	for _, n := range t.InputNodes() {
		child := tree.Append(cursor)
		Node{n}.AddToTree(tree, child)
	}
	for _, n := range t.OutputNodes() {
		child := tree.Append(cursor)
		Node{n}.AddToTree(tree, child)
	}
	for _, n := range t.ProcessingNodes() {
		child := tree.Append(cursor)
		Node{n}.AddToTree(tree, child)
	}
}

func (t SignalGraphType) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Printf("SignalGraphType.AddNewObject %T, %v\n", obj, obj)
	switch obj.(type) {
	case freesp.Node:
		n := obj.(freesp.Node)
		err := t.AddNode(n)
		if err != nil {
			log.Fatal("SignalGraphType.AddNewObject error: %s", err)
		}
		newCursor = tree.Insert(cursor)
		Node{n}.AddToTree(tree, newCursor)

	default:
		log.Fatal("SignalGraphType.AddNewObject error: wrong type %t: %v", obj, obj)
	}
	return
}

func (t SignalGraphType) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Println("SignalGraphType.RemoveObject ", cursor)
	parent := tree.Parent(cursor)
	if t.SignalGraphType != tree.Object(parent) {
		log.Fatal("SignalGraphType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case freesp.Node:
		n := obj.(freesp.Node)
		log.Println("SignalGraphType.RemoveObject remove node ", n)
		// Remove all connections first
		for _, p := range n.OutPorts() {
			//pCursor := tree.CursorAt(cursor, p)
			for _, c := range p.Connections() {
				conn := freesp.Connection{p, c}
				cCursor := tree.CursorAt(cursor, conn)
				del := Port{p}.RemoveObject(tree, cCursor)
				//removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		for _, p := range n.InPorts() {
			//pCursor := tree.CursorAt(cursor, p)
			for _, c := range p.Connections() {
				conn := freesp.Connection{c, p}
				cCursor := tree.CursorAt(cursor, conn)
				del := Port{p}.RemoveObject(tree, cCursor)
				//removed = append(removed, IdWithObject{pCursor.Path, index, conn})
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		t.RemoveNode(n)

	default:
		log.Fatal("SignalGraphType.RemoveObject error: wrong type %t: %v", obj, obj)
	}
	return
}

func init_signalgraphtype(iconPath string) (err error) {
	return
}
