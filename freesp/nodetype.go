package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts []NamedPortType
	implementation    []Implementation
	instances         []Node
}

var _ NodeType = (*nodeType)(nil)

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, nil, nil, nil, nil}
}

func (t *nodeType) AddNamedPortType(p NamedPortType) {
	pt := p.(*namedPortType)
	if p.Direction() == InPort {
		t.inPorts = append(t.inPorts, pt)
	} else {
		t.outPorts = append(t.outPorts, pt)
	}
	for _, n := range t.instances {
		if p.Direction() == InPort {
			n.(*node).addInPort(p.(*namedPortType))
		} else {
			n.(*node).addOutPort(p.(*namedPortType))
		}
	}
}

func FindNamedPortType(list []NamedPortType, elem NamedPortType) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemNamedPortType(list *[]NamedPortType, elem NamedPortType) {
	index, ok := FindNamedPortType(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *nodeType) RemoveNamedPortType(p NamedPortType) {
	for _, n := range t.instances {
		n.(*node).removePort(p.(*namedPortType))
	}
	var list *[]NamedPortType
	if p.Direction() == InPort {
		list = &t.inPorts
	} else {
		list = &t.outPorts
	}
	RemNamedPortType(list, p)
}

func (t *nodeType) Instances() []Node {
	return t.instances
}

func (t *nodeType) addInstance(n *node) {
	t.instances = append(t.instances, n)
}

func (t *nodeType) addInPort(name string, pType PortType) {
	t.inPorts = append(t.inPorts, &namedPortType{name, pType.(*portType), InPort})
}

func (t *nodeType) addOutPort(name string, pType PortType) {
	t.outPorts = append(t.outPorts, &namedPortType{name, pType.(*portType), OutPort})
}

func FindImplementation(list []Implementation, elem Implementation) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemImplementation(list *[]Implementation, elem Implementation) {
	index, ok := FindImplementation(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *nodeType) RemoveImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		gt := imp.Graph()
		for len(gt.Nodes()) > 0 {
			gt.RemoveNode(gt.Nodes()[0].(*node))
		}
	}
	RemImplementation(&t.implementation, imp)
}

func (t *nodeType) AddImplementation(imp Implementation) {
	if imp.ImplementationType() == NodeTypeGraph {
		if imp.Graph() == nil {
			log.Fatal("nodeType.AddImplementation: missing graph")
		}
		g := imp.Graph().(*signalGraphType)
		for _, p := range t.inPorts {
			st := p.SignalType()
			ntName := createInputNodeTypeName(st.TypeName())
			nt, ok := nodeTypes[ntName]
			if !ok {
				nt = NodeTypeNew(ntName, "")
				nt.addOutPort("", getPortType(st.TypeName()))
				nodeTypes[ntName] = nt
			}
			if len(nt.outPorts) == 0 {
				log.Fatal("nodeType.AddImplementation: invalid input node type")
			}
			n := NodeNew(fmt.Sprintf("in-%s", p.Name()), nt, imp.Graph())
			n.portlink = p
			err := g.addNode(n)
			if err != nil {
				log.Fatal("nodeType.AddImplementation: AddNode failed:", err)
			}
		}
		for _, p := range t.outPorts {
			st := p.SignalType()
			ntName := createOutputNodeTypeName(st.TypeName())
			nt, ok := nodeTypes[ntName]
			if !ok {
				nt = NodeTypeNew(ntName, "")
				nt.addInPort("", getPortType(st.TypeName()))
				nodeTypes[ntName] = nt
			}
			if len(nt.inPorts) == 0 {
				log.Fatal("nodeType.AddImplementation: invalid input node type")
			}
			n := NodeNew(fmt.Sprintf("out-%s", p.Name()), nt, imp.Graph())
			n.portlink = p
			err := g.addNode(n)
			if err != nil {
				log.Fatal("nodeType.AddImplementation: AddNode failed:", err)
			}
		}
	}
	t.implementation = append(t.implementation, imp)
}

func (t *nodeType) TypeName() string {
	return t.name
}

func (t *nodeType) DefinedAt() string {
	return t.definedAt
}

func (t *nodeType) InPorts() []NamedPortType {
	return t.inPorts
}

func (t *nodeType) OutPorts() []NamedPortType {
	return t.outPorts
}

func (t *nodeType) Implementation() []Implementation {
	return t.implementation
}

func createNodeTypeFromXmlNode(n backend.XmlNode, ntName string) *nodeType {
	nt := NodeTypeNew(ntName, "")
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	nodeTypes[ntName] = nt
	return nt
}

func (t *nodeType) doResolvePort(name string, dir PortDirection) *namedPortType {
	var ports []NamedPortType
	switch dir {
	case InPort:
		ports = t.inPorts
	default:
		ports = t.outPorts
	}
	for _, p := range ports {
		if name == p.Name() {
			return p.(*namedPortType)
		}
	}
	return nil
}

func createNodeTypeFromXml(n backend.XmlNodeType, filename string) *nodeType {
	nt := NodeTypeNew(n.TypeName, filename)
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	if len(n.Implementation) > 0 {
		for _, i := range n.Implementation {
			var iType ImplementationType
			if len(i.SignalGraph) == 1 {
				iType = NodeTypeGraph
			} else {
				iType = NodeTypeElement
			}
			impl := ImplementationNew(i.Name, iType)
			nt.implementation = append(nt.implementation, impl)
			switch iType {
			case NodeTypeElement:
				impl.elementName = i.Name
			default:
				var err error
				var resolvePort = func(name string, dir PortDirection) *namedPortType {
					return nt.doResolvePort(name, dir)
				}
				impl.graph, err = createSignalGraphTypeFromXml(&i.SignalGraph[0], n.TypeName, resolvePort)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nt
}
