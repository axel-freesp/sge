package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

const envIconPathKey = "SGE_ICON_PATH"
const (
	iconCol = 0
	textCol = 1
)

var (
	imageSignalGraph    *gdk.Pixbuf = nil
	imageInputNode      *gdk.Pixbuf = nil
	imageOutputNode     *gdk.Pixbuf = nil
	imageProcessingNode *gdk.Pixbuf = nil
	imageNodeType       *gdk.Pixbuf = nil
	imageNodeTypeElem   *gdk.Pixbuf = nil
	imageInputPort      *gdk.Pixbuf = nil
	imageOutputPort     *gdk.Pixbuf = nil
	imageSignalType     *gdk.Pixbuf = nil
	imageConnected      *gdk.Pixbuf = nil
	imageLibrary        *gdk.Pixbuf = nil
)

func Init() error {
	iconPath := os.Getenv(envIconPathKey)
	if len(iconPath) == 0 {
		log.Fatal("Missing Environment variable ", iconPath)
	}
	var err error
	imageSignalGraph, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test1.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test1.png", err)
	}
	imageInputNode, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/input.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageOutputNode, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/output.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageProcessingNode, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/node.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test2.png", err)
	}
	imageNodeType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/node-type.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test2.png", err)
	}
	imageNodeTypeElem, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test0.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test2.png", err)
	}
	imageInputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/inport-green.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageOutputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/outport-red.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageSignalType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/signal-type.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageConnected, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/link.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageLibrary, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test1.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	fmt.Println("models.Init: all icons intitialized")
	return nil
}

type FilesTreeStore struct {
	treestore        *gtk.TreeStore
	lookup           map[string]interface{}
	libs             map[string]freesp.Library
	currentSelection *gtk.TreeIter
}

func (s *FilesTreeStore) TreeStore() *gtk.TreeStore {
	return s.treestore
}

func FilesTreeStoreNew() (ret *FilesTreeStore, err error) {
	ts, err := gtk.TreeStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING)
	if err != nil {
		err = gtkErr("FilesTreeStoreNew", "TreeStoreNew", err)
		ret = nil
		return
	}
	ret = &FilesTreeStore{ts, make(map[string]interface{}),
		make(map[string]freesp.Library),
		nil}
	err = nil
	return
}

func (s *FilesTreeStore) GetCurrentId() (id string) {
	p, err := s.treestore.GetPath(s.currentSelection)
	if err != nil {
		id = ""
		return
	}
	id = p.String()
	return
}

func (s *FilesTreeStore) GetObjectById(id string) (ret interface{}, err error) {
	_, err = s.treestore.GetIterFromString(id)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetObjectById", "GetIterFromString()", err)
		return
	}
	ret = s.lookup[id]
	return
}

func (s *FilesTreeStore) GetValueById(id string) (ret string, err error) {
	iter, err := s.treestore.GetIterFromString(id)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetValueById", "GetIterFromString()", err)
		return
	}
	ret, err = s.GetValue(iter)
	return
}

func (s *FilesTreeStore) SetValueById(id, value string) (err error) {
	iter, err := s.treestore.GetIterFromString(id)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetValueById", "GetIterFromString()", err)
		return
	}
	err = s.setValue(value, iter)
	return
}

// Returns the string shown in textCol
func (s *FilesTreeStore) GetValue(iter *gtk.TreeIter) (ret string, err error) {
	v, err := s.treestore.GetValue(iter, textCol)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetValue", "GetValue", err)
		return
	}
	ret, err = v.GetString()
	if err != nil {
		err = gtkErr("FilesTreeStore.GetValue", "GetString", err)
		return
	}
	//s.currentSelection = iter
	return
}

// Returns the object stored with the row
func (s *FilesTreeStore) GetObject(iter *gtk.TreeIter) (ret interface{}, err error) {
	var ok bool
	ret, ok = s.getObject(iter)
	if !ok {
		err = fmt.Errorf("FilesTreeStore.GetObject: not found.")
		return
	}
	s.currentSelection = iter
	return
}

