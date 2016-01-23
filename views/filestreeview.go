package views

import (
	"fmt"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
)

type FilesTreeView struct {
	ScrolledView
	view *gtk.TreeView
}

var _ tr.TreeViewIf = (*FilesTreeView)(nil)

func FilesTreeViewNew(model *models.FilesTreeStore, width, height int) (viewer *FilesTreeView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &FilesTreeView{*v, nil}
	err = viewer.init(model)
	return
}

func (v *FilesTreeView) TreeView() *gtk.TreeView {
	return v.view
}

// Initialize view to reflect the model:
func (v *FilesTreeView) init(model *models.FilesTreeStore) error {
	renderer1, err := gtk.CellRendererTextNew()
	if err != nil {
		return fmt.Errorf("Error CellRendererTextNew:", err)
	}
	renderer2, err := gtk.CellRendererPixbufNew()
	if err != nil {
		return fmt.Errorf("Error CellRendererPixbufNew:", err)
	}
	col1, err := gtk.TreeViewColumnNewWithAttribute("Type", renderer2, "pixbuf", 0)
	if err != nil {
		return fmt.Errorf("Error col1:", err)
	}
	col2, err := gtk.TreeViewColumnNewWithAttribute("Name", renderer1, "text", 1)
	if err != nil {
		return fmt.Errorf("Error col2:", err)
	}
	v.view, err = gtk.TreeViewNewWithModel(model.TreeStore())
	if err != nil {
		return fmt.Errorf("Error TreeViewNew", err)
	}
	v.view.AppendColumn(col1)
	v.view.AppendColumn(col2)
	v.scrolled.Add(v.view)
	return nil
}

func (v *FilesTreeView) SelectId(newId string) (err error) {
	path, err := gtk.TreePathNewFromString(newId)
	if err != nil {
		return
	}
	v.TreeView().SetCursor(path, v.TreeView().GetExpanderColumn(), false)
	return
}
