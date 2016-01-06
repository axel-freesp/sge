package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

type arch struct {
	name      string
	iotypes   ioTypeList
	processes processList
}

var _ Arch = (*arch)(nil)

func ArchNew(name string) *arch {
	return &arch{name, ioTypeListInit(), processListInit()}
}

func (a *arch) createArchFromXml(xmla backend.XmlArch) (err error) {
	a.name = xmla.Name
	for _, xmlt := range xmla.IOType {
		t := IOTypeNew(xmlt.Name, ioModeMap[xmlt.Mode])
		a.iotypes.Append(t)
	}
	for _, xmlp := range xmla.Processes {
		p := ProcessNew(xmlp.Name)
		err = p.createProcessFromXml(xmlp, a.IOTypes())
		if err != nil {
			err = fmt.Errorf("arch.createArchFromXml error: %s\n", err)
		}
		a.processes.Append(p)
	}
	return
}

func (a *arch) IOTypes() []IOType {
	return a.iotypes.IoTypes()
}

func (a *arch) Processes() []Process {
	return a.processes.Processes()
}

/*
 *  Namer API
 */

func (a *arch) Name() string {
	return a.name
}

func (a *arch) SetName(newName string) {
	a.name = newName
}

/*
 *  TreeElement API
 */

func (a *arch) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolArch, a.Name(), a, mayAddObject|mayEdit|mayRemove)
	if err != nil {
		log.Fatalf("arch.AddToTree error: AddEntry failed: %s", err)
	}
	for _, t := range a.IOTypes() {
		child := tree.Append(cursor)
		t.AddToTree(tree, child)
	}
	for _, p := range a.Processes() {
		child := tree.Append(cursor)
		p.AddToTree(tree, child)
	}
}

func (a *arch) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	return
}

func (a *arch) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	return
}

/*
 *      archList
 *
 */

type archList struct {
	archs []Arch
}

func archListInit() archList {
	return archList{nil}
}

func (l *archList) Append(a Arch) {
	l.archs = append(l.archs, a)
}

func (l *archList) Remove(a Arch) {
	var i int
	for i = range l.archs {
		if a == l.archs[i] {
			break
		}
	}
	if i >= len(l.archs) {
		for _, v := range l.archs {
			log.Printf("archList.RemoveArch have Arch %v\n", v)
		}
		log.Fatalf("archList.RemoveArch error: Arch %v not in this list\n", a)
	}
	for i++; i < len(l.archs); i++ {
		l.archs[i-1] = l.archs[i]
	}
	l.archs = l.archs[:len(l.archs)-1]
}

func (l *archList) Archs() []Arch {
	return l.archs
}

func (l *archList) Find(name string) (a Arch, ok bool) {
	for _, a = range l.archs {
		if a.Name() == name {
			ok = true
			return
		}
	}
	return
}
