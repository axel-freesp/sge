package freesp

type implementation struct {
	implementationType ImplementationType
	elementName        string
	graph              SignalGraphType
}

func newImplementation(iType ImplementationType) *implementation {
	return &implementation{iType, "", nil}
}

func (n *implementation) ImplementationType() ImplementationType {
	return n.implementationType
}

func (n *implementation) ElementName() string {
	return n.elementName
}

func (n *implementation) Graph() SignalGraphType {
	return n.graph
}
