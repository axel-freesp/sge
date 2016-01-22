package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/views"
	"log"
)

type fileManagerPF struct {
	filenameFactory
	context     FilemanagerContextIf
	platformMap map[string]freesp.Platform
}

var _ freesp.FileManagerIf = (*fileManagerPF)(nil)

func FileManagerPFNew(context FilemanagerContextIf) *fileManagerPF {
	return &fileManagerPF{FilenameFactoryInit("spml"), context, make(map[string]freesp.Platform)}
}

//
//      FileManagerIf interface
//

func (f *fileManagerPF) New() (pl freesp.ToplevelTreeElement, err error) {
	filename := f.NewFilename()
	pl = freesp.PlatformNew(filename)
	f.platformMap[filename] = pl.(freesp.Platform)
	var newId string
	newId, err = f.context.FTS().AddPlatformFile(pl.(freesp.Platform))
	if err != nil {
		err = fmt.Errorf("fileManagerPF.New: %s", err)
		return
	}
	f.context.FTV().SetCursorNewId(newId)
	var pv views.GraphView
	pv, err = views.PlatformViewNew(pl.(freesp.Platform).PlatformObject(), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerPF.New: %s", err)
		return
	}
	f.context.GVC().Add(pv, filename)
	f.context.ShowAll()
	return
}

func (f *fileManagerPF) Access(name string) (pl freesp.ToplevelTreeElement, err error) {
	var ok bool
	pl, ok = f.platformMap[name]
	if ok {
		return
	}
	pl = freesp.PlatformNew(name)
	for _, try := range backend.XmlSearchPaths() {
		err = pl.ReadFile(fmt.Sprintf("%s/%s", try, name))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("fileManagerPF.Access: platform file %s not found", name)
		return
	}
	pv, err := views.PlatformViewNew(pl.(freesp.Platform).PlatformObject(), f.context)
	if err != nil {
		err = fmt.Errorf("fileManagerPF.Access: Could not create platform view.")
		return
	}
	f.context.GVC().Add(pv, name)
	log.Printf("fileManagerPF.Access: platform %s successfully loaded.\n", name)
	f.platformMap[name] = pl.(freesp.Platform)
	var newId string
	newId, err = f.context.FTS().AddPlatformFile(pl.(freesp.Platform))
	if err != nil {
		err = fmt.Errorf("fileManagerPF.Access: %s", err)
		return
	}
	f.context.FTV().SetCursorNewId(newId)
	return
}

func (f *fileManagerPF) Remove(name string) {
	// TODO: check if used in existing mappings
	pl, ok := f.platformMap[name]
	if !ok {
		log.Printf("fileManagerPF.Remove warning: platform %s not found.\n", name)
		return
	}
	for _, a := range pl.Arch() {
		for _, iot := range a.IOTypes() {
			freesp.RemoveRegisteredIOType(iot)
		}
	}
	delete(f.platformMap, name)
	cursor := f.context.FTS().Cursor(pl)
	f.context.FTS().RemoveToplevel(cursor.Path)
	f.context.GVC().RemovePlatformView(pl.PlatformObject())
}

func (f *fileManagerPF) Rename(oldName, newName string) (err error) {
	pl, ok := f.platformMap[oldName]
	if !ok {
		err = fmt.Errorf("fileManagerPF.Rename error: platform %s not found\n", oldName)
		return
	}
	_, ok = f.platformMap[newName]
	if ok {
		err = fmt.Errorf("fileManagerPF.Rename error: cannot rename platform %s to be %s: already exists\n", oldName, newName)
		return
	}
	cursor := f.context.FTS().Cursor(pl)
	f.context.FTS().SetValueById(cursor.Path, newName)
	f.context.GVC().Rename(oldName, newName)
	delete(f.platformMap, oldName)
	pl.SetFilename(newName)
	f.platformMap[newName] = pl
	return
}
