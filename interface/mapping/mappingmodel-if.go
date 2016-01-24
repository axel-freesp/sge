package mapping

import (
	"github.com/axel-freesp/sge/interface/behaviour"
	"github.com/axel-freesp/sge/interface/platform"
	"github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/interface/graph"
)

/*
 *      Mapping Model Interface
 */

type MappingIf interface {
	tree.ToplevelTreeElement
	AddMapping(n behaviour.NodeIf, p platform.ProcessIf)
	SetGraph(behaviour.SignalGraphIf)
	Graph() behaviour.SignalGraphIf
	SetPlatform(platform.PlatformIf)
	Platform() platform.PlatformIf
	Mapped(behaviour.NodeIf) (platform.ProcessIf, bool)
	MappedElement(behaviour.NodeIf) (MappedElementIf, bool)
}

type MappedElementIf interface {
	tree.TreeElement
	graph.Positioner
	Mapping() MappingIf
	Node() behaviour.NodeIf
	Process() (platform.ProcessIf, bool)
	SetProcess(platform.ProcessIf)
}

