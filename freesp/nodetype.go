package freesp

import (
	"github.com/axel-freesp/sge/backend"
	"log"
)

type nodeType struct {
	name              string
	definedAt         string
	inPorts, outPorts []NamedPortType
	implementation    Implementation
}

func newNodeType(name, definedAt string) *nodeType {
	return &nodeType{name, definedAt, nil, nil, nil}
}

func (t *nodeType) addInPort(name string, pType *portType) {
	t.inPorts = append(t.inPorts, &namedPortType{name, pType, InPort})
}

func (t *nodeType) addOutPort(name string, pType *portType) {
	t.outPorts = append(t.outPorts, &namedPortType{name, pType, OutPort})
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

func (t *nodeType) Implementation() Implementation {
	return t.implementation
}

func createNodeTypeFromXmlNode(n backend.XmlNode, ntName string) *nodeType {
	nt := newNodeType(ntName, "")
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	nodeTypes[ntName] = nt
	return nt
}

func createNodeTypeFromXml(n backend.XmlNodeType, filename string) *nodeType {
	nt := newNodeType(n.TypeName, filename)
	for _, p := range n.InPort {
		nt.addInPort(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addOutPort(p.PName, getPortType(p.PType))
	}
	nodeTypes[n.TypeName] = nt
	if len(n.Implementation) > 0 {
		for _, i := range n.Implementation {
			var iType ImplementationType
			if len(i.SignalGraph) == 1 {
				iType = NodeTypeGraph
			} else {
				iType = NodeTypeElement
			}
			nt.implementation = newImplementation(iType)
			impl := nt.implementation.(*implementation)
			switch iType {
			case NodeTypeElement:
				impl.elementName = i.Name
			default:
				var err error
				var resolvePort = func(name string, dir PortDirection) *namedPortType {
					var ports []NamedPortType
					switch dir {
					case InPort:
						ports = nt.inPorts
					default:
						ports = nt.outPorts
					}
					for _, p := range ports {
						if name == p.Name() {
							return p.(*namedPortType)
						}
					}
					return nil
				}
				log.Println("adding implementation for node type", n.TypeName)
				impl.graph, err = createSignalGraphTypeFromXml(&i.SignalGraph[0], n.TypeName, resolvePort)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return nt
}
