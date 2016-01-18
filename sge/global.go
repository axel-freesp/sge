package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
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
	win            *GoAppWindow
	jl             *jobList
	fts            *models.FilesTreeStore
	ftv            *views.FilesTreeView
	graphviewMap   map[freesp.Implementation]views.GraphView
	libraryMap     map[string]freesp.Library
	signalGraphMap map[string]freesp.SignalGraph
	platformMap    map[string]freesp.Platform
	mappingMap     map[string]freesp.Mapping
	clp            *gtk.Clipboard
}

var _ interfaces.Context = (*Global)(nil)
var _ freesp.Context = (*Global)(nil)

//
//		interfaces.Context interface
//

func (g *Global) SelectNode(node interfaces.NodeObject) {
	n := node.(freesp.Node)
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
	a := obj.(freesp.Arch)
	cursor := g.fts.Cursor(a)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectProcess(obj interfaces.ProcessObject) {
	p := obj.(freesp.Process)
	a := p.Arch()
	aCursor := g.fts.Cursor(a.(freesp.Arch))
	cursor := g.fts.CursorAt(aCursor, p)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectChannel(obj interfaces.ChannelObject) {
	c := obj.(freesp.Channel)
	pr := c.Process()
	a := pr.Arch()
	aCursor := g.fts.Cursor(a.(freesp.Arch))
	pCursor := g.fts.CursorAt(aCursor, pr.(freesp.Process))
	cursor := g.fts.CursorAt(pCursor, c)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectMapElement(obj interfaces.MapElemObject) {
	melem := obj.(freesp.MappedElement)
	cursor := g.fts.Cursor(melem)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

//
//		freesp.Context interface
//

func (g *Global) GetLibrary(libname string) (lib freesp.Library, err error) {
	var ok bool
	lib, ok = g.libraryMap[libname]
	if ok {
		return
	}
	lib = freesp.LibraryNew(libname, g)
	for _, try := range backend.XmlSearchPaths() {
		err = lib.ReadFile(fmt.Sprintf("%s/%s", try, libname))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetLibrary: library file %s not found", libname)
		return
	}
	log.Println("Global.GetLibrary: library", libname, "successfully loaded")
	g.libraryMap[libname] = lib
	_, err = g.fts.AddLibraryFile(libname, lib)
	return
}

func (g *Global) GetSignalGraph(filename string) (sg freesp.SignalGraph, err error) {
	var ok bool
	sg, ok = g.signalGraphMap[filename]
	if ok {
		return
	}
	sg = freesp.SignalGraphNew(filename, g)
	for _, try := range backend.XmlSearchPaths() {
		err = sg.ReadFile(fmt.Sprintf("%s/%s", try, filename))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetSignalGraph: graph file %s not found", filename)
		return
	}
	_, err = g.fts.AddSignalGraphFile(filename, sg)
	if err != nil {
		log.Println(err)
		return
	}
	var gv views.GraphView
	gv, err = views.SignalGraphViewNew(sg.GraphObject(), g)
	if err != nil {
		err = fmt.Errorf("Global.GetSignalGraph: Could not create graph view.")
		return
	}
	g.win.graphViews.Add(gv, filename)
	log.Println("Global.GetSignalGraph: graph", filename, "successfully loaded")
	g.signalGraphMap[filename] = sg
	return
}

func (g *Global) GetPlatform(filename string) (p freesp.Platform, err error) {
	var ok bool
	p, ok = g.platformMap[filename]
	if ok {
		return
	}
	p = freesp.PlatformNew(filename)
	for _, try := range backend.XmlSearchPaths() {
		log.Printf("Global.GetPlatform: try path %s\n", try)
		err = p.ReadFile(fmt.Sprintf("%s/%s", try, filename))
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetPlatform: platform file %s not found", filename)
		return
	}
	pv, err := views.PlatformViewNew(p.PlatformObject(), g)
	if err != nil {
		err = fmt.Errorf("Global.GetPlatform: Could not create platform view.")
		return
	}
	g.win.graphViews.Add(pv, filename)
	log.Println("Global.GetPlatform: platform", filename, "successfully loaded")
	g.platformMap[filename] = p
	_, err = g.fts.AddPlatformFile(filename, p)
	return
}

func (g *Global) GetMapping(filename string) (m freesp.Mapping, err error) {
	var ok bool
	m, ok = g.mappingMap[filename]
	if ok {
		return
	}
	m = freesp.MappingNew(filename, g)
	for _, try := range backend.XmlSearchPaths() {
		log.Printf("Global.GetMapping: try path %s\n", try)
		err = m.ReadFile(fmt.Sprintf("%s/%s", try, filename))
		if err != nil {
			log.Printf("Global.GetMapping error: %s\n", err)
		}
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetMapping: mapping file %s not found: %s", filename, err)
		return
	}
	mv, err := views.MappingViewNew(m.MappingObject(), g)
	if err != nil {
		err = fmt.Errorf("Global.GetMapping: Could not create platform view.")
		return
	}
	g.win.graphViews.Add(mv, filename)
	log.Println("Global.GetMapping: platform", filename, "successfully loaded")
	g.mappingMap[filename] = m
	_, err = g.fts.AddMappingFile(filename, m)
	return
}

func (g *Global) RemoveLibrary(name string) {
	l, ok := g.libraryMap[name]
	if !ok {
		log.Printf("Global.RemoveLibrary warning: library %s not found.\n", name)
		return
	}
	nodeTypes := l.NodeTypes()
	signalTypes := l.SignalTypes()
	for _, nt := range nodeTypes {
		if !NodeTypeIsInUse(nt) {
			CleanupNodeType(nt)
		}
	}
	for _, st := range signalTypes {
		if !SignalTypeIsInUse(st) {
			CleanupSignalType(st)
		}
	}
	delete(g.libraryMap, name)
	cursor := g.fts.Cursor(l)
	g.fts.RemoveToplevel(cursor.Path)
}

func (g *Global) RemoveSignalGraph(name string) {
	// TODO: check if used in existing mappings
	sg, ok := g.signalGraphMap[name]
	if !ok {
		log.Printf("Global.RemoveSignalGraph warning: graph %s not found.\n", name)
		return
	}
	delete(g.signalGraphMap, name)
	var nodes []freesp.Node
	for _, n := range sg.ItsType().Nodes() {
		nodes = append(nodes, n)
	}
	CleanupNodeTypesFromNodes(nodes)
	CleanupSignalTypesFromNodes(nodes)
	cursor := g.fts.Cursor(sg)
	g.fts.RemoveToplevel(cursor.Path)
	g.win.graphViews.RemoveGraphView(sg.GraphObject())
}

func (g *Global) RemovePlatform(name string) {
	// TODO: check if used in existing mappings
	p, ok := g.platformMap[name]
	if !ok {
		log.Printf("Global.RemovePlatform warning: platform %s not found.\n", name)
		return
	}
	delete(g.platformMap, name)
	cursor := g.fts.Cursor(p)
	g.fts.RemoveToplevel(cursor.Path)
	g.win.graphViews.RemovePlatformView(p.PlatformObject())
}

func (g *Global) RemoveMapping(name string) {
	m, ok := g.mappingMap[name]
	if !ok {
		log.Printf("Global.RemoveMapping warning: mapping %s not found.\n", name)
		return
	}
	delete(g.mappingMap, name)
	cursor := g.fts.Cursor(m)
	g.fts.RemoveToplevel(cursor.Path)
	g.win.graphViews.RemoveMappingView(m.MappingObject())
	// TODO: remove depending graphs and platforms if not used otherwise
}

func (g *Global) RenameLibrary(oldName, newName string) {
	lib, ok := g.libraryMap[oldName]
	if !ok {
		log.Printf("Global.RenameLibrary warning: library %s not found\n", oldName)
		return
	}
	delete(g.libraryMap, oldName)
	g.libraryMap[newName] = lib
}

func (g *Global) RenameSignalGraph(oldName, newName string) {
	sg, ok := g.signalGraphMap[oldName]
	if !ok {
		log.Printf("Global.RenameSignalGraph warning: graph %s not found\n", oldName)
		return
	}
	delete(g.signalGraphMap, oldName)
	g.signalGraphMap[newName] = sg
}

func (g *Global) RenamePlatform(oldName, newName string) {
	p, ok := g.platformMap[oldName]
	if !ok {
		log.Printf("Global.RenamePlatform warning: platform %s not found\n", oldName)
		return
	}
	delete(g.platformMap, oldName)
	g.platformMap[newName] = p
}

func (g *Global) RenameMapping(oldName, newName string) {
	m, ok := g.mappingMap[oldName]
	if !ok {
		log.Printf("Global.RenameMapping warning: mapping %s not found\n", oldName)
		return
	}
	delete(g.mappingMap, oldName)
	g.mappingMap[newName] = m
}

func (g *Global) AddNewLibrary(lib freesp.Library) {
	g.libraryMap[lib.Filename()] = lib
}

func (g *Global) AddNewSignalGraph(sg freesp.SignalGraph) {
	g.signalGraphMap[sg.Filename()] = sg
}

func (g *Global) AddNewPlatform(p freesp.Platform) {
	g.platformMap[p.Filename()] = p
}

func (g *Global) AddNewMapping(m freesp.Mapping) {
	g.mappingMap[m.Filename()] = m
}

//
//		Local functions
//

func CleanupNodeType(nt freesp.NodeType) {
	for _, impl := range nt.Implementation() {
		if impl.ImplementationType() == freesp.NodeTypeGraph {
			CleanupNodeTypesFromNodes(impl.Graph().Nodes())
		}
	}
	freesp.RemoveRegisteredNodeType(nt)
}

func CleanupNodeTypesFromNodes(nodes []freesp.Node) {
	for _, n := range nodes {
		nt := n.ItsType()
		if !NodeTypeIsInUse(nt) {
			CleanupNodeType(nt)
		}
	}
}

func CleanupSignalType(st freesp.SignalType) {
	freesp.RemoveRegisteredSignalType(st)
}

func CleanupSignalTypesFromNodes(nodes []freesp.Node) {
	for _, n := range nodes {
		for _, p := range n.InPorts() {
			st := p.SignalType()
			if !SignalTypeIsInUse(st) {
				CleanupSignalType(st)
			}
		}
		for _, p := range n.OutPorts() {
			st := p.SignalType()
			if !SignalTypeIsInUse(st) {
				CleanupSignalType(st)
			}
		}
		nt := n.ItsType()
		for _, impl := range nt.Implementation() {
			if impl.ImplementationType() == freesp.NodeTypeGraph {
				CleanupSignalTypesFromNodes(impl.Graph().Nodes())
			}
		}
	}
}

func NodeTypeIsInUse(nt freesp.NodeType) bool {
	var te freesp.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = global.fts.GetObjectById(id)
		switch te.(type) {
		case freesp.SignalGraph:
			if freesp.SignalGraphUsesNodeType(te.(freesp.SignalGraph), nt) {
				return true
			}
		case freesp.Library:
			if freesp.LibraryUsesNodeType(te.(freesp.Library), nt) {
				return true
			}
		}
	}
	return false
}

func SignalTypeIsInUse(st freesp.SignalType) bool {
	var te freesp.TreeElement
	var err error
	for i := 0; err == nil; i++ {
		id := fmt.Sprintf("%d", i)
		te, err = global.fts.GetObjectById(id)
		switch te.(type) {
		case freesp.SignalGraph:
			if freesp.SignalGraphUsesSignalType(te.(freesp.SignalGraph), st) {
				return true
			}
		case freesp.Library:
			if freesp.LibraryUsesSignalType(te.(freesp.Library), st) {
				return true
			}
		}
	}
	return false
}
