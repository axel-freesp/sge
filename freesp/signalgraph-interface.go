package freesp

import (
	interfaces "github.com/axel-freesp/sge/interface"
)

type ContextIf interface {
	SignalGraphMgr() FileManagerIf
	LibraryMgr() FileManagerIf
	PlatformMgr() FileManagerIf
	MappingMgr() FileManagerMappingIf
}

type FileManagerIf interface {
	New() (ToplevelTreeElement, error)
	Access(name string) (ToplevelTreeElement, error)
	Remove(name string)
	Rename(oldName, newName string) error
}

type FileManagerMappingIf interface {
	FileManagerIf
	SetGraphForNew(g SignalGraphIf)
	SetPlatformForNew(p PlatformIf)
}

/*
 *  Behaviour Model
 */

type SignalGraphIf interface {
	ToplevelTreeElement
	ItsType() SignalGraphTypeIf
	GraphObject() interfaces.GraphObject
}

type SignalGraphTypeIf interface {
	TreeElement
	Libraries() []LibraryIf
	Nodes() []NodeIf
	NodeByName(string) (NodeIf, bool)
	InputNodes() []NodeIf
	OutputNodes() []NodeIf
	ProcessingNodes() []NodeIf
	AddNode(NodeIf) error
	RemoveNode(n NodeIf)
	Context() ContextIf
}

type LibraryIf interface {
	ToplevelTreeElement
	SignalTypes() []SignalType
	NodeTypes() []NodeTypeIf
	AddNodeType(NodeTypeIf) error
	RemoveNodeType(t NodeTypeIf)
	AddSignalType(SignalType) bool
	RemoveSignalType(t SignalType)
}

type NodeTypeIf interface {
	TreeElement
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
	TreeElement
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
	NamedTreeElementIf
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
	TreeElement
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
	NamedTreeElementIf
	interfaces.Directioner
	SignalType() SignalType
	SetSignalType(SignalType)
}

type Port interface {
	TreeElement
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
	TreeElement
	From() Port
	To() Port
}

// TODO: Turn into interface
func GetRegisteredNodeTypes() []string {
	return registeredNodeTypes.Strings()
}

func GetRegisteredSignalTypes() []string {
	return registeredSignalTypes.Strings()
}

func GetNodeTypeByName(typeName string) (nType NodeTypeIf, ok bool) {
	nType, ok = nodeTypes[typeName]
	return
}

func GetSignalTypeByName(typeName string) (sType SignalType, ok bool) {
	sType, ok = signalTypes[typeName]
	return
}

/*
 *  PlatformIf Model
 */

type PlatformIf interface {
	ToplevelTreeElement
	PlatformId() string
	SetPlatformId(string)
	ProcessByName(string) (ProcessIf, bool)
	Arch() []ArchIf
	PlatformObject() interfaces.PlatformObject
}

type ArchIf interface {
	NamedTreeElementIf
	interfaces.ModePositioner
	Platform() PlatformIf
	IOTypes() []IOTypeIf
	Processes() []ProcessIf
}

type IOTypeIf interface {
	NamedTreeElementIf
	interfaces.IOModer
	Platform() PlatformIf
}

type ProcessIf interface {
	NamedTreeElementIf
	interfaces.ModePositioner
	Arch() ArchIf
	InChannels() []ChannelIf
	OutChannels() []ChannelIf
}

type ChannelIf interface {
	NamedTreeElementIf
	interfaces.Directioner
	interfaces.ModePositioner
	Process() ProcessIf
	IOType() IOTypeIf
	SetIOType(IOTypeIf)
	Link() ChannelIf
}

type MappingIf interface {
	ToplevelTreeElement
	MappingObject() interfaces.MappingObject
	AddMapping(n NodeIf, p ProcessIf)
	SetGraph(SignalGraphIf)
	Graph() SignalGraphIf
	SetPlatform(PlatformIf)
	Platform() PlatformIf
	Mapped(NodeIf) (ProcessIf, bool)
}

type MappedElementIf interface {
	TreeElement
	Mapping() MappingIf
	Node() NodeIf
	Process() ProcessIf
	SetProcess(ProcessIf)
}

func GetRegisteredIOTypes() []string {
	return registeredIOTypes.Strings()
}

func GetIOTypeByName(name string) (ioType IOTypeIf, ok bool) {
	ioType, ok = ioTypes[name]
	return
}

type NamedTreeElementIf interface {
	TreeElement
	interfaces.Namer
}
