package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	fd "github.com/axel-freesp/sge/interface/filedata"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

type mapping struct {
	graph      bh.SignalGraphIf
	platform   pf.PlatformIf
	context    mod.ModelContextIf
	maps       map[string]*mapelem
	filename   string
	pathPrefix string
}

var _ mp.MappingIf = (*mapping)(nil)

func MappingNew(filename string, context mod.ModelContextIf) *mapping {
	return &mapping{nil, nil, context, make(map[string]*mapelem), filename, ""}
}

func (m *mapping) CreateXml() (buf []byte, err error) {
	xmlm := CreateXmlMapping(m)
	buf, err = xmlm.Write()
	return
}

//
//		TreeElement interface
//

func (m *mapping) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var child tr.Cursor
	err := tree.AddEntry(cursor, tr.SymbolMappings, m.Filename(), m, freesp.MayAddObject)
	if err != nil {
		log.Fatalf("mapping.AddToTree error: AddEntry failed: %s\n", err)
	}
	child = tree.Append(cursor)
	m.graph.AddToTree(tree, child)
	child = tree.Append(cursor)
	m.platform.AddToTree(tree, child)
	for _, n := range m.graph.ItsType().Nodes() {
		child = tree.Append(cursor)
		melem, ok := m.maps[n.Name()]
		if !ok {
			melem = mapelemNew(n, nil, image.Point{}, m)
			m.maps[n.Name()] = melem
		}
		melem.AddToTree(tree, child)
	}
}

func (m mapping) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatalf("mapping.AddNewObject error: Nothing to add\n")
	return
}

func (m mapping) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatalf("mapping.RemoveObject error: Nothing to remove\n")
	return
}

func (m *mapping) Identify(te tr.TreeElement) bool {
	switch te.(type) {
	case *mapping:
		return te.(*mapping).Filename() == m.Filename()
	}
	return false
}

//
//		Mapping interface
//

func (m *mapping) AddMapping(n bh.NodeIf, p pf.ProcessIf) {
	m.maps[n.Name()] = mapelemNew(n, p, image.Point{}, m)
}

func (m *mapping) SetGraph(g bh.SignalGraphIf) {
	m.graph = g
}

func (m *mapping) Graph() bh.SignalGraphIf {
	return m.graph
}

func (m *mapping) SetPlatform(p pf.PlatformIf) {
	m.platform = p
}

func (m *mapping) Platform() pf.PlatformIf {
	return m.platform
}

func (m *mapping) Mapped(n bh.NodeIf) (pr pf.ProcessIf, ok bool) {
	var melem *mapelem
	melem, ok = m.maps[n.Name()]
	if ok {
		pr = melem.process
	}
	ok = pr != nil
	return
}

func (m *mapping) MappedElement(n bh.NodeIf) (melem mp.MappedElementIf, ok bool) {
	melem, ok = m.maps[n.Name()]
	return
}

//
//		Filenamer interface
//

func (m mapping) Filename() string {
	return m.filename
}

func (m *mapping) SetFilename(newFilename string) {
	m.filename = newFilename
}

func (m mapping) PathPrefix() string {
	return m.pathPrefix
}

func (m *mapping) SetPathPrefix(newP string) {
	m.pathPrefix = newP
}

func (m *mapping) Read(data []byte) (cnt int, err error) {
	xmlm := backend.XmlMappingNew(m.Graph().Filename(), m.Platform().Filename())
	cnt, err = xmlm.Read(data)
	if err != nil {
		err = fmt.Errorf("mapping.Read: %v", err)
		return
	}
	err = m.createMappingFromXml(xmlm)
	return
}

func (m *mapping) ReadFile(filepath string) error {
	xmlm := backend.XmlMappingNew("", "")
	err := xmlm.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("mapping.ReadFile: %v", err)
	}
	return m.createMappingFromXml(xmlm)
}

func (m mapping) Write() (data []byte, err error) {
	xmlm := CreateXmlMapping(&m)
	data, err = xmlm.Write()
	return
}

func (m mapping) WriteFile(filepath string) error {
	xmlm := CreateXmlMapping(&m)
	return xmlm.WriteFile(filepath)
}

func (m *mapping) RemoveFromTree(tree tr.TreeIf) {
	tree.Remove(tree.Cursor(m))
}

//
//		Local Methods and Functions
//

func (m *mapping) createMappingFromXml(xmlm *backend.XmlMapping) (err error) {
	var f fd.FileDataIf
	f, err = m.context.SignalGraphMgr().Access(xmlm.SignalGraph)
	if err != nil {
		err = fmt.Errorf("mapping.CreateMappingFromXml: graph.ReadFile failed: %s\n", err)
		return
	}
	m.graph = f.(bh.SignalGraphIf)
	f, err = m.context.PlatformMgr().Access(xmlm.Platform)
	if err != nil {
		err = fmt.Errorf("mapping.CreateMappingFromXml: platform.ReadFile failed: %s\n", err)
		return
	}
	m.platform = f.(pf.PlatformIf)
	var n bh.NodeIf
	var p pf.ProcessIf
	var ok bool
	for _, x := range xmlm.IOMappings {
		n, ok = m.graph.ItsType().NodeByName(x.Name)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: node %s not in graph %s\n", x.Name, m.graph.Filename())
			continue
		}
		p, ok = m.platform.ProcessByName(x.Process)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: process %s not in platform %s\n", x.Process, m.platform.Filename())
			continue
		}
		m.maps[n.Name()] = mapelemNew(n, p, image.Point{x.Hint.X, x.Hint.Y}, m)
	}
	for _, x := range xmlm.Mappings {
		n, ok = m.graph.ItsType().NodeByName(x.Name)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: node %s not in graph %s\n", x.Name, m.graph.Filename())
			continue
		}
		p, ok = m.platform.ProcessByName(x.Process)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: process %s not in platform %s\n", x.Process, m.platform.Filename())
			continue
		}
		m.maps[n.Name()] = mapelemNew(n, p, image.Point{x.Hint.X, x.Hint.Y}, m)
	}
	return
}
