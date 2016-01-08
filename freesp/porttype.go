package freesp

import (
	"fmt"
	"log"
	interfaces "github.com/axel-freesp/sge/interface"
)

// portType

type portType struct {
	signalType SignalType
	name       string
	direction  interfaces.PortDirection
}

var _ PortType = (*portType)(nil)

func PortTypeNew(name string, pTypeName string, dir interfaces.PortDirection) *portType {
	st, ok := signalTypes[pTypeName]
	if !ok {
		log.Fatalf("NamedPortTypeNew error: signal type '%s' not defined\n", pTypeName)
	}
	return &portType{st, name, dir}
}

func (t *portType) Name() string {
	return t.name
}

func (t *portType) SetName(newName string) {
	t.name = newName
}

func (t *portType) SignalType() SignalType {
	return t.signalType
}

func (t *portType) SetSignalType(newSignalType SignalType) {
	t.signalType = newSignalType
}

func (t *portType) Direction() interfaces.PortDirection {
	return t.direction
}

func (t *portType) SetDirection(newDir interfaces.PortDirection) {
	t.direction = newDir
}

func (t *portType) String() (s string) {
	s = fmt.Sprintf("PortType(%s, %s, %s)", t.name, t.direction, t.SignalType())
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*portType)(nil)

func (p *portType) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	if tree.Property(parentId).IsReadOnly() {
		prop = 0
	} else {
		prop = mayEdit | mayRemove | mayAddObject
	}
	var kind Symbol
	if p.Direction() == interfaces.InPort {
		kind = SymbolInputPortType
	} else {
		kind = SymbolOutputPortType
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p, prop)
	if err != nil {
		log.Fatal("PortType.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	p.SignalType().AddToTree(tree, child)
}

func (p *portType) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatal("PortType.AddNewObject - nothing to add.")
	return
}

func (p *portType) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("PortType.AddNewObject - nothing to remove.")
	return
}

/*
 *      portTypeList
 *
 */

type portTypeList struct {
	portTypes []PortType
}

func portTypeListInit() portTypeList {
	return portTypeList{nil}
}

func (l *portTypeList) Append(nt PortType) {
	l.portTypes = append(l.portTypes, nt)
}

func (l *portTypeList) Remove(nt PortType) {
	var i int
	for i = range l.portTypes {
		if nt == l.portTypes[i] {
			break
		}
	}
	if i >= len(l.portTypes) {
		for _, v := range l.portTypes {
			log.Printf("portTypeList.RemovePort have PortType %v\n", v)
		}
		log.Fatalf("portTypeList.RemovePort error: PortType %v not in this list\n", nt)
	}
	for i++; i < len(l.portTypes); i++ {
		l.portTypes[i-1] = l.portTypes[i]
	}
	l.portTypes = l.portTypes[:len(l.portTypes)-1]
}

func (l *portTypeList) PortTypes() []PortType {
	return l.portTypes
}

func (l *portTypeList) Find(name string) (p PortType, ok bool, index int) {
	ok = false
	for index, p = range l.portTypes {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}
