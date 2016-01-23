package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	interfaces "github.com/axel-freesp/sge/interface"
	"log"
)

func SignalGraphNew(filename string, context ContextIf) *signalGraph {
	return &signalGraph{filename, SignalGraphTypeNew(context)}
}

func SignalGraphUsesNodeType(s SignalGraphIf, nt NodeTypeIf) bool {
	return SignalGraphTypeUsesNodeType(s.ItsType(), nt)
}

func SignalGraphUsesSignalType(s SignalGraphIf, st SignalType) bool {
	return SignalGraphTypeUsesSignalType(s.ItsType(), st)
}

type signalGraph struct {
	filename string
	itsType  SignalGraphTypeIf
}

/*
 *  freesp.SignalGraphIf API
 */
var _ SignalGraphIf = (*signalGraph)(nil)

func (s *signalGraph) Filename() string {
	return s.filename
}

func (s *signalGraph) ItsType() SignalGraphTypeIf {
	return s.itsType
}

func (s *signalGraph) GraphObject() interfaces.GraphObject {
	return s.itsType.(*signalGraphType)
}

func (s *signalGraph) Read(data []byte) (cnt int, err error) {
	g := backend.XmlSignalGraphNew()
	cnt, err = g.Read(data)
	if err != nil {
		err = fmt.Errorf("signalGraph.Read error: %v", err)
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename, s.itsType.(*signalGraphType).context,
		func(_ string, _ interfaces.PortDirection) *portType { return nil })
	return
}

func (s *signalGraph) ReadFile(filepath string) error {
	g := backend.XmlSignalGraphNew()
	err := g.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("signalGraph.ReadFile error: %v", err)
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename, s.itsType.(*signalGraphType).context,
		func(_ string, _ interfaces.PortDirection) *portType { return nil })
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

func (s *signalGraph) RemoveFromTree(tree Tree) {
	gt := s.ItsType()
	tree.Remove(tree.Cursor(s))
	for len(gt.Nodes()) > 0 {
		gt.RemoveNode(gt.Nodes()[0].(*node))
	}
}

/*
 *  TreeElement API
 */

var _ TreeElement = (*signalGraph)(nil)

func (t *signalGraph) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolSignalGraph, t.Filename(), t, mayAddObject)
	if err != nil {
		log.Fatal("LibraryIf.AddToTree error: AddEntry failed: %s", err)
	}
	t.ItsType().AddToTree(tree, cursor)
}

func (t *signalGraph) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	return t.ItsType().AddNewObject(tree, cursor, obj)
}

func (t *signalGraph) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	return t.ItsType().RemoveObject(tree, cursor)
}
