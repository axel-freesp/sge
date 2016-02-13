package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	//"github.com/gotk3/gotk3/glib"
)

type GoAppMenu struct {
	menubar      *gtk.MenuBar
	menuFile     *gtk.Menu
	filemenu     *gtk.MenuItem
	fileNewSg    *gtk.MenuItem
	fileNewLib   *gtk.MenuItem
	fileNewPlat  *gtk.MenuItem
	fileNewMap   *gtk.MenuItem
	fileOpen     *gtk.MenuItem
	fileSave     *gtk.MenuItem
	fileSaveAs   *gtk.MenuItem
	fileClose    *gtk.MenuItem
	fileQuit     *gtk.MenuItem
	menuEdit     *gtk.Menu
	editmenu     *gtk.MenuItem
	editUndo     *gtk.MenuItem
	editRedo     *gtk.MenuItem
	editNew      *gtk.MenuItem
	editEdit     *gtk.MenuItem
	editDelete   *gtk.MenuItem
	editCopy     *gtk.MenuItem
	editPaste    *gtk.MenuItem
	menuView     *gtk.Menu
	viewmenu     *gtk.MenuItem
	viewExpand   *gtk.MenuItem
	viewCollapse *gtk.MenuItem
	menuAbout    *gtk.Menu
	aboutmenu    *gtk.MenuItem
	aboutAbout   *gtk.MenuItem
	aboutHelp    *gtk.MenuItem

	aboutdialog *gtk.AboutDialog
}

func GoAppMenuNew() *GoAppMenu {
	return &GoAppMenu{}
}

