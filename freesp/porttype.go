package freesp

import "fmt"

// portType

type portType struct {
	name string
	ref  *signalType
}

var _ PortType = (*portType)(nil)

func PortTypeNew(name string, st SignalType) *portType {
	return &portType{name, st.(*signalType)}
}

func (t *portType) TypeName() string {
	return t.name
}

func (t *portType) SignalType() SignalType {
	return t.ref
}

// namedPortType

type namedPortType struct {
	name      string
	pType     *portType
	direction PortDirection
}

var _ NamedPortType = (*namedPortType)(nil)

func NamedPortTypeNew(name string, pTypeName string, dir PortDirection) *namedPortType {
	pt := getPortType(pTypeName)
	return &namedPortType{name, pt, dir}
}

func (t *namedPortType) TypeName() string {
	return t.pType.TypeName()
}

func (t *namedPortType) Name() string {
	return t.name
}

func (t *namedPortType) Direction() PortDirection {
	return t.direction
}

func (t *namedPortType) SignalType() SignalType {
	return t.pType.SignalType()
}

func (t *namedPortType) String() (s string) {
	s = fmt.Sprintf("NamedPortType(%s, %s, %s)", t.name, t.direction, t.SignalType())
	return
}
