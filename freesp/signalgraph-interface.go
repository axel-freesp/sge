package freesp

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
}

type SignalGraph interface {
	TreeElement
	Filename() string
	ItsType() SignalGraphType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
	SetFilename(string)
}

type Library interface {
	TreeElement
	Filename() string
	SignalTypes() []SignalType
	NodeTypes() []NodeType
	Read(data []byte) error
	ReadFile(filepath string) error
	Write() (data []byte, err error)
	WriteFile(filepath string) error
	AddNodeType(NodeType) error
	RemoveNodeType(t NodeType)
	AddSignalType(SignalType)
	RemoveSignalType(t SignalType)
	SetFilename(string)
}

type NodeType interface {
	TreeElement
	TypeName() string
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
	TreeElement
	TypeName() string
	CType() string
	ChannelId() string
	Scope() Scope
	Mode() Mode
}

type PortType interface {
	TreeElement
	SignalType() SignalType
	Name() string
	Direction() PortDirection
}

type Port interface {
	TreeElement
	Name() string
	SignalType() SignalType
	Direction() PortDirection
	Connections() []Port
	Node() Node
	Connection(c *port) Connection
	AddConnection(Port) error
	RemoveConnection(c Port)
}

type PortDirection bool

const (
	InPort  PortDirection = false
	OutPort PortDirection = true
)

type Connection struct {
	From, To Port
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
