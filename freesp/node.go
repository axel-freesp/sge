package freesp

import (
	"fmt"
	"log"
)

type node struct {
	context  SignalGraphType
	name     string
	nodetype NodeType
	inPort   []Port
	outPort  []Port
	portlink NamedPortType
}

/*
 *  freesp.Node API
 */
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

/*
 *  fmt.Stringer API
 */
var _ fmt.Stringer = (*node)(nil)

func (n *node) String() (s string) {
	s = fmt.Sprintf("Node(%s: %d inports, %d outports)",
		n.name, len(n.inPort), len(n.outPort))
	return
}

/*
 *  TreeElement API
 */

func IsParentReadOnly(tree Tree, cursor Cursor) bool {
	parentId := tree.Parent(cursor)
	return tree.Property(parentId).IsReadOnly()
}

var _ TreeElement = (*node)(nil)

func (n *node) AddToTree(tree Tree, cursor Cursor) {
	var prop property
	if IsParentReadOnly(tree, cursor) {
		prop = 0
	} else {
		prop = mayEdit | mayRemove
	}
	var isImplementation bool
	var nodeType NodeType
	parentId := tree.Parent(cursor)
	parent := tree.Object(parentId)
	switch parent.(type) {
	case SignalGraphType:
		isImplementation = false
	case Implementation:
		isImplementation = true
		nodeType = tree.Object(tree.Parent(parentId)).(NodeType)
	default:
		log.Fatalf("node.AddToTree error: invalid parent type %T\n", parent)
	}
	var image Symbol
	if len(n.InPorts()) == 0 { // TODO: check if linked to parent port
		image = SymbolOutputNode
	} else if len(n.OutPorts()) == 0 {
		image = SymbolInputNode
	} else if !isImplementation {
		image = SymbolProcessingNode
	} else {
		allInputPortsLinked := len(n.InPorts()) == 1
		for _, p := range n.InPorts() {
			found := false
			for _, p2 := range nodeType.InPorts() {
				if p.PortName() == p2.Name() && n.NodeName() == fmt.Sprintf("in-%s", p2.Name()) {
					found = true
					break
				}
			}
			if !found {
				allInputPortsLinked = false
				break
			}
		}
		allOutputPortsLinked := len(n.OutPorts()) == 1
		for _, p := range n.OutPorts() {
			found := false
			for _, p2 := range nodeType.OutPorts() {
				if p.PortName() == p2.Name() && n.NodeName() == fmt.Sprintf("out-%s", p2.Name()) {
					found = true
					break
				}
			}
			if !found {
				allOutputPortsLinked = false
				break
			}
		}
		if allInputPortsLinked {
			image = SymbolInputNode
		} else if allOutputPortsLinked {
			image = SymbolOutputNode
		} else {
			image = SymbolProcessingNode
		}
	}
	err := tree.AddEntry(cursor, image, n.NodeName(), n, prop)
	if err != nil {
		log.Fatalf("node.AddToTree error: AddEntry failed: %s\n", err)
	}
	child := tree.Append(cursor)
	n.ItsType().AddToTree(tree, child)
	for _, p := range n.InPorts() {
		child := tree.Append(cursor)
		p.AddToTree(tree, child)
	}
	for _, p := range n.OutPorts() {
		child := tree.Append(cursor)
		p.AddToTree(tree, child)
	}
}

func (n *node) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case NamedPortType:
		// If this function is called, the freesp model if different from the tree model
		// cursor points to n, cursor.Position may indicate where to insert.
		// Add missing port entry to the tree, matching name and direction with obj.
		newCursor = tree.Insert(cursor)
		pt := obj.(NamedPortType)
		var list []Port
		if pt.Direction() == InPort {
			list = n.InPorts()
		} else {
			list = n.OutPorts()
		}
		for _, p := range list {
			if p.PortName() == pt.Name() {
				p.AddToTree(tree, newCursor)
				break
			}
		}
	default:
		log.Fatal("Node.AddNewObject error: invalid type %T", obj)
	}
	return
}

func (n *node) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	// TODO
	return
}

func IndexOfNodeInGraph(tree Tree, n Node) (index int) {
	nCursor := tree.Cursor(n)
	gCursor := tree.Parent(nCursor)
	index = gCursor.Position
	return
}

func IsProcessingNode(n Node) bool {
	if len(n.InPorts()) == 0 {
		return false
	}
	if len(n.OutPorts()) == 0 {
		return false
	}
	return true
}

/*
 *  freesp private functions
 */
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
