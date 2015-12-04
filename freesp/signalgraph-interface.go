package freesp

type SignalGraphType interface {
	SignalTypes() []SignalType
	Nodes() []Node
	NodeByName(string) Node
	InputNodes() []Node
	OutputNodes() []Node
	ProcessingNodes() []Node
}

type SignalGraph interface {
	ItsType() SignalGraphType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
}

type NodeType interface {
	TypeName() string
	InPorts() []NamedPortType
	OutPorts() []NamedPortType
	//Implementation() NodeImplementation
}

type NodeImplementation interface {
	ImplementationType() NodeImplementationType
	ElementName() string
	Graph() SignalGraphType
	// missing: input mapping, output mapping
}

type NodeImplementationType int

const (
	NodeTypeElement NodeImplementationType = 0
	NodeTypeGraph   NodeImplementationType = 1
)

type Node interface {
	NodeName() string
	ItsType() NodeType
	Context() SignalGraph
	InPorts() []Port
	OutPorts() []Port
}

type Scope int

const (
	Local Scope = iota
	Global
)

type SignalType interface {
	TypeName() string
	CType() string
	ChannelId() string
	Scope() Scope
}

type PortType interface {
	TypeName() string
	SignalType() SignalType
}

type NamedPortType interface {
	PortType
	Name() string
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
