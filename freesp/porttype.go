package freesp

import (
	"fmt"
	"log"
)

// portType

type portType struct {
	name string
	ref  *signalType
}

var _ PortType = (*portType)(nil)

func PortTypeNew(name string, st SignalType) *portType {
	return &portType{name, st.(*signalType)}
}

func (t *portType) TypeName() string {
	return t.name
}

func (t *portType) SignalType() SignalType {
	return t.ref
}

// namedPortType

type namedPortType struct {
	name      string
	pType     *portType
	direction PortDirection
}

var _ NamedPortType = (*namedPortType)(nil)

func NamedPortTypeNew(name string, pTypeName string, dir PortDirection) *namedPortType {
	pt := getPortType(pTypeName)
	return &namedPortType{name, pt, dir}
}

func (t *namedPortType) TypeName() string {
	return t.pType.TypeName()
}

func (t *namedPortType) Name() string {
	return t.name
}

func (t *namedPortType) Direction() PortDirection {
	return t.direction
}

func (t *namedPortType) SignalType() SignalType {
	return t.pType.SignalType()
}

func (t *namedPortType) String() (s string) {
	s = fmt.Sprintf("NamedPortType(%s, %s, %s)", t.name, t.direction, t.SignalType())
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*namedPortType)(nil)

func (p *namedPortType) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	if tree.Property(parentId).IsReadOnly() {
		prop = 0
	} else {
		prop = mayEdit | mayRemove
	}
	var kind Symbol
	if p.Direction() == InPort {
		kind = SymbolInputPortType
	} else {
		kind = SymbolOutputPortType
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p, prop)
	if err != nil {
		log.Fatal("NamedPortType.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	p.SignalType().AddToTree(tree, child)
}

func (p *namedPortType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Fatal("NamedPortType.AddNewObject - nothing to add.")
	return
}

func (p *namedPortType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("NamedPortType.AddNewObject - nothing to remove.")
	return
}