func (m *GoAppMenu) Init() {
	var err error

	m.menubar, err = gtk.MenuBarNew()
	if err != nil {
		log.Fatal("Unable to create menubar:", err)
	}

	m.menuFile, err = gtk.MenuNew()
	if err != nil {
		log.Fatal("Unable to create menuFile:", err)
	}
	m.filemenu, err = gtk.MenuItemNewWithLabel("File")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileNewSg, err = gtk.MenuItemNewWithLabel("New Signal Graph")
	if err != nil {
		log.Fatal("Unable to create fileNewSg:", err)
	}
	m.fileNewLib, err = gtk.MenuItemNewWithLabel("New Library")
	if err != nil {
		log.Fatal("Unable to create fileNewLib:", err)
	}
	m.fileNewPlat, err = gtk.MenuItemNewWithLabel("New Platform")
	if err != nil {
		log.Fatal("Unable to create fileNewPlat:", err)
	}
	m.fileNewMap, err = gtk.MenuItemNewWithLabel("New Mapping")
	if err != nil {
		log.Fatal("Unable to create fileNewPlat:", err)
	}
	m.fileOpen, err = gtk.MenuItemNewWithLabel("Open")
	if err != nil {
		log.Fatal("Unable to create fileOpen:", err)
	}
	m.fileSave, err = gtk.MenuItemNewWithLabel("Save")
	if err != nil {
		log.Fatal("Unable to create fileSave:", err)
	}
	m.fileSaveAs, err = gtk.MenuItemNewWithLabel("Save As")
	if err != nil {
		log.Fatal("Unable to create fileSaveAs:", err)
	}
	m.fileClose, err = gtk.MenuItemNewWithLabel("Close")
	if err != nil {
		log.Fatal("Unable to create fileClose:", err)
	}
	m.fileQuit, err = gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		log.Fatal("Unable to create fileQuit:", err)
	}
	m.menuFile.Append(m.fileNewSg)
	m.menuFile.Append(m.fileNewLib)
	m.menuFile.Append(m.fileNewPlat)
	m.menuFile.Append(m.fileNewMap)
	x, _ := gtk.SeparatorMenuItemNew()
	m.menuFile.Append(x)
	m.menuFile.Append(m.fileOpen)
	m.menuFile.Append(m.fileSave)
	m.menuFile.Append(m.fileSaveAs)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuFile.Append(x)
	m.menuFile.Append(m.fileClose)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuFile.Append(x)
	m.menuFile.Append(m.fileQuit)
	m.filemenu.SetSubmenu(m.menuFile)
	m.menubar.Append(m.filemenu)

	m.fileQuit.Connect("activate", func() {
		gtk.MainQuit()
	})

	m.menuEdit, err = gtk.MenuNew()
	if err != nil {
		log.Fatal("Unable to create menuEdit:", err)
	}
	m.editmenu, err = gtk.MenuItemNewWithMnemonic("_Edit")
	if err != nil {
		log.Fatal("Unable to create editmenu:", err)
	}
	m.editUndo, err = gtk.MenuItemNewWithMnemonic("_Undo")
	if err != nil {
		log.Fatal("Unable to create editUndo:", err)
	}
	m.editRedo, err = gtk.MenuItemNewWithMnemonic("_Redo")
	if err != nil {
		log.Fatal("Unable to create editRedo:", err)
	}
	m.editNew, err = gtk.MenuItemNewWithMnemonic("_New Element")
	if err != nil {
		log.Fatal("Unable to create editNew:", err)
	}
	m.editEdit, err = gtk.MenuItemNewWithMnemonic("_Edit")
	if err != nil {
		log.Fatal("Unable to create editEdit:", err)
	}
	m.editCopy, err = gtk.MenuItemNewWithMnemonic("_Copy")
	if err != nil {
		log.Fatal("Unable to create editCopy:", err)
	}
	m.editDelete, err = gtk.MenuItemNewWithMnemonic("_Delete")
	if err != nil {
		log.Fatal("Unable to create editDelete:", err)
	}
	m.editPaste, err = gtk.MenuItemNewWithMnemonic("_Paste")
	if err != nil {
		log.Fatal("Unable to create editPaste:", err)
	}
	m.menuEdit.Append(m.editUndo)
	m.menuEdit.Append(m.editRedo)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuEdit.Append(x)
	m.menuEdit.Append(m.editNew)
	m.menuEdit.Append(m.editEdit)
	m.menuEdit.Append(m.editDelete)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuEdit.Append(x)
	m.menuEdit.Append(m.editCopy)
	m.menuEdit.Append(m.editPaste)
	m.editmenu.SetSubmenu(m.menuEdit)
	m.menubar.Append(m.editmenu)

	m.menuView, err = gtk.MenuNew()
	if err != nil {
		log.Fatal("Unable to create menuView:", err)
	}
	m.viewmenu, err = gtk.MenuItemNewWithMnemonic("_View")
	if err != nil {
		log.Fatal("Unable to create viewmenu:", err)
	}
	m.viewExpand, err = gtk.MenuItemNewWithMnemonic("E_xpand")
	if err != nil {
		log.Fatal("Unable to create viewExpand:", err)
	}
	m.viewCollapse, err = gtk.MenuItemNewWithMnemonic("Co_llapse")
	if err != nil {
		log.Fatal("Unable to create viewCollapse:", err)
	}
	m.menuView.Append(m.viewExpand)
	m.menuView.Append(m.viewCollapse)
	m.viewmenu.SetSubmenu(m.menuView)
	m.menubar.Append(m.viewmenu)

	m.menuAbout, err = gtk.MenuNew()
	if err != nil {
		log.Fatal("Unable to create menuEdit:", err)
	}
	m.aboutmenu, err = gtk.MenuItemNewWithLabel("About")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.aboutAbout, err = gtk.MenuItemNewWithLabel("About...")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.aboutHelp, err = gtk.MenuItemNewWithLabel("Help")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.menuAbout.Append(m.aboutAbout)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuAbout.Append(x)
	m.menuAbout.Append(m.aboutHelp)
	m.aboutmenu.SetSubmenu(m.menuAbout)
	m.menubar.Append(m.aboutmenu)

	log.Println("GoAppMenu initialized")
}
