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
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
	"github.com/axel-freesp/sge/backend"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"fmt"
)

const (
	width  = 800
	height = 600
)

var (
	win *GoAppWindow
	xmlview *views.XmlTextView
)

func treeSelectionChangedCB(selection *gtk.TreeSelection, treeStore *models.FilesTreeStore) {
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

func main() {
	// Initialize GTK with parsing the command line arguments.
	unhandledArgs := os.Args
	gtk.Init(&unhandledArgs)
	backend.Init()
	freesp.SignalGraphInit()

	// Create a new toplevel window.
	win, err := GoAppWindowNew(width, height)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

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

	selection, err := ftv.TreeView().GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	selection.Connect("changed", treeSelectionChangedCB, fts)

	// Handle command line arguments: treat each as a filename:
	for i, p := range unhandledArgs {
		if i > 0 {
			var sg freesp.SignalGraph
			sg = freesp.SignalGraphNew()

			err := sg.ReadFile(fmt.Sprintf("%s/%s", backend.XmlRoot(), p))
			if err != nil {
				log.Println("WARNING: sg.ReadFile", p, "failed")
				continue
			}

			fts.AddBehaviourFile(p, sg)

		}
	}

	xmlview, err = views.XmlTextViewNew(width, height)
	if err != nil {
		log.Fatal("Could not create XML view.");
	}
	win.stack.AddTitled(xmlview.Widget(), "XML View", "XML View")

	win.Window().ShowAll()
	gtk.Main()
}
