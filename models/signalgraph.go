package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type SignalGraph struct {
	freesp.SignalGraph
}

var _ TreeElement = SignalGraph{}

func (g SignalGraph) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolSignalGraph, g.Filename(), g.ItsType())
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
	SignalGraphType{g.ItsType()}.AddToTree(tree, cursor)
}

func (g SignalGraph) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Fatal("SignalGraph.AddNewObject - nothing to add.")
	return
}

func (g SignalGraph) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("SignalGraph.AddNewObject - nothing to remove.")
	return
}