func (s *FilesTreeStore) AddSignalGraphFile(filename string, graph freesp.SignalGraph) (newId string, err error) {
	ts := s.treestore
	iter := ts.Append(nil)
	err = s.addEntry(iter, imageSignalGraph, filename, graph)
	if err != nil {
		err = gtkErr("FilesTreeStore.AddSignalGraphFile", "ts.addEntry(filename)", err)
		return
	}
	err = s.addSignalGraph(iter, graph)
	if err != nil {
		return
	}
	newId, err = s.getIdFromIter(iter)
	return
}

func (s *FilesTreeStore) AddLibraryFile(filename string, lib freesp.Library) (newId string, err error) {
	iter := s.treestore.Append(nil)
	err = s.addLibrary(iter, lib)
	if err != nil {
		return
	}
	newId, err = s.getIdFromIter(iter)
	return
}

func (s *FilesTreeStore) AddNewObject(parentId string, obj interface{}) (newId string, err error) {
	parent, parentIter, err := s.getObjAndIterById(parentId)
	if err != nil {
		return
	}
	if !checkParentType(obj, parent) {
		err = fmt.Errorf("FilesTreeStore.AddNewObject(%T): wrong parent type %T", obj, parent)
		return
	}
	err = doAddFreespObject(parent, obj)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.AddNewObject: doAddFreespObject failed:", err)
		return
	}

	iter := s.treestore.Append(parentIter)

	err = s.doAddTreeObject(iter, obj, parent)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.AddNewObject: doAddTreeObject failed:", err)
		return
	}

	newId, err = s.getIdFromIter(iter)
	return
}

/*
 *      Local functions
 */

func (s *FilesTreeStore) setValue(value string, iter *gtk.TreeIter) (err error) {
	err = s.treestore.SetValue(iter, textCol, value)
	if err != nil {
		err = gtkErr("FilesTreeStore.setValue", "setValue", err)
		return
	}
	return
}

func (s *FilesTreeStore) getIdFromIter(iter *gtk.TreeIter) (newId string, err error) {
	newP, err := s.treestore.GetPath(iter)
	if err != nil {
		return
	}
	newId = newP.String()
	return
}

func doAddFreespObject(parent interface{}, obj interface{}) (err error) {
	err = nil
	switch obj.(type) {
	case freesp.Implementation:
		parent.(freesp.NodeType).AddImplementation(obj.(freesp.Implementation))
	case freesp.SignalType:
		err = parent.(freesp.Library).AddSignalType(obj.(freesp.SignalType))
	case freesp.NamedPortType:
		parent.(freesp.NodeType).AddNamedPortType(obj.(freesp.NamedPortType))
	case freesp.NodeType:
		parent.(freesp.Library).AddNodeType(obj.(freesp.NodeType))
	case freesp.Node:
		switch parent.(type) {
		case freesp.SignalGraph:
			err = parent.(freesp.SignalGraph).ItsType().AddNode(obj.(freesp.Node))
		case freesp.SignalGraphType:
			err = parent.(freesp.SignalGraphType).AddNode(obj.(freesp.Node))
		}
	case freesp.Port:
		err = parent.(freesp.Port).AddConnection(obj.(freesp.Port))
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddFreespObject(1): %v", err)
			return
		}
		err = obj.(freesp.Port).AddConnection(parent.(freesp.Port))
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddFreespObject(2): %v", err)
			return
		}
	}
	return
}

