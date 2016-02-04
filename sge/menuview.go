package main

import (
	"log"
	//gr "github.com/axel-freesp/sge/interface/graph"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	//	mp "github.com/axel-freesp/sge/interface/mapping"
	tr "github.com/axel-freesp/sge/interface/tree"
	//"github.com/axel-freesp/sge/models"
	//"github.com/axel-freesp/sge/views"
	//"github.com/gotk3/gotk3/gtk"
)

func MenuViewInit(menu *GoAppMenu, g *Global) {
	menu.viewExpand.Connect("activate", func() { viewExpand(menu, g) })
	menu.viewCollapse.Connect("activate", func() { viewCollapse(menu, g) })
	menu.viewExpand.SetSensitive(false)
	menu.viewCollapse.SetSensitive(false)
}

func MenuViewPost(menu *GoAppMenu, g *Global) {
	MenuViewCurrent(menu, g)
	return

	fts := g.FTS().(tr.TreeIf)
	cursor := fts.Current()
	var obj tr.TreeElement
	if len(cursor.Path) != 0 {
		obj = fts.Object(cursor)
	}
	global.GVC().XmlTextView().Set(obj)
	global.GVC().Sync()
}

func MenuViewCurrent(menu *GoAppMenu, g *Global) {
	fts := g.FTS().(tr.TreeIf)
	cursor := fts.Current()
	menu.viewExpand.SetSensitive(false)
	menu.viewCollapse.SetSensitive(false)
	if len(cursor.Path) == 0 {
		return
	}
	obj := fts.Object(cursor)
	var n bh.NodeIf
	switch obj.(type) {
	case bh.NodeIf:
		n = obj.(bh.NodeIf)
	//case mp.MappedElementIf:
	//	n = obj.(mp.MappedElementIf).Node()
	default:
		return
	}
	impl := n.ItsType().Implementation()
	for _, i := range impl {
		if i.ImplementationType() == bh.NodeTypeGraph {
			menu.viewExpand.SetSensitive(true)
			menu.viewCollapse.SetSensitive(true)
			break
		}
	}
}

func viewExpand(menu *GoAppMenu, g *Global) {
	log.Printf("viewExpand\n")
	defer MenuViewPost(menu, g)
	fts := g.FTS().(tr.TreeIf)
	cursor := fts.Current()
	obj := fts.Object(cursor)
	v := g.GVC().CurrentView()
	v.Expand(obj)
}

func viewCollapse(menu *GoAppMenu, g *Global) {
	log.Printf("viewCollapse\n")
	defer MenuViewPost(menu, g)
	fts := g.FTS().(tr.TreeIf)
	cursor := fts.Current()
	obj := fts.Object(cursor)
	v := g.GVC().CurrentView()
	v.Collapse(obj)
}
