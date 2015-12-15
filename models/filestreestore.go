package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"strings"
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

func (s *FilesTreeStore) AddNewObject(parentId string, position int, obj interface{}) (newId string, err error) {
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
		err = fmt.Errorf("FilesTreeStore.AddNewObject: doAddFreespObject failed: %v", err)
		return
	}

	iter := s.treestore.Insert(parentIter, position)

	err = s.doAddTreeObject(iter, obj, parent)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.AddNewObject: doAddTreeObject failed: %v", err)
		return
	}

	newId, err = s.getIdFromIter(iter)
	return
}

type IdWithObject struct {
	ParentId string
	Position int
	Object   interface{}
}

// Root elements cannot be deleted. They need to be removed (e.g file->close)
func (s *FilesTreeStore) DeleteObject(id string) (deleted []IdWithObject, err error) {
	fmt.Printf("FilesTreeStore.DeleteObject: %s\n", id)
	obj, _, err := s.getObjAndIterById(id)
	if err != nil {
		return
	}
	fmt.Printf("FilesTreeStore.DeleteObject(1): %T\n", obj)
	li := strings.LastIndex(id, ":")
	if li < 0 {
		err = fmt.Errorf("Try to delete root element")
		return
	}
	parentId := id[:li]
	parent, err := s.GetObjectById(parentId)
	if err != nil {
		return
	}
	var position int
	fmt.Sscanf(id[li+1:], "%d", &position)
	toDelete := IdWithObject{parentId, position, obj}
	fmt.Printf("FilesTreeStore.DeleteObject id=%s, id[li+1:]=%s: %v\n", id, id[li+1:], toDelete)

	switch obj.(type) {
	case freesp.Implementation:
		fmt.Printf("FilesTreeStore.DeleteObject: %T\n", obj)
		parent.(freesp.NodeType).RemoveImplementation(obj.(freesp.Implementation))

	case freesp.SignalGraphType:
		fmt.Printf("FilesTreeStore.DeleteObject: %T\n", obj)
		deleted, err = s.deleteSignalGraphType(parentId, parent.(freesp.NodeType), obj.(freesp.SignalGraphType))

	case freesp.SignalType:
		err = fmt.Errorf("FilesTreeStore.DeleteObject: not yet implemented: %T\n", obj)
		return
	case freesp.NamedPortType:
		fmt.Printf("FilesTreeStore.DeleteObject: %T\n", obj)

		for _, n := range parent.(freesp.NodeType).Instances() {
			var nodeId string
			nodeId, err = s.getIdFromObject(n)
			if err != nil {
				err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced node:", err)
				return
			}
			fmt.Printf("Look at node %s: %s\n", nodeId, n)
			var portlist []freesp.Port
			var refPortTypeList []freesp.NamedPortType
			if obj.(freesp.NamedPortType).Direction() == freesp.InPort {
				portlist = n.InPorts()
				refPortTypeList = parent.(freesp.NodeType).InPorts()
			} else {
				portlist = n.OutPorts()
				refPortTypeList = parent.(freesp.NodeType).OutPorts()
			}
			var portToDelete freesp.Port
			for _, p := range portlist {
				if p.PortName() == obj.(freesp.NamedPortType).Name() {
					portToDelete = p
					break
				}
			}
			if portToDelete == nil {
				err = fmt.Errorf("FilesTreeStore.DeleteObject: error: port not found to remove with named port type.")
				return
			}
			var portTypeToDelete freesp.NamedPortType
			for _, p := range refPortTypeList {
				if p.Name() == obj.(freesp.NamedPortType).Name() {
					portTypeToDelete = p
					break
				}
			}
			if portTypeToDelete == nil {
				err = fmt.Errorf("FilesTreeStore.DeleteObject: error: referenced port type not found to remove with named port type.")
				return
			}
			var portId string
			var ok bool
			portId, ok = s.getIdFromObjectRecursive(nodeId, portToDelete)
			if !ok {
				err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced port.")
				return
			}
			fmt.Printf("portToDelete %s: %s\n", portId, portToDelete)
			var portTypeId string
			portTypeId, ok = s.getIdFromObjectRecursive(nodeId, portTypeToDelete)
			if !ok {
				err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced port type.")
				return
			}
			//deleted = append(deleted, IdWithObject{portTypeId, portTypeToDelete})
			s.deleteElement(portTypeId)
			var del []IdWithObject
			del, err = s.deletePort(portId, nodeId, n, portToDelete)
			for _, i := range del {
				deleted = append(deleted, i)
			}
		}
		parent.(freesp.NodeType).RemoveNamedPortType(obj.(freesp.NamedPortType))

	case freesp.NodeType:
		err = fmt.Errorf("FilesTreeStore.DeleteObject: not yet implemented: %T\n", obj)
		return
	case freesp.Node:
		err = fmt.Errorf("FilesTreeStore.DeleteObject: not yet implemented: %T\n", obj)
		return
	case freesp.Connection:
		fmt.Printf("FilesTreeStore.DeleteObject: %T\n", obj)
		deleted, err = s.deleteConnection(id, parentId, obj.(freesp.Connection))

	default:
		err = fmt.Errorf("FilesTreeStore.DeleteObject: invalid type %T\n", obj)
		return
	}

	s.deleteElement(id)

	deleted = append(deleted, toDelete)
	return
}

