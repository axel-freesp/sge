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
	win       *GoAppWindow
	jl        *jobList
	graphview []*views.GraphView
	fts       *models.FilesTreeStore
	ftv       *views.FilesTreeView
}

var _ views.Context = (*Global)(nil)

func (g *Global) SelectNode(node graph.GObject) {
	//log.Printf("Global.SelectNode: %v\n", node)
	n := node.(freesp.Node)
	cursor := g.fts.Cursor(n)
	path, _ := gtk.TreePathNewFromString(cursor.Path)
	g.ftv.TreeView().ExpandToPath(path)
	g.ftv.TreeView().SetCursor(path, g.ftv.TreeView().GetExpanderColumn(), false)
}

func (g *Global) EditNode(node graph.GObject) {
	log.Printf("Global.EditNode: %v\n", node)
}

var global Global

type selectionArg struct {
	treeStore *models.FilesTreeStore
	xmlview   *views.XmlTextView
	menu      *GoAppMenu
}

func treeSelectionChangedCB(selection *gtk.TreeSelection, arg *selectionArg) {
	treeStore := arg.treeStore
	xmlview := arg.xmlview
	var iter gtk.TreeIter
	var model gtk.ITreeModel
	if selection.GetSelected(&model, &iter) {
		obj, err := treeStore.GetObject(&iter)
		if err != nil {
			log.Println("treeSelectionChangedCB: Could not get object from model", err)
			obj, err = treeStore.GetObjectById("0")
			if err != nil {
				log.Fatal("treeSelectionChangedCB: Can't show root element")
			}
		}
		MenuEditCurrent(arg.menu, treeStore, global.jl)
		xmlview.Set(obj)
	}
}

func main() {
	// Initialize GTK with parsing the command line arguments.
	unhandledArgs := os.Args
	gtk.Init(&unhandledArgs)
	backend.Init()
	freesp.Init()

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

	xmlview, err := views.XmlTextViewNew(width, height)
	if err != nil {
		log.Fatal("Could not create XML view.")
	}
	global.win.stack.AddTitled(xmlview.Widget(), "XML View", "XML View")

	selection, err := global.ftv.TreeView().GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	arg := &selectionArg{global.fts, xmlview, menu}
	selection.Connect("changed", treeSelectionChangedCB, arg)

	// Handle command line arguments: treat each as a filename:
	for i, p := range unhandledArgs {
		if i > 0 {
			filepath := fmt.Sprintf("%s/%s", backend.XmlRoot(), p)
			switch tool.Suffix(p) {
			case "sml":
				var sg freesp.SignalGraph
				sg = freesp.SignalGraphNew(p)
				err1 := sg.ReadFile(filepath)
				if err1 == nil {
					log.Println("Loading signal graph", filepath)
					global.fts.AddSignalGraphFile(p, sg)
					gv, err := views.GraphViewNew(sg, width, height, &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.graphview = append(global.graphview, gv)
					global.win.stack.AddTitled(gv.Widget(), p, p)
				}
			case "alml":
				var lib freesp.Library
				lib = freesp.LibraryNew(p)
				err2 := lib.ReadFile(filepath)
				if err2 == nil {
					log.Println("Loading library file", filepath)
					global.fts.AddLibraryFile(p, lib)
					continue
				}
				log.Println("Warning: Could not read file ", filepath)
				log.Println(err2)
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
