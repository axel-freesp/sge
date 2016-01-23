package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp/mapping"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
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
	var mv views.GraphView
	mv, err = views.MappingViewNew(m.(mp.MappingIf).MappingObject(), f.context)
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
	for _, try := range backend.XmlSearchPaths() {
		log.Printf("fileManagerMap.Access: try path %s\n", try)
		err = m.ReadFile(fmt.Sprintf("%s/%s", try, name))
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
	var newId string
	newId, err = f.context.FTS().AddToplevel(m.(mp.MappingIf))
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)
	var mv views.GraphView
	mv, err = views.MappingViewNew(m.(mp.MappingIf).MappingObject(), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: Could not create platform view.")
		return
	}
	f.context.GVC().Add(mv, name)
	log.Println("fileManagerMap.Access: platform %s successfully loaded.", name)
	f.mappingMap[name] = m.(mp.MappingIf)
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
	f.context.GVC().RemoveMappingView(m.MappingObject())
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