/*
 *      Local functions
 */

func (s *FilesTreeStore) deletePortOfNode(n freesp.Node, pt freesp.NamedPortType, parent freesp.NodeType) (deleted []IdWithObject, err error) {
	var nodeId string
	nodeId, err = s.getIdFromObject(n)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced node:", err)
		return
	}
	fmt.Printf("Look at node %s: %s\n", nodeId, n)
	var portlist []freesp.Port
	var refPortTypeList []freesp.NamedPortType
	if pt.Direction() == freesp.InPort {
		portlist = n.InPorts()
		refPortTypeList = parent.InPorts()
	} else {
		portlist = n.OutPorts()
		refPortTypeList = parent.OutPorts()
	}
	var portToDelete freesp.Port
	for _, p := range portlist {
		if p.PortName() == pt.Name() {
			portToDelete = p
			break
		}
	}
	if portToDelete == nil {
		err = fmt.Errorf("FilesTreeStore.DeleteObject: error: port not found to remove with named port type.")
		return
	}
	var portTypeToDelete freesp.NamedPortType
	for _, p := range refPortTypeList {
		if p.Name() == pt.Name() {
			portTypeToDelete = p
			break
		}
	}
	if portTypeToDelete == nil {
		err = fmt.Errorf("FilesTreeStore.DeleteObject: error: referenced port type not found to remove with named port type.")
		return
	}
	var portId string
	var ok bool
	portId, ok = s.getIdFromObjectRecursive(nodeId, portToDelete)
	if !ok {
		err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced port.")
		return
	}
	fmt.Printf("portToDelete %s: %s\n", portId, portToDelete)
	var portTypeId string
	portTypeId, ok = s.getIdFromObjectRecursive(nodeId, portTypeToDelete)
	if !ok {
		err = fmt.Errorf("FilesTreeStore.DeleteObject: error: no id for referenced port type.")
		return
	}
	//deleted = append(deleted, IdWithObject{portTypeId, portTypeToDelete})
	s.deleteElement(portTypeId)
	var del []IdWithObject
	del, err = s.deletePort(portId, nodeId, n, portToDelete)
	for _, i := range del {
		deleted = append(deleted, i)
	}
	return
}

func (s *FilesTreeStore) deletePort(id, parentId string, parent freesp.Node, port freesp.Port) (deleted []IdWithObject, err error) {
	for _, c := range port.Connections() {
		var connId string
		var ok bool
		conn := freesp.Connection{port, c}
		connId, ok = s.getIdFromObjectRecursive(id, conn)
		if !ok {
			err = fmt.Errorf("FilesTreeStore.deletePort error: cannot find connection: ", err)
			return
		}
		_, err = s.deleteConnection(connId, id, conn)
		if err != nil {
			return
		}
		var position int
		fmt.Sscanf(connId[len(id)+1:], "%d", &position)
		deleted = append(deleted, IdWithObject{id, position, conn})
		s.deleteElement(connId)
	}
	s.deleteElement(id)
	return
}

