package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp/behaviour"
	"github.com/axel-freesp/sge/freesp/mapping"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"log"
)

type fileManagerMap struct {
	filenameFactory
	context        FilemanagerContextIf
	mappingMap     map[string]mp.MappingIf
	graphForNew    bh.SignalGraphIf
	platformForNew pf.PlatformIf
	// TODO: Semaphore to lock against concurrent New() calls
}

var _ mod.FileManagerMappingIf = (*fileManagerMap)(nil)

func FileManagerMapNew(context FilemanagerContextIf) *fileManagerMap {
	return &fileManagerMap{FilenameFactoryInit("mml"), context, make(map[string]mp.MappingIf), nil, nil}
}

func (f *fileManagerMap) SetGraphForNew(g interface{}) {
	f.graphForNew = g.(bh.SignalGraphIf)
}

func (f *fileManagerMap) SetPlatformForNew(p interface{}) {
	f.platformForNew = p.(pf.PlatformIf)
}

//
//      FileManagerIf interface
//

func (f *fileManagerMap) New() (m tr.ToplevelTreeElement, err error) {
	if f.graphForNew == nil || f.platformForNew == nil {
		err = fmt.Errorf("fileManagerMap.New: not completely prepared.\n")
		return
	}
	filename := f.NewFilename()
	m = mapping.MappingNew(filename, f.context)
	m.(mp.MappingIf).SetGraph(f.graphForNew)
	m.(mp.MappingIf).SetPlatform(f.platformForNew)
	f.graphForNew = nil
	f.platformForNew = nil
	var newId string
	newId, err = f.context.FTS().AddToplevel(m.(mp.MappingIf))
	if err != nil {
		err = fmt.Errorf("fileManagerMap.New: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)
	var mv views.GraphViewIf
	mv, err = views.MappingViewNew(m.(mp.MappingIf), f.context)
	if err != nil {
		log.Fatal("fileNewMap: Could not create graph view.")
	}
	f.context.GVC().Add(mv, filename)
	f.context.ShowAll()
	f.mappingMap[filename] = m.(mp.MappingIf)
	return
}

func (f *fileManagerMap) Access(name string) (m tr.ToplevelTreeElement, err error) {
	var ok bool
	m, ok = f.mappingMap[name]
	if ok {
		return
	}
	m = mapping.MappingNew(name, f.context)
	var filedir string
	for _, filedir = range backend.XmlSearchPaths() {
		err = m.ReadFile(fmt.Sprintf("%s/%s", filedir, name))
		if err != nil {
			log.Printf("fileManagerMap.Access error: %s\n", err)
		}
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: mapping file %s not found.", name)
		return
	}
	var filename string
	if len(filedir) > 0 {
		filename = fmt.Sprintf("%s/%s", filedir, name)
	} else {
		filename = name
	}
	m.SetPathPrefix(filedir)
	// TODO: check consistency with graph/platform
	//mapping.AddMapping(n, nId, nil)
	//melem, _ = mapping.MappedElement(nId)
	err = f.Consistent(m.(mp.MappingIf))
	if err != nil {
		return
	}

	hint := backend.XmlMappingHintNew(name)
	hintfilename := f.HintFilename(filename)
	var buf []byte
	buf, err = tool.ReadFile(hintfilename)
	if err == nil {
		_, err = hint.Read(buf)
		if err != nil {
			return
		}
		err = mapping.MappingApplyHints(m.(mp.MappingIf), hint)
		if err != nil {
			err = fmt.Errorf("fileManagerMap.Access: %s", err)
			return
		}
	}
	var newId string
	newId, err = f.context.FTS().AddToplevel(m.(mp.MappingIf))
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)
	var mv views.GraphViewIf
	mv, err = views.MappingViewNew(m.(mp.MappingIf), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: Could not create platform view.")
		return
	}
	f.context.GVC().Add(mv, name)
	log.Printf("fileManagerMap.Access: platform %s successfully loaded.\n", name)
	f.mappingMap[name] = m.(mp.MappingIf)
	return
}

func findNodeInIdList(n bh.NodeIf, parent bh.NodeIdIf, idlist []bh.NodeIdIf) bool {
	nId := behaviour.NodeIdNew(parent, n.Name())
	log.Printf("findNodeInIdList: n=%s, parent=%v\n", n.Name(), parent)
	for _, id := range idlist {
		if nId.String() == id.String() {
			return true
		}
	}
	for _, impl := range n.ItsType().Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, nn := range impl.Graph().ProcessingNodes() {
				if !findNodeInIdList(nn, nId, idlist) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func addExpandedMappings(m mp.MappingIf, n bh.NodeIf, parentId bh.NodeIdIf) {
	idlist := m.MappedIds()
	nId := behaviour.NodeIdNew(parentId, n.Name())
	for _, id := range idlist {
		if nId.String() == id.String() {
			return
		}
	}
	melem := m.AddMapping(n, nId, nil)
	melem.SetExpanded(true)
	for _, impl := range n.ItsType().Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, nn := range impl.Graph().ProcessingNodes() {
				addExpandedMappings(m, nn, nId)
			}
		}
	}
}

func (f *fileManagerMap) Consistent(m mp.MappingIf) (err error) {
	nodelist := m.Graph().ItsType().Nodes()
	idlist := m.MappedIds()
	for _, n := range nodelist {
		// check that every node has a mapping
		if !findNodeInIdList(n, behaviour.NodeIdFromString("", m.Graph().Filename()), idlist){
			// TODO: add missing mapping(s)
			err = fmt.Errorf("fileManagerMap.Consistent error: node %s is not fully mapped\n", n.Name())
			return
		}
		addExpandedMappings(m, n, behaviour.NodeIdFromString("", m.Graph().Filename()))
	}
	for _, id := range idlist {
		// check that every mapping has a node
		melem, ok := m.MappedElement(id)
		if !ok {
			// Should never occur, nothing to recover though...
			err = fmt.Errorf("fileManagerMap.Consistent error: invalid maplist\n")
			return
		}
		n := melem.Node()
		if n == nil {
			// TODO: remove this mapping from mappings
			err = fmt.Errorf("fileManagerMap.Consistent error: mapping %s has no node\n", id)
			return
		}
	}
	return
}

func (f *fileManagerMap) Remove(name string) {
	m, ok := f.mappingMap[name]
	if !ok {
		log.Printf("fileManagerMap.Remove warning: mapping %s not found.\n", name)
		return
	}
	delete(f.mappingMap, name)
	id, _ := f.context.FTS().GetToplevelId(m)
	f.context.FTS().RemoveToplevel(id)
	f.context.GVC().RemoveMappingView(m)
	// TODO: remove depending graphs and platforms if not used otherwise
}

func (f *fileManagerMap) Rename(oldName, newName string) (err error) {
	m, ok := f.mappingMap[oldName]
	if !ok {
		err = fmt.Errorf("fileManagerMap.Rename error: mapping %s not found\n", oldName)
		return
	}
	_, ok = f.mappingMap[newName]
	if ok {
		err = fmt.Errorf("fileManagerMap.Rename error: cannot rename mapping %s to be %s: already exists\n", oldName, newName)
		return
	}
	id, _ := f.context.FTS().GetToplevelId(m)
	f.context.FTS().SetValueById(id, newName)
	f.context.GVC().Rename(oldName, newName)
	delete(f.mappingMap, oldName)
	m.SetFilename(newName)
	f.mappingMap[newName] = m
	return
}

func (f *fileManagerMap) Store(name string) (err error) {
	m, ok := f.mappingMap[name]
	if !ok {
		err = fmt.Errorf("fileManagerMap.Store error: mapping %s not found.\n", name)
		return
	}
	context := m.Graph().ItsType().Context()
	gMgr := context.SignalGraphMgr()
	err = gMgr.Store(m.Graph().Filename())
	if err != nil {
		log.Printf("fileManagerMap.Store WARNING: could not save graph file %s\n", m.Graph().Filename())
	}
	pMgr := context.PlatformMgr()
	err = pMgr.Store(m.Platform().Filename())
	if err != nil {
		log.Printf("fileManagerMap.Store WARNING: could not save platform file %s\n", m.Platform().Filename())
	}

	var filename string
	if len(m.PathPrefix()) == 0 {
		filename = m.Filename()
	} else {
		filename = fmt.Sprintf("%s/%s", m.PathPrefix(), m.Filename())
	}
	err = m.WriteFile(filename)
	if err != nil {
		return
	}
	hint := mapping.CreateXmlMappingHint(m)
	hintfilename := f.HintFilename(filename)
	var buf []byte
	buf, err = hint.Write()
	if err != nil {
		return
	}
	err = tool.WriteFile(hintfilename, buf)
	return
}
