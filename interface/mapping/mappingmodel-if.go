package mapping

import (
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/interface/behaviour"
	"github.com/axel-freesp/sge/interface/platform"
	"github.com/axel-freesp/sge/interface/tree"
)

/*
 *      Mapping Model Interface
 */

type MappingIf interface {
	tree.ToplevelTreeElement
	MappingObject() interfaces.MappingObject
	AddMapping(n behaviour.NodeIf, p platform.ProcessIf)
	SetGraph(behaviour.SignalGraphIf)
	Graph() behaviour.SignalGraphIf
	SetPlatform(platform.PlatformIf)
	Platform() platform.PlatformIf
	Mapped(behaviour.NodeIf) (platform.ProcessIf, bool)
	MappedElement(behaviour.NodeIf) MappedElementIf
}

type MappedElementIf interface {
	tree.TreeElement
	interfaces.Positioner
	Mapping() MappingIf
	Node() behaviour.NodeIf
	Process() platform.ProcessIf
	SetProcess(platform.ProcessIf)
}
