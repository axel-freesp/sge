package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/tool"
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

func (l *library) Filename() string {
	return l.filename
}

func (l *library) SignalTypes() []SignalType {
	return l.signalTypes
}

func (l *library) NodeTypes() []NodeType {
	return l.nodeTypes
}

func (s *library) Read(data []byte) error {
	l := backend.XmlLibraryNew()
	err := l.Read(data)
	if err != nil {
		return fmt.Errorf("library.Read: %v", err)
	}
	for _, st := range l.SignalTypes {
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
		err := s.AddSignalType(sType)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	for _, n := range l.NodeTypes {
		nType := createNodeTypeFromXml(n, s.Filename())
		err := s.AddNodeType(nType)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	return nil
}

func (s *library) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("library.ReadFile: %v", err)
	}
	err = s.Read(data)
	if err != nil {
		return fmt.Errorf("library.ReadFile: %v", err)
	}
	return err
}

func (s *library) Write() (data []byte, err error) {
	xmllib := CreateXmlLibrary(s)
	data, err = xmllib.Write()
	return
}

func (s *library) WriteFile(filepath string) error {
	data, err := s.Write()
	if err != nil {
		return err
	}
	err = tool.WriteFile(filepath, data)
	return nil
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
