package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

func SignalGraphNew(filename string, context mod.ModelContextIf) *signalGraph {
	return &signalGraph{filename, "", SignalGraphTypeNew(context)}
}

func SignalGraphUsesNodeType(s bh.SignalGraphIf, nt bh.NodeTypeIf) bool {
	return SignalGraphTypeUsesNodeType(s.ItsType(), nt)
}

func SignalGraphUsesSignalType(s bh.SignalGraphIf, st bh.SignalTypeIf) bool {
	return SignalGraphTypeUsesSignalType(s.ItsType(), st)
}

func SignalGraphApplyHints(s bh.SignalGraphIf, xmlhints *backend.XmlGraphHint) (err error) {
	if s.Filename() != xmlhints.Ref {
		err = fmt.Errorf("SignalGraphApplyHints error: filename mismatch\n")
		return
	}
	g := s.ItsType()
	for _, xmln := range xmlhints.InputNode {
		n, ok := g.NodeByName(xmln.Name)
		if ok {
			n.SetExpanded(false)
			PathModePositionerApplyHints(n, xmln.XmlModeHint)
		}
	}
	for _, xmln := range xmlhints.OutputNode {
		n, ok := g.NodeByName(xmln.Name)
		if ok {
			n.SetExpanded(false)
			PathModePositionerApplyHints(n, xmln.XmlModeHint)
		}
	}
	NodesApplyHints("", g, xmlhints.ProcessingNode)
	return
}

func PathModePositionerApplyHints(pmp gr.PathModePositioner, xmln backend.XmlModeHint) {
	for _, m := range xmln.Entry {
		path, modestring := gr.SeparatePathMode(m.Mode)
		mode, ok := freesp.ModeFromString[string(modestring)]
		if !ok {
			log.Printf("PathModePositionerApplyHints Warning: hint mode %s not defined\n", m.Mode)
			continue
		}
		pos := image.Point{m.X, m.Y}
		pmp.SetPathModePosition(path, mode, pos)
	}
}

func ModePositionerApplyHints(mp gr.ModePositioner, xmln backend.XmlModeHint) {
	for _, m := range xmln.Entry {
		mode, ok := freesp.ModeFromString[m.Mode]
		if !ok {
			log.Printf("ModePositionerApplyHints Warning: hint mode %s not defined\n", m.Mode)
			continue
		}
		pos := image.Point{m.X, m.Y}
		mp.SetModePosition(mode, pos)
	}
}

func NodesApplyHints(path string, g bh.SignalGraphTypeIf, xmlh []backend.XmlNodePosHint) {
	for _, n := range g.ProcessingNodes() {
		var p string
		if len(path) == 0 {
			p = n.Name()
		} else {
			p = fmt.Sprintf("%s/%s", path, n.Name())
		}
		ok := false
		var xmln backend.XmlNodePosHint
		for _, xmln = range xmlh {
			if p == xmln.Name {
				ok = true
				break
			}
		}
		if ok {
			n.SetExpanded(xmln.Expanded)
			PathModePositionerApplyHints(n, xmln.XmlModeHint)
			for i, p := range n.InPorts() {
				ok = p.Name() == xmln.InPorts[i].Name
				if !ok {
					log.Printf("NodesApplyHints FIXME: name mismatch\n")
					continue
				}
				ModePositionerApplyHints(p, xmln.InPorts[i].XmlModeHint)
			}
			for i, p := range n.OutPorts() {
				ok = p.Name() == xmln.OutPorts[i].Name
				if !ok {
					log.Printf("NodesApplyHints FIXME: name mismatch\n")
					continue
				}
				ModePositionerApplyHints(p, xmln.OutPorts[i].XmlModeHint)
			}
		}
		nt := n.ItsType()
		for _, impl := range nt.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				NodesApplyHints(p, impl.Graph(), xmlh)
				break
			}
		}
	}
}

type signalGraph struct {
	filename   string
	pathPrefix string
	itsType    bh.SignalGraphTypeIf
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

func (s signalGraph) PathPrefix() string {
	return s.pathPrefix
}

func (s *signalGraph) SetPathPrefix(newP string) {
	s.pathPrefix = newP
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
	prop := freesp.PropertyNew(true, false, false)
	err := tree.AddEntry(cursor, tr.SymbolSignalGraph, t.Filename(), t, prop)
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
