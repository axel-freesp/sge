package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
)

type connection struct {
	from, to bh.PortIf
}

var _ bh.ConnectionIf = (*connection)(nil)

func ConnectionNew(from, to bh.PortIf) *connection {
	return &connection{from, to}
}

func (c *connection) From() bh.PortIf {
	return c.from
}

func (c *connection) To() bh.PortIf {
	return c.to
}

func (c *connection) CreateXml() (buf []byte, err error) {
	xmlconn := CreateXmlConnection(c)
	buf, err = xmlconn.Write()
	return
}

//
//  fmt.Stringer API
//

var _ fmt.Stringer = (*connection)(nil)

func (c *connection) String() (s string) {
	s = fmt.Sprintf("Connection(%s/%s -> %s/%s)",
		c.from.Node().Name(), c.from.Name(),
		c.to.Node().Name(), c.to.Name())
	return
}

//
//  tr.TreeElementIf API
//

var _ tr.TreeElementIf = (*connection)(nil)

func (c *connection) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	text := fmt.Sprintf("%s/%s -> %s/%s", c.from.Node().Name(), c.from.Name(),
		c.to.Node().Name(), c.to.Name())
	prop := freesp.PropertyNew(false, false, true)
	err := tree.AddEntry(cursor, tr.SymbolConnection, text, c, prop)
	if err != nil {
		log.Fatalf("connection.AddToTree error: AddEntry failed: %s\n", err)
	}
}

func (c *connection) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	log.Fatal("Connection.AddNewObject - nothing to add.")
	return
}

func (c *connection) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	log.Fatal("Connection.AddNewObject - nothing to remove.")
	return
}

func (c *connection) Identify(te tr.TreeElementIf) bool {
	switch te.(type) {
	case *connection:
		return te.(*connection) == c
	}
	return false
}
