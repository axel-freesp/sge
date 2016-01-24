package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

type arch struct {
	name      string
	iotypes   ioTypeList
	processes processList
	platform  pf.PlatformIf
	position  map[gr.PositionMode]image.Point
	archports []pf.ArchPortIf
}

var _ pf.ArchIf = (*arch)(nil)

func ArchNew(name string, platform pf.PlatformIf) *arch {
	return &arch{name, ioTypeListInit(), processListInit(), platform,
		make(map[gr.PositionMode]image.Point), nil}
}

func createArchFromXml(xmla backend.XmlArch, platform pf.PlatformIf) (a *arch, err error) {
	a = ArchNew(xmla.Name, platform)
	for _, xmlt := range xmla.IOType {
		var t pf.IOTypeIf
		t, err = IOTypeNew(xmlt.Name, ioModeMap[xmlt.Mode], a.platform)
		if err != nil {
			return
		}
		a.iotypes.Append(t.(*iotype))
	}
	for _, xmlp := range xmla.Processes {
		var pr *process
		pr, err = createProcessFromXml(xmlp, a)
		if err != nil {
			return
		}
		a.processes.Append(pr)
	}
	for _, xmlh := range xmla.Entry {
		mode, ok := ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("createArchFromXml Warning: hint mode %s not defined\n", xmlh.Mode)
			continue
		}
		a.position[mode] = image.Point{xmlh.X, xmlh.Y}
	}
	return
}

func (a *arch) Platform() pf.PlatformIf {
	return a.platform
}

func (a *arch) IOTypes() []pf.IOTypeIf {
	return a.iotypes.IoTypes()
}

func (a *arch) Processes() []pf.ProcessIf {
	return a.processes.Processes()
}

func (a *arch) ArchPorts() []pf.ArchPortIf {
	return a.archports
}

func (a *arch) AddArchPort(ch pf.ChannelIf) (p *archPort) {
	// TODO: Do all the checks...
	p = archPortNew(ch)
	a.archports = append(a.archports, p)
	ch.(*channel).archport = p
	return
}

func (a *arch) CreateXml() (buf []byte, err error) {
	xmla := CreateXmlArch(a)
	buf, err = xmla.Write()
	return
}

//
//  Namer API
//

func (a *arch) Name() string {
	return a.name
}

func (a *arch) SetName(newName string) {
	a.name = newName
}

//
//      ModePositioner API
//

func (a *arch) ModePosition(mode gr.PositionMode) (p image.Point) {
	p = a.position[mode]
	return
}

func (a *arch) SetModePosition(mode gr.PositionMode, p image.Point) {
	a.position[mode] = p
}

//
//  fmt.Stringer API
//

func (a *arch) String() string {
	return fmt.Sprintf("Arch(%s)", a.name)
}

//
//  tr.TreeElement API
//

func (a *arch) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	//log.Printf("arch.AddToTree: %s\n", a.Name())
	err := tree.AddEntry(cursor, tr.SymbolArch, a.Name(), a, MayAddObject|MayEdit|MayRemove)
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

func (a *arch) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("arch.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case pf.IOTypeIf:
		t := obj.(pf.IOTypeIf)
		_, ok := a.iotypes.Find(t.Name())
		if ok {
			err = fmt.Errorf("arch.AddNewObject warning: duplicate ioType name %s (abort)\n", t.Name())
			return
		}
		a.iotypes.Append(t.(*iotype))
		cursor.Position = len(a.IOTypes()) - 1
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	case pf.ProcessIf:
		p := obj.(pf.ProcessIf)
		_, ok := a.processes.Find(p.Name())
		if ok {
			err = fmt.Errorf("arch.AddNewObject warning: duplicate process name %s (abort)\n", p.Name())
			return
		}
		a.processes.Append(p.(*process))
		newCursor = tree.Insert(cursor)
		p.AddToTree(tree, newCursor)
		//log.Printf("arch.AddNewObject: successfully added process %v\n", p)

	default:
		log.Fatalf("arch.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (a *arch) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if a != tree.Object(parent) {
		log.Printf("arch.RemoveObject error: not removing child of mine.")
		return
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case pf.IOTypeIf:
		t := obj.(pf.IOTypeIf)
		_, ok := a.iotypes.Find(t.Name())
		if ok {
			a.iotypes.Remove(t)
		} else {
			log.Printf("arch.RemoveObject error: iotype to be removed not found.\n")
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, t})

	case pf.ProcessIf:
		p := obj.(pf.ProcessIf)
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
			removed = append(removed, tr.IdWithObject{prefix, index, cc})
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
			removed = append(removed, tr.IdWithObject{prefix, index, cc})
		}
		for _, c := range p.OutChannels() {
			p.(*process).outChannels.Remove(c)
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, p})
		//log.Printf("arch.RemoveObject: successfully removed process %v\n", p)

	default:
		log.Fatalf("arch.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (a *arch) Identify(te tr.TreeElement) bool {
	switch te.(type) {
	case *arch:
		return te.(*arch).Name() == a.Name()
	}
	return false
}

//
//      archList
//

type archList struct {
	archs []pf.ArchIf
}

func archListInit() archList {
	return archList{nil}
}

func (l *archList) Append(a *arch) {
	l.archs = append(l.archs, a)
}

func (l *archList) Remove(a pf.ArchIf) {
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

func (l *archList) Archs() []pf.ArchIf {
	return l.archs
}

func (l *archList) Find(name string) (a pf.ArchIf, ok bool) {
	for _, a = range l.archs {
		if a.Name() == name {
			ok = true
			return
		}
	}
	return
}
