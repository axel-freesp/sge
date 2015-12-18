package freesp

import (
	"log"
)

// signalType

type signalType struct {
	name, ctype, msgid string
	scope              Scope
	mode               Mode
}

var _ SignalType = (*signalType)(nil)

func SignalTypeNew(name, ctype, msgid string, scope Scope, mode Mode) *signalType {
	return &signalType{name, ctype, msgid, scope, mode}
}

func (t *signalType) TypeName() string {
	return t.name
}

func (t *signalType) CType() string {
	return t.ctype
}

func (t *signalType) ChannelId() string {
	return t.msgid
}

func (t *signalType) Scope() Scope {
	return t.scope
}

func (t *signalType) Mode() Mode {
	return t.mode
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*signalType)(nil)

func (t *signalType) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolSignalType, t.TypeName(), t)
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
}

func (t *signalType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Fatal("SignalType.AddNewObject - nothing to add.")
	return
}

func (t *signalType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("SignalType.AddNewObject - nothing to remove.")
	return
}
