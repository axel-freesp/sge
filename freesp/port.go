package freesp

import (
	"fmt"
	"unsafe"
)

type port struct {
	name        string
	itsType     PortType
	direction   PortDirection
	connections []Port
	node        Node
}

var _ Port = (*port)(nil)

func newPort(name string, pt *portType, dir PortDirection, n *node) *port {
	return &port{name, pt, dir, nil, n}
}

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

func (p *port) String() (s string) {
	s = fmt.Sprintf("%sPort of node %s(%s, %d connections)",
		p.direction, p.Node().NodeName(), p.name, len(p.connections))
	return
}

func findPort(list []Port, prt *port) bool {
	for _, p := range list {
		if unsafe.Pointer(p.(*port)) == unsafe.Pointer(prt) {
			return true
		}
	}
	return false
}

func (p *port) AddConnection(c Port) error {
	return PortConnect(p, c)
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

func (p *port) RemoveConnection(c Port) {
	i, ok := p.indexOfConnection(c)
	if ok {
		for j := i + 1; j < len(p.connections); j++ {
			p.connections[j-1] = p.connections[j]
		}
		p.connections = p.connections[:len(p.connections)-1]
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
