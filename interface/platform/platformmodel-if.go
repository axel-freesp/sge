package platform

import (
	"github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/interface/tree"
)

/*
 *  Platform Model Interface
 */

type PlatformIf interface {
	tree.ToplevelTreeElement
	PlatformId() string
	SetPlatformId(string)
	ProcessByName(string) (ProcessIf, bool)
	Arch() []ArchIf
}

type ArchIf interface {
	tree.NamedTreeElementIf
	graph.ModePositioner
	Platform() PlatformIf
	IOTypes() []IOTypeIf
	Processes() []ProcessIf
	ArchPorts() []ArchPortIf
}

type ArchPortIf interface {
	graph.ModePositioner
	Channel() ChannelIf
}

type IOTypeIf interface {
	tree.NamedTreeElementIf
	graph.IOModer
	Platform() PlatformIf
}

type ProcessIf interface {
	tree.NamedTreeElementIf
	graph.ModePositioner
	Arch() ArchIf
	InChannels() []ChannelIf
	OutChannels() []ChannelIf
}

type ChannelIf interface {
	tree.NamedTreeElementIf
	graph.Directioner
	graph.ModePositioner
	Process() ProcessIf
	IOType() IOTypeIf
	SetIOType(IOTypeIf)
	Link() ChannelIf
	ArchPort() ArchPortIf
}
