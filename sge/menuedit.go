package main

import (
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func MenuEditInit(menu *GoAppMenu) {
	fts := global.fts
	jl := global.jl
	ftv := global.ftv
	menu.editUndo.Connect("activate", func() { editUndo(menu, fts, jl, ftv) })
	menu.editRedo.Connect("activate", func() { editRedo(menu, fts, jl, ftv) })
	menu.editNew.Connect("activate", func() { editNew(menu, fts, jl, ftv) })
	menu.editCopy.Connect("activate", func() { editCopy(fts) })
	menu.editDelete.Connect("activate", func() { editDelete(menu, fts, jl, ftv) })
}

func MenuEditCurrent(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList) {
	cursor := freesp.Cursor{fts.GetCurrentId(), freesp.AppendCursor}
	prop := fts.Property(cursor)
	menu.editNew.SetSensitive(prop.MayAddObject())
	menu.editDelete.SetSensitive(prop.MayRemove())
	menu.editUndo.SetSensitive(jl.CanUndo())
	menu.editRedo.SetSensitive(jl.CanRedo())
	for _, v := range global.graphview {
		v.Sync()
	}
}

func editUndo(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editUndo")
	defer MenuEditCurrent(menu, fts, jl)
	state, ok := jl.Undo()
	if ok {
		path, err := gtk.TreePathNewFromString(state.(string))
		if err != nil {
			log.Println("editNew error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().ExpandToPath(path)
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}

func editRedo(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editRedo")
	defer MenuEditCurrent(menu, fts, jl)
	state, ok := jl.Redo()
	if ok {
		path, err := gtk.TreePathNewFromString(state.(string))
		if err != nil {
			log.Println("editNew error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().ExpandToPath(path)
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}

func editNew(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editNew")
	defer MenuEditCurrent(menu, fts, jl)
	dialog, err := NewElementDialogNew(fts)
	if err != nil {
		log.Println("editNew error: ", err)
		return
	}
	job, ok := dialog.Run(fts)
	if ok {
		state, ok := jl.Apply(job)
		if ok {
			path, err := gtk.TreePathNewFromString(state.(string))
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

func editDelete(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	log.Println("editDelete")
	defer MenuEditCurrent(menu, fts, jl)
	job := DeleteObjectJobNew(fts.GetCurrentId())
	state, ok := jl.Apply(EditorJobNew(JobDeleteObject, job))
	if ok {
		path, err := gtk.TreePathNewFromString(state.(string))
		if err != nil {
			log.Println("editNew error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().ExpandToPath(path)
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}
