package freesp

type nodeType struct {
	name              string
	inPorts, outPorts []NamedPortType
	implementation    NodeImplementation
}

func newNodeType(name string) *nodeType {
	return &nodeType{name, nil, nil, nil}
}

func (t *nodeType) addInPort(name string, pType *portType) {
	t.inPorts = append(t.inPorts, &namedPortType{name, pType})
}

func (t *nodeType) addOutPort(name string, pType *portType) {
	t.outPorts = append(t.outPorts, &namedPortType{name, pType})
}

func (t *nodeType) TypeName() string {
	return t.name
}

func (t *nodeType) InPorts() []NamedPortType {
	return t.inPorts
}

func (t *nodeType) OutPorts() []NamedPortType {
	return t.outPorts
}

func (t *nodeType) Implementation() NodeImplementation {
	return t.implementation
}
