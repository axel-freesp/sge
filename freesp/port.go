package freesp

import (
	"fmt"
	"log"
	"unsafe"
)

type port struct {
	itsType     PortType
	connections portList
	node        Node
}

// TODO: Create namedPortType object first and hand it here?
func newPort(pt PortType, n Node) *port {
	return &port{pt, portListInit(), n}
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
 *      Connection(c *port) Connection
 *  	AddConnection(Port) error
 *  	RemoveConnection(c Port)
 *  }
 *
 *  func PortConnect(port1, port2 Port) error
 */
var _ Port = (*port)(nil)

func (p *port) Name() string {
	return p.itsType.Name()
}

func (p *port) SignalType() SignalType {
	return p.itsType.SignalType()
}

func (p *port) Direction() PortDirection {
	return p.itsType.Direction()
}

func (p *port) Connections() []Port {
	return p.connections.Ports()
}

func (p *port) Node() Node {
	return p.node
}

func (p *port) AddConnection(c Port) error {
	return PortConnect(p, c)
}

func (p *port) RemoveConnection(c Port) {
	p.connections.Remove(c)
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

func PortConnect(port1, port2 Port) error {
	p1, p2 := port1.(*port), port2.(*port)
	if port1.SignalType().TypeName() != port2.SignalType().TypeName() {
		return fmt.Errorf("type mismatch")
	}
	if port1.Direction() == port2.Direction() {
		return fmt.Errorf("direction mismatch")
	}
	p1.connections.Append(port2)
	p2.connections.Append(port1)
	return nil
}

/*
 *  fmt.Stringer API
 */
var _ fmt.Stringer = (*port)(nil)

func (p *port) String() (s string) {
	s = fmt.Sprintf("%sPort of node %s(%s, %d connections)",
		p.Direction(), p.Node().Name(), p.Name(), len(p.Connections()))
	return
}

var _ fmt.Stringer = PortDirection(false)

func (d PortDirection) String() (s string) {
	if d == InPort {
		s = "Input"
	} else {
		s = "Output"
	}
	return
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*port)(nil)

func (p *port) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	parentId := tree.Parent(cursor)
	if tree.Property(parentId).IsReadOnly() {
		prop = 0
	} else {
		prop = mayAddObject
	}
	var kind Symbol
	if p.Direction() == InPort {
		kind = SymbolInputPort
	} else {
		kind = SymbolOutputPort
	}
	err := tree.AddEntry(cursor, kind, p.Name(), p, prop)
	if err != nil {
		log.Fatalf("port.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	t := p.SignalType()
	t.AddToTree(tree, child)
	for _, c := range p.Connections() {
		child = tree.Append(cursor)
		p.Connection(c.(*port)).AddToTree(tree, child)
	}
	return
}

func (p *port) treeAddNewObject(tree Tree, cursor Cursor, conn Connection, otherPort Port) (newCursor Cursor) {
	newCursor = tree.Insert(cursor)
	conn.AddToTree(tree, newCursor)
	contextCursor := tree.Parent(tree.Parent(cursor))
	cCursor := tree.CursorAt(contextCursor, otherPort)
	cChild := tree.Append(cCursor)
	conn.AddToTree(tree, cChild)
	return
}

func (p *port) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case Connection:
		conn := obj.(Connection)
		var thisPort, otherPort Port
		if p.Direction() == InPort {
			otherPort = conn.From
			thisPort = conn.To
		} else {
			otherPort = conn.To
			thisPort = conn.From
		}
		if p != thisPort {
			log.Println("p =", p)
			log.Println("thisPort =", thisPort)
			log.Fatal("port.AddNewObject error: invalid connection ", conn)
		}
		sgt := thisPort.Node().Context()
		if sgt != otherPort.Node().Context() {
			log.Fatal("port.AddNewObject error: nodes have different context: ", conn)
		}
		contextCursor := tree.Parent(tree.Parent(cursor))

		thisPort.AddConnection(otherPort)
		otherPort.AddConnection(thisPort)
		newCursor = p.treeAddNewObject(tree, cursor, conn, otherPort)

		context := tree.Object(contextCursor)
		switch context.(type) {
		case SignalGraph:
		case Implementation:
			// propagate new edge to all instances of embracing type
			nt := tree.Object(tree.Parent(contextCursor))
			for _, nn := range nt.(NodeType).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, context)
				pCursor := tree.CursorAt(tCursor, p)
				pCursor.Position = cursor.Position

				p.treeAddNewObject(tree, pCursor, conn, otherPort)
			}

		default:
			log.Fatalf("port.AddNewObject error: wrong context type %T: %v\n", context, context)
		}

	default:
		log.Fatalf("port.AddNewObject error: invalid type %T: %v\n", obj, obj)
	}
	return
}

func (p *port) treeRemoveObject(tree Tree, cursor Cursor, conn Connection, otherPort Port) (removed []IdWithObject) {
	contextCursor := tree.Parent(tree.Parent(tree.Parent(cursor)))
	pCursor := tree.CursorAt(contextCursor, otherPort)
	otherCursor := tree.CursorAt(pCursor, conn)
	prefix, index := tree.Remove(cursor)
	removed = append(removed, IdWithObject{prefix, index, conn})
	tree.Remove(otherCursor)
	return
}

func (p *port) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("port.RemoveObject error: not removing child of mine.")
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
				log.Fatal("port.AddNewObject error: invalid connection ", conn)
			}
		} else {
			otherPort = conn.To
			thisPort = conn.From
			if p != thisPort {
				log.Fatal("port.AddNewObject error: invalid connection ", conn)
			}
		}
		contextCursor := tree.Parent(tree.Parent(tree.Parent(cursor)))
		removed = p.treeRemoveObject(tree, cursor, conn, otherPort)
		context := tree.Object(contextCursor)
		switch context.(type) {
		case SignalGraphType:
		case Implementation:
			// propagate removed edge to all instances of embracing type
			nt := tree.Object(tree.Parent(contextCursor))
			for _, nn := range nt.(NodeType).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, context)
				pCursor := tree.CursorAt(tCursor, p)
				cCursor := tree.CursorAt(pCursor, conn)
				p.treeRemoveObject(tree, cCursor, conn, otherPort)
			}

		default:
			log.Fatalf("port.AddNewObject error: wrong context type %t: %v\n", contextCursor, contextCursor)
		}
		p.RemoveConnection(otherPort)
		otherPort.RemoveConnection(p)

	default:
		log.Fatalf("Port.RemoveObject error: invalid type %T: %v\n", obj, obj)
	}
	return
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

/*
 *      portList
 *
 */

type portList struct {
	ports []Port
}

func portListInit() portList {
	return portList{nil}
}

func (l *portList) Append(p Port) {
	l.ports = append(l.ports, p)
}

func (l *portList) Remove(p Port) {
	var i int
	for i = range l.ports {
		if p == l.ports[i] {
			break
		}
	}
	if i >= len(l.ports) {
		for _, v := range l.ports {
			log.Printf("portList.RemovePort have Port %v\n", v)
		}
		log.Fatalf("portList.RemovePort error: Port %v not in this list\n", p)
	}
	for i++; i < len(l.ports); i++ {
		l.ports[i-1] = l.ports[i]
	}
	l.ports = l.ports[:len(l.ports)-1]
}

func (l *portList) Ports() []Port {
	return l.ports
}

func (l *portList) Find(name string) (p Port, ok bool) {
	ok = false
	for _, p = range l.ports {
		if p.Name() == name {
			ok = true
			return
		}
	}
	return
}
