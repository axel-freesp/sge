package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/tool"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"strings"
)

var global Global

func treeSelectionChangedCB(selection *gtk.TreeSelection, menu *GoAppMenu) {
	treeStore := global.fts
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool
	model, iter, ok = selection.GetSelected()
	if ok {
		var err error
		var tpath *gtk.TreePath
		var path string
		var obj tr.TreeElementIf
		tpath, err = model.(*gtk.TreeModel).GetPath(iter)
		if err != nil {
			log.Println("treeSelectionChangedCB: Could not get path from model", err)
			return
		}
		path = tpath.String()
		obj, err = treeStore.GetObject(iter) // This one updates treeStore.Current...
		if err != nil {
			log.Println("treeSelectionChangedCB: Could not get object from model", err)
			obj, err = treeStore.GetObjectById("0")
			if err != nil {
				log.Fatal("treeSelectionChangedCB: Can't show root element")
			}
		}
		MenuEditCurrent(menu, treeStore, global.jl)
		MenuViewCurrent(menu, &global)
		global.win.graphViews.XmlTextView().Set(obj)
		switch obj.(type) {
		case bh.ImplementationIf:
			impl := obj.(bh.ImplementationIf)
			if impl.ImplementationType() == bh.NodeTypeGraph {
				gv, ok := global.graphviewMap[impl]
				if !ok {
					cursor := treeStore.Cursor(obj)
					ntCursor := treeStore.Parent(cursor)
					nt := treeStore.Object(ntCursor).(bh.NodeTypeIf)
					_, err = global.libraryMgr.Access(nt.DefinedAt())
					log.Printf("treeSelectionChangedCB: Need library %s: %v\n", nt.DefinedAt(), ok)
					if err != nil {
						log.Printf("%s\n", err)
						return
					}
					gv, err = views.SignalGraphViewNewFromType(impl.Graph(), &global)
					if err != nil {
						log.Fatal("Could not create graph view.")
					}
					global.win.graphViews.Add(gv, nt.TypeName())
					global.graphviewMap[impl] = gv
					global.ShowAll()
				}
				gv.Sync()
			}
		case bh.NodeIf:
			log.Printf("treeSelectionChangedCB(NodeIf %v)\n", nodeIdFromPath(treeStore, path))
			global.win.graphViews.Select2(obj, nodeIdFromPath(treeStore, path))
		case bh.PortIf:
			p := strings.Split(path, ":")
			npath := strings.Join(p[:len(p)-1], ":")
			global.win.graphViews.Select2(obj, nodeIdFromPath(treeStore, npath))
		case bh.ConnectionIf:
			global.win.graphViews.Select(obj)
		case pf.ArchIf, pf.ProcessIf, pf.ChannelIf, mp.MappedElementIf:
			log.Printf("treeSelectionChangedCB(%T)\n", obj)
			global.win.graphViews.Select(obj)
		}
	}
}

func nodeIdFromPath(fts tr.TreeMgrIf, path string) string {
	p := strings.Split(path, ":")
	obj, err := fts.GetObjectById(path)
	if err != nil {
		log.Panicf("nodeIdFromPath: could not get object of %s\n", path)
	}
	last := obj.(bh.NodeIf).Name()
	//log.Printf("nodeIdFromPath(%s): last=%s\n", path, last)
	if len(p)/3 == 0 {
		return last
	} else {
		parentPath := strings.Join(p[:len(p)-3], ":")
		//log.Printf("nodeIdFromPath(%s): last=%s, return %s/%s\n", path, last, nodeIdFromPath(fts, parentPath), last)
		return fmt.Sprintf("%s/%s", nodeIdFromPath(fts, parentPath), last)
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

	japp := jobApplierNew(global.fts)
	global.jl = jobListNew(japp)

	MenuFileInit(menu)
	MenuEditInit(menu)
	MenuViewInit(menu, &global)
	MenuAboutInit(menu)

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

	global.win.Window().ShowAll()
	gtk.Main()
}
