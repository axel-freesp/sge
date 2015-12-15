package freesp

import "fmt"

type node struct {
	context  SignalGraphType
	name     string
	nodetype NodeType
	inPort   []Port
	outPort  []Port
	portlink NamedPortType
}

var _ Node = (*node)(nil)

func NodeNew(name string, ntype NodeType, context SignalGraphType) *node {
	ret := &node{context.(*signalGraphType), name, ntype.(*nodeType), nil, nil, nil}
	for _, p := range ntype.InPorts() {
		ret.addInPort(p.(*namedPortType))
	}
	for _, p := range ntype.OutPorts() {
		ret.addOutPort(p.(*namedPortType))
	}
	ntype.(*nodeType).addInstance(ret)
	return ret
}

func (n *node) String() (s string) {
	s = fmt.Sprintf("Node(%s: %d inports, %d outports)",
		n.name, len(n.inPort), len(n.outPort))
	return
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
		if len(list) == 0 {
			err = newSignalGraphError("no port found")
			ret = nil
			return
		}
		if len(list) > 1 {
			var listtext string
			for _, l := range list {
				listtext = fmt.Sprintf("%s(%v)\n", listtext, l)
			}
			err = newSignalGraphError(fmt.Sprintf("ambiguous port, list = %v", listtext))
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

func (n *node) removePort(pt *namedPortType) {
	var list []Port
	if pt.Direction() == InPort {
		list = n.inPort
	} else {
		list = n.outPort
	}
	var i int
	for i = 0; i < len(list); i++ {
		if list[i].PortName() == pt.Name() {
			break
		}
	}
	toRemove := list[i]
	for _, c := range toRemove.(*port).connections {
		c.RemoveConnection(toRemove)
	}
	for j := i + 1; j < len(list); j++ {
		list[j-1] = list[j]
	}
	list = list[:len(list)-2]
}

func (n *node) NodeName() string {
	return n.name
}

func (n *node) ItsType() NodeType {
	return n.nodetype
}

func (n *node) Context() SignalGraphType {
	return n.context
}

func (n *node) InPorts() []Port {
	return n.inPort
}

func (n *node) OutPorts() []Port {
	return n.outPort
}
