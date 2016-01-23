package behaviour

import (
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/interface/model"
)

/*
 *  Behaviour Model
 */

type SignalGraphIf interface {
	tree.ToplevelTreeElement
	ItsType() SignalGraphTypeIf
	GraphObject() interfaces.GraphObject
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
	SignalTypes() []SignalType
	NodeTypes() []NodeTypeIf
	AddNodeType(NodeTypeIf) error
	RemoveNodeType(t NodeTypeIf)
	AddSignalType(SignalType) bool
	RemoveSignalType(t SignalType)
}

type NodeTypeIf interface {
	tree.TreeElement
	TypeName() string
	SetTypeName(string)
	DefinedAt() string
	InPorts() []PortType
	OutPorts() []PortType
	Implementation() []ImplementationIf
	Instances() []NodeIf
	AddNamedPortType(PortType)
	RemoveNamedPortType(PortType)
	AddImplementation(ImplementationIf)
	RemoveImplementation(ImplementationIf)
}

type ImplementationIf interface {
	tree.TreeElement
	ImplementationType() ImplementationType
	ElementName() string
	SetElemName(string)
	Graph() SignalGraphTypeIf
	GraphObject() interfaces.GraphObject
}

type ImplementationType int

const (
	NodeTypeElement ImplementationType = 0
	NodeTypeGraph   ImplementationType = 1
)

type NodeIf interface {
	tree.NamedTreeElementIf
	interfaces.Positioner
	interfaces.Porter
	ItsType() NodeTypeIf
	InPorts() []Port
	OutPorts() []Port
	Context() SignalGraphTypeIf
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

type SignalType interface {
	tree.TreeElement
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

type PortType interface {
	tree.NamedTreeElementIf
	interfaces.Directioner
	SignalType() SignalType
	SetSignalType(SignalType)
}

type Port interface {
	tree.TreeElement
	interfaces.Directioner
	Name() string
	SignalType() SignalType
	Connections() []Port
	Node() NodeIf
	Connection(Port) Connection
	AddConnection(Connection) error
	RemoveConnection(c Port)
}

type Connection interface {
	tree.TreeElement
	From() Port
	To() Port
}
