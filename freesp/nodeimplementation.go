package freesp

type nodeImplementation struct {
	implementationType          NodeImplementationType
	elementName                 string
	graph                       SignalGraphType
	inputMapping, outputMapping map[Port]Node
}

func (n *nodeImplementation) ImplementationType() NodeImplementationType {
	return n.implementationType
}

func (n *nodeImplementation) ElementName() string {
	return n.elementName
}

func (n *nodeImplementation) Graph() SignalGraphType {
	return n.graph
}
