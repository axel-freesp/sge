package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/views"
	"log"
)

type fileManagerMap struct {
	filenameFactory
	context    FilemanagerContextIf
	mappingMap map[string]freesp.Mapping
}

var _ freesp.FileManagerIf = (*fileManagerMap)(nil)

func FileManagerMapNew(context FilemanagerContextIf) *fileManagerMap {
	return &fileManagerMap{FilenameFactoryInit("mml"), context, make(map[string]freesp.Mapping)}
}

//
//      FileManagerIf interface
//

func (f *fileManagerMap) New() (m freesp.FileDataIf, err error) {
	filename := f.NewFilename()
	m = freesp.MappingNew(filename, f.context)
	var graph freesp.SignalGraph
	var platform freesp.Platform
	// TODO: open dialog to choose graph and platform
	_ = graph
	_ = platform
	log.Printf("fileNewMap: not implemented\n")
	return

	var cursor freesp.Cursor
	cursor, err = f.context.FTS().AddMappingFile(m.(freesp.Mapping))
	if err != nil {
		err = fmt.Errorf("fileManagerMap.New: %s", err)
		return
	}
	f.context.FTV().SetCursorNewId(cursor.Path)
	var mv views.GraphView
	mv, err = views.MappingViewNew(m.(freesp.Mapping).MappingObject(), f.context)
	if err != nil {
		log.Fatal("fileNewMap: Could not create graph view.")
	}
	f.context.GVC().Add(mv, filename)
	f.context.ShowAll()
	f.mappingMap[filename] = m.(freesp.Mapping)
	return
}

func (f *fileManagerMap) Access(name string) (m freesp.FileDataIf, err error) {
	var ok bool
	m, ok = f.mappingMap[name]
	if ok {
		return
	}
	m = freesp.MappingNew(name, f.context)
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
	_, err = f.context.FTS().AddMappingFile(m.(freesp.Mapping))
	var mv views.GraphView
	mv, err = views.MappingViewNew(m.(freesp.Mapping).MappingObject(), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerMap.Access: Could not create platform view.")
		return
	}
	f.context.GVC().Add(mv, name)
	log.Println("fileManagerMap.Access: platform %s successfully loaded.", name)
	f.mappingMap[name] = m.(freesp.Mapping)
	return
}

func (f *fileManagerMap) Remove(name string) {
	m, ok := f.mappingMap[name]
	if !ok {
		log.Printf("fileManagerMap.Remove warning: mapping %s not found.\n", name)
		return
	}
	delete(f.mappingMap, name)
	cursor := f.context.FTS().Cursor(m)
	f.context.FTS().RemoveToplevel(cursor.Path)
	f.context.GVC().RemoveMappingView(m.MappingObject())
	// TODO: remove depending graphs and platforms if not used otherwise
}

func (f *fileManagerMap) Rename(oldName, newName string) {
	m, ok := f.mappingMap[oldName]
	if !ok {
		log.Printf("fileManagerMap.Rename warning: mapping %s not found\n", oldName)
		return
	}
	delete(f.mappingMap, oldName)
	m.SetFilename(newName)
	f.mappingMap[newName] = m
	// TODO: update GVC
}
