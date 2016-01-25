package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
)

type node struct {
	context  bh.SignalGraphTypeIf
	name     string
	nodetype bh.NodeTypeIf
	inPort   portList
	outPort  portList
	portlink bh.PortTypeIf
	position map[gr.PositionMode]image.Point
	expanded bool
}

/*
 *  freesp.bh.NodeIf API
 */
var _ bh.NodeIf = (*node)(nil)

func NodeNew(name string, ntype bh.NodeTypeIf, context bh.SignalGraphTypeIf) (ret *node, err error) {
	for _, n := range context.Nodes() {
		if n.Name() == name {
			err = fmt.Errorf("NodeNew error: node '%s' already exists in context.", name)
			return
		}
	}
	if len(ntype.InPorts())+len(ntype.OutPorts()) == 0 {
		err = fmt.Errorf("NodeNew error: type '%s' has no ports.", ntype.TypeName())
		return
	}
	ret = &node{context, name, ntype, portListInit(), portListInit(), nil, make(map[gr.PositionMode]image.Point), false}
	for _, p := range ntype.InPorts() {
		ret.addInPort(p)
	}
	for _, p := range ntype.OutPorts() {
		ret.addOutPort(p)
	}
	ntype.(*nodeType).addInstance(ret)
	return
}

func InputNodeNew(name string, stypeName string, context bh.SignalGraphTypeIf) (ret *node, err error) {
	st, ok := freesp.GetSignalTypeByName(stypeName)
	if !ok {
		err = fmt.Errorf("InputNodeNew error: signaltype %s not defined", stypeName)
		return
	}
	ntName := createInputNodeTypeName(stypeName)
	nt, ok := freesp.GetNodeTypeByName(ntName)
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.(*nodeType).addOutPort("", st)
		freesp.RegisterNodeType(nt)
	}
	return NodeNew(name, nt, context)
}

func OutputNodeNew(name string, stypeName string, context bh.SignalGraphTypeIf) (ret *node, err error) {
	st, ok := freesp.GetSignalTypeByName(stypeName)
	if !ok {
		err = fmt.Errorf("InputNodeNew error: signaltype %s not defined", stypeName)
		return
	}
	ntName := createOutputNodeTypeName(stypeName)
	nt, ok := freesp.GetNodeTypeByName(ntName)
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.(*nodeType).addInPort("", st)
		freesp.RegisterNodeType(nt)
	}
	return NodeNew(name, nt, context)
}

func (n *node) ItsType() bh.NodeTypeIf {
	return n.nodetype
}

func (n *node) InPorts() []bh.PortIf {
	return n.inPort.Ports()
}

func (n *node) OutPorts() []bh.PortIf {
	return n.outPort.Ports()
}

func (n *node) Context() bh.SignalGraphTypeIf {
	return n.context
}

func (n *node) Expanded() bool {
	return n.expanded
}

func (n *node) SetExpanded(xp bool) {
	n.expanded = xp
}

func (n *node) CreateXml() (buf []byte, err error) {
	if len(n.InPorts()) == 0 {
		xmlnode := CreateXmlInputNode(n)
		buf, err = xmlnode.Write()
	} else if len(n.OutPorts()) == 0 {
		xmlnode := CreateXmlOutputNode(n)
		buf, err = xmlnode.Write()
	} else {
		xmlnode := CreateXmlProcessingNode(n)
		buf, err = xmlnode.Write()
	}
	return
}

/*
 *  fmt.Stringer API
 */
var _ fmt.Stringer = (*node)(nil)

func (n *node) String() (s string) {
	s = fmt.Sprintf("bh.NodeIf(%s: %d inports, %d outports)",
		n.name, len(n.inPort.Ports()), len(n.outPort.Ports()))
	return
}

/*
 *  tr.TreeElement API
 */

var _ tr.TreeElement = (*node)(nil)

func (n *node) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	var prop tr.Property
	if isParentReadOnly(tree, cursor) {
		prop = freesp.PropertyNew(false, false, false)
	} else {
		prop = freesp.PropertyNew(true, true, true)
	}
	var image tr.Symbol
	if len(n.InPorts()) == 0 {
		image = tr.SymbolInputNode
	} else if len(n.OutPorts()) == 0 {
		image = tr.SymbolOutputNode
	} else {
		image = tr.SymbolProcessingNode
	}
	err := tree.AddEntry(cursor, image, n.Name(), n, prop)
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

func (n *node) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	log.Fatal("node.AddNewObject - nothing to add.")
	return
}

