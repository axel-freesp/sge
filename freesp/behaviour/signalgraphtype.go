package behaviour

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/freesp"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	fd "github.com/axel-freesp/sge/interface/filedata"
	gr "github.com/axel-freesp/sge/interface/graph"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"image"
	"log"
	"strings"
)

type signalGraphType struct {
	context                                  mod.ModelContextIf
	libraries                                []bh.LibraryIf
	nodes                                    nodeList
	inputNodes, outputNodes, processingNodes []bh.NodeIf
}

/*
 *  freesp.bh.SignalGraphTypeIf API
 */

var _ bh.SignalGraphTypeIf = (*signalGraphType)(nil)

func SignalGraphTypeNew(context mod.ModelContextIf) *signalGraphType {
	return &signalGraphType{context, nil, nodeListInit(), nil, nil, nil}
}

func SignalGraphTypeUsesNodeType(t bh.SignalGraphTypeIf, nt bh.NodeTypeIf) bool {
	for _, n := range t.Nodes() {
		if n.ItsType().TypeName() == nt.TypeName() {
			return true
		}
		for _, impl := range n.ItsType().Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				if SignalGraphTypeUsesNodeType(impl.Graph(), nt) {
					return true
				}
			}
		}
	}
	return false
}

func SignalGraphTypeUsesSignalType(t bh.SignalGraphTypeIf, st bh.SignalTypeIf) bool {
	for _, n := range t.Nodes() {
		for _, p := range n.InPorts() {
			if p.SignalType() == st {
				return true
			}
		}
		for _, p := range n.OutPorts() {
			if p.SignalType() == st {
				return true
			}
		}
		for _, impl := range n.ItsType().Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				if SignalGraphTypeUsesSignalType(impl.Graph(), st) {
					return true
				}
			}
		}
	}
	return false
}

func (t *signalGraphType) Nodes() []bh.NodeIf {
	return t.nodes.Nodes()
}

func (t *signalGraphType) NodeByName(name string) (n bh.NodeIf, ok bool) {
	for _, n = range t.Nodes() {
		if n.Name() == name {
			ok = true
			return
		}
	}
	return
}

func (t *signalGraphType) Libraries() []bh.LibraryIf {
	return t.libraries
}

func (t *signalGraphType) InputNodes() []bh.NodeIf {
	return t.inputNodes
}

func (t *signalGraphType) OutputNodes() []bh.NodeIf {
	return t.outputNodes
}

func (t *signalGraphType) ProcessingNodes() []bh.NodeIf {
	return t.processingNodes
}

func (t *signalGraphType) AddNode(n bh.NodeIf) error {
	nType := n.ItsType()
	if !isAutoType(nType) {
		libname := nType.DefinedAt()
		if len(libname) == 0 {
			return fmt.Errorf("signalGraphType.AddNode error: node type %s has no DefinedAt...", nType.TypeName())
		}
		if !t.containsLibRef(libname) {
			lib, ok := freesp.GetLibraryByName(libname)
			if !ok {
				return fmt.Errorf("signalGraphType.AddNode error: library %s not registered", libname)
			}
			t.libraries = append(t.libraries, lib)
		}
	}
	return t.addNode(n)
}

func (t *signalGraphType) RemoveNode(n bh.NodeIf) {
	for _, p := range n.(*node).inPort.Ports() {
		for _, c := range p.Connections() {
			c.RemoveConnection(p)
		}
	}
	t.nodes.Remove(n)
	RemNode(&t.inputNodes, n.(*node))
	RemNode(&t.outputNodes, n.(*node))
	RemNode(&t.processingNodes, n.(*node))
	n.ItsType().(*nodeType).removeInstance(n.(*node))
}

func (t *signalGraphType) Context() mod.ModelContextIf {
	return t.context
}

func (t *signalGraphType) containsLibRef(libname string) bool {
	for _, l := range t.libraries {
		if l.Filename() == libname {
			return true
		}
	}
	return false
}

func FindNode(list []bh.NodeIf, elem bh.NodeIf) (index int, ok bool) {
	for index = 0; index < len(list); index++ {
		if elem == list[index] {
			break
		}
	}
	ok = (index < len(list))
	return
}

