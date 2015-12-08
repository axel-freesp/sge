package freesp

type SignalGraphType interface {
	Name() string
	Libraries() []Library
	SignalTypes() []SignalType
	Nodes() []Node
	NodeByName(string) Node
	InputNodes() []Node
	OutputNodes() []Node
	ProcessingNodes() []Node
}

type SignalGraph interface {
	Filename() string
	ItsType() SignalGraphType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
}

type Library interface {
	Filename() string
	SignalTypes() []SignalType
	NodeTypes() []NodeType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
}

type NodeType interface {
	TypeName() string
	DefinedAt() string
	InPorts() []NamedPortType
	OutPorts() []NamedPortType
	Implementation() Implementation
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
