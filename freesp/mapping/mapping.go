package mapping

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	//	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	fd "github.com/axel-freesp/sge/interface/filedata"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
)

type mapping struct {
	graph      bh.SignalGraphIf
	platform   pf.PlatformIf
	context    mod.ModelContextIf
	maps       map[string]*mapelem
	filename   string
	pathPrefix string
	maplist    behaviour.NodeIdList
}

var _ mp.MappingIf = (*mapping)(nil)

func MappingNew(filename string, context mod.ModelContextIf) *mapping {
	return &mapping{nil, nil, context, make(map[string]*mapelem), filename, "", behaviour.NodeIdListInit()}
}

func findMapHint(xmlhints *backend.XmlMappingHint, nId string) (xmlh backend.XmlNodePosHint, ok bool) {
	for _, xmlh = range xmlhints.MappedNodes {
		if xmlh.Name == nId {
			ok = true
			break
		}
	}
	return
}

func MappingApplyHints(m mp.MappingIf, xmlhints *backend.XmlMappingHint) (err error) {
	if m.Filename() != xmlhints.Ref {
		err = fmt.Errorf("MappingApplyHints error: filename mismatch\n")
		return
	}
	for _, nId := range m.MappedIds() {
		melem, ok := m.MappedElement(nId)
		if !ok {
			log.Fatal("MappingApplyHints internal error: inconsistent maplist\n")
		}
		xmlh, ok := findMapHint(xmlhints, nId.String())
		if ok {
			melem.SetExpanded(xmlh.Expanded)
			freesp.PathModePositionerApplyHints(melem, xmlh.XmlModeHint)
			log.Printf("MappingApplyHints(%s) applying: expanded=%v, pos=%v\n", nId.String(), xmlh.Expanded, melem.ModePosition(gr.PositionModeMapping))
			// TODO: ports
			for i, xmlp := range xmlh.InPorts {
				freesp.PathModePositionerApplyHints(&melem.(*mapelem).inports[i], xmlp.XmlModeHint)
			}
			for i, xmlp := range xmlh.OutPorts {
				freesp.PathModePositionerApplyHints(&melem.(*mapelem).outports[i], xmlp.XmlModeHint)
			}
		}
	}
	return

	/*
		g := m.Graph().ItsType()
		for _, xmln := range xmlhints.InputNodes {
			n, ok := g.NodeByName(xmln.Name)
			if ok {
				melem, ok := m.MappedElement(behaviour.NodeIdFromString(n.Name(), m.Graph().Filename()))
				if ok {
					melem.SetExpanded(false)
					freesp.ModePositionerApplyHints(melem, xmln.XmlModeHint)
				}
			}
		}
		for _, xmln := range xmlhints.OutputNodes {
			n, ok := g.NodeByName(xmln.Name)
			if ok {
				melem, ok := m.MappedElement(behaviour.NodeIdFromString(n.Name(), m.Graph().Filename()))
				if ok {
					melem.SetExpanded(false)
					freesp.ModePositionerApplyHints(melem, xmln.XmlModeHint)
				}
			}
		}
		NodesApplyHints("", m, g, xmlhints.MappedNodes)
		return
	*/
}

/*
func NodesApplyHints(path string, m mp.MappingIf, g bh.SignalGraphTypeIf, xmlh []backend.XmlNodePosHint) {
	for _, n := range g.ProcessingNodes() {
		var p string
		if len(path) == 0 {
			p = n.Name()
		} else {
			p = fmt.Sprintf("%s/%s", path, n.Name())
		}
		ok := false
		var xmln backend.XmlNodePosHint
		for _, xmln = range xmlh {
			if p == xmln.Name {
				ok = true
				break
			}
		}
		if ok {
			melem, ok := m.MappedElement(behaviour.NodeIdFromString(p, m.Graph().Filename()))
			if ok {
				melem.SetExpanded(xmln.Expanded)
				freesp.ModePositionerApplyHints(melem, xmln.XmlModeHint)
			}
		}
		nt := n.ItsType()
		for _, impl := range nt.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				NodesApplyHints(p, m, impl.Graph(), xmlh)
				break
			}
		}
	}
}
*/

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
	for _, nId := range m.MappedIds() {
		log.Printf("mapping.AddToTree: id=%s\n", nId.String())
		melem, ok := m.MappedElement(nId)
		if !ok {
			log.Fatal("mapping) AddToTree internal error: inconsistent maplist")
		}
		child = tree.Append(cursor)
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

func (m *mapping) AddMapping(n bh.NodeIf, nId bh.NodeIdIf, p pf.ProcessIf) {
	m.maps[nId.String()] = mapelemNew(n, nId, p, m)
	m.maplist.Append(nId)
}

func (m *mapping) SetGraph(g bh.SignalGraphIf) {
	m.graph = g
	log.Printf("mapping.SetGraph: TODO: any checks?\n")
}

func (m *mapping) Graph() bh.SignalGraphIf {
	return m.graph
}

func (m *mapping) SetPlatform(p pf.PlatformIf) {
	m.platform = p
	log.Printf("mapping.SetPlatform: TODO: any checks?\n")
}

func (m *mapping) Platform() pf.PlatformIf {
	return m.platform
}

func (m *mapping) Mapped(nId string) (pr pf.ProcessIf, ok bool) {
	var melem *mapelem
	melem, ok = m.maps[nId]
	if ok {
		pr = melem.process
	}
	ok = pr != nil
	return
}

func (m *mapping) MappedElement(nId bh.NodeIdIf) (melem mp.MappedElementIf, ok bool) {
	melem, ok = m.maps[nId.String()]
	return
}

func (m mapping) MappedIds() []bh.NodeIdIf {
	return m.maplist.NodeIds()
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
		nId := behaviour.NodeIdFromString(x.Name, m.Graph().Filename())
		m.maps[n.Name()] = mapelemNew(n, nId, p, m)
		m.maplist.Append(nId)
	}
	for _, x := range xmlm.Mappings {
		n, ok = m.graph.ItsType().NodeByPath(x.Name)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: node %s not in graph %s\n", x.Name, m.graph.Filename())
			continue
		}
		p, ok = m.platform.ProcessByName(x.Process)
		if !ok {
			err = fmt.Errorf("mapping.CreateMappingFromXml FIXME: process %s not in platform %s\n", x.Process, m.platform.Filename())
			continue
		}
		nId := behaviour.NodeIdFromString(x.Name, m.Graph().Filename())
		m.maps[n.Name()] = mapelemNew(n, nId, p, m)
		m.maplist.Append(nId)
	}
	return
}
