package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	//"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

const (
	width  = 800
	height = 600
)

var (
	win *GoAppWindow
	jl  *jobList
)

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
		MenuEditCurrent(arg.menu, treeStore, jl)
		xmlview.Set(obj)
	}
}

func main() {
	// Initialize GTK with parsing the command line arguments.
	unhandledArgs := os.Args
	gtk.Init(&unhandledArgs)
	backend.Init()
	freesp.Init()

	// Create a new toplevel window.
	win, err := GoAppWindowNew(width, height)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	menu := GoAppMenuNew()
	menu.Init()
	win.layout_box.Add(menu.menubar)

	err = models.Init()
	if err != nil {
		log.Fatal("Unable to initialize models:", err)
	}

	fts, err := models.FilesTreeStoreNew()
	if err != nil {
		log.Fatal("Unable to create FilesTreeStore:", err)
	}
	ftv, err := views.FilesTreeViewNew(fts, width/2, height)
	if err != nil {
		log.Fatal("Unable to create FilesTreeView:", err)
	}
	win.navigation_box.Add(ftv.Widget())

	xmlview, err := views.XmlTextViewNew(width, height)
	if err != nil {
		log.Fatal("Could not create XML view.")
	}
	win.stack.AddTitled(xmlview.Widget(), "XML View", "XML View")

	selection, err := ftv.TreeView().GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	arg := &selectionArg{fts, xmlview, menu}
	selection.Connect("changed", treeSelectionChangedCB, arg)

	if len(unhandledArgs) < 2 {
		_, err := fts.AddSignalGraphFile("new-file.sml", freesp.SignalGraphNew("new-file.sml"))
		if err != nil {
			log.Fatal("ftv.AddSignalGraphFile('new-file.sml') failed.")
		}
	}
	// Handle command line arguments: treat each as a filename:
	for i, p := range unhandledArgs {
		if i > 0 {
			filepath := fmt.Sprintf("%s/%s", backend.XmlRoot(), p)
			var sg freesp.SignalGraph
			sg = freesp.SignalGraphNew(p)
			err1 := sg.ReadFile(filepath)
			if err1 == nil {
				log.Println("Loading signal graph", filepath)
				fts.AddSignalGraphFile(p, sg)
				graphview, err := views.GraphViewNew(sg, width, height)
				if err != nil {
					log.Fatal("Could not create graph view.")
				}
				win.stack.AddTitled(graphview.Widget(), "Graph View", "Graph View")
				continue
			}
			var lib freesp.Library
			lib = freesp.LibraryNew(p)
			err2 := lib.ReadFile(filepath)
			if err2 == nil {
				log.Println("Loading library file", filepath)
				fts.AddLibraryFile(p, lib)
				continue
			}
			log.Println("Warning: Could not read file ", filepath)
			log.Println(err1)
			log.Println(err2)
		}
	}

	japp := jobApplierNew(fts)
	jl = jobListNew(japp)

	MenuFileInit(menu, fts, ftv)
	MenuEditInit(menu, fts, jl, ftv)
	MenuAboutInit(menu)

	win.Window().ShowAll()
	gtk.Main()
}
