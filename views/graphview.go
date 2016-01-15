package views

import (
	"fmt"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type Context interface {
	SelectNode(node interfaces.NodeObject)          // single click selection
	EditNode(node interfaces.NodeObject)            // double click selection
	SelectPort(port interfaces.PortObject)          // single click selection
	SelectConnect(conn interfaces.ConnectionObject) // single click selection
	SelectArch(arch interfaces.ArchObject)
	SelectProcess(arch interfaces.ProcessObject)
	SelectChannel(arch interfaces.ChannelObject)
}

type GraphView interface {
	Widget() *gtk.Widget
	Sync()
	Select(obj interfaces.GraphElement)
	IdentifyGraph(interfaces.GraphObject) bool
	IdentifyPlatform(interfaces.PlatformObject) bool
}

type GraphViewCollection interface {
	Add(gv GraphView, title string)
	RemoveGraphView(g interfaces.GraphObject)
	RemovePlatformView(p interfaces.PlatformObject)
	Rename(old, new string)
	Widget() *gtk.Widget
	XmlTextView() *XmlTextView
	Sync()
	Select(obj interfaces.GraphElement)
}

type graphViewCollection struct {
	graphview []GraphView
	xmlview   *XmlTextView
	box       *gtk.Box
	header    *gtk.HeaderBar
	tabs      *gtk.StackSwitcher
	stack     *gtk.Stack
}

var _ GraphViewCollection = (*graphViewCollection)(nil)

func GraphViewCollectionNew(width, height int) (gvc *graphViewCollection, err error) {
	gvc = &graphViewCollection{}
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

func (gvc graphViewCollection) XmlTextView() *XmlTextView {
	return gvc.xmlview
}

func (gvc *graphViewCollection) Add(gv GraphView, title string) {
	gvc.graphview = append(gvc.graphview, gv)
	gvc.stack.AddTitled(gv.Widget(), title, title)
	gv.Widget().ShowAll()
}

func (gvc *graphViewCollection) Rename(old, new string) {
	widget, err := gvc.stack.GetChildByName(old)
	if err != nil {
		log.Printf("graphViewCollection.Rename warning: stack child %s not found\n", old)
		return
	}
	if widget == gvc.stack.GetVisibleChild() {
		gvc.stack.SetVisibleChild(gvc.xmlview.Widget())
		gvc.stack.Remove(widget)
		gvc.stack.AddTitled(widget, new, new)
		gvc.stack.SetVisibleChildName(new)
	} else {
		gvc.stack.Remove(widget)
		gvc.stack.AddTitled(widget, new, new)
	}
}

func (gvc *graphViewCollection) RemoveGraphView(g interfaces.GraphObject) {
	var tmp []GraphView
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

func (gvc *graphViewCollection) RemovePlatformView(p interfaces.PlatformObject) {
	var tmp []GraphView
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

func (gvc *graphViewCollection) Widget() *gtk.Widget {
	return &gvc.box.Widget
}

func (gvc *graphViewCollection) Sync() {
	for _, v := range gvc.graphview {
		v.Sync()
	}
}

func (gvc *graphViewCollection) Select(obj interfaces.GraphElement) {
	for _, v := range gvc.graphview {
		v.Select(obj)
	}
}
