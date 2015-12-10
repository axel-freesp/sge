package main

import (
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"log"
	"strings"
)

func fileNewSg(fts *models.FilesTreeStore) {
	log.Println("fileNewSg")
	err := fts.AddSignalGraphFile("new-file.sml", freesp.SignalGraphNew("new-file.sml"))
	if err != nil {
		log.Println("Warning: ftv.AddSignalGraphFile('new-file.sml') failed.")
	}
}

func fileNewLib(fts *models.FilesTreeStore) {
	log.Println("fileNewLib")
	err := fts.AddLibraryFile("new-file.alml", freesp.LibraryNew("new-file.alml"))
	if err != nil {
		log.Println("Warning: ftv.AddLibraryFile('new-file.alml') failed.")
	}
}

func fileOpen(fts *models.FilesTreeStore) {
	log.Println("fileOpen")
}

func fileSaveAs(fts *models.FilesTreeStore) {
	log.Println("fileSaveAs")
}

func fileSave(fts *models.FilesTreeStore) {
	log.Println("fileSave")
	var id0, filename string
	p := fts.GetCurrentId()
	if p == "" {
		log.Fatal("fileSave error: fts.GetCurrentId() failed")
		return
	}
	log.Println("Current selection id: ", p)
	spl := strings.Split(p, ":")
	id0 = spl[0] // TODO: move to function in fts...

	obj, err := fts.GetObjectById(id0)
	if err != nil {
		log.Fatal("fileSave error: fts.GetObjectByPath failed:", err)
	}
	switch obj.(type) {
	case freesp.SignalGraph:
		filename = obj.(freesp.SignalGraph).Filename()
	case freesp.Library:
		filename = obj.(freesp.Library).Filename()
	default:
		log.Fatal("fileSave error: wrong type '%T' of toplevel object (%v)", obj, obj)
		return
	}
	log.Println("fileSave: filename =", filename)
}

func MenuFileInit(menu *GoAppMenu, fts *models.FilesTreeStore) {
	menu.fileNewSg.Connect("activate", func() { fileNewSg(fts) })
	menu.fileNewLib.Connect("activate", func() { fileNewLib(fts) })
	menu.fileOpen.Connect("activate", func() { fileOpen(fts) })
	menu.fileSave.Connect("activate", func() { fileSave(fts) })
	menu.fileSaveAs.Connect("activate", func() { fileSaveAs(fts) })

}
