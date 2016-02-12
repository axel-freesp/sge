package mapping

import (
	"github.com/axel-freesp/sge/interface/behaviour"
	"github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/interface/platform"
	"github.com/axel-freesp/sge/interface/tree"
)

/*
 *      Mapping Model Interface
 */

type MappingIf interface {
	tree.ToplevelTreeElementIf
	SetGraph(behaviour.SignalGraphIf)
	Graph() behaviour.SignalGraphIf
	SetPlatform(platform.PlatformIf)
	Platform() platform.PlatformIf
	AddMapping(n behaviour.NodeIf, nId behaviour.NodeIdIf, p platform.ProcessIf) MappedElementIf
	Mapped(string) (platform.ProcessIf, bool)
	MappedElement(behaviour.NodeIdIf) (MappedElementIf, bool)
	MappedIds() []behaviour.NodeIdIf
}

type MappedElementIf interface {
	tree.TreeElementIf
	graph.ModePositioner
	graph.Expander
	Mapping() MappingIf
	NodeId() behaviour.NodeIdIf
	Node() behaviour.NodeIf
	Process() (platform.ProcessIf, bool)
	SetProcess(platform.ProcessIf)
}
