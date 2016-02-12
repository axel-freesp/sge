package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	tr "github.com/axel-freesp/sge/interface/tree"
	//"image"
	"log"
	"unsafe"
)

type port struct {
	gr.PathModePositionerObject
	itsType   bh.PortTypeIf
	connected portList
	conn      []bh.ConnectionIf
	node      bh.NodeIf
}

var _ bh.PortIf = (*port)(nil)

// TODO: Create namedPortType object first and hand it here?
func newPort(pt bh.PortTypeIf, n bh.NodeIf) *port {
	return &port{*gr.PathModePositionerObjectNew(), pt, portListInit(), nil, n}
}

/*
 *  freesp.bh.PortIf API
 *
 *  type bh.PortIf interface {
 *  	Name() string
 *  	ItsType() bh.PortTypeIf
 *  	Direction() PortDirection
 *  	Connections() []bh.PortIf
 *  	bh.NodeIf() bh.NodeIf
 *      bh.ConnectionIf(c *port) bh.ConnectionIf
 *  	AddConnection(bh.PortIf) error
 *  	RemoveConnection(c bh.PortIf)
 *  }
 *
 *  func PortConnect(port1, port2 bh.PortIf) error
 */
var _ bh.PortIf = (*port)(nil)

func (p *port) Name() string {
	return p.itsType.Name()
}

func (p *port) SetName(newName string) {
	log.Panicf("port.SetName: not allowed.\n")
}

func (p *port) SignalType() bh.SignalTypeIf {
	return p.itsType.SignalType()
}

func (p *port) Direction() gr.PortDirection {
	return p.itsType.Direction()
}

func (p *port) SetDirection(gr.PortDirection) {
	log.Panicf("port.SetDirection() is not allowed!\n")
}

func (p *port) Connections() []bh.PortIf {
	return p.connected.Ports()
}

func (p *port) Node() bh.NodeIf {
	return p.node
}

func (p *port) AddConnection(c bh.ConnectionIf) error {
	return portConnect(p, c.(*connection))
}

func (p *port) RemoveConnection(c bh.PortIf) {
	_, ok, index := p.connected.Find(c.Node().Name(), c.Name())
	if !ok {
		log.Fatalf("port.bh.ConnectionIf error: port %v not in connected list of port %v\n", c, p)
	}
	p.connected.Remove(c)
	for index++; index < len(p.conn); index++ {
		p.conn[index-1] = p.conn[index]
	}
	p.conn = p.conn[:len(p.conn)-1]
}

func (p *port) Connection(c bh.PortIf) bh.ConnectionIf {
	_, ok, index := p.connected.Find(c.Node().Name(), c.Name())
	if !ok {
		log.Fatalf("port.bh.ConnectionIf error: port %v not in connected list of port %v\n", c, p)
	}
	return p.conn[index]
}

func portConnect(port1 bh.PortIf, c *connection) error {
	p1 := port1.(*port)
	var port2 bh.PortIf
	var p2 *port
	if c.from.(*port) == p1 {
		port2 = c.to
	} else if c.to.(*port) == p1 {
		port2 = c.from
	}
	p2 = port2.(*port)
	if port1.SignalType().TypeName() != port2.SignalType().TypeName() {
		return fmt.Errorf("type mismatch")
	}
	if port1.Direction() == port2.Direction() {
		return fmt.Errorf("direction mismatch")
	}
	p1.connected.Append(port2.(*port))
	p1.conn = append(p1.conn, c)
	p2.connected.Append(port1.(*port))
	p2.conn = append(p2.conn, c)
	return nil
}

func (p *port) CreateXml() (buf []byte, err error) {
	if p.Direction() == gr.OutPort {
		xmlport := CreateXmlOutPort(p)
		buf, err = xmlport.Write()
	} else {
		xmlport := CreateXmlInPort(p)
		buf, err = xmlport.Write()
	}
	return
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

/*
 *  tr.TreeElementIf API
 */

var _ tr.TreeElementIf = (*port)(nil)

func (p *port) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var prop tr.Property
	parentId := tree.Parent(cursor)
	if tree.Property(parentId).IsReadOnly() {
		prop = freesp.PropertyNew(false, false, false)
	} else {
		prop = freesp.PropertyNew(true, false, false)
	}
	var kind tr.Symbol
	if p.Direction() == gr.InPort {
		kind = tr.SymbolInputPort
	} else {
		kind = tr.SymbolOutputPort
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
		p.Connection(c).AddToTree(tree, child)
	}
	return
}

