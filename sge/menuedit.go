package main

import (
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func editUndo(fts *models.FilesTreeStore) {
	log.Println("editUndo")
}

func editRedo(fts *models.FilesTreeStore) {
	log.Println("editRedo")
}

func editNew(fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editNew")
	dialog, err := NewElementDialogNew(fts)
	if err != nil {
		log.Println("editNew error: ", err)
		return
	}
	job, ok := dialog.Run(fts)
	if ok {
		if jl.Apply(job) {
			id, err := gtk.TreePathNewFromString(job.newElement.newId)
			if err != nil {
				log.Println("editNew error: TreePathNewFromString failed:", err)
				return
			}
			ftv.TreeView().SetCursor(id, ftv.TreeView().GetExpanderColumn(), false)
		}
	}
	log.Println("editNew finished")
}

func editCopy(fts *models.FilesTreeStore) {
	log.Println("editCopy")
}

func editDelete(fts *models.FilesTreeStore) {
	log.Println("editDelete")
}

func MenuEditInit(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	menu.editUndo.Connect("activate", func() { editUndo(fts) })
	menu.editRedo.Connect("activate", func() { editRedo(fts) })
	menu.editNew.Connect("activate", func() { editNew(fts, jl, ftv) })
	menu.editCopy.Connect("activate", func() { editCopy(fts) })
	menu.editDelete.Connect("activate", func() { editDelete(fts) })
}