func RemNode(list *[]bh.NodeIf, elem bh.NodeIf) {
	index, ok := FindNode(*list, elem)
	if !ok {
		return
	}
	for j := index + 1; j < len(*list); j++ {
		(*list)[j-1] = (*list)[j]
	}
	(*list) = (*list)[:len(*list)-1]
}

func (t *signalGraphType) addNode(n bh.NodeIf) error {
	if len(n.InPorts()) > 0 {
		if len(n.OutPorts()) > 0 {
			t.processingNodes = append(t.processingNodes, n.(*node))
		} else {
			t.outputNodes = append(t.outputNodes, n.(*node))
		}
	} else {
		if len(n.OutPorts()) > 0 {
			t.inputNodes = append(t.inputNodes, n.(*node))
		} else {
			return fmt.Errorf("signalGraphType.AddNode error: node has no ports")
		}
	}
	t.nodes.Append(n.(*node))
	return nil
}

func createSignalGraphTypeFromXml(g *backend.XmlSignalGraph, name string, context mod.ModelContextIf,
	resolvePort func(portname string, dir gr.PortDirection) *portType) (t *signalGraphType, err error) {
	t = SignalGraphTypeNew(context)
	for _, ref := range g.Libraries {
		l, ok := freesp.GetLibraryByName(ref.Name)
		if !ok {
			var f fd.FileDataIf
			f, err = t.context.LibraryMgr().Access(ref.Name)
			if err != nil {
				err = fmt.Errorf("createSignalGraphTypeFromXml error: referenced library file %s not found", ref.Name)
				return
			}
			l = f.(bh.LibraryIf)
			freesp.RegisterLibrary(l)
			fmt.Println("createSignalGraphTypeFromXml: library", ref.Name, "successfully loaded")
		}
		log.Printf("createSignalGraphTypeFromXml: ok=%v, l=%v\n", ok, l)
		t.libraries = append(t.libraries, l)
	}
	for _, n := range g.InputNodes {
		var nnode *node
		nnode, err = t.createInputNodeFromXml(n, resolvePort)
		if err != nil {
			return
		}
		t.inputNodes = append(t.inputNodes, nnode)
		t.nodes.Append(nnode)
	}
	for _, n := range g.OutputNodes {
		var nnode *node
		nnode, err = t.createOutputNodeFromXml(n, resolvePort)
		if err != nil {
			return
		}
		t.outputNodes = append(t.outputNodes, nnode)
		t.nodes.Append(nnode)
	}
	for _, n := range g.ProcessingNodes {
		nnode := t.createNodeFromXml(n.XmlNode)
		t.processingNodes = append(t.processingNodes, nnode)
		t.nodes.Append(nnode)
	}
	for i, c := range g.Connections {
		n1, ok := t.NodeByName(c.From)
		if !ok {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.From, dump))
		}
		n2, ok := t.NodeByName(c.To)
		if !ok {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d: node %s not found\n%s", i, c.To, dump))
		}
		p1, err := n1.(*node).outPortFromName(c.FromPort)
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Printf("edge = %v\n", c)
			log.Printf("node = %v, missing port = %s\n", n1, c.FromPort)
			log.Fatal(fmt.Sprintf("invalid edge %d outPortFromName failed: %s\n%s", i, err, dump))
		}
		p2, err := n2.(*node).inPortFromName(c.ToPort)
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d inPortFromName failed: %s\n%s", i, err, dump))
		}
		err = p1.AddConnection(ConnectionNew(p1, p2))
		if err != nil {
			dump, _ := g.Write()
			log.Println("createSignalGraphTypeFromXml error:")
			log.Fatal(fmt.Sprintf("invalid edge %d PortConnect failed: %s\n%s", i, err, dump))
		}
	}
	return
}

func createNodeTypeName(n backend.XmlNode) string {
	ntName := n.NType
	if len(ntName) == 0 {
		ntName = fmt.Sprintf("autoTypeOfNode-%s", n.NName)
	}
	return ntName
}

