package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
)

var (
	currentDir string
)

func MenuFileInit(menu *GoAppMenu) {
	menu.fileNewSg.Connect("activate", func() { fileNewSg(global.fts, global.ftv) })
	menu.fileNewLib.Connect("activate", func() { fileNewLib(global.fts, global.ftv) })
	menu.fileOpen.Connect("activate", func() { fileOpen(global.fts, global.ftv) })
	menu.fileSave.Connect("activate", func() { fileSave(global.fts) })
	menu.fileSaveAs.Connect("activate", func() { fileSaveAs(global.fts) })

}

/*
 *		Callbacks
 */

var sgFilenameIndex = 0

func newSGFilename() string {
	ret := fmt.Sprintf("new-file-%d.sml", sgFilenameIndex)
	sgFilenameIndex++
	return ret
}

func fileNewSg(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	log.Println("fileNewSg")
	filename := newSGFilename()
	sg := freesp.SignalGraphNew(filename, &global)
	newId, err := fts.AddSignalGraphFile(sg.Filename(), sg)
	if err != nil {
		log.Printf("Warning: ftv.AddSignalGraphFile('%s') failed.\n", filename)
	}
	setCursorNewId(ftv, newId)
	gv, err := views.GraphViewNew(sg.ItsType(), width, height, &global)
	if err != nil {
		log.Fatal("fileNewSg: Could not create graph view.")
	}
	global.graphview = append(global.graphview, gv)
	global.win.stack.AddTitled(gv.Widget(), filename, filename)
	global.win.Window().ShowAll()
}

func setCursorNewId(ftv *views.FilesTreeView, newId string) {
	path, err := gtk.TreePathNewFromString(newId)
	if err != nil {
		log.Println("expandNewId error: TreePathNewFromString failed:", err)
		return
	}
	ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
}

func fileNewLib(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	log.Println("fileNewLib")
	cursor, err := fts.AddLibraryFile("new-file.alml", freesp.LibraryNew("new-file.alml", &global))
	if err != nil {
		log.Println("Warning: ftv.AddLibraryFile('new-file.alml') failed.")
	}
	setCursorNewId(ftv, cursor.Path)
}

func fileOpen(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	log.Println("fileOpen")
	filename, ok := getFilenameToOpen()
	if !ok {
		return
	}
	switch tool.Suffix(filename) {
	case "sml":
		sg := freesp.SignalGraphNew(filenameToShow(filename), &global)
		err := sg.ReadFile(filename)
		if err != nil {
			log.Println(err)
			return
		}
		//sg.SetFilename(filenameToShow(filename))
		newId, err := fts.AddSignalGraphFile(filename, sg)
		if err != nil {
			log.Println(err)
			return
		}
		setCursorNewId(ftv, newId)
		gv, err := views.GraphViewNew(sg.ItsType(), width, height, &global)
		if err != nil {
			log.Fatal("fileOpen: Could not create graph view.")
		}
		global.graphview = append(global.graphview, gv)
		global.win.stack.AddTitled(gv.Widget(), filenameToShow(filename), filenameToShow(filename))
		global.win.Window().ShowAll()
	case "alml":
		lib, err := global.GetLibrary(filenameToShow(filename))
		if err != nil {
			log.Println(err)
			return
		}
		cursor := fts.Cursor(lib)
		setCursorNewId(ftv, cursor.Path)
	case "spml":
		p := freesp.PlatformNew(filenameToShow(filename))
		err := p.ReadFile(filename)
		if err != nil {
			log.Println(err)
			return
		}
		newId, err := fts.AddPlatformFile(filename, p)
		if err != nil {
			log.Println(err)
			return
		}
		setCursorNewId(ftv, newId)
	default:
	}
}

func fileSaveAs(fts *models.FilesTreeStore) {
	log.Println("fileSaveAs")
	filename, ok := getFilenameToSave(getFilenameProposal(fts))
	if !ok {
		return
	}
	obj := getCurrentTopObject(fts)
	switch obj.(type) {
	case freesp.SignalGraph:
		obj.(freesp.SignalGraph).SetFilename(filenameToShow(filename))
	case freesp.Library:
		oldName := obj.(freesp.Library).Filename()
		obj.(freesp.Library).SetFilename(filenameToShow(filename))
		global.RenameLibrary(oldName, filenameToShow(filename))
	default:
		log.Fatalf("fileSaveAs error: wrong type '%T' of toplevel object (%v)\n", obj, obj)
	}
	err := doSave(fts, filename, obj)
	if err != nil {
		log.Println(err)
	}
}

