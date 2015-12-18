package models

import (
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type SignalType struct {
	freesp.SignalType
}

func (t SignalType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolSignalType, t.TypeName(), t)
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
}