func createInputNodeTypeName(name string) string {
	return fmt.Sprintf("autoInputNodeType-%s", name)
}

func createOutputNodeTypeName(name string) string {
	return fmt.Sprintf("autoOutputNodeType-%s", name)
}

func isAutoType(nt bh.NodeTypeIf) bool {
	if strings.HasPrefix(nt.TypeName(), "autoInputNodeType-") {
		return true
	}
	if strings.HasPrefix(nt.TypeName(), "autoOutputNodeType-") {
		return true
	}
	return false
}

func (t *signalGraphType) createNodeFromXml(xmln backend.XmlNode) (nd *node) {
	nName := xmln.NName
	ntName := xmln.NType
	if len(ntName) == 0 {
		ntName = createNodeTypeName(xmln)
	}
	nt, ok := freesp.GetNodeTypeByName(ntName)
	if !ok {
		nt = createNodeTypeFromXmlNode(xmln, ntName)
	}
	var err error
	nd, err = NodeNew(nName, nt, t)
	if err != nil {
		log.Fatal("signalGraphType.createNodeFromXml: TODO: error handling")
	}
	for _, xmlh := range xmln.Entry {
		mode, ok := freesp.ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("signalGraphType.createNodeFromXml Warning: hint mode %s not defined\n",
				xmlh.Mode)
			continue
		}
		nd.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
	}
	nd.expanded = xmln.Expanded
	for _, p := range nd.InPorts() {
		pname := p.Name()
		for _, xmlp := range xmln.InPort {
			if xmlp.PName == pname {
				for _, m := range xmlp.Entry {
					p.SetModePosition(freesp.ModeFromString[m.Mode], image.Point{m.X, m.Y})
				}
				break
			}
		}
	}
	for _, p := range nd.OutPorts() {
		pname := p.Name()
		for _, xmlp := range xmln.OutPort {
			if xmlp.PName == pname {
				for _, m := range xmlp.Entry {
					p.SetModePosition(freesp.ModeFromString[m.Mode], image.Point{m.X, m.Y})
				}
				break
			}
		}
	}
	if nd.expanded {
		for _, impl := range nt.Implementation() {
			if impl.ImplementationType() == bh.NodeTypeGraph {
				for i, chn := range impl.Graph().ProcessingNodes() {
					xmlch := xmln.Children[i]
					for _, xmlh := range xmlch.Entry {
						mode, ok := freesp.ModeFromString[xmlh.Mode]
						if !ok {
							log.Printf("signalGraphType.createNodeFromXml Warning: hint mode %s not defined\n",
								xmlh.Mode)
							continue
						}
						chn.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
					}
				}
				break
			}
		} 
	}
	return
}

func (t *signalGraphType) createInputNodeFromXml(n backend.XmlInputNode,
	resolvePort func(portname string, dir gr.PortDirection) *portType) (ret *node, err error) {
	nName := n.NName
	ntName := createInputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	ret, err = NodeNew(nName, nt, t)
	if err != nil {
		err = fmt.Errorf("signalGraphType.createInputNodeFromXml: %s", err)
		return
	}
	pt := resolvePort(n.NPort, gr.InPort)
	if pt != nil {
		ret.portlink = pt
	}
	for _, xmlh := range n.Entry {
		mode, ok := freesp.ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("signalGraphType.createNodeFromXml Warning: hint mode %s not defined\n",
				xmlh.Mode)
			continue
		}
		ret.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
	}
	return
}

func (t *signalGraphType) createOutputNodeFromXml(n backend.XmlOutputNode,
	resolvePort func(portname string, dir gr.PortDirection) *portType) (ret *node, err error) {
	nName := n.NName
	ntName := createOutputNodeTypeName(nName)
	nt := createNodeTypeFromXmlNode(n.XmlNode, ntName)
	ret, err = NodeNew(nName, nt, t)
	if err != nil {
		err = fmt.Errorf("signalGraphType.createOutputNodeFromXml: %s", err)
		return
	}
	pt := resolvePort(n.NPort, gr.OutPort) // matches also empty names
	if pt != nil {
		ret.portlink = pt
	}
	for _, xmlh := range n.Entry {
		mode, ok := freesp.ModeFromString[xmlh.Mode]
		if !ok {
			log.Printf("signalGraphType.createNodeFromXml Warning: hint mode %s not defined\n",
				xmlh.Mode)
			continue
		}
		ret.SetModePosition(mode, image.Point{xmlh.X, xmlh.Y})
	}
	return
}

