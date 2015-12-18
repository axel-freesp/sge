package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type SignalGraph struct {
	freesp.SignalGraph
}

var _ TreeElement = SignalGraph{}

func (g SignalGraph) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, imageSignalGraph, g.Filename(), g.ItsType())
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
	SignalGraphType{g.ItsType()}.AddToTree(tree, cursor)
}

func (g SignalGraph) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Fatal("SignalGraph.AddNewObject - nothing to add.")
	return
}

func (g SignalGraph) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("SignalGraph.AddNewObject - nothing to remove.")
	return
}

var (
	imageSignalGraph *gdk.Pixbuf = nil
)

func init_signalgraph(iconPath string) (err error) {
	imageSignalGraph, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test1.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_signaltype error loading test1.png: %s", err)
	}
	return
}
