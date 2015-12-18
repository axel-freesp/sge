package freesp

import (
	"fmt"
	"log"
)

/*
 *  TreeElement API
 */

var _ TreeElement = Connection{}

func (c Connection) AddToTree(tree Tree, cursor Cursor) {
	text := fmt.Sprintf("%s/%s -> %s/%s", c.From.Node().NodeName(), c.From.PortName(),
		c.To.Node().NodeName(), c.To.PortName())
	err := tree.AddEntry(cursor, SymbolConnection, text, c)
	if err != nil {
		log.Fatal("Connection.AddToTree error: AddEntry failed: %s", err)
	}
}

func (c Connection) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Fatal("Connection.AddNewObject - nothing to add.")
	return
}

func (c Connection) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("Connection.AddNewObject - nothing to remove.")
	return
}