func (g *signalGraphType) addInputNodeFromPortType(p bh.PortTypeIf) {
	st := p.SignalType()
	ntName := createInputNodeTypeName(st.TypeName())
	nt, ok := freesp.GetNodeTypeByName(ntName)
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.(*nodeType).addOutPort("", st)
		freesp.RegisterNodeType(nt)
	}
	if len(nt.(*nodeType).outPorts.PortTypes()) == 0 {
		log.Fatal("signalGraphType.addInputNodeFromNamedPortType: invalid input node type")
	}
	n, err := NodeNew(fmt.Sprintf("in-%s", p.Name()), nt, g)
	if err != nil {
		log.Fatal("signalGraphType.addInputNodeFromPortType: TODO: error handling")
	}
	n.portlink = p
	err = g.addNode(n)
	if err != nil {
		log.Fatal("signalGraphType.addInputNodeFromNamedPortType: AddNode failed:", err)
	}
}

func (g *signalGraphType) addOutputNodeFromPortType(p bh.PortTypeIf) {
	st := p.SignalType()
	ntName := createOutputNodeTypeName(st.TypeName())
	nt, ok := freesp.GetNodeTypeByName(ntName)
	if !ok {
		nt = NodeTypeNew(ntName, "")
		nt.(*nodeType).addInPort("", st)
		freesp.RegisterNodeType(nt)
	}
	if len(nt.(*nodeType).inPorts.PortTypes()) == 0 {
		log.Fatal("signalGraphType.addOutputNodeFromNamedPortType: invalid output node type")
	}
	n, err := NodeNew(fmt.Sprintf("out-%s", p.Name()), nt, g)
	if err != nil {
		log.Fatal("signalGraphType.addOutputNodeFromPortType: TODO: error handling")
	}
	n.portlink = p
	err = g.addNode(n)
	if err != nil {
		log.Fatal("signalGraphType.addOutputNodeFromNamedPortType: AddNode failed:", err)
	}
}

func (g *signalGraphType) removeInputNodeFromPortType(p bh.PortTypeIf) {
	for _, n := range g.InputNodes() {
		if n.Name() == fmt.Sprintf("in-%s", p.Name()) {
			g.RemoveNode(n)
			return
		}
	}
}

func (g *signalGraphType) removeOutputNodeFromPortType(p bh.PortTypeIf) {
	for _, n := range g.OutputNodes() {
		if n.Name() == fmt.Sprintf("out-%s", p.Name()) {
			g.RemoveNode(n)
			return
		}
	}
}

func (g *signalGraphType) findInputNodeFromPortType(p bh.PortTypeIf) bh.NodeIf {
	for _, n := range g.InputNodes() {
		if n.Name() == fmt.Sprintf("in-%s", p.Name()) {
			return n
		}
	}
	return nil
}

func (g *signalGraphType) findOutputNodeFromPortType(p bh.PortTypeIf) bh.NodeIf {
	for _, n := range g.OutputNodes() {
		if n.Name() == fmt.Sprintf("out-%s", p.Name()) {
			return n
		}
	}
	return nil
}

func (t *signalGraphType) CreateXml() (buf []byte, err error) {
	xmlsignalgraph := CreateXmlSignalGraphType(t)
	buf, err = xmlsignalgraph.Write()
	return
}

/*
 *  tr.TreeElement API
 */

var _ tr.TreeElement = (*signalGraphType)(nil)

func (t *signalGraphType) AddToTree(tree tr.TreeIf, cursor tr.Cursor) {
	for _, n := range t.InputNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
	for _, n := range t.OutputNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
	for _, n := range t.ProcessingNodes() {
		child := tree.Append(cursor)
		n.AddToTree(tree, child)
	}
}

