package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
)

// signalType

type signalType struct {
	name, ctype, msgid string
	scope              bh.Scope
	mode               bh.Mode
	definedAt          string
}

/*
 *  freesp.bh.SignalType API
 */

var _ bh.SignalTypeIf = (*signalType)(nil)

func SignalTypeNew(name, ctype, msgid string, scope bh.Scope, mode bh.Mode, definedAt string) (t *signalType, err error) {
	newT := &signalType{name, ctype, msgid, scope, mode, definedAt}
	sType, ok := freesp.GetSignalTypeByName(name)
	if ok {
		if !signalTypeCpmpatible(newT, sType) {
			err = fmt.Errorf(`SignalTypeNew error: adding existing signal
				type %s, which is incompatible`, name)
			return
		}
		log.Printf(`SignalTypeNew: warning: adding existing signal type
			definition %s (taking the existing)`, name)
		t = sType.(*signalType)
	} else {
		t = newT
		freesp.RegisterSignalType(t)
	}
	return
}

func signalTypeCpmpatible(s1, s2 bh.SignalTypeIf) bool {
	if s1.TypeName() != s2.TypeName() {
		return true
	}
	if s1.CType() != s2.CType() {
		return false
	}
	if s1.Scope() != s2.Scope() {
		return false
	}
	if s1.ChannelId() != s2.ChannelId() {
		return false
	}
	if s1.Mode() != s2.Mode() {
		return false
	}
	return true
}

func (t signalType) DefinedAt() string {
	return t.definedAt
}

func SignalTypeDestroy(t bh.SignalTypeIf) {
	freesp.RemoveRegisteredSignalType(t)
}

func (t *signalType) TypeName() string {
	return t.name
}

func (t *signalType) SetTypeName(newName string) {
	// TODO: how to make this consistent in all cases?
	log.Println("signalType.SetTypeName WARNING: this is not yet implemented!")
}

func (t *signalType) CType() string {
	return t.ctype
}

func (t *signalType) SetCType(newCType string) {
	t.ctype = newCType
}

func (t *signalType) ChannelId() string {
	return t.msgid
}

func (t *signalType) SetChannelId(newChannelId string) {
	t.msgid = newChannelId
}

func (t *signalType) Scope() bh.Scope {
	return t.scope
}

func (t *signalType) SetScope(newScope bh.Scope) {
	t.scope = newScope
}

func (t *signalType) Mode() bh.Mode {
	return t.mode
}

func (t *signalType) SetMode(newMode bh.Mode) {
	t.mode = newMode
}

func (t *signalType) CreateXml() (buf []byte, err error) {
	if t != nil {
		xmlsignaltype := CreateXmlSignalType(t)
		buf, err = xmlsignaltype.Write()
	}
	return
}

/*
 *	fmt.Stringer API
 */

func (t *signalType) String() string {
	return fmt.Sprintf("bh.SignalType(%s, %s, %s, %v, %v)", t.name, t.ctype, t.msgid, t.scope, t.mode)
}

/*
 *  tr.TreeElementIf API
 */

var _ tr.TreeElementIf = (*signalType)(nil)

func (t *signalType) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var prop tr.Property
	parentId := tree.Parent(cursor)
	parent := tree.Object(parentId)
	switch parent.(type) {
	case bh.LibraryIf:
		prop = freesp.PropertyNew(true, true, true)
	case bh.PortIf, bh.PortTypeIf:
		prop = freesp.PropertyNew(false, false, false)
	default:
		log.Fatalf("signalType.AddToTree error: invalid parent type %T\n", parent)
	}
	err := tree.AddEntry(cursor, tr.SymbolSignalType, t.TypeName(), t, prop)
	if err != nil {
		log.Fatalf("signalType.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (t *signalType) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	log.Fatal("signalType.AddNewObject - nothing to add.")
	return
}

func (t *signalType) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatal("signalType.AddNewObject - nothing to remove.")
	return
}

/*
 *      signalTypeList
 *
 */

type signalTypeList struct {
	signalTypes []bh.SignalTypeIf
}

func signalTypeListInit() signalTypeList {
	return signalTypeList{nil}
}

func (l *signalTypeList) Append(st bh.SignalTypeIf) {
	l.signalTypes = append(l.signalTypes, st)
}

func (l *signalTypeList) Remove(st bh.SignalTypeIf) {
	var i int
	for i = range l.signalTypes {
		if st == l.signalTypes[i] {
			break
		}
	}
	if i >= len(l.signalTypes) {
		for _, v := range l.signalTypes {
			log.Printf("signalTypeList.RemoveNodeType have bh.SignalType %v\n", v)
		}
		log.Fatalf("signalTypeList.RemoveNodeType error: bh.SignalType %v not in this list\n", st)
	}
	for i++; i < len(l.signalTypes); i++ {
		l.signalTypes[i-1] = l.signalTypes[i]
	}
	l.signalTypes = l.signalTypes[:len(l.signalTypes)-1]
}

func (l *signalTypeList) SignalTypes() []bh.SignalTypeIf {
	return l.signalTypes
}
