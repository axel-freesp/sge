package main

import (
	"fmt"
	"github.com/axel-freesp/sge/filemanager"
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
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
	graphviewMap                            map[freesp.ImplementationIf]views.GraphView
	clp                                     *gtk.Clipboard
	signalGraphMgr, libraryMgr, platformMgr freesp.FileManagerIf
	mappingMgr                              freesp.FileManagerMappingIf
}

var _ interfaces.Context = (*Global)(nil)
var _ freesp.ContextIf = (*Global)(nil)
var _ filemanager.FilemanagerContextIf = (*Global)(nil)

func GlobalInit(g *Global) {
	g.graphviewMap = make(map[freesp.ImplementationIf]views.GraphView)
	g.signalGraphMgr = filemanager.FileManagerSGNew(g)
	g.libraryMgr = filemanager.FileManagerLibNew(g)
	g.platformMgr = filemanager.FileManagerPFNew(g)
	g.mappingMgr = filemanager.FileManagerMapNew(g)
}

//
//		freesp.Context interface
//

func (g *Global) SignalGraphMgr() freesp.FileManagerIf {
	return g.signalGraphMgr
}

func (g *Global) LibraryMgr() freesp.FileManagerIf {
	return g.libraryMgr
}

func (g *Global) PlatformMgr() freesp.FileManagerIf {
	return g.platformMgr
}

func (g *Global) MappingMgr() freesp.FileManagerMappingIf {
	return g.mappingMgr
}

func (g *Global) FileMgr(obj freesp.TreeElement) (mgr freesp.FileManagerIf) {
	switch obj.(type) {
	case freesp.SignalGraphIf:
		mgr = global.SignalGraphMgr()
	case freesp.LibraryIf:
		mgr = global.LibraryMgr()
	case freesp.PlatformIf:
		mgr = global.PlatformMgr()
	case freesp.MappingIf:
		mgr = global.MappingMgr()
	default:
		log.Panicf("Global.FileMgr error: wrong type '%T' (%v)\n", obj, obj)
	}
	return
}

//
//		filemanager.FilemanagerContextIf interface
//

func (g *Global) FTS() freesp.TreeMgrIf {
	return g.fts
}

func (g *Global) FTV() freesp.TreeViewIf {
	return g.ftv
}

func (g *Global) GVC() views.GraphViewCollection {
	return g.win.graphViews
}

func (g *Global) ShowAll() {
	g.win.Window().ShowAll()
}

func (g *Global) CleanupNodeTypesFromNodes(nodes []freesp.NodeIf) {
	for _, n := range nodes {
		nt := n.ItsType()
		if !g.NodeTypeIsInUse(nt) {
			g.CleanupNodeType(nt)
		}
	}
}

func (g *Global) CleanupSignalTypesFromNodes(nodes []freesp.NodeIf) {
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
			if impl.ImplementationType() == freesp.NodeTypeGraph {
				g.CleanupSignalTypesFromNodes(impl.Graph().Nodes())
			}
		}
	}
}

func (g *Global) NodeTypeIsInUse(nt freesp.NodeTypeIf) bool {
	var te freesp.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = g.fts.GetObjectById(id)
		switch te.(type) {
		case freesp.SignalGraphIf:
			if freesp.SignalGraphUsesNodeType(te.(freesp.SignalGraphIf), nt) {
				return true
			}
		case freesp.LibraryIf:
			if freesp.LibraryUsesNodeType(te.(freesp.LibraryIf), nt) {
				return true
			}
		}
	}
	return false
}

func (g *Global) CleanupNodeType(nt freesp.NodeTypeIf) {
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == freesp.NodeTypeGraph {
			g.CleanupNodeTypesFromNodes(impl.Graph().Nodes())
		}
	}
	freesp.RemoveRegisteredNodeType(nt)
}

func (g *Global) SignalTypeIsInUse(st freesp.SignalType) bool {
	var te freesp.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = g.fts.GetObjectById(id)
		switch te.(type) {
		case freesp.SignalGraphIf:
			if freesp.SignalGraphUsesSignalType(te.(freesp.SignalGraphIf), st) {
				return true
			}
		case freesp.LibraryIf:
			if freesp.LibraryUsesSignalType(te.(freesp.LibraryIf), st) {
				return true
			}
		}
	}
	return false
}

func (g *Global) CleanupSignalType(st freesp.SignalType) {
	freesp.RemoveRegisteredSignalType(st)
}

//
//		interfaces.Context interface
//

func (g *Global) SelectNode(node interfaces.NodeObject) {
	n := node.(freesp.NodeIf)
	cursor := g.fts.Cursor(n)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) EditNode(node interfaces.NodeObject) {
	log.Printf("Global.EditNode: %v\n", node)
}

func (g *Global) SelectPort(port interfaces.PortObject) {
	p := port.(freesp.Port)
	n := p.Node()
	cursor := g.fts.Cursor(n)
	pCursor := g.fts.CursorAt(cursor, p)
	path, _ := gtk.TreePathNewFromString(pCursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectConnect(conn interfaces.ConnectionObject) {
	c := conn.(freesp.Connection)
	p := c.From()
	n := p.Node()
	cursor := g.fts.Cursor(n)
	pCursor := g.fts.CursorAt(cursor, p)
	cCursor := g.fts.CursorAt(pCursor, c)
	path, _ := gtk.TreePathNewFromString(cCursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectArch(obj interfaces.ArchObject) {
	a := obj.(freesp.ArchIf)
	cursor := g.fts.Cursor(a)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectProcess(obj interfaces.ProcessObject) {
	p := obj.(freesp.ProcessIf)
	a := p.Arch()
	aCursor := g.fts.Cursor(a)
	cursor := g.fts.CursorAt(aCursor, p)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectChannel(obj interfaces.ChannelObject) {
	c := obj.(freesp.ChannelIf)
	pr := c.Process()
	a := pr.Arch()
	aCursor := g.fts.Cursor(a)
	pCursor := g.fts.CursorAt(aCursor, pr)
	cursor := g.fts.CursorAt(pCursor, c)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectMapElement(obj interfaces.MapElemObject) {
	melem := obj.(freesp.MappedElementIf)
	cursor := g.fts.Cursor(melem)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}
