package main

import (
	"fmt"
	"github.com/axel-freesp/sge/filemanager"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/freesp/behaviour"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	mod "github.com/axel-freesp/sge/interface/model"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
)

const (
	width  = 800
	height = 600
)

type Global struct {
	win                                     *GoAppWindow
	jl                                      *jobList
	fts                                     *models.FilesTreeStore
	ftv                                     *views.FilesTreeView
	graphviewMap                            map[bh.ImplementationIf]views.GraphViewIf
	clp                                     *gtk.Clipboard
	signalGraphMgr, libraryMgr, platformMgr mod.FileManagerIf
	mappingMgr                              mod.FileManagerMappingIf
}

var _ views.ContextIf = (*Global)(nil)
var _ mod.ModelContextIf = (*Global)(nil)
var _ filemanager.FilemanagerContextIf = (*Global)(nil)

func GlobalInit(g *Global) {
	g.graphviewMap = make(map[bh.ImplementationIf]views.GraphViewIf)
	g.signalGraphMgr = filemanager.FileManagerSGNew(g)
	g.libraryMgr = filemanager.FileManagerLibNew(g)
	g.platformMgr = filemanager.FileManagerPFNew(g)
	g.mappingMgr = filemanager.FileManagerMapNew(g)
}

//
//		freesp.Context interface
//

func (g *Global) SignalGraphMgr() mod.FileManagerIf {
	return g.signalGraphMgr
}

func (g *Global) LibraryMgr() mod.FileManagerIf {
	return g.libraryMgr
}

func (g *Global) PlatformMgr() mod.FileManagerIf {
	return g.platformMgr
}

func (g *Global) MappingMgr() mod.FileManagerMappingIf {
	return g.mappingMgr
}

func (g *Global) FileMgr(obj tr.TreeElement) (mgr mod.FileManagerIf) {
	switch obj.(type) {
	case bh.SignalGraphIf:
		mgr = global.SignalGraphMgr()
	case bh.LibraryIf:
		mgr = global.LibraryMgr()
	case pf.PlatformIf:
		mgr = global.PlatformMgr()
	case mp.MappingIf:
		mgr = global.MappingMgr()
	default:
		log.Panicf("Global.FileMgr error: wrong type '%T' (%v)\n", obj, obj)
	}
	return
}

//
//		filemanager.FilemanagerContextIf interface
//

func (g *Global) FTS() tr.TreeMgrIf {
	return g.fts
}

func (g *Global) FTV() tr.TreeViewIf {
	return g.ftv
}

func (g *Global) GVC() views.GraphViewCollectionIf {
	return g.win.graphViews
}

func (g *Global) ShowAll() {
	g.win.Window().ShowAll()
}

func (g *Global) CleanupNodeTypesFromNodes(nodes []bh.NodeIf) {
	for _, n := range nodes {
		nt := n.ItsType()
		if !g.NodeTypeIsInUse(nt) {
			g.CleanupNodeType(nt)
		}
	}
}

func (g *Global) CleanupSignalTypesFromNodes(nodes []bh.NodeIf) {
	for _, n := range nodes {
		for _, p := range n.InPorts() {
			st := p.SignalType()
			if !g.SignalTypeIsInUse(st) {
				g.CleanupSignalType(st)
			}
		}
		for _, p := range n.OutPorts() {
			st := p.SignalType()
			if !g.SignalTypeIsInUse(st) {
				g.CleanupSignalType(st)
			}
		}
		nt := n.ItsType()
		for _, impl := range nt.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				g.CleanupSignalTypesFromNodes(impl.Graph().Nodes())
			}
		}
	}
}

func (g *Global) NodeTypeIsInUse(nt bh.NodeTypeIf) bool {
	var te tr.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = g.fts.GetObjectById(id)
		switch te.(type) {
		case bh.SignalGraphIf:
			if behaviour.SignalGraphUsesNodeType(te.(bh.SignalGraphIf), nt) {
				return true
			}
		case bh.LibraryIf:
			if behaviour.LibraryUsesNodeType(te.(bh.LibraryIf), nt) {
				return true
			}
		}
	}
	return false
}

func (g *Global) CleanupNodeType(nt bh.NodeTypeIf) {
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			g.CleanupNodeTypesFromNodes(impl.Graph().Nodes())
		}
	}
	freesp.RemoveRegisteredNodeType(nt)
}

func (g *Global) SignalTypeIsInUse(st bh.SignalTypeIf) bool {
	var te tr.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = g.fts.GetObjectById(id)
		switch te.(type) {
		case bh.SignalGraphIf:
			if behaviour.SignalGraphUsesSignalType(te.(bh.SignalGraphIf), st) {
				return true
			}
		case bh.LibraryIf:
			if behaviour.LibraryUsesSignalType(te.(bh.LibraryIf), st) {
				return true
			}
		}
	}
	return false
}

func (g *Global) CleanupSignalType(st bh.SignalTypeIf) {
	freesp.RemoveRegisteredSignalType(st)
}

//
//		views.Context interface
//

func (g *Global) SelectNode(n bh.NodeIf, id bh.NodeIdIf) {
	//log.Printf("Global.SelectNode(%v)\n", id)
	cursor := g.nodePath(n, g.fts.Cursor(n), id)
	path, _ := gtk.TreePathNewFromString(g.fts.Parent(cursor).Path)
	g.ftv.TreeView().ExpandToPath(path)
	path, _ = gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) nodePath(n bh.NodeIf, nCursor tr.Cursor, selectId bh.NodeIdIf) (cursor tr.Cursor) {
	ids := strings.Split(selectId.String(), "/")
	if len(ids) == 1 {
		cursor = nCursor
		return
	}
	nt := n.ItsType()
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == bh.NodeTypeGraph {
			for _, nn := range impl.Graph().ProcessingNodes() {
				if nn.Name() == ids[1] {
					cursor = g.nodePath(nn, g.fts.CursorAt(nCursor, nn), behaviour.NodeIdFromString(strings.Join(ids[1:], "/")))
					break
				}
			}
			break
		}
	}
	return
}

func (g *Global) EditNode(node bh.NodeIf, id bh.NodeIdIf) {
	log.Printf("Global.EditNode: %v\n", node)
}

func (g *Global) SelectPort(p bh.PortIf, n bh.NodeIf, id bh.NodeIdIf) {
	cursor := g.nodePath(n, g.fts.Cursor(n), id)
	pCursor := g.fts.CursorAt(cursor, p)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	path, _ = gtk.TreePathNewFromString(pCursor.Path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectConnect(conn bh.ConnectionIf) {
	c := conn.(bh.ConnectionIf)
	p := c.From()
	n := p.Node()
	cursor := g.fts.Cursor(n)
	pCursor := g.fts.CursorAt(cursor, p)
	cCursor := g.fts.CursorAt(pCursor, c)
	path, _ := gtk.TreePathNewFromString(cCursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectArch(obj pf.ArchIf) {
	a := obj.(pf.ArchIf)
	cursor := g.fts.Cursor(a)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectProcess(obj pf.ProcessIf) {
	p := obj.(pf.ProcessIf)
	a := p.Arch()
	aCursor := g.fts.Cursor(a)
	cursor := g.fts.CursorAt(aCursor, p)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectChannel(obj pf.ChannelIf) {
	c := obj.(pf.ChannelIf)
	pr := c.Process()
	a := pr.Arch()
	aCursor := g.fts.Cursor(a)
	pCursor := g.fts.CursorAt(aCursor, pr)
	cursor := g.fts.CursorAt(pCursor, c)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectMapElement(melem mp.MappedElementIf) {
	cursor := g.fts.Cursor(melem)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}
