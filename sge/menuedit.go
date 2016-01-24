package main

import (
	//"github.com/axel-freesp/sge/freesp"
	gr "github.com/axel-freesp/sge/interface/graph"
	tr "github.com/axel-freesp/sge/interface/tree"
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

func MenuEditPost(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList) {
	MenuEditCurrent(menu, fts, jl)
	cursor := fts.Current()
	var obj tr.TreeElement
	if len(cursor.Path) != 0 {
		obj = fts.Object(cursor)
	}
	global.GVC().XmlTextView().Set(obj)
	global.GVC().Sync()
}

func MenuEditCurrent(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList) {
	menu.editUndo.SetSensitive(jl.CanUndo())
	menu.editRedo.SetSensitive(jl.CanRedo())
	var prop tr.Property
	cursor := fts.Current()
	if len(cursor.Path) != 0 {
		prop = fts.Property(cursor)
		menu.editNew.SetSensitive(prop.MayAddObject())
		menu.editDelete.SetSensitive(prop.MayRemove())
		menu.editEdit.SetSensitive(prop.MayEdit())
	} else {
		menu.editNew.SetSensitive(false)
		menu.editDelete.SetSensitive(false)
		menu.editEdit.SetSensitive(false)
	}
}

func editUndo(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	defer MenuEditPost(menu, fts, jl)
	state, ok := jl.Undo()
	if ok {
		global.win.graphViews.Sync()
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
	defer MenuEditPost(menu, fts, jl)
	state, ok := jl.Redo()
	if ok {
		global.win.graphViews.Sync()
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
	defer MenuEditPost(menu, fts, jl)
	dialog, err := NewElementDialogNew(fts)
	if err != nil {
		log.Println("editNew error: ", err)
		return
	}
	job, ok := dialog.Run(fts)
	if ok {
		state, ok := jl.Apply(job)
		if ok {
			global.win.graphViews.Sync()
			path, err := gtk.TreePathNewFromString(state.(string))
			if err != nil {
				log.Println("editNew error: TreePathNewFromString failed:", err)
				return
			}
			ftv.TreeView().ExpandToPath(path)
			ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
		}
	}
}

func editEdit(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView) {
	defer MenuEditPost(menu, fts, jl)
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
}

func editCopy(fts *models.FilesTreeStore, clp *gtk.Clipboard) {
	var buf []byte
	obj, err := fts.GetObjectById(fts.GetCurrentId())
	if err != nil {
		log.Printf("editCopy error: %s\n", err)
		return
	}
	buf, err = obj.(gr.XmlCreator).CreateXml()
	if err != nil {
		log.Printf("editCopy error: %s\n", err)
		return
	}
	clp.SetText(string(buf))
}

func editPaste(menu *GoAppMenu, fts *models.FilesTreeStore, jl IJobList, ftv *views.FilesTreeView, clp *gtk.Clipboard) {
	defer MenuEditPost(menu, fts, jl)
	if !clp.WaitIsTextAvailable() {
		log.Println("editPaste: Nothing on clipboard")
		return
	}
	text, err := clp.WaitForText()
	if err != nil {
		log.Println("editPaste: Failed to get data from clipboard")
		return
	}
	job, err := ParseText(text, fts)
	if err != nil {
		log.Printf("editPaste: error parsing data from clipboard: %s\n", err)
		return
	}
	state, ok := jl.Apply(job)
	if ok {
		global.win.graphViews.Sync()
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
	defer MenuEditPost(menu, fts, jl)
	job := DeleteObjectJobNew(fts.GetCurrentId())
	state, ok := jl.Apply(EditorJobNew(JobDeleteObject, job))
	if ok {
		global.win.graphViews.Sync()
		path, err := gtk.TreePathNewFromString(state.(string))
		if err != nil {
			log.Println("editNew error: TreePathNewFromString failed:", err)
			return
		}
		ftv.TreeView().ExpandToPath(path)
		ftv.TreeView().SetCursor(path, ftv.TreeView().GetExpanderColumn(), false)
	}
}
