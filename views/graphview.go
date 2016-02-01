package views

import (
	"fmt"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type graphViewCollection struct {
	graphview []GraphViewIf
	gmap      map[string]GraphViewIf
	xmlview   *XmlTextView
	box       *gtk.Box
	header    *gtk.HeaderBar
	tabs      *gtk.StackSwitcher
	stack     *gtk.Stack
}

var _ GraphViewCollectionIf = (*graphViewCollection)(nil)

func GraphViewCollectionNew(width, height int) (gvc *graphViewCollection, err error) {
	gvc = &graphViewCollection{}
	gvc.gmap = make(map[string]GraphViewIf)
	gvc.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		err = fmt.Errorf("GraphViewCollectionNew: Unable to create box:", err)
		return
	}
	gvc.header, err = gtk.HeaderBarNew()
	if err != nil {
		err = fmt.Errorf("GraphViewCollectionNew: Unable to create bar:", err)
		return
	}
	gvc.tabs, err = gtk.StackSwitcherNew()
	if err != nil {
		err = fmt.Errorf("GraphViewCollectionNew: Unable to create stackswitcher:", err)
		return
	}
	gvc.stack, err = gtk.StackNew()
	if err != nil {
		err = fmt.Errorf("GraphViewCollectionNew: Unable to create Stack:", err)
		return
	}
	gvc.box.PackStart(gvc.header, false, true, 0)
	gvc.header.Add(gvc.tabs)
	gvc.tabs.SetStack(gvc.stack)
	gvc.box.Add(gvc.stack)

	gvc.xmlview, err = XmlTextViewNew(width, height)
	if err != nil {
		err = fmt.Errorf("GraphViewCollectionNew: Could not create XML view.")
		return
	}
	gvc.stack.AddTitled(gvc.xmlview.Widget(), "XML View", "XML View")
	return
}

func (gvc graphViewCollection) XmlTextView() XmlTextViewIf {
	return gvc.xmlview
}

func (gvc *graphViewCollection) Add(gv GraphViewIf, title string) {
	gvc.gmap[title] = gv
	gvc.graphview = append(gvc.graphview, gv)
	gvc.stack.AddTitled(gv.Widget(), title, title)
	gv.Widget().ShowAll()
}

func (gvc *graphViewCollection) Rename(old, new string) {
	v, ok := gvc.gmap[old]
	if !ok {
		log.Printf("graphViewCollection.Rename warning: old %s not found\n", old)
		return
	}
	widget, err := gvc.stack.GetChildByName(old)
	if err != nil {
		log.Printf("graphViewCollection.Rename warning: stack child %s not found\n", old)
		return
	}
	// TODO ~~~~~~~~~~~~~~~~~~~~~~~~~~~
	if old == gvc.stack.GetVisibleChildName() {
		gvc.stack.SetVisibleChild(gvc.xmlview.Widget())
		gvc.stack.Remove(widget)
		gvc.stack.AddTitled(widget, new, new)
		gvc.stack.SetVisibleChildName(new)
	} else {
		gvc.stack.Remove(widget)
		gvc.stack.AddTitled(widget, new, new)
	}
	delete(gvc.gmap, old)
	gvc.gmap[new] = v
}

func (gvc *graphViewCollection) RemoveGraphView(g bh.SignalGraphIf) {
	var tmp []GraphViewIf
	for _, v := range gvc.graphview {
		if v.IdentifyGraph(g) {
			gvc.stack.SetVisibleChild(gvc.xmlview.Widget())
			gvc.stack.Remove(v.Widget())
			//views.SignalGraphViewDestroy(v)
		} else {
			tmp = append(tmp, v)
		}
	}
	gvc.graphview = tmp
}

func (gvc *graphViewCollection) RemovePlatformView(p pf.PlatformIf) {
	var tmp []GraphViewIf
	for _, v := range gvc.graphview {
		if v.IdentifyPlatform(p) {
			gvc.stack.SetVisibleChild(gvc.xmlview.Widget())
			gvc.stack.Remove(v.Widget())
			//views.SignalGraphViewDestroy(v)
		} else {
			tmp = append(tmp, v)
		}
	}
	gvc.graphview = tmp
}

func (gvc *graphViewCollection) RemoveMappingView(m mp.MappingIf) {
	var tmp []GraphViewIf
	for _, v := range gvc.graphview {
		if v.IdentifyMapping(m) {
			gvc.stack.SetVisibleChild(gvc.xmlview.Widget())
			gvc.stack.Remove(v.Widget())
			//views.MappingViewDestroy(v)
		} else {
			tmp = append(tmp, v)
		}
	}
	gvc.graphview = tmp
}

func (gvc *graphViewCollection) Widget() *gtk.Widget {
	return &gvc.box.Widget
}

func (gvc *graphViewCollection) Sync() {
	for _, v := range gvc.graphview {
		v.Sync()
	}
}

func (gvc *graphViewCollection) Select(obj interface{}) {
	for _, v := range gvc.graphview {
		v.Select(obj)
	}
}

func (gvc *graphViewCollection) Select2(obj interface{}, id string) {
	for _, v := range gvc.graphview {
		v.Select2(obj, id)
	}
}

func (gvc *graphViewCollection) CurrentView() (v GraphViewIf) {
	name := gvc.stack.GetVisibleChildName()
	var ok bool
	v, ok = gvc.gmap[name]
	if !ok {
		log.Fatalf("FIXME: graphViewCollection.CurrentView: no visible child\n")
	}
	return
}
