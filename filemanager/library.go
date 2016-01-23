package filemanager

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"log"
)

type fileManagerLib struct {
	filenameFactory
	context    FilemanagerContextIf
	libraryMap map[string]freesp.LibraryIf
}

var _ freesp.FileManagerIf = (*fileManagerLib)(nil)

func FileManagerLibNew(context FilemanagerContextIf) *fileManagerLib {
	return &fileManagerLib{FilenameFactoryInit("alml"), context, make(map[string]freesp.LibraryIf)}
}

//
//      FileManagerIf interface
//

func (f *fileManagerLib) New() (lib freesp.ToplevelTreeElement, err error) {
	filename := f.NewFilename()
	lib = freesp.LibraryNew(filename, f.context)
	f.libraryMap[filename] = lib.(freesp.LibraryIf)
	var newId string
	newId, err = f.context.FTS().AddToplevel(lib.(freesp.LibraryIf))
	if err != nil {
		err = fmt.Errorf("fileManagerLib.New: %s.\n", err)
		return
	}
	f.context.FTV().SelectId(newId)
	return
}

func (f *fileManagerLib) Access(name string) (lib freesp.ToplevelTreeElement, err error) {
	var ok bool
	lib, ok = f.libraryMap[name]
	if ok {
		return
	}
	lib = freesp.LibraryNew(name, f.context)
	for _, try := range backend.XmlSearchPaths() {
		err = lib.ReadFile(fmt.Sprintf("%s/%s", try, name))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("fileManagerLib.Access: library file %s not found", name)
		return
	}
	f.libraryMap[name] = lib.(freesp.LibraryIf)
	var newId string
	newId, err = f.context.FTS().AddToplevel(lib.(freesp.LibraryIf))
	if err != nil {
		err = fmt.Errorf("fileManagerLib.Access: %s", err)
		return
	}
	f.context.FTV().SelectId(newId)
	if err != nil {
		err = fmt.Errorf("fileManagerLib.Access: %s", err)
		return
	}
	log.Printf("fileManagerLib.Access: library %s successfully loaded\n", name)
	return
}

func (f *fileManagerLib) Remove(name string) {
	lib, ok := f.libraryMap[name]
	if !ok {
		log.Printf("fileManagerLib.Remove warning: library %s not found.\n", name)
		return
	}
	nodeTypes := lib.NodeTypes()
	signalTypes := lib.SignalTypes()
	for _, nt := range nodeTypes {
		if !f.context.NodeTypeIsInUse(nt) {
			f.context.CleanupNodeType(nt)
		}
	}
	for _, st := range signalTypes {
		if !f.context.SignalTypeIsInUse(st) {
			f.context.CleanupSignalType(st)
		}
	}
	delete(f.libraryMap, name)
	id, _ := f.context.FTS().GetToplevelId(lib)
	f.context.FTS().RemoveToplevel(id)
	log.Printf("fileManagerLib.Access: library %s successfully unloaded\n", name)
}

func (f *fileManagerLib) Rename(oldName, newName string) (err error) {
	lib, ok := f.libraryMap[oldName]
	if !ok {
		err = fmt.Errorf("fileManagerLib.Remove error: library %s not found\n", oldName)
		return
	}
	_, ok = f.libraryMap[newName]
	if ok {
		err = fmt.Errorf("fileManagerLib.Rename error: cannot rename library %s to be %s: already exists\n", oldName, newName)
		return
	}
	id, _ := f.context.FTS().GetToplevelId(lib)
	f.context.FTS().SetValueById(id, newName)
	delete(f.libraryMap, oldName)
	lib.SetFilename(newName)
	f.libraryMap[newName] = lib
	log.Printf("fileManagerLib.Access: library %s successfully renamed to %d\n", oldName, newName)
	return
}