func fileSave(fts *models.FilesTreeStore) {
	log.Println("fileSave")
	if len(currentDir) == 0 {
		currentDir = backend.XmlRoot()
	}
	var filename string
	obj := getCurrentTopObject(fts)
	ok := true
	switch obj.(type) {
	case freesp.SignalGraph:
		filename = obj.(freesp.SignalGraph).Filename()
		if filename == "new-file.sml" {
			filename, ok = getFilenameToSave(fmt.Sprintf("%s/%s", currentDir, filename))
		} else if !tool.IsSubPath("/", filename) {
			filename = fmt.Sprintf("%s/%s", backend.XmlRoot(), filename)
		}
	case freesp.Library:
		filename = obj.(freesp.Library).Filename()
		if filename == "new-file.alml" {
			filename, ok = getFilenameToSave(fmt.Sprintf("%s/%s", currentDir, filename))
		} else if !tool.IsSubPath("/", filename) {
			filename = fmt.Sprintf("%s/%s", backend.XmlRoot(), filename)
		}
	default:
		log.Fatalf("fileSave error: wrong type '%T' of toplevel object (%v)\n", obj, obj)
	}
	if !ok {
		return
	}
	err := doSave(fts, filename, obj)
	if err != nil {
		log.Println(err)
	}
	return
}

/*
 *		Local functions
 */

func doSave(fts *models.FilesTreeStore, filename string, obj interface{}) (err error) {
	relpath := tool.RelPath(backend.XmlRoot(), filename)
	log.Println("doSave: filename =", filename)
	log.Println("doSave: relative filename =", relpath)
	switch obj.(type) {
	case freesp.SignalGraph:
		err = obj.(freesp.SignalGraph).WriteFile(filename)
		if err != nil {
			return
		}
		obj.(freesp.SignalGraph).SetFilename(relpath)
	case freesp.Library:
		err = obj.(freesp.Library).WriteFile(filename)
		if err != nil {
			return
		}
		obj.(freesp.Library).SetFilename(relpath)
	}
	setCurrentTopValue(fts, relpath)
	if tool.IsSubPath(backend.XmlRoot(), filename) {
		currentDir = tool.Dirname(filename)
	}
	return
}

func getFilenameToSave(proposed string) (filename string, ok bool) {
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Choose file to save",
		nil,
		gtk.FILE_CHOOSER_ACTION_SAVE,
		"Cancel",
		gtk.RESPONSE_CANCEL,
		"Save",
		gtk.RESPONSE_OK)
	if err != nil {
		log.Fatal(err)
	}
	dialog.SetCurrentFolder(tool.Dirname(proposed))
	dialog.SetCurrentName(tool.Basename(proposed))
	response := dialog.Run()
	ok = (gtk.ResponseType(response) == gtk.RESPONSE_OK)
	if ok {
		filename = dialog.GetFilename()
	}
	dialog.Destroy()
	return
}

func getFilenameToOpen() (filename string, ok bool) {
	if len(currentDir) == 0 {
		currentDir = backend.XmlRoot()
	}
	dialog, err := gtk.FileChooserDialogNewWith2Buttons("Choose file to open",
		nil,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel",
		gtk.RESPONSE_CANCEL,
		"Open",
		gtk.RESPONSE_OK)
	if err != nil {
		log.Fatal(err)
	}
	dialog.SetCurrentFolder(currentDir)
	response := dialog.Run()
	ok = (gtk.ResponseType(response) == gtk.RESPONSE_OK)
	if ok {
		filename = dialog.GetFilename()
	}
	dialog.Destroy()
	return
}

func getToplevelId(fts *models.FilesTreeStore) string {
	p := fts.GetCurrentId()
	if p == "" {
		log.Fatal("fileSave error: fts.GetCurrentId() failed")
	}
	spl := strings.Split(p, ":")
	return spl[0] // TODO: move to function in fts...
}

func setCurrentTopValue(fts *models.FilesTreeStore, value string) {
	id0 := getToplevelId(fts)
	fts.SetValueById(id0, value)
}

func getCurrentTopObject(fts *models.FilesTreeStore) interface{} {
	id0 := getToplevelId(fts)
	obj, err := fts.GetObjectById(id0)
	if err != nil {
		log.Fatal("fileSave error: fts.GetObjectByPath failed:", err)
	}
	return obj
}

func getFilenameProposal(fts *models.FilesTreeStore) (filename string) {
	if len(currentDir) == 0 {
		currentDir = backend.XmlRoot()
	}
	obj := getCurrentTopObject(fts)
	switch obj.(type) {
	case freesp.SignalGraph:
		filename = obj.(freesp.SignalGraph).Filename()
	case freesp.Library:
		filename = obj.(freesp.Library).Filename()
	default:
		log.Fatal("fileSave error: wrong type '%T' of toplevel object (%v)", obj, obj)
	}
	if !tool.IsSubPath("/", filename) {
		filename = fmt.Sprintf("%s/%s", backend.XmlRoot(), filename)
	}
	return
}

func filenameToShow(filepath string) (filename string) {
	if tool.IsSubPath(currentDir, filepath) {
		filename = tool.RelPath(currentDir, filepath)
	} else {
		filename = filepath
	}
	return
}
