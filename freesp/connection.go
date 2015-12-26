package freesp

import (
	"fmt"
	"log"
)

/*
 *  fmt.Stringer API
 */

var _ fmt.Stringer = Connection{}

func (c Connection) String() (s string) {
	s = fmt.Sprintf("Connection(%s/%s -> %s/%s)",
		c.From.Node().Name(), c.From.PortName(),
		c.To.Node().Name(), c.To.PortName())
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = Connection{}

func (c Connection) AddToTree(tree Tree, cursor Cursor) {
	text := fmt.Sprintf("%s/%s -> %s/%s", c.From.Node().Name(), c.From.PortName(),
		c.To.Node().Name(), c.To.PortName())
	err := tree.AddEntry(cursor, SymbolConnection, text, c, mayRemove)
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
