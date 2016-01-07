package freesp

/*
 *  Behaviour Model
 */

type Context interface {
	GetLibrary(libname string) (lib Library, err error)
}

type SignalGraphType interface {
	TreeElement
	Libraries() []Library
	Nodes() []Node
	NodeByName(string) Node
	InputNodes() []Node
	OutputNodes() []Node
	ProcessingNodes() []Node
	AddNode(Node) error
	RemoveNode(n Node)
	Context() Context
}

type SignalGraph interface {
	TreeElement
	Filenamer
	ItsType() SignalGraphType
}

type Library interface {
	TreeElement
	Filenamer
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
	// missing: input mapping, output mapping
}

type ImplementationType int

const (
	NodeTypeElement ImplementationType = 0
	NodeTypeGraph   ImplementationType = 1
)

type Node interface {
	TreeElement
	Namer
	Positioner
	Porter
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
	Namer
	Directioner
	SignalType() SignalType
	SetSignalType(SignalType)
}

type Port interface {
	TreeElement
	Directioner
	Name() string
	SignalType() SignalType
	Connections() []Port
	Node() Node
	Connection(Port) Connection
	AddConnection(Connection) error
	RemoveConnection(c Port)
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

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
	TreeElement
	Filenamer
	PlatformId() string
	SetPlatformId(string)
	Arch() []Arch
}

type Arch interface {
	TreeElement
	Namer
	Platform() Platform
	IOTypes() []IOType
	Processes() []Process
}

type IOType interface {
	TreeElement
	Namer
	Mode() IOMode
	SetMode(IOMode)
	Platform() Platform
}

type IOMode string

const (
	IOModeShmem IOMode = "Shared Memory"
	IOModeAsync IOMode = "Asynchronous"
	IOModeSync  IOMode = "Isochronous"
)

type Process interface {
	TreeElement
	Namer
	Arch() Arch
	InChannels() []Channel
	OutChannels() []Channel
}

type Channel interface {
	TreeElement
	Namer
	Directioner
	Process() Process
	IOType() IOType
	SetIOType(IOType)
	Link() Channel
}

func GetRegisteredIOTypes() []string {
	return registeredIOTypes.Strings()
}

func GetIOTypeByName(name string) (ioType IOType, ok bool) {
	ioType, ok = ioTypes[name]
	return
}
