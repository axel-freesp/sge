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
	var path0 string
	if fts.CurrentSelection == nil {
		log.Fatal("fileSave error: CurrentSelection = nil")
	}
	path, err := fts.TreeStore().GetPath(fts.CurrentSelection)
	if err != nil {
		log.Fatal("fileSave error: iter.GetPath failed:", err)
		return
	}
	p := path.String()
	log.Println("Current selection: ", p)
	spl := strings.Split(p, ":")
	path0 = spl[0]
	iter, err := fts.TreeStore().GetIterFromString(path0)
	if err != nil {
		log.Fatal("fileSave error: fts.TreeStore().GetIterFromString failed:", err)
	}
	filename, err := fts.GetValue(iter)
	if err != nil {
		log.Fatal("fileSave error: fts.GetValue failed:", err)
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
