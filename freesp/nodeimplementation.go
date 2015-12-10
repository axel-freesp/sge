package freesp

type implementation struct {
	implementationType ImplementationType
	elementName        string
	graph              SignalGraphType
}

func ImplementationNew(iName string, iType ImplementationType) *implementation {
	ret := &implementation{iType, iName, nil}
	if iType == NodeTypeGraph {
		ret.graph = SignalGraphTypeNew()
	}
	return ret
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