func (n *node) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parentId := tree.Parent(cursor)
	if n != tree.Object(parentId) {
		log.Fatal("node.RemoveObject error: not removing child of mine.")
	}
	nt := n.ItsType()
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.PortIf:
		p := obj.(bh.PortIf)
		for index, c := range p.Connections() {
			conn := p.Connection(c)
			removed = append(removed, tr.IdWithObject{cursor.Path, index, conn})
		}
		var list portTypeList
		if p.Direction() == gr.InPort {
			list = nt.(*nodeType).inPorts
		} else {
			list = nt.(*nodeType).outPorts
		}
		_, ok, index := list.Find(p.Name())
		if !ok {
			log.Println("node.RemoveObject: saving removed port", p)
			removed = append(removed, tr.IdWithObject{parentId.Path, index, obj})
		}
		tree.Remove(cursor)

	default:
		log.Fatal("bh.NodeIf.RemoveObject error: invalid type %T", obj)
	}
	return
}

/*
 *  freesp private functions
 */
func (n *node) inPortFromName(name string) (p bh.PortIf, err error) {
	return portFromName(n.inPort.Ports(), name)
}

func (n *node) outPortFromName(name string) (p bh.PortIf, err error) {
	return portFromName(n.outPort.Ports(), name)
}

func portFromName(list []bh.PortIf, name string) (ret bh.PortIf, err error) {
	switch {
	case len(name) == 0:
		if len(list) == 0 {
			err = fmt.Errorf("portFromName: no port found (1)")
			ret = nil
			return
		}
		if len(list) > 1 {
			var listtext string
			for _, l := range list {
				listtext = fmt.Sprintf("%s(%v)\n", listtext, l)
			}
			err = fmt.Errorf("portFromName: ambiguous port, list = %v", listtext)
			ret = nil
			return
		}
		ret = list[0].(*port)
		err = nil
		return
	default:
		for _, p := range list {
			if p.Name() == name {
				ret = p
				err = nil
				return
			}
		}
	}
	err = fmt.Errorf("portFromName: no port found (2)")
	ret = nil
	return
}

func (n *node) addInPort(pt bh.PortTypeIf) {
	n.inPort.Append(newPort(pt, n))
}

func (n *node) addOutPort(pt bh.PortTypeIf) {
	n.outPort.Append(newPort(pt, n))
}

func (n *node) removePort(pt bh.PortTypeIf) {
	var list *portList
	if pt.Direction() == gr.InPort {
		list = &n.inPort
	} else {
		list = &n.outPort
	}
	var i int
	for i = 0; i < len(list.Ports()); i++ {
		if list.Ports()[i].Name() == pt.Name() {
			break
		}
	}
	toRemove := list.Ports()[i]
	for _, c := range toRemove.Connections() {
		c.RemoveConnection(toRemove)
	}
	list.Remove(toRemove)
}

func isParentReadOnly(tree tr.TreeIf, cursor tr.Cursor) bool {
	parentId := tree.Parent(cursor)
	return tree.Property(parentId).IsReadOnly()
}

func IsProcessingNode(n bh.NodeIf) bool {
	if len(n.InPorts()) == 0 {
		return false
	}
	if len(n.OutPorts()) == 0 {
		return false
	}
	return true
}

/*
 *      ModePositioner API
 */

func (n *node) ModePosition(mode gr.PositionMode) (pos image.Point) {
	pos = n.position[mode]
	return
}

func (n *node) SetModePosition(mode gr.PositionMode, pos image.Point) {
	n.position[mode] = pos
}

/*
 *      Namer API
 */

func (n *node) Name() string {
	return n.name
}

func (n *node) SetName(newName string) {
	n.name = newName
}

var _ gr.Namer = (*node)(nil)

/*
 *      Porter API
 */

func (n *node) InPortIndex(portName string) int {
	_, ok, index := n.inPort.Find(n.Name(), portName)
	if ok {
		return index
	}
	return -1
}

func (n *node) OutPortIndex(portName string) int {
	_, ok, index := n.outPort.Find(n.Name(), portName)
	if ok {
		return index
	}
	return -1
}

/*
 *      nodeList
 *
 */

type nodeList struct {
	nodes []bh.NodeIf
}

func nodeListInit() nodeList {
	return nodeList{nil}
}

func (l *nodeList) Append(n *node) {
	l.nodes = append(l.nodes, n)
}

func (l *nodeList) Remove(n bh.NodeIf) {
	var i int
	for i = range l.nodes {
		if n == l.nodes[i] {
			break
		}
	}
	if i >= len(l.nodes) {
		for _, v := range l.nodes {
			log.Printf("nodeList.RemoveNode have bh.NodeIf %v\n", v)
		}
		log.Fatalf("nodeList.RemoveNode error: bh.NodeIf %v not in this list\n", n)
	}
	for i++; i < len(l.nodes); i++ {
		l.nodes[i-1] = l.nodes[i]
	}
	l.nodes = l.nodes[:len(l.nodes)-1]
}

func (l *nodeList) Nodes() []bh.NodeIf {
	return l.nodes
}

func (l *nodeList) Find(node bh.NodeIf) (ok bool, index int) {
	ok = false
	var n bh.NodeIf
	for index, n = range l.nodes {
		if n == node {
			ok = true
			return
		}
	}
	return
}
