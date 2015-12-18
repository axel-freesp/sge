package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type SignalType struct {
	freesp.SignalType
}

func (t SignalType) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	err := tree.AddEntry(cursor, imageSignalType, t.TypeName(), t)
	if err != nil {
		log.Fatal("SignalType.AddToTree error: AddEntry failed: %s", err)
	}
}

var (
	imageSignalType *gdk.Pixbuf = nil
)

func init_signaltype(iconPath string) (err error) {
	imageSignalType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/signal-type.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_signaltype error loading signal-type.png: %s", err)
	}
	return
}
