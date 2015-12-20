package freesp

import (
	"fmt"
	"log"
)

// signalType

type signalType struct {
	name, ctype, msgid string
	scope              Scope
	mode               Mode
}

/*
 *  freesp.SignalType API
 */

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
 *	fmt.Stringer API
 */

func (t *signalType) String() string {
	return fmt.Sprintf("SignalType(%s, %s, %s, %v, %v)", t.name, t.ctype, t.msgid, t.scope, t.mode)
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*signalType)(nil)

func (t *signalType) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	parent := tree.Object(parentId)
	switch parent.(type) {
	case Library:
		prop = mayAddObject | mayEdit | mayRemove
	case Port, NamedPortType:
		prop = 0
	default:
		log.Fatalf("signalType.AddToTree error: invalid parent type %T\n", parent)
	}
	err := tree.AddEntry(cursor, SymbolSignalType, t.TypeName(), t, prop)
	if err != nil {
		log.Fatal("signalType.AddToTree error: AddEntry failed: %s", err)
	}
}

func (t *signalType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Fatal("signalType.AddNewObject - nothing to add.")
	return
}

func (t *signalType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("signalType.AddNewObject - nothing to remove.")
	return
}

/*
 *      signalTypeList
 *
 */

type signalTypeList struct {
	signalTypes []SignalType
}

func signalTypeListInit() signalTypeList {
	return signalTypeList{nil}
}

func (l *signalTypeList) Append(st SignalType) {
	l.signalTypes = append(l.signalTypes, st)
}

func (l *signalTypeList) Remove(st SignalType) {
	var i int
	for i = range l.signalTypes {
		if st == l.signalTypes[i] {
			break
		}
	}
	if i >= len(l.signalTypes) {
		for _, v := range l.signalTypes {
			log.Printf("signalTypeList.RemoveNodeType have SignalType %v\n", v)
		}
		log.Fatalf("signalTypeList.RemoveNodeType error: SignalType %v not in this list\n", st)
	}
	for i++; i < len(l.signalTypes); i++ {
		l.signalTypes[i-1] = l.signalTypes[i]
	}
	l.signalTypes = l.signalTypes[:len(l.signalTypes)-1]
}

func (l *signalTypeList) SignalTypes() []SignalType {
	return l.signalTypes
}
