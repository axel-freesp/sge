package main

import (
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func MenuEditInit(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	menu.editUndo.Connect("activate", func() { editUndo(fts, jl) })
	menu.editRedo.Connect("activate", func() { editRedo(fts, jl) })
	menu.editNew.Connect("activate", func() { editNew(fts, jl, ftv) })
	menu.editCopy.Connect("activate", func() { editCopy(fts) })
	menu.editDelete.Connect("activate", func() { editDelete(fts, jl, ftv) })
}

func editUndo(fts *models.FilesTreeStore, jl IJobList) {
	log.Println("editUndo")
	jl.Undo()
}

func editRedo(fts *models.FilesTreeStore, jl IJobList) {
	log.Println("editRedo")
	jl.Redo()
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
			path, err := gtk.TreePathNewFromString(job.newElement.newId)
			if err != nil {
				log.Println("editNew error: TreePathNewFromString failed:", err)
				return
			}
			ftv.TreeView().ExpandToPath(path)
			ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
		}
	}
	log.Println("editNew finished")
}

func editCopy(fts *models.FilesTreeStore) {
	log.Println("editCopy")
}

func editDelete(fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editDelete")
	job := DeleteObjectJobNew(fts.GetCurrentId())
	if jl.Apply(EditorJobNew(JobDeleteObject, job)) {
		log.Printf("Deleted %d objects\n", len(job.deletedObjects))
		var parentId string
		for _, d := range job.deletedObjects {
			log.Println(d)
			parentId = d.ParentId
		}
		path, err := gtk.TreePathNewFromString(parentId)
		if err != nil {
			log.Println("editNew error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}
