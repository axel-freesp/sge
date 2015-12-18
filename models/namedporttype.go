package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type NamedPortType struct {
	freesp.NamedPortType
}

var _ TreeElement = NamedPortType{}

func (p NamedPortType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	var kind *gdk.Pixbuf
	if p.Direction() == freesp.InPort {
		kind = imageInputPortType
	} else {
		kind = imageOutputPortType
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p.NamedPortType)
	if err != nil {
		log.Fatal("NamedPortType.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	SignalType{p.SignalType()}.AddToTree(tree, child)
}

func (p NamedPortType) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Fatal("NamedPortType.AddNewObject - nothing to add.")
	return
}

func (p NamedPortType) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("NamedPortType.AddNewObject - nothing to remove.")
	return
}

var (
	imageInputPortType  *gdk.Pixbuf = nil
	imageOutputPortType *gdk.Pixbuf = nil
)

func init_namedporttype(iconPath string) (err error) {
	imageInputPortType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/inport-green.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_port error loading inport-green.png: %s", err)
		return
	}
	imageOutputPortType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/outport-red.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_port error loading outport-red.png: %s", err)
	}
	return
}
