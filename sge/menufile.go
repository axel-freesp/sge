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
	menu.fileNewPlat.Connect("activate", func() { fileNewPlat(global.fts, global.ftv) })
	menu.fileNewMap.Connect("activate", func() { fileNewMap(global.fts, global.ftv) })
	menu.fileOpen.Connect("activate", func() { fileOpen(global.fts, global.ftv) })
	menu.fileSave.Connect("activate", func() { fileSave(global.fts) })
	menu.fileSaveAs.Connect("activate", func() { fileSaveAs(global.fts) })
	menu.fileClose.Connect("activate", func() { fileClose(menu, global.fts, global.ftv, global.jl) })

}

/*
 *		Callbacks
 */

func fileNewSg(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.SignalGraphMgr().New()
	if err != nil {
		log.Printf("fileNewSg: %s\n", err)
	}
}

func fileNewLib(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.LibraryMgr().New()
	if err != nil {
		log.Printf("fileNewLib: %s\n", err)
	}
}

func fileNewPlat(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.PlatformMgr().New()
	if err != nil {
		log.Printf("fileNewPlat: %s\n", err)
	}
}

func fileNewMap(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	_, err := global.MappingMgr().New()
	if err != nil {
		log.Printf("fileNewMap: %s\n", err)
	}
}

func fileOpen(fts *models.FilesTreeStore, ftv *views.FilesTreeView) {
	log.Println("fileOpen")
	filename, ok := getFilenameToOpen()
	if !ok {
		return
	}
	var err error
	switch tool.Suffix(filename) {
	case "sml":
		_, err = global.SignalGraphMgr().Access(filenameToShow(filename))
		if err != nil {
			log.Printf("fileOpen FIXME: could not get signal graph %s. Try full path!\n", filenameToShow(filename))
		}
	case "alml":
		_, err = global.LibraryMgr().Access(filenameToShow(filename))
	case "spml":
		_, err = global.PlatformMgr().Access(filenameToShow(filename))
	case "mml":
		_, err = global.MappingMgr().Access(filenameToShow(filename))
	default:
		return
	}
	if err != nil {
		log.Printf("fileOpen: %s\n", err)
	}
}

func fileSaveAs(fts *models.FilesTreeStore) {
	log.Println("fileSaveAs")
	filename, ok := getFilenameToSave(getFilenameProposal(fts))
	if !ok {
		return
	}
	obj := getCurrentTopObject(fts)
	oldName := obj.(freesp.FileDataIf).Filename()
	global.FileMgr(obj).Rename(oldName, filenameToShow(filename))
	err := doSave(fts, filename, obj)
	if err != nil {
		log.Printf("fileSaveAs: %s\n", err)
	}
}

func isGeneratedFilename(filename string) bool {
	return strings.HasPrefix(filename, "new-file-")
}

func fileSave(fts *models.FilesTreeStore) {
	log.Println("fileSave")
	if len(currentDir) == 0 {
		currentDir = backend.XmlRoot()
	}
	obj := getCurrentTopObject(fts)
	filename := obj.(freesp.FileDataIf).Filename()
	if isGeneratedFilename(filename) {
		oldName := filename
		ok := true
		filename, ok = getFilenameToSave(fmt.Sprintf("%s/%s", currentDir, filename))
		if !ok {
			return
		}
		global.FileMgr(obj).Rename(oldName, filenameToShow(filename))
	} else if !tool.IsSubPath("/", filename) {
		filename = fmt.Sprintf("%s/%s", backend.XmlRoot(), filename)
	}
	err := doSave(fts, filename, obj)
	if err != nil {
		log.Println(err)
	}
	return
}

func fileClose(menu *GoAppMenu, fts *models.FilesTreeStore, ftv *views.FilesTreeView, jl IJobList) {
	path := fts.GetCurrentId()
	if strings.Contains(path, ":") {
		path = strings.Split(path, ":")[0]
	}
	obj, err := fts.GetObjectById(path)
	if err != nil {
		return
	}
	global.FileMgr(obj).Remove(obj.(freesp.Filenamer).Filename())
	jl.Reset()
	MenuEditPost(menu, fts, jl)
}

/*
 *		Local functions
 */

func doSave(fts *models.FilesTreeStore, filename string, obj interface{}) (err error) {
	relpath := tool.RelPath(backend.XmlRoot(), filename)
	log.Println("doSave: filename =", filename)
	log.Println("doSave: relative filename =", relpath)
	err = obj.(freesp.FileDataIf).WriteFile(filename)
	if err != nil {
		return
	}
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

func getCurrentTopObject(fts *models.FilesTreeStore) freesp.TreeElement {
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
	case freesp.Platform:
		filename = obj.(freesp.Platform).Filename()
	case freesp.Mapping:
		filename = obj.(freesp.Mapping).Filename()
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
