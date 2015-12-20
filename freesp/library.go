package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

func LibraryNew(filename string) *library {
	ret := &library{filename, signalTypeListInit(), nodeTypeListInit()}
	libraries[filename] = ret
	return ret
}

type library struct {
	filename    string
	signalTypes signalTypeList
	nodeTypes   nodeTypeList
}

var _ Library = (*library)(nil)

func (l *library) Filename() string {
	return l.filename
}

func (l *library) SignalTypes() []SignalType {
	return l.signalTypes.SignalTypes()
}

func (l *library) NodeTypes() []NodeType {
	return l.nodeTypes.NodeTypes()
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
		l.AddSignalType(sType)
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

func (l *library) AddNodeType(t NodeType) error {
	nType, ok := nodeTypes[t.TypeName()]
	if ok {
		log.Printf(`library.AddNodeType: warning: adding existing node
			type definition %s (taking the existing one)`, t.TypeName())
	} else {
		nType = t.(*nodeType)
		nodeTypes[t.TypeName()] = nType
		registeredNodeTypes.Append(t.TypeName())
		log.Println("library.AddNodeType: registered ", t.TypeName())
	}
	for _, nt := range l.nodeTypes.NodeTypes() {
		if nt.TypeName() == t.TypeName() {
			return fmt.Errorf(`adding duplicate node type definition
				%s (ignored)`, t.TypeName())
		}
	}
	l.nodeTypes.Append(nType)
	return nil
}

func (l *library) RemoveNodeType(nt NodeType) {
	for _, n := range nt.(*nodeType).instances {
		n.(*node).context.RemoveNode(n)
	}
	delete(nodeTypes, nt.TypeName())
	registeredNodeTypes.Remove(nt.TypeName())
	l.nodeTypes.Remove(nt)
}

func (l *library) AddSignalType(s SignalType) {
	sType := signalTypes[s.TypeName()]
	if sType != nil {
		log.Printf(`library.AddSignalType: warning: adding existing
			signal type definition %s (taking the existing)`, s.TypeName())
	} else {
		sType = s.(*signalType)
		signalTypes[s.TypeName()] = sType
		registeredSignalTypes.Append(s.TypeName())
	}
	for _, st := range l.signalTypes.SignalTypes() {
		if st.TypeName() == s.TypeName() {
			log.Printf(`library.AddSignalType: warning: adding
				duplicate signal type definition %s (ignored)`, s.TypeName())
			return
		}
	}
	l.signalTypes.Append(sType)
}

func (l *library) RemoveSignalType(st SignalType) {
	for _, ntName := range registeredNodeTypes.Strings() {
		nt := nodeTypes[ntName]
		for _, p := range nt.InPorts() {
			if p.SignalType() == st {
				log.Printf(`library.RemoveSignalType warning:
					SignalType %v is still in use\n`, st)
				return
			}
		}
		for _, p := range nt.OutPorts() {
			if p.SignalType() == st {
				log.Printf(`library.RemoveSignalType warning:
					SignalType %v is still in use\n`, st)
				return
			}
		}
	}
	delete(portTypes, st.TypeName())
	delete(signalTypes, st.TypeName())
	registeredSignalTypes.Remove(st.TypeName())
	l.signalTypes.Remove(st)
}

var _ TreeElement = (*library)(nil)

func (l *library) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolLibrary, l.Filename(), l, mayAddObject|mayEdit)
	if err != nil {
		log.Fatal("Library.AddToTree error: AddEntry failed: %s", err)
	}
	for _, t := range l.SignalTypes() {
		child := tree.Append(cursor)
		t.AddToTree(tree, child)
	}
	for _, t := range l.NodeTypes() {
		child := tree.Append(cursor)
		t.AddToTree(tree, child)
	}
}

func (l *library) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor) {
	switch obj.(type) {
	case SignalType:
		t := obj.(SignalType)
		l.AddSignalType(t)
		cursor.Position = len(l.SignalTypes()) - 1
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	case NodeType:
		t := obj.(NodeType)
		err := l.AddNodeType(t)
		if err != nil {
			log.Fatalf("library.AddNewObject error: AddNodeType failed: %s\n", err)
		}
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	default:
		log.Fatalf("library.AddNewObject error: invalid type %T\n", obj)
	}
	return
}

func (l *library) RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject) {
	parent := tree.Parent(cursor)
	if l != tree.Object(parent) {
		log.Fatal("library.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case SignalType:
		st := tree.Object(cursor).(SignalType)
		for _, nt := range l.NodeTypes() {
			for _, p := range nt.InPorts() {
				if p.SignalType() == st {
					log.Printf(`library.RemoveObject warning:
						SignalType %v is still in use\n`, st)
					return
				}
			}
			for _, p := range nt.OutPorts() {
				if p.SignalType() == st {
					log.Printf(`library.RemoveObject warning:
						SignalType %v is still in use\n`, st)
					return
				}
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		l.RemoveSignalType(st)

	case NodeType:
		nt := tree.Object(cursor).(NodeType)
		if len(nt.Instances()) > 0 {
			log.Printf(`library.RemoveObject warning: The following nodes
				are still instances of NodeType %s:\n`, nt.TypeName())
			for _, n := range nt.Instances() {
				log.Printf("	%s\n", n.NodeName())
			}
			return
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		l.RemoveNodeType(nt)

	default:
		log.Fatalf("library.RemoveObject error: invalid type %T\n", obj)
	}
	return
}