func (s *FilesTreeStore) deleteSignalGraphType(parentId string, parent freesp.NodeType, obj freesp.SignalGraphType) (deleted []IdWithObject, err error) {
	var imp freesp.Implementation
	for _, i := range parent.Implementation() {
		if i.ImplementationType() == freesp.NodeTypeGraph && i.Graph() == obj {
			imp = i
		}
	}
	if imp == nil {
		return
	}
	for _, n := range imp.Graph().ProcessingNodes() {
		nid, ok := s.getIdFromObjectRecursive(parentId, n)
		if ok {
			var position int
			fmt.Sscanf(nid[len(parentId)+1:], "%d", &position)
			deleted = append(deleted, IdWithObject{parentId, position, n})
		}
	}
	parent.RemoveImplementation(imp)
	return
}

func (s *FilesTreeStore) deleteConnection(id, parentId string, obj freesp.Connection) (deleted []IdWithObject, err error) {
	from, to := obj.From, obj.To
	from.RemoveConnection(to)
	to.RemoveConnection(from)
	var fromId, toId string
	fromId, err = s.getIdFromObject(from)
	if err != nil {
		return
	}
	toId, err = s.getIdFromObject(to)
	if err != nil {
		return
	}
	var otherId string
	if fromId == parentId {
		otherId = toId
	} else {
		otherId = fromId
	}
	otherObj := s.lookup[otherId]
	ok := true
	var i int
	for i = 0; ok; i++ {
		e, ok := s.lookup[fmt.Sprintf("%s:%d", otherId, i)]
		if ok {
			if e == obj {
				break
			}
		}
	}
	if !ok {
		fmt.Println("FilesTreeStore.DeleteObject error: Could not find remote connection object")
	} else {
		s.deleteElement(fmt.Sprintf("%s:%d", otherId, i))
	}
	deleted = append(deleted, IdWithObject{otherId, i, otherObj})
	return
}

func (s *FilesTreeStore) deleteElement(id string) {
	s.deleteSubtree(id)
	prefix := id[:strings.LastIndex(id, ":")]
	suffix := id[strings.LastIndex(id, ":")+1:]
	var index int
	fmt.Sscanf(suffix, "%d", &index)
	for ok := true; ok; index++ {
		ok = s.shiftLookup(prefix, fmt.Sprintf("%d", index+1), fmt.Sprintf("%d", index))
	}
}

func (s *FilesTreeStore) deleteSubtree(id string) {
	_, iter, err := s.getObjAndIterById(id)
	if err != nil {
		return
	}
	ok := true
	for i := 0; ok; i++ {
		child := fmt.Sprintf("%s:%d", id, i)
		_, ok = s.lookup[child]
		if ok {
			s.deleteSubtree(child)
		}
	}
	s.treestore.Remove(iter)
}

