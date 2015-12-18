package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

func SignalGraphNew(filename string) *signalGraph {
	return &signalGraph{filename, SignalGraphTypeNew()}
}

type signalGraph struct {
	filename string
	itsType  SignalGraphType
}

var _ SignalGraph = (*signalGraph)(nil)

func (s *signalGraph) Filename() string {
	return s.filename
}

func (s *signalGraph) ItsType() SignalGraphType {
	return s.itsType
}

func (s *signalGraph) Read(data []byte) error {
	g := backend.XmlSignalGraphNew()
	err := g.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.Read: %v", err))
	}
	fmt.Println("signalGraph.Read: call createSignalGraphTypeFromXml")
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename,
		func(_ string, _ PortDirection) *namedPortType { return nil })
	return err
}

func (s *signalGraph) ReadFile(filepath string) error {
	g := backend.XmlSignalGraphNew()
	err := g.ReadFile(filepath)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.ReadFile: %v", err))
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename,
		func(_ string, _ PortDirection) *namedPortType { return nil })
	return err
}

func (s *signalGraph) Write() (data []byte, err error) {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	data, err = xmlsignalgraph.Write()
	return
}

func (s *signalGraph) WriteFile(filepath string) error {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	return xmlsignalgraph.WriteFile(filepath)
}

func (s *signalGraph) SetFilename(filename string) {
	s.filename = filename
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*signalGraph)(nil)

func (t *signalGraph) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolSignalGraph, t.Filename(), t.ItsType())
	if err != nil {
		log.Fatal("Library.AddToTree error: AddEntry failed: %s", err)
	}
	t.ItsType().AddToTree(tree, cursor)
}

func (t *signalGraph) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	log.Fatal("signalGraph.AddNewObject - nothing to add.")
	return
}

func (t *signalGraph) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	log.Fatal("signalGraph.RemoveObject - nothing to remove.")
	return
}

//------------------------------

type signalGraphError struct {
	reason string
}

func (e *signalGraphError) Error() string {
	return fmt.Sprintf("signal graph error: %s", e.reason)
}

func newSignalGraphError(reason string) *signalGraphError {
	return &signalGraphError{reason}
}
