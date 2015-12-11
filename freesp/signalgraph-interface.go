package freesp

type SignalGraphType interface {
	Libraries() []Library
	Nodes() []Node
	NodeByName(string) Node
	InputNodes() []Node
	OutputNodes() []Node
	ProcessingNodes() []Node
	AddNode(Node) error
}

type SignalGraph interface {
	Filename() string
	ItsType() SignalGraphType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
	SetFilename(string)
}

type Library interface {
	Filename() string
	SignalTypes() []SignalType
	NodeTypes() []NodeType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
	AddNodeType(NodeType) error
	AddSignalType(SignalType) error
	SetFilename(string)
}

type NodeType interface {
	TypeName() string
	DefinedAt() string
	InPorts() []NamedPortType
	OutPorts() []NamedPortType
	Implementation() []Implementation
	AddNamedPortType(NamedPortType)
	AddImplementation(Implementation)
}

type Implementation interface {
	ImplementationType() ImplementationType
	ElementName() string
	Graph() SignalGraphType
	// missing: input mapping, output mapping
}

type ImplementationType int

const (
	NodeTypeElement ImplementationType = 0
	NodeTypeGraph   ImplementationType = 1
)

type Node interface {
	NodeName() string
	ItsType() NodeType
	Context() SignalGraphType
	InPorts() []Port
	OutPorts() []Port
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
	TypeName() string
	CType() string
	ChannelId() string
	Scope() Scope
	Mode() Mode
}

type PortType interface {
	TypeName() string
	SignalType() SignalType
}

type NamedPortType interface {
	PortType
	Name() string
	Direction() PortDirection
}

type Port interface {
	PortName() string
	ItsType() PortType
	Direction() PortDirection
	Connections() []Port
	Node() Node
	AddConnection(Port) error
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

type Connection struct {
	From, To Port
}

func GetRegisteredNodeTypes() []string {
	return registeredNodeTypes
}

func GetRegisteredSignalTypes() []string {
	return registeredSignalTypes
}

func GetNodeTypeByName(typeName string) (nType NodeType, ok bool) {
	nType, ok = nodeTypes[typeName]
	return
}

func GetSignalTypeByName(typeName string) (sType SignalType, ok bool) {
	sType, ok = signalTypes[typeName]
	return
}
