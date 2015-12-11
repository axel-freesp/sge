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
}

func NodeTypeNew(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, nil, nil, nil}
}

func (t *nodeType) AddNamedPortType(p NamedPortType) {
	pt := p.(*namedPortType)
	if p.Direction() == InPort {
		t.inPorts = append(t.inPorts, pt)
	} else {
		t.outPorts = append(t.outPorts, pt)
	}
}

func (t *nodeType) AddInPort(name string, pType PortType) {
	t.inPorts = append(t.inPorts, &namedPortType{name, pType.(*portType), InPort})
}

func (t *nodeType) AddOutPort(name string, pType PortType) {
	t.outPorts = append(t.outPorts, &namedPortType{name, pType.(*portType), OutPort})
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
				nt.AddOutPort("", getPortType(st.TypeName()))
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
				nt.AddInPort("", getPortType(st.TypeName()))
				nodeTypes[ntName] = nt
			}
			if len(nt.inPorts) == 0 {
				log.Fatal("nodeType.AddImplementation: invalid input node type")
			}
			n := NodeNew(fmt.Sprintf("in-%s", p.Name()), nt, imp.Graph())
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
		nt.AddInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.AddOutPort(p.PName, getPortType(p.PType))
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
		nt.AddInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.AddOutPort(p.PName, getPortType(p.PType))
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
