package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

func LibraryNew(filename string) *library {
	ret := &library{filename, nil, nil}
	libraries[filename] = ret
	return ret
}

type library struct {
	filename    string
	signalTypes []SignalType
	nodeTypes   []NodeType
}

var _ Library = (*library)(nil)

func (l *library) Filename() string {
	return l.filename
}

func (l *library) SignalTypes() []SignalType {
	return l.signalTypes
}

func (l *library) NodeTypes() []NodeType {
	return l.nodeTypes
}

func (l *library) createLibFromXml(xmlLib *backend.XmlLibrary) error {
	for _, st := range xmlLib.SignalTypes {
		var scope Scope
		var mode Mode
		switch st.Scope {
		case "local":
			scope = Local
		default:
			scope = Global
		}
		switch st.Mode {
		case "sync":
			mode = Synchronous
		default:
			mode = Asynchronous
		}
		sType := SignalTypeNew(st.Name, st.Ctype, st.Msgid, scope, mode)
		err := l.AddSignalType(sType)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	for _, n := range xmlLib.NodeTypes {
		nType := createNodeTypeFromXml(n, l.Filename())
		err := l.AddNodeType(nType)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	return nil
}

func (s *library) Read(data []byte) error {
	l := backend.XmlLibraryNew()
	err := l.Read(data)
	if err != nil {
		return fmt.Errorf("library.Read: %v", err)
	}
	return s.createLibFromXml(l)
}

func (s *library) ReadFile(filepath string) error {
	l := backend.XmlLibraryNew()
	err := l.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("library.Read: %v", err)
	}
	return s.createLibFromXml(l)
}

func (s *library) Write() (data []byte, err error) {
	xmllib := CreateXmlLibrary(s)
	data, err = xmllib.Write()
	return
}

func (s *library) WriteFile(filepath string) error {
	xmllib := CreateXmlLibrary(s)
	return xmllib.WriteFile(filepath)
}

func (s *library) SetFilename(filename string) {
	s.filename = filename
}

func (s *library) AddNodeType(t NodeType) error {
	nType, ok := nodeTypes[t.TypeName()]
	if ok {
		log.Printf("library.AddNodeType: warning: adding existing node type definition %s (taking the existing one)", t.TypeName())
	} else {
		nType = t.(*nodeType)
		nodeTypes[t.TypeName()] = nType
		registeredNodeTypes = append(registeredNodeTypes, t.TypeName())
		log.Println("library.AddNodeType: registered ", t.TypeName())
	}
	for _, nt := range s.nodeTypes {
		if nt.TypeName() == t.TypeName() {
			return fmt.Errorf("adding duplicate node type definition %s (ignored)", t.TypeName())
		}
	}
	s.nodeTypes = append(s.nodeTypes, nType)
	return nil
}

func FindNodeType(list []NodeType, elem NodeType) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemoveNodeType(list []NodeType, elem NodeType) {
	index, ok := FindNodeType(list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(list); j++ {
		list[j-1] = list[j]
	}
	list = list[:len(list)-1]
}

func (l *library) RemoveNodeType(t NodeType) {
	for _, n := range t.(*nodeType).instances {
		n.(*node).context.RemoveNode(n)
	}
	RemoveNodeType(l.nodeTypes, t)
}

func (l *library) AddSignalType(s SignalType) error {
	sType := signalTypes[s.TypeName()]
	if sType != nil {
		log.Printf("library.AddSignalType: warning: adding existing signal type definition %s (taking the existing one)", s.TypeName())
	} else {
		sType = s.(*signalType)
		signalTypes[s.TypeName()] = sType
		registeredSignalTypes = append(registeredSignalTypes, s.TypeName())
	}
	for _, st := range l.signalTypes {
		if st.TypeName() == s.TypeName() {
			return fmt.Errorf("adding duplicate signal type definition %s (ignored)", s.TypeName())
		}
	}
	l.signalTypes = append(l.signalTypes, sType)
	return nil
}

func (l *library) RemoveSignalType(s SignalType) {
}
