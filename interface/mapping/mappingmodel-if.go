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
	tree.ToplevelTreeElement
	SetGraph(behaviour.SignalGraphIf)
	Graph() behaviour.SignalGraphIf
	SetPlatform(platform.PlatformIf)
	Platform() platform.PlatformIf
	AddMapping(n behaviour.NodeIf, nId behaviour.NodeIdIf, p platform.ProcessIf)
	Mapped(string) (platform.ProcessIf, bool)
	MappedElement(behaviour.NodeIdIf) (MappedElementIf, bool)
	MappedIds() []behaviour.NodeIdIf
}

type MappedElementIf interface {
	tree.TreeElement
	graph.PathModePositioner
	graph.Expander
	Mapping() MappingIf
	NodeId() behaviour.NodeIdIf
	Process() (platform.ProcessIf, bool)
	SetProcess(platform.ProcessIf)
}
