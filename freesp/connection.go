package freesp

import (
	"fmt"
	"log"
)

type connection struct {
	from, to Port
}

func ConnectionNew(from, to Port) *connection {
	return &connection{from, to}
}

func (c *connection) From() Port {
	return c.from
}

func (c *connection) To() Port {
	return c.to
}

/*
 *  fmt.Stringer API
 */

var _ fmt.Stringer = (*connection)(nil)

func (c *connection) String() (s string) {
	s = fmt.Sprintf("Connection(%s/%s -> %s/%s)",
		c.from.Node().Name(), c.from.Name(),
		c.to.Node().Name(), c.to.Name())
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*connection)(nil)

func (c *connection) AddToTree(tree Tree, cursor Cursor) {
	text := fmt.Sprintf("%s/%s -> %s/%s", c.from.Node().Name(), c.from.Name(),
		c.to.Node().Name(), c.to.Name())
	err := tree.AddEntry(cursor, SymbolConnection, text, c, mayRemove)
	if err != nil {
		log.Fatalf("connection.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (c *connection) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	log.Fatal("Connection.AddNewObject - nothing to add.")
	return
}

func (c *connection) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("Connection.AddNewObject - nothing to remove.")
	return
}
