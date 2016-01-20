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
	libraryMap map[string]freesp.Library
}

var _ freesp.FileManagerIf = (*fileManagerLib)(nil)

func FileManagerLibNew(context FilemanagerContextIf) *fileManagerLib {
	return &fileManagerLib{FilenameFactoryInit("alml"), context, make(map[string]freesp.Library)}
}

//
//      FileManagerIf interface
//

func (f *fileManagerLib) New() (lib freesp.FileDataIf, err error) {
	filename := f.NewFilename()
	lib = freesp.LibraryNew(filename, f.context)
	f.libraryMap[filename] = lib.(freesp.Library)
	var cursor freesp.Cursor
	cursor, err = f.context.FTS().AddLibraryFile(lib.(freesp.Library))
	if err != nil {
		err = fmt.Errorf("fileManagerLib.New: %s.\n", err)
		return
	}
	f.context.FTV().SetCursorNewId(cursor.Path)
	return
}

func (f *fileManagerLib) Access(name string) (lib freesp.FileDataIf, err error) {
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
	f.libraryMap[name] = lib.(freesp.Library)
	_, err = f.context.FTS().AddLibraryFile(lib.(freesp.Library))
	if err != nil {
		err = fmt.Errorf("fileManagerLib.Access: %s", err)
		return
	}
	log.Println("fileManagerLib.Access: library", name, "successfully loaded")
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
	cursor := f.context.FTS().Cursor(lib)
	f.context.FTS().RemoveToplevel(cursor.Path)
}

func (f *fileManagerLib) Rename(oldName, newName string) {
	lib, ok := f.libraryMap[oldName]
	if !ok {
		log.Printf("fileManagerLib.Remove warning: library %s not found\n", oldName)
		return
	}
	delete(f.libraryMap, oldName)
	lib.SetFilename(newName)
	f.libraryMap[newName] = lib
}
