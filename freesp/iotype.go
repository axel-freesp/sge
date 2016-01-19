package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	interfaces "github.com/axel-freesp/sge/interface"
	"log"
)

var ioModeMap = map[backend.XmlIOMode]interfaces.IOMode{
	backend.IOModeShmem: interfaces.IOModeShmem,
	backend.IOModeAsync: interfaces.IOModeAsync,
	backend.IOModeSync:  interfaces.IOModeSync,
}

var ioXmlModeMap = map[interfaces.IOMode]backend.XmlIOMode{
	interfaces.IOModeShmem: backend.IOModeShmem,
	interfaces.IOModeAsync: backend.IOModeAsync,
	interfaces.IOModeSync:  backend.IOModeSync,
}

type iotype struct {
	name     string
	mode     interfaces.IOMode
	platform Platform
}

var _ IOType = (*iotype)(nil)
var _ interfaces.IOTypeObject = (*iotype)(nil)

func IOTypeNew(name string, mode interfaces.IOMode, platform Platform) (t *iotype, err error) {
	newT := &iotype{name, mode, platform}
	ioType := ioTypes[name]
	if ioType != nil {
		if (*newT) != (*ioType) {
			err = fmt.Errorf("IOTypeNew error: adding existing io type %s, which is incompatible.", name)
			err = fmt.Errorf("%s\nexisting: %v - new: %v\n", err, ioType, newT)
			return
		}
		t = ioType
	} else {
		t = newT
		ioTypes[name] = t
		registeredIOTypes.Append(name)
	}
	return
}

func (t *iotype) IOMode() interfaces.IOMode {
	return t.mode
}

func (t *iotype) SetIOMode(newMode interfaces.IOMode) {
	t.mode = newMode
}

func (t *iotype) Platform() Platform {
	return t.platform
}

//
//  Namer API
//

func (t *iotype) Name() string {
	return t.name
}

func (t *iotype) SetName(newName string) {
	if ioTypes[newName] != nil {
		log.Printf("iotype.SetName error: cannot rename to existing iotype.\n")
		return
	}
	registeredIOTypes.Remove(t.name)
	delete(ioTypes, t.name)
	t.name = newName
	ioTypes[t.name] = t
	registeredIOTypes.Append(t.name)
}

//
//  TreeElement API
//

func (t *iotype) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolIOType, t.Name(), t, mayEdit|mayAddObject|mayRemove)
	if err != nil {
		log.Fatalf("iotype.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (t *iotype) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatalf("iotype.AddNewObject error: nothing to add\n")
	return
}

func (t *iotype) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatalf("iotype.RemoveObject error: nothing to remove\n")
	return
}

func (t *iotype) Identify(te TreeElement) bool {
	switch te.(type) {
	case *iotype:
		return te.(*iotype).Name() == t.Name()
	}
	return false
}

//
//      ioTypeList
//

type ioTypeList struct {
	ioTypes []IOType
	exports []interfaces.IOTypeObject
}

func ioTypeListInit() ioTypeList {
	return ioTypeList{nil, nil}
}

func (l *ioTypeList) Append(st *iotype) {
	l.ioTypes = append(l.ioTypes, st)
	l.exports = append(l.exports, st)
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
		l.exports[i-1] = l.exports[i]
	}
	l.ioTypes = l.ioTypes[:len(l.ioTypes)-1]
	l.exports = l.exports[:len(l.exports)-1]
}

func (l *ioTypeList) IoTypes() []IOType {
	return l.ioTypes
}

func (l *ioTypeList) Exports() []interfaces.IOTypeObject {
	return l.exports
}

func (l *ioTypeList) Find(name string) (t IOType, ok bool) {
	for _, t = range l.ioTypes {
		if t.Name() == name {
			ok = true
			return
		}
	}
	return
}
