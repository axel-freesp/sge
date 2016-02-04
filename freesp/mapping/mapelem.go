package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	//"github.com/axel-freesp/sge/tool"
	//"image"
	"log"
)

type mapelem struct {
	gr.PathModePositionerObject
	node              bh.NodeIf
	nodeId            bh.NodeIdIf
	process           pf.ProcessIf
	mapping           mp.MappingIf
	expanded          bool
	inports, outports []gr.PathModePositionerObject
}

func mapelemNew(n bh.NodeIf, nId bh.NodeIdIf, p pf.ProcessIf, mapping mp.MappingIf) (m *mapelem) {
	m = &mapelem{*gr.PathModePositionerObjectNew(), n, nId, p, mapping, false,
		make([]gr.PathModePositionerObject, len(n.InPorts())),
		make([]gr.PathModePositionerObject, len(n.OutPorts()))}
	return
}

var _ mp.MappedElementIf = (*mapelem)(nil)

/*
 *      Expander API
 */

func (m mapelem) Expanded() bool {
	return m.expanded
}

func (m *mapelem) SetExpanded(xp bool) {
	m.expanded = xp
}

func (m *mapelem) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var s tr.Symbol
	if m.process == nil {
		s = tr.SymbolUnmapped
	} else {
		s = tr.SymbolMapped
	}
	err := tree.AddEntry(cursor, s, m.nodeId.String(), m, freesp.MayEdit)
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

func (m mapelem) NodeId() bh.NodeIdIf {
	return m.nodeId
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
	if len(m.node.InPorts()) > 0 && len(m.node.OutPorts()) > 0 {
		xmlm := CreateXmlNodeMap(m.nodeId.String(), pname)
		buf, err = xmlm.Write()
	} else {
		xmlm := CreateXmlIOMap(m.nodeId.String(), pname)
		buf, err = xmlm.Write()
	}
	return
}
