package views

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/gotk3/gotk3/gtk"
)

type ContextIf interface {
	// node is toplevel node:
	SelectNode(bh.NodeIf, bh.NodeIdIf)            // single click selection
	EditNode(bh.NodeIf, bh.NodeIdIf)              // double click selection
	SelectPort(bh.PortIf, bh.NodeIf, bh.NodeIdIf) // single click selection
	SelectConnect(bh.ConnectionIf)                // single click selection
	SelectArch(pf.ArchIf)
	SelectProcess(pf.ProcessIf)
	SelectChannel(pf.ChannelIf)
	SelectMapElement(mp.MappedElementIf)
}

type GraphViewCollectionIf interface {
	Add(gv GraphViewIf, title string)
	RemoveGraphView(bh.SignalGraphIf)
	RemovePlatformView(pf.PlatformIf)
	RemoveMappingView(mp.MappingIf)
	Rename(old, new string)
	Widget() *gtk.Widget
	XmlTextView() XmlTextViewIf
	Sync()
	Select(obj interface{})
	Select2(obj interface{}, id string)
	CurrentView() GraphViewIf
}

type GraphViewIf interface {
	Widget() *gtk.Widget
	Sync()
	Select(obj interface{})
	Select2(obj interface{}, id string)
	Expand(obj interface{})
	Collapse(obj interface{})
	IdentifyGraph(bh.SignalGraphIf) bool
	IdentifyPlatform(pf.PlatformIf) bool
	IdentifyMapping(mp.MappingIf) bool
}

type XmlTextViewIf interface {
	Set(gr.XmlCreator) error
}
