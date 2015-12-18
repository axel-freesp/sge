package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type Library struct {
	freesp.Library
}

var _ TreeElement = Library{}

func (l Library) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolLibrary, l.Filename(), l.Library)
	if err != nil {
		log.Fatal("Library.AddToTree error: AddEntry failed: %s", err)
	}
	for _, t := range l.SignalTypes() {
		child := tree.Append(cursor)
		SignalType{t}.AddToTree(tree, child)
	}
	for _, t := range l.NodeTypes() {
		child := tree.Append(cursor)
		NodeType{t}.AddToTree(tree, child)
	}
}

func (l Library) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	switch obj.(type) {
	case freesp.SignalType:
		t := obj.(freesp.SignalType)
		err := l.AddSignalType(t)
		if err != nil {
			log.Fatal("Library.AddNewObject error: AddSignalType failed: %s", err)
		}
		newCursor = tree.Insert(cursor)
		SignalType{t}.AddToTree(tree, newCursor)

	case freesp.NodeType:
		t := obj.(freesp.NodeType)
		err := l.AddNodeType(t)
		if err != nil {
			log.Fatal("Library.AddNewObject error: AddNodeType failed: %s", err)
		}
		newCursor = tree.Insert(cursor)
		NodeType{t}.AddToTree(tree, newCursor)

	default:
		log.Fatal("Library.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (l Library) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	// TODO
	return
}
