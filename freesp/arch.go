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
	platform  Platform
}

var _ Arch = (*arch)(nil)

func ArchNew(name string, platform Platform) *arch {
	return &arch{name, ioTypeListInit(), processListInit(), platform}
}

func (a *arch) createArchFromXml(xmla backend.XmlArch) (err error) {
	a.name = xmla.Name
	for _, xmlt := range xmla.IOType {
		var t IOType
		t, err = IOTypeNew(xmlt.Name, ioModeMap[xmlt.Mode], a.platform)
		if err != nil {
			return
		}
		a.iotypes.Append(t)
	}
	for _, xmlp := range xmla.Processes {
		p := ProcessNew(xmlp.Name, a)
		err = p.createProcessFromXml(xmlp, a.IOTypes())
		if err != nil {
			err = fmt.Errorf("arch.createArchFromXml error: %s\n", err)
		}
		a.processes.Append(p)
	}
	return
}

func (a *arch) Platform() Platform {
	return a.platform
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
	//log.Printf("arch.AddToTree: %s\n", a.Name())
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
	if obj == nil {
		err = fmt.Errorf("arch.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case IOType:
		t := obj.(IOType)
		_, ok := a.iotypes.Find(t.Name())
		if ok {
			err = fmt.Errorf("arch.AddNewObject warning: duplicate ioType name %s (abort)\n", t.Name())
			return
		}
		a.iotypes.Append(t)
		cursor.Position = len(a.IOTypes()) - 1
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	case Process:
		p := obj.(Process)
		_, ok := a.processes.Find(p.Name())
		if ok {
			err = fmt.Errorf("arch.AddNewObject warning: duplicate process name %s (abort)\n", p.Name())
			return
		}
		a.processes.Append(p)
		newCursor = tree.Insert(cursor)
		p.AddToTree(tree, newCursor)
		//log.Printf("arch.AddNewObject: successfully added process %v\n", p)

	default:
		log.Fatalf("arch.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (a *arch) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if a != tree.Object(parent) {
		log.Printf("arch.RemoveObject error: not removing child of mine.")
		return
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case IOType:
		t := obj.(IOType)
		_, ok := a.iotypes.Find(t.Name())
		if ok {
			a.iotypes.Remove(t)
		} else {
			log.Printf("arch.RemoveObject error: iotype to be removed not found.\n")
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, t})

	case Process:
		p := obj.(Process)
		if p.Arch() != a {
			log.Printf("arch.RemoveObject error: process to be removed is no child of mine.")
		}
		_, ok := a.processes.Find(p.Name())
		if ok {
			a.processes.Remove(p)
		} else {
			log.Printf("arch.RemoveObject error: process to be removed not found.\n")
		}
		for _, c := range p.InChannels() {
			cc := c.Link()
			pp := cc.Process()
			ppCursor := tree.Cursor(pp) // TODO: better search over platform...
			ccCursor := tree.CursorAt(ppCursor, cc)
			//log.Printf("arch.RemoveObject: remove %v\n", cc)
			pp.(*process).outChannels.Remove(cc)
			prefix, index := tree.Remove(ccCursor)
			removed = append(removed, IdWithObject{prefix, index, cc})
		}
		for _, c := range p.InChannels() {
			p.(*process).inChannels.Remove(c)
		}
		for _, c := range p.OutChannels() {
			cc := c.Link()
			pp := cc.Process()
			ppCursor := tree.Cursor(pp) // TODO: better search over platform...
			ccCursor := tree.CursorAt(ppCursor, cc)
			//log.Printf("arch.RemoveObject: remove %v\n", cc)
			pp.(*process).inChannels.Remove(cc)
			prefix, index := tree.Remove(ccCursor)
			removed = append(removed, IdWithObject{prefix, index, cc})
		}
		for _, c := range p.OutChannels() {
			p.(*process).outChannels.Remove(c)
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, p})
		//log.Printf("arch.RemoveObject: successfully removed process %v\n", p)

	default:
		log.Fatalf("arch.AddNewObject error: invalid type %T\n", obj)
	}
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
