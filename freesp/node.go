package freesp

type node struct {
	context  SignalGraph
	name     string
	nodetype NodeType
	inPort   []Port
	outPort  []Port
}

func newNode(name string, ntype *nodeType, context *signalGraph) *node {
	return &node{context, name, ntype, nil, nil}
}

func (n *node) inPortFromName(name string) (p Port, err error) {
	return portFromName(n.inPort, name)
}

func (n *node) outPortFromName(name string) (p Port, err error) {
	return portFromName(n.outPort, name)
}

func portFromName(list []Port, name string) (ret Port, err error) {
	switch {
	case len(name) == 0:
		if len(list) > 1 {
			err = newSignalGraphError("ambiguous port")
			ret = nil
			return
		}
		ret = list[0].(*port)
		err = nil
		return
	default:
		for _, p := range list {
			if p.(*port).name == name {
				ret = p
				err = nil
				return
			}
		}
	}
	err = newSignalGraphError("port not found")
	ret = nil
	return
}

func (n *node) addInPort(pt *namedPortType) {
	n.inPort = append(n.inPort, newPort(pt.name, pt.pType, InPort, n))
}

func (n *node) addOutPort(pt *namedPortType) {
	n.outPort = append(n.outPort, newPort(pt.name, pt.pType, OutPort, n))
}

func (n *node) NodeName() string {
	return n.name
}

func (n *node) ItsType() NodeType {
	return n.nodetype
}

func (n *node) Context() SignalGraph {
	return n.context
}

func (n *node) InPorts() []Port {
	return n.inPort
}

func (n *node) OutPorts() []Port {
	return n.outPort
}
