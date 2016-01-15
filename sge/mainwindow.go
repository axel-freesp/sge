package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
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

var _ views.Context = (*Global)(nil)
var _ freesp.Context = (*Global)(nil)

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

var global Global

func treeSelectionChangedCB(selection *gtk.TreeSelection, menu *GoAppMenu) {
	treeStore := global.fts
	var iter gtk.TreeIter
	var model gtk.ITreeModel
	if selection.GetSelected(&model, &iter) {
		obj, err := treeStore.GetObject(&iter) // This one updates treeStore.Current...
		if err != nil {
			log.Println("treeSelectionChangedCB: Could not get object from model", err)
			obj, err = treeStore.GetObjectById("0")
			if err != nil {
				log.Fatal("treeSelectionChangedCB: Can't show root element")
			}
		}
		MenuEditCurrent(menu, treeStore, global.jl)
		global.win.graphViews.XmlTextView().Set(obj)
		switch obj.(type) {
		case freesp.Node, freesp.Port, freesp.Connection, freesp.Arch, freesp.Process, freesp.Channel:
			global.win.graphViews.Select(obj)
		case freesp.Implementation:
			impl := obj.(freesp.Implementation)
			if impl.ImplementationType() == freesp.NodeTypeGraph {
				gv, ok := global.graphviewMap[impl]
				if !ok {
					cursor := treeStore.Cursor(obj)
					ntCursor := treeStore.Parent(cursor)
					nt := treeStore.Object(ntCursor).(freesp.NodeType)
					_, ok := global.libraryMap[nt.DefinedAt()]
					log.Printf("treeSelectionChangedCB: Need library %s: %v\n", nt.DefinedAt(), ok)
					if !ok {
						return
					}
					gv, err = views.SignalGraphViewNew(impl.GraphObject(), &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.win.graphViews.Add(gv, nt.TypeName())
					global.graphviewMap[impl] = gv
					gv.Widget().ShowAll()
				}
				gv.Sync()
			}
		}
	}
}

func main() {
	// Initialize GTK with parsing the command line arguments.
	unhandledArgs := os.Args
	gtk.Init(&unhandledArgs)
	backend.Init()
	freesp.Init()
	global.graphviewMap = make(map[freesp.Implementation]views.GraphView)
	global.libraryMap = make(map[string]freesp.Library)
	global.signalGraphMap = make(map[string]freesp.SignalGraph)
	global.platformMap = make(map[string]freesp.Platform)
	global.mappingMap = make(map[string]freesp.Mapping)

	var err error
	iconPath := os.Getenv("SGE_ICON_PATH")
	if len(iconPath) == 0 {
		log.Println("WARNING: Missing environment variable SGE_ICON_PATH")
	} else {
		err = gtk.WindowSetDefaultIconFromFile(fmt.Sprintf("%s/%s", iconPath, "sge-logo.png"))
		if err != nil {
			log.Printf("WARNING: Failed to set default icon: %s.\n", err)
		}
	}

	atom := gdk.GdkAtomIntern("CLIPBOARD", false)
	global.clp, err = gtk.ClipboardGet(atom)
	if err != nil {
		log.Println("WARNING: Failed to get clipboard CLIPBOARD")
	}

	// Create a new toplevel window.
	global.win, err = GoAppWindowNew(width, height)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	menu := GoAppMenuNew()
	menu.Init()
	global.win.layout_box.Add(menu.menubar)

	err = models.Init()
	if err != nil {
		log.Fatal("Unable to initialize models:", err)
	}

	global.fts, err = models.FilesTreeStoreNew()
	if err != nil {
		log.Fatal("Unable to create FilesTreeStore:", err)
	}
	global.ftv, err = views.FilesTreeViewNew(global.fts, width/2, height)
	if err != nil {
		log.Fatal("Unable to create FilesTreeView:", err)
	}
	global.win.navigation_box.Add(global.ftv.Widget())

	selection, err := global.ftv.TreeView().GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	selection.Connect("changed", treeSelectionChangedCB, menu)

	// Handle command line arguments: treat each as a filename:
	for i, p := range unhandledArgs {
		if i > 0 {
			switch tool.Suffix(p) {
			case "sml":
				_, err := global.GetSignalGraph(p)
				if err != nil {
					log.Println(err)
					continue
				}
			case "alml":
				_, err = global.GetLibrary(p)
			case "spml":
				_, err := global.GetPlatform(p)
				if err != nil {
					log.Println(err)
					continue
				}
			case "mml":
				_, err = global.GetMapping(p)
				if err != nil {
					log.Println(err)
					continue
				}
			default:
				log.Println("Warning: unknown suffix", tool.Suffix(p))
			}
		}
	}

	japp := jobApplierNew(global.fts)
	global.jl = jobListNew(japp)

	MenuFileInit(menu)
	MenuEditInit(menu)
	MenuAboutInit(menu)

	global.win.Window().ShowAll()
	gtk.Main()
}

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
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetMapping: platform file %s not found: %s", filename, err)
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

func (g *Global) AddNewLibrary(lib freesp.Library) {
	g.libraryMap[lib.Filename()] = lib
}

func (g *Global) RenameLibrary(oldName, newName string) {
	lib, ok := g.libraryMap[oldName]
	if !ok {
		log.Fatalf("Global.RenameLibrary error: library %s not found\n", oldName)
	}
	delete(g.libraryMap, oldName)
	g.libraryMap[newName] = lib
}

func (g *Global) RemoveLibrary(name string) {
	l := g.libraryMap[name]
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
}

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
