package freesp

// portType

type portType struct {
	name string
	ref  *signalType
}

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
