package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/interface/graph"
	"github.com/axel-freesp/sge/interface/model"
	"github.com/axel-freesp/sge/interface/tree"
)

/*
 *  Behaviour Model
 */

type SignalGraphIf interface {
	tree.ToplevelTreeElement
	ItsType() SignalGraphTypeIf
	Nodes() []NodeIf
}

type SignalGraphTypeIf interface {
	tree.TreeElement
	Libraries() []LibraryIf
	Nodes() []NodeIf
	NodeByName(string) (NodeIf, bool)
	InputNodes() []NodeIf
	OutputNodes() []NodeIf
	ProcessingNodes() []NodeIf
	AddNode(NodeIf) error
	RemoveNode(n NodeIf)
	Context() model.ModelContextIf
}

type LibraryIf interface {
	tree.ToplevelTreeElement
	SignalTypes() []SignalTypeIf
	NodeTypes() []NodeTypeIf
	AddNodeType(NodeTypeIf) error
	RemoveNodeType(t NodeTypeIf)
	AddSignalType(SignalTypeIf) bool
	RemoveSignalType(SignalTypeIf)
}

type NodeTypeIf interface {
	tree.TreeElement
	TypeName() string
	SetTypeName(string)
	DefinedAt() string
	InPorts() []PortTypeIf
	OutPorts() []PortTypeIf
	Implementation() []ImplementationIf
	Instances() []NodeIf
	AddNamedPortType(PortTypeIf)
	RemoveNamedPortType(PortTypeIf)
	AddImplementation(ImplementationIf)
	RemoveImplementation(ImplementationIf)
}

type ImplementationIf interface {
	tree.TreeElement
	ImplementationType() ImplementationType
	ElementName() string
	SetElemName(string)
	Graph() SignalGraphTypeIf
}

type ImplementationType int

const (
	NodeTypeElement ImplementationType = 0
	NodeTypeGraph   ImplementationType = 1
)

type NodeIf interface {
	tree.NamedTreeElementIf
	graph.PathModePositioner
	ItsType() NodeTypeIf
	InPorts() []PortIf
	OutPorts() []PortIf
	InPortIndex(portname string) int
	OutPortIndex(portname string) int
	Context() SignalGraphTypeIf
	PortLink() (string, bool)
	Expanded() bool
	SetExpanded(bool)
	SubNode(ownId, childId NodeIdIf) (NodeIf, bool)
	Children(ownId NodeIdIf) []NodeIdIf
}

type NodeIdIf interface {
	fmt.Stringer
	Parent() NodeIdIf
	IsAncestor(NodeIdIf) bool
}

type SignalTypeIf interface {
	tree.TreeElement
	DefinedAt() string
	TypeName() string
	SetTypeName(string)
	CType() string
	SetCType(string)
	ChannelId() string
	SetChannelId(string)
	Scope() Scope
	SetScope(Scope)
	Mode() Mode
	SetMode(Mode)
}

type Scope int
type Mode int

const (
	Local Scope = iota
	Global
)

const (
	Synchronous Mode = iota
	Asynchronous
)

type PortTypeIf interface {
	tree.NamedTreeElementIf
	graph.Directioner
	graph.ModePositioner
	SignalType() SignalTypeIf
	SetSignalType(SignalTypeIf)
}

type PortIf interface {
	tree.TreeElement
	graph.Directioner
	graph.ModePositioner
	Name() string
	SignalType() SignalTypeIf
	Connections() []PortIf
	Node() NodeIf
	Connection(PortIf) ConnectionIf
	AddConnection(ConnectionIf) error
	RemoveConnection(c PortIf)
}

type ConnectionIf interface {
	tree.TreeElement
	From() PortIf
	To() PortIf
}
