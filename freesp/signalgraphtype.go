package freesp

type signalGraphType struct {
	signalTypes                                     []SignalType
	nodes, inputNodes, outputNodes, processingNodes []Node
}

func newSignalGraphType() *signalGraphType {
	return &signalGraphType{nil, nil, nil, nil, nil}
}

func (t *signalGraphType) Nodes() []Node {
	return t.nodes
}

func (t *signalGraphType) SignalTypes() []SignalType {
	return t.signalTypes
}

func (t *signalGraphType) NodeByName(name string) Node {
	for _, n := range t.nodes {
		if n.NodeName() == name {
			return n
		}
	}
	return nil
}

func (t *signalGraphType) InputNodes() []Node {
	return t.inputNodes
}

func (t *signalGraphType) OutputNodes() []Node {
	return t.outputNodes
}

func (t *signalGraphType) ProcessingNodes() []Node {
	return t.processingNodes
}
