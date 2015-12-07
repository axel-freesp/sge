package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	//"github.com/gotk3/gotk3/glib"
)

type GoAppMenu struct {
	menubar    *gtk.MenuBar
	menuFile   *gtk.Menu
	filemenu   *gtk.MenuItem
	fileNewSg  *gtk.MenuItem
	fileNewLib *gtk.MenuItem
	fileOpen   *gtk.MenuItem
	fileSave   *gtk.MenuItem
	fileSaveAs *gtk.MenuItem
	fileQuit   *gtk.MenuItem
	menuEdit   *gtk.Menu
	editmenu   *gtk.MenuItem
	editUndo   *gtk.MenuItem
	editRedo   *gtk.MenuItem
	editNew    *gtk.MenuItem
	editCopy   *gtk.MenuItem
	editDelete *gtk.MenuItem
	menuAbout  *gtk.Menu
	aboutmenu  *gtk.MenuItem
	aboutAbout *gtk.MenuItem
	aboutHelp  *gtk.MenuItem
}

func GoAppMenuNew() *GoAppMenu {
	return &GoAppMenu{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
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
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileNewLib, err = gtk.MenuItemNewWithLabel("New Library")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileOpen, err = gtk.MenuItemNewWithLabel("Open")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileSave, err = gtk.MenuItemNewWithLabel("Save")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileSaveAs, err = gtk.MenuItemNewWithLabel("Save As")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.fileQuit, err = gtk.MenuItemNewWithLabel("Quit")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.menuFile.Append(m.fileNewSg)
	m.menuFile.Append(m.fileNewLib)
	x, _ := gtk.SeparatorMenuItemNew()
	m.menuFile.Append(x)
	m.menuFile.Append(m.fileOpen)
	m.menuFile.Append(m.fileSave)
	m.menuFile.Append(m.fileSaveAs)
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
	m.editmenu, err = gtk.MenuItemNewWithLabel("Edit")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.editUndo, err = gtk.MenuItemNewWithLabel("Undo")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.editRedo, err = gtk.MenuItemNewWithLabel("Redo")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.editNew, err = gtk.MenuItemNewWithLabel("New")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.editCopy, err = gtk.MenuItemNewWithLabel("Copy")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.editDelete, err = gtk.MenuItemNewWithLabel("Delete")
	if err != nil {
		log.Fatal("Unable to create filemenu:", err)
	}
	m.menuEdit.Append(m.editUndo)
	m.menuEdit.Append(m.editRedo)
	x, _ = gtk.SeparatorMenuItemNew()
	m.menuEdit.Append(x)
	m.menuEdit.Append(m.editNew)
	m.menuEdit.Append(m.editCopy)
	m.menuEdit.Append(m.editDelete)
	m.editmenu.SetSubmenu(m.menuEdit)
	m.menubar.Append(m.editmenu)

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
