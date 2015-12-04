package views

import (
	"fmt"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
)

type NodeTypeListView struct {
	ScrolledView
	view *gtk.TreeView
}

func NodeTypeListViewNew(model *models.NodeTypeTreeStore, width, height int) (viewer *NodeTypeListView, err error) {
	v, err := ScrolledViewNew(width, height)
	if err != nil {
		viewer = nil
		return
	}
	viewer = &NodeTypeListView{*v, nil}
	err = viewer.init(model)
	return
}

func (v *NodeTypeListView) TreeView() *gtk.TreeView {
	return v.view
}

// Initialize view to render text:
func (v *NodeTypeListView) init(model *models.NodeTypeTreeStore) error {
	renderer, err := gtk.CellRendererTextNew()
	if err != nil {
		return fmt.Errorf("Error CellRendererTextNew:", err)
	}
	col1, err := gtk.TreeViewColumnNewWithAttribute("Type", renderer, "text", 0)
	if err != nil {
		return fmt.Errorf("Error col1:", err)
	}
	col2, err := gtk.TreeViewColumnNewWithAttribute("Name", renderer, "text", 1)
	if err != nil {
		return fmt.Errorf("Error col2:", err)
	}
	v.view, err = gtk.TreeViewNewWithModel(model.TreeStore())
	if err != nil {
		return fmt.Errorf("Error TreeViewNew", err)
	}
	v.view.AppendColumn(col1)
	v.view.AppendColumn(col2)
	//v.view.ExpandAll()
	v.scrolled.Add(v.view)
	return nil
}
