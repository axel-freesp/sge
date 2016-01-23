package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	interfaces "github.com/axel-freesp/sge/interface"
	"image"
	"log"
)

type mapping struct {
	graph    SignalGraphIf
	platform PlatformIf
	context  ContextIf
	maps     map[string]*mapelem
	filename string
}

// Ensure node and process are copies:
type mapelem struct {
	node     NodeIf
	process  ProcessIf
	position image.Point
	mapping  MappingIf
}

func mapelemNew(n NodeIf, p ProcessIf, nodePos image.Point, mapping MappingIf) (m *mapelem) {
	m = &mapelem{n, p, nodePos, mapping}
	return
}

var _ MappedElementIf = (*mapelem)(nil)

func (m mapelem) Position() image.Point {
	return m.position
}

func (m *mapelem) SetPosition(pos image.Point) {
	m.position = pos
}

func (m *mapelem) AddToTree(tree Tree, cursor Cursor) {
	var s Symbol
	if m.process == nil {
		s = SymbolUnmapped
	} else {
		s = SymbolMapped
	}
	err := tree.AddEntry(cursor, s, m.node.Name(), m, mayEdit)
	if err != nil {
		log.Fatalf("mapping.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (m mapelem) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatalf("mapping.AddNewObject error: Nothing to add\n")
	return
}

func (m mapelem) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatalf("mapping.RemoveObject error: Nothing to remove\n")
	return
}

func (m mapelem) NodeObject() interfaces.NodeObject {
	return m.node.(*node)
}

func (m mapelem) ProcessObject() (p interfaces.ProcessObject, ok bool) {
	ok = m.process != nil
	if ok {
		p = m.process.(*process)
	}
	return
}

func (m mapelem) Mapping() MappingIf {
	return m.mapping
}

func (m mapelem) Node() NodeIf {
	return m.node
}

func (m mapelem) Process() ProcessIf {
	return m.process
}

func (m *mapelem) SetProcess(p ProcessIf) {
	m.process = p
}

var _ interfaces.MappingObject = (*mapping)(nil)
var _ MappingIf = (*mapping)(nil)

func MappingNew(filename string, context ContextIf) *mapping {
	return &mapping{nil, nil, context, make(map[string]*mapelem), filename}
}

//
//		TreeElement interface
//

func (m *mapping) AddToTree(tree Tree, cursor Cursor) {
	var child Cursor
	err := tree.AddEntry(cursor, SymbolMappings, m.Filename(), m, mayAddObject)
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

func (m mapping) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatalf("mapping.AddNewObject error: Nothing to add\n")
	return
}

func (m mapping) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatalf("mapping.RemoveObject error: Nothing to remove\n")
	return
}

func (m *mapping) Identify(te TreeElement) bool {
	switch te.(type) {
	case *mapping:
		return te.(*mapping).Filename() == m.Filename()
	}
	return false
}

//
//		interfaces.MappingObject interface
//

func (m mapping) GraphObject() interfaces.GraphObject {
	return m.graph.GraphObject()
}

func (m mapping) PlatformObject() interfaces.PlatformObject {
	return m.platform.PlatformObject()
}

func (m mapping) MappedObject(n interfaces.NodeObject) (p interfaces.ProcessObject, ok bool) {
	var melem *mapelem
	melem, ok = m.maps[n.Name()]
	if ok {
		ok = melem.process != nil
		if ok {
			p = melem.process.(*process)
		}
	}
	return
}

func (m mapping) MapElemObject(n interfaces.NodeObject) (melem interfaces.MapElemObject, ok bool) {
	melem, ok = m.maps[n.Name()]
	return
}

//
//		Mapping interface
//

func (m *mapping) MappingObject() interfaces.MappingObject {
	return m
}

func (m *mapping) AddMapping(n NodeIf, p ProcessIf) {
	m.maps[n.Name()] = mapelemNew(n, p, image.Point{}, m)
}

func (m *mapping) SetGraph(g SignalGraphIf) {
	m.graph = g
}

func (m *mapping) Graph() SignalGraphIf {
	return m.graph
}

func (m *mapping) SetPlatform(p PlatformIf) {
	m.platform = p
}

func (m *mapping) Platform() PlatformIf {
	return m.platform
}

func (m *mapping) Mapped(n NodeIf) (pr ProcessIf, ok bool) {
	var melem *mapelem
	melem, ok = m.maps[n.Name()]
	if ok {
		pr = melem.process
	}
	ok = pr != nil
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

func (m *mapping) RemoveFromTree(tree Tree) {
	tree.Remove(tree.Cursor(m))
}

//
//		Local Methods and Functions
//

func (m *mapping) createMappingFromXml(xmlm *backend.XmlMapping) (err error) {
	var f FileDataIf
	f, err = m.context.SignalGraphMgr().Access(xmlm.SignalGraph)
	if err != nil {
		err = fmt.Errorf("mapping.CreateMappingFromXml: graph.ReadFile failed: %s\n", err)
		return
	}
	m.graph = f.(SignalGraphIf)
	f, err = m.context.PlatformMgr().Access(xmlm.Platform)
	if err != nil {
		err = fmt.Errorf("mapping.CreateMappingFromXml: platform.ReadFile failed: %s\n", err)
		return
	}
	m.platform = f.(PlatformIf)
	var n NodeIf
	var p ProcessIf
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