func (t *signalGraphType) treeAddNewObject(tree tr.TreeIf, cursor tr.Cursor, n bh.NodeIf) (newCursor tr.Cursor) {
	newCursor = tree.Insert(cursor)
	n.AddToTree(tree, newCursor)
	return
}

func (t *signalGraphType) AddNewObject(tree tr.TreeIf, cursor tr.Cursor, obj tr.TreeElement) (newCursor tr.Cursor, err error) {
	switch obj.(type) {
	case bh.NodeIf:
		// TODO: Check if IO node and exists: copy position only and return
		n := obj.(bh.NodeIf)
		err = t.AddNode(n)
		if err != nil {
			err = fmt.Errorf("signalGraphType.AddNewObject error: %s", err)
			nt := n.ItsType().(*nodeType)
			if nt != nil {
				ok, _ := nt.instances.Find(n)
				if ok {
					nt.instances.Remove(n)
				}
			}
			return
		}
		newCursor = t.treeAddNewObject(tree, cursor, n)

		parent := tree.Object(cursor)
		switch parent.(type) {
		case bh.SignalGraphIf:
		case bh.ImplementationIf:
			// propagate new node to all instances of embracing type
			pCursor := tree.Parent(cursor)
			nt := tree.Object(pCursor)
			for _, nn := range nt.(bh.NodeTypeIf).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, parent)
				tCursor.Position = cursor.Position
				t.treeAddNewObject(tree, tCursor, n)
			}

		default:
			log.Fatalf("signalGraphType.AddNewObject error: wrong parent type %T: %v\n", parent, parent)
		}

	case bh.ConnectionIf:
		conn := obj.(bh.ConnectionIf)
		var n bh.NodeIf
		var p bh.PortIf
		for _, n = range t.Nodes() {
			if n.Name() == conn.From().Node().Name() {
				nCursor := tree.CursorAt(cursor, n)
				for _, p = range n.OutPorts() {
					if conn.From().Name() == p.Name() {
						pCursor := tree.CursorAt(nCursor, p)
						return p.AddNewObject(tree, pCursor, obj)
					}
				}
			}
		}
	default:
		log.Fatalf("signalGraphType.AddNewObject error: wrong type %t: %v\n", obj, obj)
	}
	return
}

func (t *signalGraphType) RemoveObject(tree tr.TreeIf, cursor tr.Cursor) (removed []tr.IdWithObject) {
	obj := tree.Object(cursor)
	switch obj.(type) {
	case bh.NodeIf:
		n := obj.(bh.NodeIf)
		// Remove all connections first
		for _, p := range n.OutPorts() {
			for _, c := range p.Connections() {
				conn := p.Connection(c)
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		for _, p := range n.InPorts() {
			for _, c := range p.Connections() {
				conn := p.Connection(c)
				cCursor := tree.CursorAt(cursor, conn)
				del := p.RemoveObject(tree, cCursor)
				for _, d := range del {
					removed = append(removed, d)
				}
			}
		}
		parentCursor := tree.Parent(cursor)
		parent := tree.Object(parentCursor)
		switch parent.(type) {
		case bh.SignalGraphIf:
		case bh.ImplementationIf:
			// propagate new node to all instances of embracing type
			pCursor := tree.Parent(parentCursor)
			nt := tree.Object(pCursor)
			for _, nn := range nt.(bh.NodeTypeIf).Instances() {
				nCursor := tree.Cursor(nn)
				tCursor := tree.CursorAt(nCursor, parent)
				tree.Remove(tree.CursorAt(tCursor, n))
			}

		default:
			log.Fatalf("signalGraphType.RemoveObject error: wrong parent type %t: %v\n", parent, parent)
		}
		prefix, index := tree.Remove(cursor)
		removed = append(removed, tr.IdWithObject{prefix, index, obj})
		t.RemoveNode(n)

	default:
		log.Fatalf("signalGraphType.RemoveObject error: wrong type %t: %v", obj, obj)
	}
	return
}
