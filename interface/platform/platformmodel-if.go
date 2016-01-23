package platform

import (
	interfaces "github.com/axel-freesp/sge/interface"
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
	PlatformObject() interfaces.PlatformObject
}

type ArchIf interface {
	tree.NamedTreeElementIf
	interfaces.ModePositioner
	Platform() PlatformIf
	IOTypes() []IOTypeIf
	Processes() []ProcessIf
}

type IOTypeIf interface {
	tree.NamedTreeElementIf
	interfaces.IOModer
	Platform() PlatformIf
}

type ProcessIf interface {
	tree.NamedTreeElementIf
	interfaces.ModePositioner
	Arch() ArchIf
	InChannels() []ChannelIf
	OutChannels() []ChannelIf
}

type ChannelIf interface {
	tree.NamedTreeElementIf
	interfaces.Directioner
	interfaces.ModePositioner
	Process() ProcessIf
	IOType() IOTypeIf
	SetIOType(IOTypeIf)
	Link() ChannelIf
}