func (s *FilesTreeStore) shiftLookup(parent, from, to string) (ok bool) {
	fromId, toId := fmt.Sprintf("%s:%s", parent, from), fmt.Sprintf("%s:%s", parent, to)
	obj, ok := s.lookup[fromId]
	if ok {
		s.lookup[toId] = obj
		ok2 := true
		for i := 0; ok2; i++ {
			ok2 = s.shiftLookup(parent, fmt.Sprintf("%s:%d", from, i), fmt.Sprintf("%s:%d", to, i))
		}
	} else {
		delete(s.lookup, toId)
	}
	return
}

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
		var theGraph freesp.SignalGraphType
		switch parent.(type) {
		case freesp.SignalGraph:
			theGraph = parent.(freesp.SignalGraph).ItsType()
		case freesp.SignalGraphType:
			theGraph = parent.(freesp.SignalGraphType)
		}
		err = theGraph.AddNode(obj.(freesp.Node))
	case freesp.Connection:
		var port1, port2 freesp.Port
		var node1, node2 freesp.Node
		node1 = obj.(freesp.Connection).To.Node()
		for _, p := range node1.InPorts() {
			if p.PortName() == obj.(freesp.Connection).To.PortName() {
				port1 = p
				break
			}
		}
		node2 = obj.(freesp.Connection).From.Node()
		for _, p := range node2.OutPorts() {
			if p.PortName() == obj.(freesp.Connection).From.PortName() {
				port2 = p
				break
			}
		}
		if port1 == nil || port2 == nil {
			err = fmt.Errorf("FilesTreeStore.doAddFreespObject error: one of referenced ports (%s) not found\n", obj.(freesp.Connection))
			return
		}
		err = port1.AddConnection(port2)
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddFreespObject(1): %v", err)
			return
		}
		err = port2.AddConnection(port1)
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddFreespObject(2): %v", err)
			return
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
		if err != nil {
			return
		}
		for _, n := range parent.(freesp.NodeType).Instances() {
			var i *gtk.TreeIter
			i, err = s.getIterFromObject(n)
			if err != nil {
				return
			}
			var p freesp.Port
			var plist []freesp.Port
			if obj.(freesp.NamedPortType).Direction() == freesp.InPort {
				plist = n.InPorts()
			} else {
				plist = n.OutPorts()
			}
			for _, p = range plist {
				if p.PortName() == obj.(freesp.NamedPortType).Name() {
					break
				}
			}
			fmt.Println("Adding port ", p)
			i = s.treestore.Append(i)
			err = s.addPort(i, p, icon)
			var nodeId string
			nodeId, err = s.getIdFromObject(n)
			if err != nil {
				return
			}
			id, ok := s.getIdFromObjectRecursive(nodeId, parent.(freesp.NodeType))
			if !ok {
				return
			}
			i, err = s.treestore.GetIterFromString(id)
			if err != nil {
				return
			}
			j := s.treestore.Append(i)
			err = s.addNamedPortType(j, obj.(freesp.NamedPortType), icon)
			if err != nil {
				return
			}
		}
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
	case freesp.Connection:
		fmt.Println("Adding ", obj.(freesp.Connection))
		var port1, port2 freesp.Port
		var node1, node2 freesp.Node
		node1 = obj.(freesp.Connection).To.Node()
		for _, port1 = range node1.InPorts() {
			if port1.PortName() == obj.(freesp.Connection).To.PortName() {
				break
			}
		}
		node2 = obj.(freesp.Connection).From.Node()
		for _, port2 = range node2.OutPorts() {
			if port2.PortName() == obj.(freesp.Connection).From.PortName() {
				break
			}
		}
		fmt.Println("node1 =", node1)
		fmt.Println("port1 =", port1)
		fmt.Println("node2 =", node2)
		fmt.Println("port2 =", port2)
		var iter1, iter2 *gtk.TreeIter
		iter1, err = s.getIterFromObject(port1)
		if err != nil {
			return
		}
		iter2, err = s.getIterFromObject(port2)
		if err != nil {
			return
		}
		iter = s.treestore.Append(iter1)
		err = s.addConnection(iter, port1, port2)
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddTreeObject: addConnection(3) failed")
			return
		}
		iter = s.treestore.Append(iter2)
		err = s.addConnection(iter, port2, port1)
		if err != nil {
			err = fmt.Errorf("FilesTreeStore.doAddTreeObject: addConnection(4) failed")
			return
		}
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
	case freesp.Port, freesp.Connection:
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

func (s *FilesTreeStore) getIdFromObject(obj interface{}) (id string, err error) {
	iter, err := s.getIterFromObject(obj)
	if err != nil {
		return
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		return
	}
	id = path.String()
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

func (s *FilesTreeStore) addPorts(iter *gtk.TreeIter, ports []freesp.Port, kind *gdk.Pixbuf) (err error) {
	ts := s.treestore
	for _, p := range ports {
		i := ts.Append(iter)
		err = s.addPort(i, p, kind)
		if err != nil {
			return
		}
	}
	return
}

func (s *FilesTreeStore) addPort(iter *gtk.TreeIter, p freesp.Port, kind *gdk.Pixbuf) (err error) {
	err = s.addEntry(iter, kind, p.PortName(), p)
	if err != nil {
		err = gtkErr("FilesTreeStore.addPorts", "s.addEntry()", err)
	}
	err = s.addPortDetails(iter, p)
	return
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
