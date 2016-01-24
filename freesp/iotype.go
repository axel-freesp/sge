package freesp

import (
	"fmt"
	"log"
	"github.com/axel-freesp/sge/backend"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	gr "github.com/axel-freesp/sge/interface/graph"
)

var ioModeMap = map[backend.XmlIOMode]gr.IOMode{
	backend.IOModeShmem: gr.IOModeShmem,
	backend.IOModeAsync: gr.IOModeAsync,
	backend.IOModeSync:  gr.IOModeSync,
}

var ioXmlModeMap = map[gr.IOMode]backend.XmlIOMode{
	gr.IOModeShmem: backend.IOModeShmem,
	gr.IOModeAsync: backend.IOModeAsync,
	gr.IOModeSync:  backend.IOModeSync,
}

type iotype struct {
	name     string
	mode     gr.IOMode
	platform pf.PlatformIf
}

var _ pf.IOTypeIf = (*iotype)(nil)

func IOTypeNew(name string, mode gr.IOMode, platform pf.PlatformIf) (t *iotype, err error) {
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

func (t *iotype) IOMode() gr.IOMode {
	return t.mode
}

func (t *iotype) SetIOMode(newMode gr.IOMode) {
	t.mode = newMode
}

func (t *iotype) Platform() pf.PlatformIf {
	return t.platform
}

func (t *iotype) CreateXml() (buf []byte, err error) {
	xmlt := CreateXmlIOType(t)
	buf, err = xmlt.Write()
	return
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

func (t *iotype) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	err := tree.AddEntry(cursor, tr.SymbolIOType, t.Name(), t, MayEdit|MayAddObject|MayRemove)
	if err != nil {
		log.Fatalf("iotype.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (t *iotype) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatalf("iotype.AddNewObject error: nothing to add\n")
	return
}

func (t *iotype) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatalf("iotype.RemoveObject error: nothing to remove\n")
	return
}

func (t *iotype) Identify(te tr.TreeElement) bool {
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
	ioTypes []pf.IOTypeIf
}

func ioTypeListInit() ioTypeList {
	return ioTypeList{}
}

func (l *ioTypeList) Append(st *iotype) {
	l.ioTypes = append(l.ioTypes, st)
}

func (l *ioTypeList) Remove(st pf.IOTypeIf) {
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

func (l *ioTypeList) IoTypes() []pf.IOTypeIf {
	return l.ioTypes
}

func (l *ioTypeList) Find(name string) (t pf.IOTypeIf, ok bool) {
	for _, t = range l.ioTypes {
		if t.Name() == name {
			ok = true
			return
		}
	}
	return
}
