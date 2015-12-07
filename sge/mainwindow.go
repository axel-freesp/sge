/*
 * Copyright (c) 2013-2014 Conformal Systems <info@conformal.com>
 *
 * This file originated from: http://opensource.conformal.com/
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"strings"
)

const (
	width  = 800
	height = 600
)

var (
	win *GoAppWindow
)

type selectionArg struct {
	treeStore *models.FilesTreeStore
	xmlview   *views.XmlTextView
}

func treeSelectionChangedCB(selection *gtk.TreeSelection, arg *selectionArg) {
	treeStore := arg.treeStore
	xmlview := arg.xmlview
	var iter gtk.TreeIter
	var model gtk.ITreeModel
	if selection.GetSelected(&model, &iter) {
		obj, err := treeStore.GetObject(&iter)
		if err != nil {
			log.Fatal("treeSelectionChangedCB: Could not get object from model", err)
		}
		xmlview.Set(obj)
	}
}

func fileNewSg(fts *models.FilesTreeStore) {
	log.Println("fileNewSg")
	err := fts.AddSignalGraphFile("new-file.sml", freesp.SignalGraphNew("new-file.sml"))
	if err != nil {
		log.Println("Warning: ftv.AddSignalGraphFile('new-file.sml') failed.")
	}
}

func fileNewLib(fts *models.FilesTreeStore) {
	log.Println("fileNewLib")
	err := fts.AddLibraryFile("new-file.alml", freesp.LibraryNew("new-file.alml"))
	if err != nil {
		log.Println("Warning: ftv.AddLibraryFile('new-file.alml') failed.")
	}
}

func fileOpen(fts *models.FilesTreeStore) {
	log.Println("fileOpen")
}

func fileSaveAs(fts *models.FilesTreeStore) {
	log.Println("fileSaveAs")
}

func fileSave(fts *models.FilesTreeStore) {
	log.Println("fileSave")
	var path0 string
	if fts.CurrentSelection == nil {
		log.Fatal("fileSave error: CurrentSelection = nil")
	}
	path, err := fts.TreeStore().GetPath(fts.CurrentSelection)
	if err != nil {
		log.Fatal("fileSave error: iter.GetPath failed:", err)
		return
	}
	p := path.String()
	log.Println("Current selection: ", p)
	spl := strings.Split(p, ":")
	path0 = spl[0]
	iter, err := fts.TreeStore().GetIterFromString(path0)
	if err != nil {
		log.Fatal("fileSave error: fts.TreeStore().GetIterFromString failed:", err)
	}
	filename, err := fts.GetValue(iter)
	if err != nil {
		log.Fatal("fileSave error: fts.GetValue failed:", err)
	}
	log.Println("fileSave: filename =", filename)
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
	arg := &selectionArg{fts, xmlview}
	selection.Connect("changed", treeSelectionChangedCB, arg)

	if len(unhandledArgs) < 2 {
		err := fts.AddSignalGraphFile("new-file.sml", freesp.SignalGraphNew("new-file.sml"))
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
				continue
			}
			var lib freesp.Library
			lib = freesp.LibraryNew(filepath)
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

	menu.fileNewSg.Connect("activate", func() { fileNewSg(fts) })
	menu.fileNewLib.Connect("activate", func() { fileNewLib(fts) })
	menu.fileOpen.Connect("activate", func() { fileOpen(fts) })
	menu.fileSave.Connect("activate", func() { fileSave(fts) })
	menu.fileSaveAs.Connect("activate", func() { fileSaveAs(fts) })

	win.Window().ShowAll()
	gtk.Main()
}
