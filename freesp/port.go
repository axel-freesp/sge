package freesp

import (
	"fmt"
	"log"
	"unsafe"
)

type port struct {
	name        string
	itsType     PortType
	direction   PortDirection
	connections []Port
	node        Node
}

func newPort(name string, pt *portType, dir PortDirection, n *node) *port {
	return &port{name, pt, dir, nil, n}
}

/*
 *  freesp.Port API
 *
 *  type Port interface {
 *  	PortName() string
 *  	ItsType() PortType
 *  	Direction() PortDirection
 *  	Connections() []Port
 *  	Node() Node
 *  	AddConnection(Port) error
 *  	RemoveConnection(c Port)
 *  }
 *
 *  func PortConnect(port1, port2 Port) error
 */
var _ Port = (*port)(nil)

type portError struct {
	reason string
}

func (p *port) PortName() string {
	return p.name
}

func (p *port) ItsType() PortType {
	return p.itsType
}

func (p *port) Direction() PortDirection {
	return p.direction
}

func (p *port) Connections() []Port {
	return p.connections
}

func (p *port) Node() Node {
	return p.node
}

func (p *port) AddConnection(c Port) error {
	return PortConnect(p, c)
}

func (p *port) RemoveConnection(c Port) {
	i, ok := p.indexOfConnection(c)
	if ok {
		for j := i + 1; j < len(p.connections); j++ {
			p.connections[j-1] = p.connections[j]
		}
		p.connections = p.connections[:len(p.connections)-1]
	} else {
		fmt.Printf("port(%s).RemoveConnection error: could not find port %s\n", p, c)
	}
}

func PortConnect(port1, port2 Port) error {
	p1, p2 := port1.(*port), port2.(*port)
	if p1.itsType != p2.itsType {
		return fmt.Errorf("type mismatch")
	}
	if p1.direction == p2.direction {
		return fmt.Errorf("direction mismatch")
	}
	if !findPort(p1.connections, p2) {
		p1.connections = append(p1.connections, port2)
	}
	if !findPort(p2.connections, p1) {
		p2.connections = append(p2.connections, port1)
	}
	return nil
}

/*
 *  fmt.Stringer API
 */
var _ fmt.Stringer = (*port)(nil)

func (p *port) String() (s string) {
	s = fmt.Sprintf("%sPort of node %s(%s, %d connections)",
		p.direction, p.Node().NodeName(), p.name, len(p.connections))
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*port)(nil)

func (p *port) AddToTree(tree Tree, cursor Cursor) {
	var kind Symbol
	if p.Direction() == InPort {
		kind = SymbolInputPort
	} else {
		kind = SymbolOutputPort
	}
	err := tree.AddEntry(cursor, kind, p.PortName(), p)
	if err != nil {
		log.Fatal("Port.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	t := p.ItsType().SignalType()
	t.AddToTree(tree, child)
	for _, c := range p.Connections() {
		child = tree.Append(cursor)
		p.Connection(c.(*port)).AddToTree(tree, child)
	}
	return
}

func (p *port) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Connection:
		conn := obj.(Connection)
		log.Println("Port.AddNewObject: conn =", conn)
		var thisPort, otherPort Port
		if p.Direction() == InPort {
			otherPort = conn.From
			thisPort = conn.To
		} else {
			otherPort = conn.To
			thisPort = conn.From
		}
		if p != thisPort {
			log.Println("p.Port =", p)
			log.Println("thisPort =", thisPort)
			log.Fatal("Port.AddNewObject error: invalid connection ", conn)
		}
		thisPort.AddConnection(otherPort)
		otherPort.AddConnection(thisPort)
		newCursor = tree.Insert(cursor)
		conn.AddToTree(tree, newCursor)
		cCursor := tree.Cursor(otherPort) // TODO: faster search!
		cChild := tree.Append(cCursor)
		conn.AddToTree(tree, cChild)

	default:
		log.Fatalf("Port.AddNewObject error: invalid type %T: %v", obj, obj)
	}
	return
}

func (p *port) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case Connection:
		conn := obj.(Connection)
		var thisPort, otherPort Port
		if p.Direction() == InPort {
			otherPort = conn.From
			thisPort = conn.To
			if p != thisPort {
				log.Fatal("Port.AddNewObject error: invalid connection ", conn)
			}
		} else {
			otherPort = conn.To
			thisPort = conn.From
			if p != thisPort {
				log.Fatal("Port.AddNewObject error: invalid connection ", conn)
			}
		}
		pCursor := tree.Cursor(otherPort)
		otherCursor := tree.CursorAt(pCursor, conn)
		thisPort.RemoveConnection(otherPort)
		otherPort.RemoveConnection(thisPort)
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, conn})
		prefix, index = tree.Remove(otherCursor)

	default:
		log.Fatalf("Port.RemoveObject error: invalid type %T: %v", obj, obj)
	}
	return
}

func (p *port) Connection(c *port) Connection {
	var from, to Port
	if p.Direction() == InPort {
		from = c
		to = p
	} else {
		from = p
		to = c
	}
	return Connection{from, to}
}

/*
 *  freesp private functions
 */
func findPort(list []Port, prt *port) bool {
	for _, p := range list {
		if unsafe.Pointer(p.(*port)) == unsafe.Pointer(prt) {
			return true
		}
	}
	return false
}

func (p *port) indexOfConnection(c Port) (index int, ok bool) {
	for index = 0; index < len(p.connections); index++ {
		if c == p.connections[index] {
			break
		}
	}
	ok = (index < len(p.connections))
	return
}
