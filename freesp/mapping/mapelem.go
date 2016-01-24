package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

type mapelem struct {
	node     bh.NodeIf
	process  pf.ProcessIf
	position image.Point
	mapping  mp.MappingIf
}

func mapelemNew(n bh.NodeIf, p pf.ProcessIf, nodePos image.Point, mapping mp.MappingIf) (m *mapelem) {
	m = &mapelem{n, p, nodePos, mapping}
	return
}

var _ mp.MappedElementIf = (*mapelem)(nil)

func (m mapelem) Position() image.Point {
	return m.position
}

func (m *mapelem) SetPosition(pos image.Point) {
	m.position = pos
}

func (m *mapelem) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var s tr.Symbol
	if m.process == nil {
		s = tr.SymbolUnmapped
	} else {
		s = tr.SymbolMapped
	}
	err := tree.AddEntry(cursor, s, m.node.Name(), m, freesp.MayEdit)
	if err != nil {
		log.Fatalf("mapping.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (m mapelem) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatalf("mapping.AddNewObject error: Nothing to add\n")
	return
}

func (m mapelem) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatalf("mapping.RemoveObject error: Nothing to remove\n")
	return
}

func (m mapelem) Mapping() mp.MappingIf {
	return m.mapping
}

func (m mapelem) Node() bh.NodeIf {
	return m.node
}

func (m mapelem) Process() (p pf.ProcessIf, ok bool) {
	ok = m.process != nil
	if ok {
		p = m.process
	}
	return
}

func (m *mapelem) SetProcess(p pf.ProcessIf) {
	m.process = p
}

func (m mapelem) CreateXml() (buf []byte, err error) {
	var pname string
	p, ok := m.Process()
	if ok {
		pname = fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
	}
	if len(m.Node().InPorts()) > 0 && len(m.Node().OutPorts()) > 0 {
		xmlm := CreateXmlNodeMap(m.Node().Name(), pname, m.Position())
		buf, err = xmlm.Write()
	} else {
		xmlm := CreateXmlIOMap(m.Node().Name(), pname, m.Position())
		buf, err = xmlm.Write()
	}
	return
}
