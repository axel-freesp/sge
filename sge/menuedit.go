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
	clp := global.clp
	menu.editUndo.Connect("activate", func() { editUndo(menu, fts, jl, ftv) })
	menu.editRedo.Connect("activate", func() { editRedo(menu, fts, jl, ftv) })
	menu.editNew.Connect("activate", func() { editNew(menu, fts, jl, ftv) })
	menu.editEdit.Connect("activate", func() { editEdit(menu, fts, jl, ftv) })
	menu.editDelete.Connect("activate", func() { editDelete(menu, fts, jl, ftv) })
	menu.editCopy.Connect("activate", func() { editCopy(fts, clp) })
	menu.editPaste.Connect("activate", func() { editPaste(menu, fts, jl, ftv, clp) })
}

func MenuEditCurrent(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList) {
	cursor := fts.Current()
	if len(cursor.Path) == 0 {
		return
	}
	prop := fts.Property(cursor)
	obj := fts.Object(cursor)
	menu.editNew.SetSensitive(prop.MayAddObject())
	menu.editDelete.SetSensitive(prop.MayRemove())
	menu.editEdit.SetSensitive(prop.MayEdit())
	menu.editUndo.SetSensitive(jl.CanUndo())
	menu.editRedo.SetSensitive(jl.CanRedo())
	for _, v := range global.graphview {
		v.Sync()
	}
	global.xmlview.Set(obj)
}

func editUndo(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	//log.Println("editUndo")
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
	//log.Println("editRedo")
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
	//log.Println("editNew")
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
	//log.Println("editNew finished")
}

func editEdit(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	//log.Println("editEdit")
	defer MenuEditCurrent(menu, fts, jl)
	dialog, err := EditDialogNew(fts)
	if err != nil {
		log.Println("editEdit error: ", err)
		return
	}
	job, ok := dialog.Run(fts)
	if ok {
		state, ok := jl.Apply(job)
		if ok {
			path, err := gtk.TreePathNewFromString(state.(string))
			if err != nil {
				log.Println("editEdit error: TreePathNewFromString failed:", err)
				return
			}
			ftv.TreeView().ExpandToPath(path)
			ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
		}
	}
	//log.Println("editEdit finished")
}

func editCopy(fts *models.FilesTreeStore, clp *gtk.Clipboard) {
	log.Println("editCopy")
	var buf []byte
	obj, err := fts.GetObjectById(fts.GetCurrentId())
	if err != nil {
		log.Printf("editCopy error: %s\n", err)
		return
	}
	buf, err = freesp.CreateXML(obj)
	if err != nil {
		log.Printf("editCopy error: %s\n", err)
		return
	}
	clp.SetText(string(buf))
}

func editPaste(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView, clp *gtk.Clipboard) {
	defer MenuEditCurrent(menu, fts, jl)
	if !clp.WaitIsTextAvailable() {
		log.Println("editPaste: Nothing on clipboard")
		return
	}
	text, err := clp.WaitForText()
	if err != nil {
		log.Println("editPaste: Failed to get data from clipboard")
		return
	}
	log.Printf("editPaste: receive from clipboard: %s\n", text)
	job, err := ParseText(text, fts)
	if err != nil {
		log.Printf("editPaste: error parsing data from clipboard: %s\n", err)
		return
	}
	state, ok := jl.Apply(job)
	if ok {
		path, err := gtk.TreePathNewFromString(state.(string))
		if err != nil {
			log.Println("editPaste error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().ExpandToPath(path)
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}

func editDelete(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	//log.Println("editDelete")
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
