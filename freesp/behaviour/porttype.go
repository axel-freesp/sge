package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

// portType

type portType struct {
	signalType bh.SignalTypeIf
	name       string
	direction  gr.PortDirection
	position   map[gr.PositionMode]image.Point
}

var _ bh.PortTypeIf = (*portType)(nil)

func PortTypeNew(name string, pTypeName string, dir gr.PortDirection) *portType {
	st, ok := freesp.GetSignalTypeByName(pTypeName)
	if !ok {
		log.Fatalf("NamedPortTypeNew error: FIXME: signal type '%s' not defined\n", pTypeName)
	}
	return &portType{st, name, dir, make(map[gr.PositionMode]image.Point)}
}

func (t *portType) Name() string {
	return t.name
}

func (t *portType) SetName(newName string) {
	t.name = newName
}

func (t *portType) SignalType() bh.SignalTypeIf {
	return t.signalType
}

func (t *portType) SetSignalType(newSignalType bh.SignalTypeIf) {
	t.signalType = newSignalType
}

func (t *portType) Direction() gr.PortDirection {
	return t.direction
}

func (t *portType) SetDirection(newDir gr.PortDirection) {
	t.direction = newDir
}

func (t *portType) String() (s string) {
	s = fmt.Sprintf("bh.PortTypeIf(%s, %s, %s)", t.name, t.direction, t.SignalType())
	return
}

/*
 *      ModePositioner API
 */

func (t *portType) ModePosition(mode gr.PositionMode) (p image.Point) {
	p = t.position[mode]
	return
}

func (t *portType) SetModePosition(mode gr.PositionMode, p image.Point) {
	t.position[mode] = p
}

func (t *portType) CreateXml() (buf []byte, err error) {
	if t.Direction() == gr.InPort {
		xmlporttype := CreateXmlNamedInPort(t)
		buf, err = xmlporttype.Write()
	} else {
		xmlporttype := CreateXmlNamedOutPort(t)
		buf, err = xmlporttype.Write()
	}
	return
}

/*
 *  tr.TreeElement API
 */

var _ tr.TreeElement = (*portType)(nil)

func (p *portType) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var prop tr.Property
	parentId := tree.Parent(cursor)
	if tree.Property(parentId).IsReadOnly() {
		prop = freesp.PropertyNew(false, false, false)
	} else {
		prop = freesp.PropertyNew(true, true, true)
	}
	var kind tr.Symbol
	if p.Direction() == gr.InPort {
		kind = tr.SymbolInputPortType
	} else {
		kind = tr.SymbolOutputPortType
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p, prop)
	if err != nil {
		log.Fatal("bh.PortTypeIf.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	p.SignalType().AddToTree(tree, child)
}

func (p *portType) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatal("bh.PortTypeIf.AddNewObject - nothing to add.")
	return
}

func (p *portType) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatal("bh.PortTypeIf.AddNewObject - nothing to remove.")
	return
}

/*
 *      portTypeList
 *
 */

type portTypeList struct {
	portTypes []bh.PortTypeIf
}

func portTypeListInit() portTypeList {
	return portTypeList{nil}
}

func (l *portTypeList) Append(nt bh.PortTypeIf) {
	l.portTypes = append(l.portTypes, nt)
}

func (l *portTypeList) Remove(nt bh.PortTypeIf) {
	var i int
	for i = range l.portTypes {
		if nt == l.portTypes[i] {
			break
		}
	}
	if i >= len(l.portTypes) {
		for _, v := range l.portTypes {
			log.Printf("portTypeList.RemovePort have bh.PortTypeIf %v\n", v)
		}
		log.Fatalf("portTypeList.RemovePort error: bh.PortTypeIf %v not in this list\n", nt)
	}
	for i++; i < len(l.portTypes); i++ {
		l.portTypes[i-1] = l.portTypes[i]
	}
	l.portTypes = l.portTypes[:len(l.portTypes)-1]
}

func (l *portTypeList) PortTypes() []bh.PortTypeIf {
	return l.portTypes
}

func (l *portTypeList) Find(name string) (p bh.PortTypeIf, ok bool, index int) {
	ok = false
	for index, p = range l.portTypes {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}
