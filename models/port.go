package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type Port struct {
	freesp.Port
}

var _ TreeElement = Port{}

func (p Port) AddToTree(tree *FilesTreeStore, cursor Cursor) {
	var kind *gdk.Pixbuf
	if p.Direction() == freesp.InPort {
		kind = imageInputPort
	} else {
		kind = imageOutputPort
	}
	err := tree.AddEntry(cursor, kind, p.PortName(), p.Port)
	if err != nil {
		log.Fatal("Port.AddToTree: FilesTreeStore.AddEntry() failed: %s\n", err)
	}
	child := tree.Append(cursor)
	t := SignalType{p.ItsType().SignalType()}
	t.AddToTree(tree, child)
	for _, c := range p.Connections() {
		child = tree.Append(cursor)
		Connection{p.Connection(c)}.AddToTree(tree, child)
	}
	return
}

func (p Port) AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor) {
	switch obj.(type) {
	case freesp.Connection:
		conn := obj.(freesp.Connection)
		log.Println("Port.AddNewObject: conn =", conn)
		var thisPort, otherPort freesp.Port
		if p.Direction() == freesp.InPort {
			otherPort = conn.From
			thisPort = conn.To
		} else {
			otherPort = conn.To
			thisPort = conn.From
		}
		if p.Port != thisPort {
			log.Println("p.Port =", p.Port)
			log.Println("thisPort =", thisPort)
			log.Fatal("Port.AddNewObject error: invalid connection ", conn)
		}
		thisPort.AddConnection(otherPort)
		otherPort.AddConnection(thisPort)
		newCursor = tree.Insert(cursor)
		Connection{conn}.AddToTree(tree, newCursor)
		cCursor := tree.Cursor(otherPort) // TODO: faster search!
		cChild := tree.Append(cCursor)
		Connection{conn}.AddToTree(tree, cChild)

	default:
		fmt.Printf("Port.AddNewObject error: invalid type %T: %v", obj, obj)
		log.Fatal()
	}
	return
}

func (p Port) RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if p.Port != tree.Object(parent) {
		log.Fatal("NodeType.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case freesp.Connection:
		conn := obj.(freesp.Connection)
		var thisPort, otherPort freesp.Port
		if p.Direction() == freesp.InPort {
			otherPort = conn.From
			thisPort = conn.To
			if p.Port != thisPort {
				log.Fatal("Port.AddNewObject error: invalid connection ", conn)
			}
		} else {
			otherPort = conn.To
			thisPort = conn.From
			if p.Port != thisPort {
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
		fmt.Printf("Port.RemoveObject error: invalid type %T: %v", obj, obj)
		log.Fatal()
	}
	return
}

func (p Port) Connection(c freesp.Port) freesp.Connection {
	var from, to freesp.Port
	if p.Direction() == freesp.InPort {
		from = c
		to = p.Port
	} else {
		from = p.Port
		to = c
	}
	return freesp.Connection{from, to}
}

var (
	imageInputPort  *gdk.Pixbuf = nil
	imageOutputPort *gdk.Pixbuf = nil
)

func init_port(iconPath string) (err error) {
	imageInputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/inport-green.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_port error loading inport-green.png: %s", err)
		return
	}
	imageOutputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/outport-red.png", iconPath))
	if err != nil {
		err = fmt.Errorf("init_port error loading outport-red.png: %s", err)
	}
	return
}
