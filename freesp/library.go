package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/tool"
)

func LibraryNew(filename string) *library {
	return &library{filename, nil, nil}
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

func createNodeTypeFromXml(n backend.XmlNodeType) *nodeType {
	nt := newNodeType(n.TypeName)
	for _, p := range n.InPort {
		nt.addNamedInPortType(p.PName, getPortType(p.PType))
	}
	for _, p := range n.OutPort {
		nt.addNamedOutPortType(p.PName, getPortType(p.PType))
	}
	nodeTypes[n.TypeName] = nt
	// TODO: evaluate <implementation>
	return nt
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
		sType := newSignalType(st.Name, st.Ctype, st.Msgid, scope, mode)
		s.signalTypes = append(s.signalTypes, sType)
	}
	for _, n := range l.NodeTypes {
		nType := createNodeTypeFromXml(n)
		s.nodeTypes = append(s.nodeTypes, nType)
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
	// TODO
	data = nil
	err = fmt.Errorf("library.Write() interface not implemented")
	return
}

func (s *library) WriteFile(filepath string) error {
	// TODO
	return fmt.Errorf("library.WriteFile() interface not implemented")
}
