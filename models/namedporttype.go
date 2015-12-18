package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type NamedPortType struct {
	freesp.NamedPortType
}

var _ TreeElement = NamedPortType{}

func (p NamedPortType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	var kind Symbol
	if p.Direction() == freesp.InPort {
		kind = SymbolInputPortType
	} else {
		kind = SymbolOutputPortType
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p.NamedPortType)
	if err != nil {
		log.Fatal("NamedPortType.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	SignalType{p.SignalType()}.AddToTree(tree, child)
}

func (p NamedPortType) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Fatal("NamedPortType.AddNewObject - nothing to add.")
	return
}

func (p NamedPortType) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("NamedPortType.AddNewObject - nothing to remove.")
	return
}
