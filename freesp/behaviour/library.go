package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"log"
)

type library struct {
	filename    string
	pathPrefix  string
	signalTypes signalTypeList
	nodeTypes   nodeTypeList
	context     mod.ModelContextIf
}

var _ bh.LibraryIf = (*library)(nil)

func LibraryNew(filename string, context mod.ModelContextIf) *library {
	ret := &library{filename, "", signalTypeListInit(), nodeTypeListInit(), context}
	freesp.RegisterLibrary(ret)
	return ret
}

func LibraryUsesNodeType(l bh.LibraryIf, nt bh.NodeTypeIf) bool {
	for _, t := range l.NodeTypes() {
		if t.TypeName() == nt.TypeName() {
			return true
		}
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				if SignalGraphTypeUsesNodeType(impl.Graph(), nt) {
					return true
				}
			}
		}
	}
	return false
}

func LibraryUsesSignalType(l bh.LibraryIf, st bh.SignalTypeIf) bool {
	for _, t := range l.SignalTypes() {
		if t.TypeName() == st.TypeName() {
			return true
		}
	}
	for _, t := range l.NodeTypes() {
		for _, impl := range t.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
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
	// TODO: check if new name is existing
	freesp.RemoveRegisteredLibrary(l)
	l.filename = filename
	for _, t := range l.NodeTypes() {
		t.(*nodeType).definedAt = filename
	}
	freesp.RegisterLibrary(l)
}

func (l library) PathPrefix() string {
	return l.pathPrefix
}

func (l *library) SetPathPrefix(newP string) {
	l.pathPrefix = newP
}

func (l library) SignalTypes() []bh.SignalTypeIf {
	return l.signalTypes.SignalTypes()
}

func (l library) NodeTypes() []bh.NodeTypeIf {
	return l.nodeTypes.NodeTypes()
}

func (l *library) createLibFromXml(xmlLib *backend.XmlLibrary) error {
	libMgr := l.context.LibraryMgr()
	for _, ref := range xmlLib.Libraries {
		_, err := libMgr.Access(ref.Name)
		if err != nil {
			log.Println("library.Read warning:", err)
		}
	}
	for _, st := range xmlLib.SignalTypes {
		var scope bh.Scope
		var mode bh.Mode
		switch st.Scope {
		case "local":
			scope = bh.Local
		default:
			scope = bh.Global
		}
		switch st.Mode {
		case "sync":
			mode = bh.Synchronous
		default:
			mode = bh.Asynchronous
		}
		sType, err := SignalTypeNew(st.Name, st.Ctype, st.Msgid, scope, mode, l.Filename())
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

func (l *library) RemoveFromTree(tree tr.TreeIf) {
	tree.Remove(tree.Cursor(l))
	freesp.RemoveRegisteredLibrary(l)
}

func (l *library) AddNodeType(t bh.NodeTypeIf) error {
	nType, ok := freesp.GetNodeTypeByName(t.TypeName())
	if ok {
		log.Printf(`library.AddNodeType: warning: adding existing node
			type definition %s (taking the existing one)`, t.TypeName())
	} else {
		nType = t.(*nodeType)
		freesp.RegisterNodeType(nType)
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

func (l *library) RemoveNodeType(nt bh.NodeTypeIf) {
	for _, n := range nt.(*nodeType).instances.Nodes() {
		n.(*node).context.RemoveNode(n)
	}
	freesp.RemoveRegisteredNodeType(nt)
	l.nodeTypes.Remove(nt)
}

func (l *library) AddSignalType(s bh.SignalTypeIf) (ok bool) {
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

func (l *library) RemoveSignalType(st bh.SignalTypeIf) {
	for _, ntName := range freesp.GetRegisteredNodeTypes() {
		nt, ok := freesp.GetNodeTypeByName(ntName)
		if !ok {
			log.Panicf("library.RemoveSignalType internal error: node type %s\n", ntName)
		}
		for _, p := range nt.InPorts() {
			if p.SignalType().TypeName() == st.TypeName() {
				log.Printf(`library.RemoveSignalType warning:
					bh.SignalTypeIf %v is still in use\n`, st)
				return
			}
		}
		for _, p := range nt.OutPorts() {
			if p.SignalType().TypeName() == st.TypeName() {
				log.Printf(`library.RemoveSignalType warning:
					bh.SignalTypeIf %v is still in use\n`, st)
				return
			}
		}
	}
	SignalTypeDestroy(st)
	l.signalTypes.Remove(st)
}

func (l *library) CreateXml() (buf []byte, err error) {
	xmllib := CreateXmlLibrary(l)
	buf, err = xmllib.Write()
	return
}

var _ tr.TreeElementIf = (*library)(nil)

//
//		tr.TreeElementIf interface
//

func (l *library) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	prop := freesp.PropertyNew(true, false, false)
	err := tree.AddEntry(cursor, tr.SymbolLibrary, l.Filename(), l, prop)
	if err != nil {
		log.Fatalf("bh.LibraryIf.AddToTree error: AddEntry failed: %s\n", err)
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

func (l *library) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElementIf) (newCursor tr.Cursor, err error) {
	if obj == nil {
		err = fmt.Errorf("library.AddNewObject error: nil object")
		return
	}
	switch obj.(type) {
	case bh.SignalTypeIf:
		t := obj.(bh.SignalTypeIf)
		ok := l.AddSignalType(t)
		if !ok {
			err = fmt.Errorf("library.AddNewObject warning: duplicate")
			return
		}
		cursor.Position = len(l.SignalTypes()) - 1
		newCursor = tree.Insert(cursor)
		t.AddToTree(tree, newCursor)

	case bh.NodeTypeIf:
		t := obj.(bh.NodeTypeIf)
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

func (l *library) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	parent := tree.Parent(cursor)
	if l != tree.Object(parent) {
		log.Fatal("library.RemoveObject error: not removing child of mine.")
	}
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.SignalTypeIf:
		st := tree.Object(cursor).(bh.SignalTypeIf)
		for _, nt := range l.NodeTypes() {
			for _, p := range nt.InPorts() {
				if p.SignalType().TypeName() == st.TypeName() {
					log.Printf(`library.RemoveObject warning:
						bh.SignalTypeIf %v is still in use\n`, st)
					return
				}
			}
			for _, p := range nt.OutPorts() {
				if p.SignalType().TypeName() == st.TypeName() {
					log.Printf(`library.RemoveObject warning:
						bh.SignalTypeIf %v is still in use\n`, st)
					return
				}
			}
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, obj})
		l.RemoveSignalType(st)

	case bh.NodeTypeIf:
		nt := tree.Object(cursor).(bh.NodeTypeIf)
		log.Printf("library.RemoveObject: nt=%v\n", nt)
		if len(nt.Instances()) > 0 {
			log.Printf(`library.RemoveObject warning: The following nodes
				are still instances of bh.NodeTypeIf %s:\n`, nt.TypeName())
			for _, n := range nt.Instances() {
				log.Printf("	%s\n", n.Name())
			}
			return
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, obj})
		l.RemoveNodeType(nt)

	default:
		log.Fatalf("library.RemoveObject error: invalid type %T\n", obj)
	}
	return
}

func (l *library) Identify(te tr.TreeElementIf) bool {
	switch te.(type) {
	case *library:
		return te.(*library).Filename() == l.Filename()
	}
	return false
}
