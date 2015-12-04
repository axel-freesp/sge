package models

import (
	//"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type NodeTypeTreeStore struct {
	treestore *gtk.TreeStore
	obj       freesp.NodeType
}

func NodeTypeTreeStoreNew(obj freesp.NodeType) (ret *NodeTypeTreeStore, err error) {
	ts, err := gtk.TreeStoreNew(glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		err = gtkErr("NodeTypeTreeStoreNew", "TreeStoreNew", err)
		ret = nil
		return
	}
	ret = &NodeTypeTreeStore{ts, obj}
	iter := ts.Append(nil)
	err = ts.SetValue(iter, 0, obj.TypeName())
	if err != nil {
		err = gtkErr("NodeTypeTreeStoreNew", "ts.SetValue", err)
		ret = nil
		return
	}
	setPortList(ts, iter, obj.InPorts(), "Input ports")
	setPortList(ts, iter, obj.OutPorts(), "Output ports")
	return
}

func setPortList(ts *gtk.TreeStore, iter *gtk.TreeIter, list []freesp.NamedPortType, listname string) error {
	iter = ts.Append(iter)
	err := ts.SetValue(iter, 0, listname)
	if err != nil {
		return gtkErr("setPortList", "ts.SetValue", err)
	}
	for _, p := range list {
		child := ts.Append(iter)
		err = ts.SetValue(child, 1, p.Name())
		if err != nil {
			return gtkErr("setPortList", "ts.SetValue", err)
		}
		err = ts.SetValue(child, 0, p.TypeName())
		if err != nil {
			return gtkErr("setPortList", "ts.SetValue", err)
		}
	}
	return nil
}

func (s *NodeTypeTreeStore) TreeStore() *gtk.TreeStore {
	return s.treestore
}