func (s *FilesTreeStore) doAddTreeObject(iter *gtk.TreeIter, obj interface{}, parent interface{}) (err error) {
	var icon *gdk.Pixbuf
	err = nil
	switch obj.(type) {
	case freesp.Implementation:
		err = s.addImplementation(iter, obj.(freesp.Implementation))
		if err != nil {
			return
		}
	case freesp.SignalType:
		err = s.addSignalType(iter, obj.(freesp.SignalType))
	case freesp.NamedPortType:
		if obj.(freesp.NamedPortType).Direction() == freesp.InPort {
			icon = imageInputPort
		} else {
			icon = imageOutputPort
		}
		err = s.addNamedPortType(iter, obj.(freesp.NamedPortType), icon)
	case freesp.NodeType:
		err = s.addNodeType(iter, obj.(freesp.NodeType))
	case freesp.Node:
		if len(obj.(freesp.Node).InPorts()) > 0 {
			if len(obj.(freesp.Node).OutPorts()) > 0 {
				icon = imageProcessingNode
			} else {
				icon = imageOutputNode
			}
		} else {
			icon = imageInputNode
		}
		err = s.addNode(iter, obj.(freesp.Node), icon)
	case freesp.Port:
		err = s.addConnection(iter, parent.(freesp.Port), obj.(freesp.Port))
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddTreeObject: addConnection(1) failed")
			return
		}
		var parentIter2 *gtk.TreeIter
		parentIter2, err = s.getIterFromObject(obj.(freesp.Port))
		if err != nil {
			return
		}
		iter = s.treestore.Append(parentIter2)
		err = s.addConnection(iter, obj.(freesp.Port), parent.(freesp.Port))
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddTreeObject: addConnection(2) failed")
			return
		}
	}
	return
}

func checkParentType(obj interface{}, parent interface{}) bool {
	switch obj.(type) {
	case freesp.Implementation:
		switch parent.(type) {
		case freesp.NodeType:
			return true
		default:
			return false
		}
	case freesp.SignalType:
		switch parent.(type) {
		case freesp.Library:
			return true
		default:
			return false
		}
	case freesp.NamedPortType:
		switch parent.(type) {
		case freesp.NodeType:
			return true
		default:
			return false
		}
	case freesp.NodeType:
		switch parent.(type) {
		case freesp.Library:
			return true
		default:
			return false
		}
	case freesp.Node:
		switch parent.(type) {
		case freesp.SignalGraphType, freesp.SignalGraph:
			return true
		default:
			return false
		}
	case freesp.Port:
		switch parent.(type) {
		case freesp.Port:
			return true
		default:
			return false
		}
	}
	return false
}

func (s *FilesTreeStore) getObjAndIterById(id string) (obj interface{}, iter *gtk.TreeIter, err error) {
	obj, err = s.GetObjectById(id)
	if err != nil {
		return
	}
	iter, err = s.treestore.GetIterFromString(id)
	return
}

func (s *FilesTreeStore) getIterFromObject(obj interface{}) (iter *gtk.TreeIter, err error) {
	var id string
	var o interface{}
	ok := false
	for i := 0; !ok; i++ {
		id = fmt.Sprintf("%d", i)
		o, ok = s.lookup[id]
		if !ok {
			break
		} else if o != obj {
			id, ok = s.getIdFromObjectRecursive(id, obj)
		} else {
			fmt.Println("match: ", ok, id)
		}
	}
	if !ok {
		err = fmt.Errorf("GetIterFromObject error: object not found.")
		return
	}
	iter, err = s.treestore.GetIterFromString(id)
	return
}

func (s *FilesTreeStore) getIdFromObjectRecursive(parent string, obj interface{}) (id string, ok bool) {
	ok = false
	var o interface{}
	for i := 0; !ok; i++ {
		id = fmt.Sprintf("%s:%d", parent, i)
		o, ok = s.lookup[id]
		if !ok {
			break
		} else if o != obj {
			id, ok = s.getIdFromObjectRecursive(id, obj)
		} else {
			fmt.Println("match: ", id)
		}
	}
	return
}

func (s *FilesTreeStore) getObject(iter *gtk.TreeIter) (obj interface{}, ok bool) {
	ok = false
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		return
	}
	obj, ok = s.lookup[path.String()]
	return
}

