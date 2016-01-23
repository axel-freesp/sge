package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"log"
)

type library struct {
	filename    string
	signalTypes signalTypeList
	nodeTypes   nodeTypeList
	context     ContextIf
}

var _ LibraryIf = (*library)(nil)

func LibraryNew(filename string, context ContextIf) *library {
	ret := &library{filename, signalTypeListInit(), nodeTypeListInit(), context}
	libraries[filename] = ret
	return ret
}

func LibraryUsesNodeType(l LibraryIf, nt NodeTypeIf) bool {
	for _, t := range l.NodeTypes() {
		if t.TypeName() == nt.TypeName() {
			return true
		}
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				if SignalGraphTypeUsesNodeType(impl.Graph(), nt) {
					return true
				}
			}
		}
	}
	return false
}

func LibraryUsesSignalType(l LibraryIf, st SignalType) bool {
	for _, t := range l.SignalTypes() {
		if t.TypeName() == st.TypeName() {
			return true
		}
	}
	for _, t := range l.NodeTypes() {
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == NodeTypeGraph {
				if SignalGraphTypeUsesSignalType(impl.Graph(), st) {
					return true
				}
			}
		}
	}
	return false
}

func (l library) Filename() string {
	return l.filename
}

func (l *library) SetFilename(filename string) {
	delete(libraries, l.filename)
	l.filename = filename
	for _, t := range l.NodeTypes() {
		t.(*nodeType).definedAt = filename
	}
	libraries[filename] = l
}

func (l library) SignalTypes() []SignalType {
	return l.signalTypes.SignalTypes()
}

func (l library) NodeTypes() []NodeTypeIf {
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
		sType, err := SignalTypeNew(st.Name, st.Ctype, st.Msgid, scope, mode)
		if err != nil {
			return err
		}
		l.AddSignalType(sType)
	}
	for _, n := range xmlLib.NodeTypes {
		nType := createNodeTypeFromXml(n, l.Filename(), l.context)
		err := l.AddNodeType(nType)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	return nil
}

func (l *library) Read(data []byte) (cnt int, err error) {
	xmllib := backend.XmlLibraryNew()
	cnt, err = xmllib.Read(data)
	if err != nil {
		err = fmt.Errorf("library.Read: %v", err)
		return
	}
	err = l.createLibFromXml(xmllib)
	return
}

func (l *library) ReadFile(filepath string) error {
	xmllib := backend.XmlLibraryNew()
	err := xmllib.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("library.Read: %v", err)
	}
	return l.createLibFromXml(xmllib)
}

func (l library) Write() (data []byte, err error) {
	xmllib := CreateXmlLibrary(&l)
	data, err = xmllib.Write()
	return
}

func (l library) WriteFile(filepath string) error {
	xmllib := CreateXmlLibrary(&l)
	return xmllib.WriteFile(filepath)
}

func (l *library) RemoveFromTree(tree Tree) {
	tree.Remove(tree.Cursor(l))
	delete(libraries, l.filename)
}

func (l *library) AddNodeType(t NodeTypeIf) error {
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

func (l *library) RemoveNodeType(nt NodeTypeIf) {
	for _, n := range nt.(*nodeType).instances.Nodes() {
		n.(*node).context.RemoveNode(n)
	}
	delete(nodeTypes, nt.TypeName())
	registeredNodeTypes.Remove(nt.TypeName())
	l.nodeTypes.Remove(nt)
}

func (l *library) AddSignalType(s SignalType) (ok bool) {
	for _, st := range l.signalTypes.SignalTypes() {
		if st.TypeName() == s.TypeName() {
			log.Printf(`library.AddSignalType: warning: adding
				duplicate signal type definition %s (ignored)`, s.TypeName())
			return
		}
	}
	ok = true
	l.signalTypes.Append(s)
	return
}

func (l *library) RemoveSignalType(st SignalType) {
	for _, ntName := range registeredNodeTypes.Strings() {
		nt := nodeTypes[ntName]
		for _, p := range nt.InPorts() {
			if p.SignalType().TypeName() == st.TypeName() {
				log.Printf(`library.RemoveSignalType warning:
					SignalType %v is still in use\n`, st)
				return
			}
		}
		for _, p := range nt.OutPorts() {
			if p.SignalType().TypeName() == st.TypeName() {
				log.Printf(`library.RemoveSignalType warning:
					SignalType %v is still in use\n`, st)
				return
			}
		}
	}
	SignalTypeDestroy(st)
	l.signalTypes.Remove(st)
}

var _ TreeElement = (*library)(nil)

//
//		TreeElement interface
//

func (l *library) AddToTree(tree Tree, cursor Cursor) {
	err := tree.AddEntry(cursor, SymbolLibrary, l.Filename(), l, mayAddObject)
	if err != nil {
		log.Fatalf("LibraryIf.AddToTree error: AddEntry failed: %s\n", err)
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

func (l *library) AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("library.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case SignalType:
		t := obj.(SignalType)
		ok := l.AddSignalType(t)
		if !ok {
			err = fmt.Errorf("library.AddNewObject warning: duplicate")
			return
		}
		cursor.Position = len(l.SignalTypes()) - 1
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	case NodeTypeIf:
		t := obj.(NodeTypeIf)
		err = l.AddNodeType(t)
		if err != nil {
			err = fmt.Errorf("library.AddNewObject error: AddNodeType failed: %s", err)
			return
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
				if p.SignalType().TypeName() == st.TypeName() {
					log.Printf(`library.RemoveObject warning:
						SignalType %v is still in use\n`, st)
					return
				}
			}
			for _, p := range nt.OutPorts() {
				if p.SignalType().TypeName() == st.TypeName() {
					log.Printf(`library.RemoveObject warning:
						SignalType %v is still in use\n`, st)
					return
				}
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, IdWithObject{prefix, index, obj})
		l.RemoveSignalType(st)

	case NodeTypeIf:
		nt := tree.Object(cursor).(NodeTypeIf)
		log.Printf("library.RemoveObject: nt=%v\n", nt)
		if len(nt.Instances()) > 0 {
			log.Printf(`library.RemoveObject warning: The following nodes
				are still instances of NodeTypeIf %s:\n`, nt.TypeName())
			for _, n := range nt.Instances() {
				log.Printf("	%s\n", n.Name())
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

func (l *library) Identify(te TreeElement) bool {
	switch te.(type) {
	case *library:
		return te.(*library).Filename() == l.Filename()
	}
	return false
}
