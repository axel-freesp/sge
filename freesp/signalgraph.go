package freesp

import (
	"fmt"
	"log"
	"github.com/axel-freesp/sge/backend"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	gr "github.com/axel-freesp/sge/interface/graph"
)

func SignalGraphNew(filename string, context mod.ModelContextIf) *signalGraph {
	return &signalGraph{filename, SignalGraphTypeNew(context)}
}

func SignalGraphUsesNodeType(s bh.SignalGraphIf, nt bh.NodeTypeIf) bool {
	return SignalGraphTypeUsesNodeType(s.ItsType(), nt)
}

func SignalGraphUsesSignalType(s bh.SignalGraphIf, st bh.SignalTypeIf) bool {
	return SignalGraphTypeUsesSignalType(s.ItsType(), st)
}

type signalGraph struct {
	filename string
	itsType  bh.SignalGraphTypeIf
}

/*
 *  freesp.bh.SignalGraphIf API
 */
var _ bh.SignalGraphIf = (*signalGraph)(nil)

func (s *signalGraph) Filename() string {
	return s.filename
}

func (s *signalGraph) ItsType() bh.SignalGraphTypeIf {
	return s.itsType
}

func (s *signalGraph) Read(data []byte) (cnt int, err error) {
	g := backend.XmlSignalGraphNew()
	cnt, err = g.Read(data)
	if err != nil {
		err = fmt.Errorf("signalGraph.Read error: %v", err)
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename, s.itsType.(*signalGraphType).context,
		func(_ string, _ gr.PortDirection) *portType { return nil })
	return
}

func (s *signalGraph) ReadFile(filepath string) error {
	g := backend.XmlSignalGraphNew()
	err := g.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("signalGraph.ReadFile error: %v", err)
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename, s.itsType.(*signalGraphType).context,
		func(_ string, _ gr.PortDirection) *portType { return nil })
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

func (s *signalGraph) RemoveFromTree(tree tr.TreeIf) {
	gt := s.ItsType()
	tree.Remove(tree.Cursor(s))
	for len(gt.Nodes()) > 0 {
		gt.RemoveNode(gt.Nodes()[0].(*node))
	}
}

func (s *signalGraph) CreateXml() (buf []byte, err error) {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	buf, err = xmlsignalgraph.Write()
	return
}

func (s *signalGraph) Nodes() []bh.NodeIf {
	return s.ItsType().Nodes()
}

/*
 *  tr.TreeElement API
 */

var _ tr.TreeElement = (*signalGraph)(nil)

func (t *signalGraph) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	err := tree.AddEntry(cursor, tr.SymbolSignalGraph, t.Filename(), t, MayAddObject)
	if err != nil {
		log.Fatal("LibraryIf.AddToTree error: AddEntry failed: %s", err)
	}
	t.ItsType().AddToTree(tree, cursor)
}

func (t *signalGraph) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	return t.ItsType().AddNewObject(tree, cursor, obj)
}

func (t *signalGraph) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	return t.ItsType().RemoveObject(tree, cursor)
}
