package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/axel-freesp/sge/views/graph"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

const (
	width  = 800
	height = 600
)

type Global struct {
	win          *GoAppWindow
	jl           *jobList
	graphview    []*views.GraphView
	xmlview      *views.XmlTextView
	fts          *models.FilesTreeStore
	ftv          *views.FilesTreeView
	graphviewMap map[freesp.Implementation]*views.GraphView
	libraryMap   map[string]freesp.Library
}

var _ views.Context = (*Global)(nil)

func (g *Global) SelectNode(node graph.NodeObject) {
	n := node.(freesp.Node)
	cursor := g.fts.Cursor(n)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) EditNode(node graph.NodeObject) {
	log.Printf("Global.EditNode: %v\n", node)
}

func (g *Global) SelectPort(port freesp.Port) {
	n := port.Node()
	cursor := g.fts.Cursor(n)
	pCursor := g.fts.CursorAt(cursor, port)
	path, _ := gtk.TreePathNewFromString(pCursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) SelectConnect(conn freesp.Connection) {
	p := conn.From()
	n := p.Node()
	cursor := g.fts.Cursor(n)
	pCursor := g.fts.CursorAt(cursor, p)
	cCursor := g.fts.CursorAt(pCursor, conn)
	path, _ := gtk.TreePathNewFromString(cCursor.Path)
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
		global.xmlview.Set(obj)
		switch obj.(type) {
		case freesp.Node, freesp.Port, freesp.Connection:
			for _, v := range global.graphview {
				v.Select(obj)
			}
		case freesp.Implementation:
			impl := obj.(freesp.Implementation)
			if impl.ImplementationType() == freesp.NodeTypeGraph {
				log.Println("treeSelectionChangedCB: Have graph implementation to show")
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
					gv, err = views.GraphViewNew(impl.Graph(), width, height, &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.graphview = append(global.graphview, gv)
					global.win.stack.AddTitled(gv.Widget(), nt.TypeName(), nt.TypeName())
					global.graphviewMap[impl] = gv
					gv.Widget().ShowAll()
					log.Println("treeSelectionChangedCB: Created graphview for implementation to show")
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
	global.graphviewMap = make(map[freesp.Implementation]*views.GraphView)
	global.libraryMap = make(map[string]freesp.Library)

	var err error
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

	global.xmlview, err = views.XmlTextViewNew(width, height)
	if err != nil {
		log.Fatal("Could not create XML view.")
	}
	global.win.stack.AddTitled(global.xmlview.Widget(), "XML View", "XML View")

	selection, err := global.ftv.TreeView().GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	selection.Connect("changed", treeSelectionChangedCB, menu)

	// Handle command line arguments: treat each as a filename:
	for i, p := range unhandledArgs {
		if i > 0 {
			filepath := fmt.Sprintf("%s/%s", backend.XmlRoot(), p)
			switch tool.Suffix(p) {
			case "sml":
				var sg freesp.SignalGraph
				sg = freesp.SignalGraphNew(p, &global)
				err1 := sg.ReadFile(filepath)
				if err1 == nil {
					log.Println("Loading signal graph", filepath)
					global.fts.AddSignalGraphFile(p, sg)
					gv, err := views.GraphViewNew(sg.ItsType(), width, height, &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.graphview = append(global.graphview, gv)
					global.win.stack.AddTitled(gv.Widget(), p, p)
				}
			case "alml":
				_, err = global.GetLibrary(p)
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
	lib = freesp.LibraryNew(libname)
	for _, try := range backend.XmlSearchPaths() {
		err = lib.ReadFile(fmt.Sprintf("%s/%s", try, libname), g)
		if err == nil {
			break
		}
	}
	if err != nil {
		err = fmt.Errorf("Global.GetLibrary: library file %s not found", libname)
		return
	}
	log.Println("Global.GetLibrary: library", libname, "successfully loaded")
	if err == nil {
		g.libraryMap[libname] = lib
		_, err = g.fts.AddLibraryFile(libname, lib)
	}
	return
}
