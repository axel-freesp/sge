package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type Connection struct {
	freesp.Connection
}

var _ TreeElement = Connection{}

func (c Connection) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	text := fmt.Sprintf("%s/%s -> %s/%s", c.From.Node().NodeName(), c.From.PortName(),
		c.To.Node().NodeName(), c.To.PortName())
	err := tree.AddEntry(cursor, imageConnected, text, c.Connection)
	if err != nil {
		log.Fatal("Connection.AddToTree error: AddEntry failed: %s", err)
	}
}

func (c Connection) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	log.Fatal("Connection.AddNewObject - nothing to add.")
	return
}

func (c Connection) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("Connection.AddNewObject - nothing to remove.")
	return
}

var (
	imageConnected *gdk.Pixbuf = nil
)

func init_connection(iconPath string) (err error) {
	imageConnected, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/link.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_connection error loading link.png: %s", err)
	}
	return
}
