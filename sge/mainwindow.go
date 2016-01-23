package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

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
		case freesp.ImplementationIf:
			impl := obj.(freesp.ImplementationIf)
			if impl.ImplementationType() == freesp.NodeTypeGraph {
				gv, ok := global.graphviewMap[impl]
				if !ok {
					cursor := treeStore.Cursor(obj)
					ntCursor := treeStore.Parent(cursor)
					nt := treeStore.Object(ntCursor).(freesp.NodeTypeIf)
					_, err = global.libraryMgr.Access(nt.DefinedAt())
					log.Printf("treeSelectionChangedCB: Need library %s: %v\n", nt.DefinedAt(), ok)
					if err != nil {
						log.Printf("%s\n", err)
						return
					}
					gv, err = views.SignalGraphViewNew(impl.GraphObject(), &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.win.graphViews.Add(gv, nt.TypeName())
					global.graphviewMap[impl] = gv
					global.ShowAll()
				}
				gv.Sync()
			}
		case freesp.NodeIf, freesp.Port, freesp.Connection, freesp.ArchIf, freesp.ProcessIf, freesp.ChannelIf, freesp.MappedElementIf:
			global.win.graphViews.Select(obj)
		}
	}
}

func main() {
	// Initialize GTK with parsing the command line arguments.
	unhandledArgs := os.Args
	gtk.Init(&unhandledArgs)
	backend.Init()
	freesp.Init()
	GlobalInit(&global)

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
				_, err = global.SignalGraphMgr().Access(p)
			case "alml":
				_, err = global.LibraryMgr().Access(p)
			case "spml":
				_, err = global.PlatformMgr().Access(p)
			case "mml":
				_, err = global.MappingMgr().Access(p)
			default:
				log.Println("Warning: unknown suffix", tool.Suffix(p))
			}
			if err != nil {
				log.Println(err)
				continue
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
