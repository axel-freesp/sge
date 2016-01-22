package freesp

import (
	interfaces "github.com/axel-freesp/sge/interface"
)

type Context interface {
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
	SetGraphForNew(g SignalGraph)
	SetPlatformForNew(p Platform)
}

/*
 *  Behaviour Model
 */

type SignalGraphType interface {
	TreeElement
	Libraries() []Library
	Nodes() []Node
	NodeByName(string) (Node, bool)
	InputNodes() []Node
	OutputNodes() []Node
	ProcessingNodes() []Node
	AddNode(Node) error
	RemoveNode(n Node)
	Context() Context
}

type SignalGraph interface {
	ToplevelTreeElement
	ItsType() SignalGraphType
	GraphObject() interfaces.GraphObject
}

type Library interface {
	ToplevelTreeElement
	SignalTypes() []SignalType
	NodeTypes() []NodeType
	AddNodeType(NodeType) error
	RemoveNodeType(t NodeType)
	AddSignalType(SignalType) bool
	RemoveSignalType(t SignalType)
}

type NodeType interface {
	TreeElement
	TypeName() string
	SetTypeName(string)
	DefinedAt() string
	InPorts() []PortType
	OutPorts() []PortType
	Implementation() []Implementation
	Instances() []Node
	AddNamedPortType(PortType)
	RemoveNamedPortType(PortType)
	AddImplementation(Implementation)
	RemoveImplementation(Implementation)
}

type Implementation interface {
	TreeElement
	ImplementationType() ImplementationType
	ElementName() string
	SetElemName(string)
	Graph() SignalGraphType
	GraphObject() interfaces.GraphObject
}

type ImplementationType int

const (
	NodeTypeElement ImplementationType = 0
	NodeTypeGraph   ImplementationType = 1
)

type Node interface {
	TreeElement
	interfaces.Namer
	interfaces.Positioner
	interfaces.Porter
	ItsType() NodeType
	InPorts() []Port
	OutPorts() []Port
	Context() SignalGraphType
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
	TreeElement
	interfaces.Namer
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
	Node() Node
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

func GetNodeTypeByName(typeName string) (nType NodeType, ok bool) {
	nType, ok = nodeTypes[typeName]
	return
}

func GetSignalTypeByName(typeName string) (sType SignalType, ok bool) {
	sType, ok = signalTypes[typeName]
	return
}

/*
 *  Platform Model
 */

type Platform interface {
	ToplevelTreeElement
	interfaces.Shaper
	PlatformId() string
	SetPlatformId(string)
	ProcessByName(string) (Process, bool)
	Arch() []Arch
	PlatformObject() interfaces.PlatformObject
}

type Arch interface {
	TreeElement
	interfaces.Namer
	interfaces.ModePositioner
	Platform() Platform
	IOTypes() []IOType
	Processes() []Process
}

type IOType interface {
	TreeElement
	interfaces.Namer
	interfaces.IOModer
	Platform() Platform
}

type Process interface {
	TreeElement
	interfaces.Namer
	interfaces.ModePositioner
	Arch() Arch
	InChannels() []Channel
	OutChannels() []Channel
}

type Channel interface {
	TreeElement
	interfaces.Namer
	interfaces.Directioner
	interfaces.ModePositioner
	Process() Process
	IOType() IOType
	SetIOType(IOType)
	Link() Channel
}

type Mapping interface {
	ToplevelTreeElement
	MappingObject() interfaces.MappingObject
	AddMapping(n Node, p Process)
	SetGraph(SignalGraph)
	Graph() SignalGraph
	SetPlatform(Platform)
	Platform() Platform
	Mapped(Node) (Process, bool)
}

type MappedElement interface {
	TreeElement
	Mapping() Mapping
	Node() Node
	Process() Process
	SetProcess(Process)
}

func GetRegisteredIOTypes() []string {
	return registeredIOTypes.Strings()
}

func GetIOTypeByName(name string) (ioType IOType, ok bool) {
	ioType, ok = ioTypes[name]
	return
}

type Filenamer interface {
	Filename() string
	SetFilename(string)
}

type Remover interface {
	Remove(Tree)
}

type FileDataIf interface {
	Filenamer
	Remover
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
}