func (s *FilesTreeStore) addEntry(iter *gtk.TreeIter, icon *gdk.Pixbuf, text string, data interface{}) error {
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "GetPath", err)
		return err
	}
	s.lookup[path.String()] = data
	err = s.treestore.SetValue(iter, iconCol, icon)
	if err != nil {
		return gtkErr("FilesTreeStore.addEntry", "SetValue(iconCol)", err)
	}
	err = s.treestore.SetValue(iter, textCol, text)
	if err != nil {
		return gtkErr("FilesTreeStore.addEntry", "SetValue(textCol)", err)
	}
	return nil
}

func (s *FilesTreeStore) addSignalGraph(iter *gtk.TreeIter, graph freesp.SignalGraph) error {
	return s.addSignalGraphType(iter, graph.ItsType())
}

func (s *FilesTreeStore) addSignalGraphType(iter *gtk.TreeIter, gtype freesp.SignalGraphType) error {
	var err error
	err = s.addNodes(iter, gtype.InputNodes(), imageInputNode)
	if err != nil {
		return err
	}
	err = s.addNodes(iter, gtype.OutputNodes(), imageOutputNode)
	if err != nil {
		return err
	}
	err = s.addNodes(iter, gtype.ProcessingNodes(), imageProcessingNode)
	if err != nil {
		return err
	}
	return nil
}