func (p *port) treeAddNewObject(tree tr.TreeIf, cursor tr.Cursor, conn bh.ConnectionIf, otherPort bh.PortIf) (newCursor tr.Cursor) {
	newCursor = tree.Insert(cursor)
	conn.AddToTree(tree, newCursor)
	contextCursor := tree.Parent(tree.Parent(cursor))
	cCursor := tree.CursorAt(contextCursor, otherPort)
	cChild := tree.Append(cCursor)
	conn.AddToTree(tree, cChild)
	return
}

func (p *port) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	switch obj.(type) {
	case bh.ConnectionIf:
		conn := obj.(bh.ConnectionIf)
		var thisPort, otherPort bh.PortIf
		if p.Direction() == gr.InPort {
			otherPort = conn.From()
			thisPort = conn.To()
		} else {
			otherPort = conn.To()
			thisPort = conn.From()
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

		thisPort.AddConnection(conn)
		newCursor = p.treeAddNewObject(tree, cursor, conn, otherPort)

		context := tree.Object(contextCursor)
		switch context.(type) {
		case bh.SignalGraphIf:
		case bh.ImplementationIf:
			// propagate new edge to all instances of embracing type
			nt := tree.Object(tree.Parent(contextCursor))
			for _, nn := range nt.(bh.NodeTypeIf).Instances() {
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

func (p *port) treeRemoveObject(tree tr.TreeIf, cursor tr.Cursor, conn bh.ConnectionIf, otherPort bh.PortIf) (removed []tr.IdWithObject) {
	contextCursor := tree.Parent(tree.Parent(tree.Parent(cursor)))
	pCursor := tree.CursorAt(contextCursor, otherPort)
	otherCursor := tree.CursorAt(pCursor, conn)
	prefix, index := tree.Remove(cursor)
	removed = append(removed, tr.IdWithObject{prefix, index, conn})
	tree.Remove(otherCursor)
	return
}

func (p *port) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if p != tree.Object(parent) {
		log.Fatal("port.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.ConnectionIf:
		conn := obj.(bh.ConnectionIf)
		var thisPort, otherPort bh.PortIf
		if p.Direction() == gr.InPort {
			otherPort = conn.From()
			thisPort = conn.To()
			if p != thisPort {
				log.Fatal("port.RemoveObject error: invalid connection ", conn)
			}
		} else {
			otherPort = conn.To()
			thisPort = conn.From()
			if p != thisPort {
				log.Fatal("port.RemoveObject error: invalid connection ", conn)
			}
		}
		contextCursor := tree.Parent(tree.Parent(tree.Parent(cursor)))
		removed = p.treeRemoveObject(tree, cursor, conn, otherPort)
		context := tree.Object(contextCursor)
		switch context.(type) {
		case bh.SignalGraphIf:
		case bh.ImplementationIf:
			// propagate removed edge to all instances of embracing type
			nt := tree.Object(tree.Parent(contextCursor))
			for _, nn := range nt.(bh.NodeTypeIf).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, context)
				pCursor := tree.CursorAt(tCursor, p)
				cCursor := tree.CursorAt(pCursor, conn)
				p.treeRemoveObject(tree, cCursor, conn, otherPort)
			}

		default:
			log.Fatalf("port.RemoveObject error: wrong context type %T: %v\n", context, context)
		}
		p.RemoveConnection(otherPort)
		otherPort.RemoveConnection(p)

	default:
		log.Fatalf("bh.PortIf.RemoveObject error: invalid type %T: %v\n", obj, obj)
	}
	return
}

/*
 *  freesp private functions
 */
func findPort(list []bh.PortIf, prt *port) bool {
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
	ports []bh.PortIf
}

func portListInit() portList {
	return portList{nil}
}

func (l *portList) Append(p *port) {
	l.ports = append(l.ports, p)
}

func (l *portList) Remove(p bh.PortIf) {
	var i int
	for i = range l.ports {
		if p.Node().Name() == l.ports[i].Node().Name() && p.Name() == l.ports[i].Name() {
			break
		}
	}
	if i >= len(l.ports) {
		for _, v := range l.ports {
			log.Printf("portList.RemovePort have bh.PortIf %v\n", v)
		}
		log.Fatalf("portList.RemovePort error: bh.PortIf %v not in this list\n", p)
	}
	for i++; i < len(l.ports); i++ {
		l.ports[i-1] = l.ports[i]
	}
	l.ports = l.ports[:len(l.ports)-1]
}

func (l *portList) Ports() []bh.PortIf {
	return l.ports
}

func (l *portList) Find(nodeName, portName string) (p bh.PortIf, ok bool, index int) {
	ok = false
	for index, p = range l.ports {
		if p.Node().Name() == nodeName && p.Name() == portName {
			ok = true
			return
		}
	}
	return
}
