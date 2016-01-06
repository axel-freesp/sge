package freesp

import (
	//"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

var ioModeMap = map[backend.XmlIOMode]IOMode{
	backend.IOModeShmem: IOModeShmem,
	backend.IOModeAsync: IOModeAsync,
	backend.IOModeSync:  IOModeSync,
}

var ioXmlModeMap = map[IOMode]backend.XmlIOMode{
	IOModeShmem: backend.IOModeShmem,
	IOModeAsync: backend.IOModeAsync,
	IOModeSync:  backend.IOModeSync,
}

type iotype struct {
	name string
	mode IOMode
}

var _ IOType = (*iotype)(nil)

func IOTypeNew(name string, mode IOMode) *iotype {
	return &iotype{name, mode}
}

func (t *iotype) Mode() IOMode {
	return t.mode
}

func (t *iotype) SetMode(newMode IOMode) {
	t.mode = newMode
}

/*
 *  Namer API
 */

func (t *iotype) Name() string {
	return t.name
}

func (t *iotype) SetName(string) {
}

/*
 *  TreeElement API
 */

func (t *iotype) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolIOType, t.Name(), t, mayEdit|mayRemove)
	if err != nil {
		log.Fatalf("iotype.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (t *iotype) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	return
}

func (t *iotype) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	return
}

/*
 *      ioTypeList
 *
 */

type ioTypeList struct {
	ioTypes []IOType
}

func ioTypeListInit() ioTypeList {
	return ioTypeList{nil}
}

func (l *ioTypeList) Append(st IOType) {
	l.ioTypes = append(l.ioTypes, st)
}

func (l *ioTypeList) Remove(st IOType) {
	var i int
	for i = range l.ioTypes {
		if st == l.ioTypes[i] {
			break
		}
	}
	if i >= len(l.ioTypes) {
		for _, v := range l.ioTypes {
			log.Printf("ioTypeList.RemoveNodeType have IoType %v\n", v)
		}
		log.Fatalf("ioTypeList.RemoveNodeType error: IoType %v not in this list\n", st)
	}
	for i++; i < len(l.ioTypes); i++ {
		l.ioTypes[i-1] = l.ioTypes[i]
	}
	l.ioTypes = l.ioTypes[:len(l.ioTypes)-1]
}

func (l *ioTypeList) IoTypes() []IOType {
	return l.ioTypes
}