func (s *FilesTreeStore) addLibrary(iter *gtk.TreeIter, lib freesp.Library) error {
	ts := s.treestore
	err := s.addEntry(iter, imageLibrary, lib.Filename(), lib)
	if err != nil {
		return gtkErr("FilesTreeStore.addLibrary", "ts.addEntry(filename)", err)
	}
	if s.libs[lib.Filename()] == nil {
		s.libs[lib.Filename()] = lib
		i := ts.Append(iter)
		err = s.addEntry(i, imageSignalType, "Signal Types", nil)
		if err != nil {
			return gtkErr("FilesTreeStore.addLibrary", "ts.addEntry(SignalTypes dummy)", err)
		}
		for _, st := range lib.SignalTypes() {
			j := ts.Append(i)
			err = s.addSignalType(j, st)
			if err != nil {
				return gtkErr("FilesTreeStore.addLibrary", "ts.addSignalType(SignalTypes)", err)
			}
		}
		err = s.addNodeTypes(iter, lib.NodeTypes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addSignalType(iter *gtk.TreeIter, sType freesp.SignalType) error {
	err := s.addEntry(iter, imageSignalType, sType.TypeName(), sType)
	if err != nil {
		return gtkErr("FilesTreeStore.addSignalType", "ts.addEntry()", err)
	}
	return nil
}

func (s *FilesTreeStore) addNodeTypes(iter *gtk.TreeIter, types []freesp.NodeType) error {
	ts := s.treestore
	for _, t := range types {
		i := ts.Append(iter)
		err := s.addNodeType(i, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addNodeType(iter *gtk.TreeIter, ntype freesp.NodeType) error {
	ts := s.treestore
	err := s.addEntry(iter, imageNodeType, ntype.TypeName(), ntype)
	if err != nil {
		return gtkErr("FilesTreeStore.addNodeType", "ts.addEntry()", err)
	}
	for _, impl := range ntype.Implementation() {
		j := ts.Append(iter)
		err = s.addImplementation(j, impl)
		if err != nil {
			return err
		}
	}
	if len(ntype.InPorts()) > 0 {
		err = s.addNamedPortTypes(iter, ntype.InPorts(), imageInputPort)
		if err != nil {
			return err
		}
	}
	if len(ntype.OutPorts()) > 0 {
		err = s.addNamedPortTypes(iter, ntype.OutPorts(), imageOutputPort)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addImplementation(iter *gtk.TreeIter, impl freesp.Implementation) (err error) {
	if impl.ImplementationType() == freesp.NodeTypeElement {
		err = s.addEntry(iter, imageNodeTypeElem, impl.ElementName(), impl)
	} else {
		if impl.Graph() == nil {
			err = fmt.Errorf("FilesTreeStore.addImplementation: missing Graph().")
			return
		}
		err = s.addEntry(iter, imageSignalGraph, "Sub-graph", impl.Graph())
		if err != nil {
			return
		}
		err = s.addSignalGraphType(iter, impl.Graph())
	}
	return
}

func (s *FilesTreeStore) addNamedPortTypes(iter *gtk.TreeIter, ports []freesp.NamedPortType, kind *gdk.Pixbuf) error {
	for _, p := range ports {
		i := s.treestore.Append(iter)
		err := s.addNamedPortType(i, p, kind)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addNamedPortType(iter *gtk.TreeIter, p freesp.NamedPortType, kind *gdk.Pixbuf) error {
	err := s.addEntry(iter, kind, p.Name(), p)
	if err != nil {
		return gtkErr("FilesTreeStore.addNamedPortType", "s.addEntry(1)", err)
	}
	j := s.treestore.Append(iter)
	err = s.addEntry(j, imageSignalType, p.TypeName(), p.SignalType())
	if err != nil {
		return gtkErr("FilesTreeStore.addNamedPortType", "s.addEntry(2)", err)
	}
	return nil
}

func (s *FilesTreeStore) addNodes(iter *gtk.TreeIter, nodes []freesp.Node, kind *gdk.Pixbuf) error {
	for _, n := range nodes {
		nodeIter := s.treestore.Append(iter)
		err := s.addNode(nodeIter, n, kind)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addNode(nodeIter *gtk.TreeIter, n freesp.Node, kind *gdk.Pixbuf) (err error) {
	ts := s.treestore
	err = s.addEntry(nodeIter, kind, n.NodeName(), n)
	if err != nil {
		err = gtkErr("FilesTreeStore.addNode", "ts.addEntry(1)", err)
		return
	}
	j := ts.Append(nodeIter)
	t := n.ItsType()
	err = s.addNodeType(j, t)
	if err != nil {
		err = gtkErr("FilesTreeStore.addNode", "ts.addEntry(2)", err)
		return
	}
	if len(n.InPorts()) > 0 {
		err = s.addPorts(nodeIter, n.InPorts(), imageInputPort)
		if err != nil {
			return
		}
	}
	if len(n.OutPorts()) > 0 {
		err = s.addPorts(nodeIter, n.OutPorts(), imageOutputPort)
		if err != nil {
			return
		}
	}
	return
}

func (s *FilesTreeStore) addConnection(iter *gtk.TreeIter, thisPort, p freesp.Port) (err error) {
	var from, to freesp.Port
	if thisPort.(freesp.Port).Direction() == freesp.OutPort {
		from = thisPort.(freesp.Port)
		to = p
	} else {
		from = p
		to = thisPort.(freesp.Port)
	}
	conn := freesp.Connection{from, to}
	connText := fmt.Sprintf("%s/%s -> %s/%s",
		from.Node().NodeName(), from.PortName(),
		to.Node().NodeName(), to.PortName())
	err = s.addEntry(iter, imageConnected, connText, conn)
	if err != nil {
		err = gtkErr("FilesTreeStore.addConnection", "s.addEntry()", err)
		return
	}
	return
}

func (s *FilesTreeStore) addPorts(iter *gtk.TreeIter, ports []freesp.Port, kind *gdk.Pixbuf) error {
	ts := s.treestore
	for _, p := range ports {
		i := ts.Append(iter)
		err := s.addEntry(i, kind, p.PortName(), p)
		if err != nil {
			return gtkErr("FilesTreeStore.addPorts", "s.addEntry()", err)
		}
		err = s.addPortDetails(i, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FilesTreeStore) addPortDetails(iter *gtk.TreeIter, port freesp.Port) error {
	ts := s.treestore
	i := ts.Append(iter)
	t := port.ItsType()
	err := s.addEntry(i, imageSignalType, port.ItsType().TypeName(), t)
	if err != nil {
		return gtkErr("FilesTreeStore.addPortDetails", "s.addEntry(1)", err)
	}
	for _, c := range port.Connections() {
		i = ts.Append(iter)
		err := s.addConnection(i, port, c)
		if err != nil {
			return gtkErr("FilesTreeStore.addPortDetails", "s.addEntry(2)", err)
		}
	}
	return nil
}
