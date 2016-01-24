package views

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/gotk3/gotk3/gtk"
)

type Context interface {
	SelectNode(bh.NodeIf)          // single click selection
	EditNode(bh.NodeIf)            // double click selection
	SelectPort(bh.PortIf)          // single click selection
	SelectConnect(bh.ConnectionIf) // single click selection
	SelectArch(pf.ArchIf)
	SelectProcess(pf.ProcessIf)
	SelectChannel(pf.ChannelIf)
	SelectMapElement(mp.MappedElementIf)
}

type GraphViewIf interface {
	Widget() *gtk.Widget
	Sync()
	Select(obj interface{})
	Expand(obj interface{})
	Collapse(obj interface{})
	IdentifyGraph(bh.SignalGraphIf) bool
	IdentifyPlatform(pf.PlatformIf) bool
	IdentifyMapping(mp.MappingIf) bool
}

type XmlTextViewIf interface {
	Set(gr.XmlCreator) error
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
}
