package main

import (
	"github.com/axel-freesp/sge/models"
	"log"
)

func editUndo(fts *models.FilesTreeStore) {
	log.Println("editUndo")
}

func editRedo(fts *models.FilesTreeStore) {
	log.Println("editRedo")
}

func editNew(fts *models.FilesTreeStore) {
	log.Println("editNew")
	dialog, err := NewElementDialogNew(fts)
	if err != nil {
		log.Println("editNew error: ", err)
		return
	}
	dialog.Run(fts)
	log.Println("editNew finished")
}

func editCopy(fts *models.FilesTreeStore) {
	log.Println("editCopy")
}

func editDelete(fts *models.FilesTreeStore) {
	log.Println("editDelete")
}

func MenuEditInit(menu *GoAppMenu, fts *models.FilesTreeStore) {
	menu.editUndo.Connect("activate", func() { editUndo(fts) })
	menu.editRedo.Connect("activate", func() { editRedo(fts) })
	menu.editNew.Connect("activate", func() { editNew(fts) })
	menu.editCopy.Connect("activate", func() { editCopy(fts) })
	menu.editDelete.Connect("activate", func() { editDelete(fts) })
}
