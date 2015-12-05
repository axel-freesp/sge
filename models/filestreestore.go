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

type storedType int

const (
	TInvalid storedType = iota
	TNone
	TSignalGraph
	TNode
	TNodeType
	TPort
	TPortType
	TNamedPortType
	TSignalType
	TConnection
)

var (
	imageBehaviour      *gdk.Pixbuf = nil
	imageInputNode      *gdk.Pixbuf = nil
	imageOutputNode     *gdk.Pixbuf = nil
	imageProcessingNode *gdk.Pixbuf = nil
	imageNodeType       *gdk.Pixbuf = nil
	imageInputPort      *gdk.Pixbuf = nil
	imageOutputPort     *gdk.Pixbuf = nil
	imagePortType       *gdk.Pixbuf = nil
	imageConnected      *gdk.Pixbuf = nil
)

func Init() error {
	iconPath := os.Getenv(envIconPathKey)
	if len(iconPath) == 0 {
		log.Fatal("Missing Environment variable ", iconPath)
	}
	var err error
	imageBehaviour, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/test1.png", iconPath))
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
	imageInputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/inport-green.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageOutputPort, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/outport-red.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imagePortType, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/signal-type.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	imageConnected, err = gdk.PixbufNewFromFile(fmt.Sprintf("%s/link.png", iconPath))
	if err != nil {
		return gtkErr("FilesTreeStore.Init", "test0.png", err)
	}
	fmt.Println("models.Init: all icons intitialized")
	return nil
}

type behaviourFile struct {
	filename string
	graph    freesp.SignalGraph
}

type FilesTreeStore struct {
	treestore *gtk.TreeStore
	//library []freesp.Library
	behaviour []behaviourFile
	//mappings [] freesp.Mapping
    lookup map[string]interface{}
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
	ret = &FilesTreeStore{ts, nil, make(map[string]interface{})}
	err = nil
	return
}

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
	return
}

func (s *FilesTreeStore) GetObject(iter *gtk.TreeIter) (ret interface{}, err error) {
    path, err := s.treestore.GetPath(iter)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetObject", "iter.GetPath", err)
		return
	}
    ret = s.lookup[path.String()]
	return
}

func (ts *FilesTreeStore) addEntry(iter *gtk.TreeIter, icon *gdk.Pixbuf, text string, data interface{}) error {
    path, err := ts.treestore.GetPath(iter)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "iter.GetPath", err)
		return err
	}
    ts.lookup[path.String()] = data
	err = ts.treestore.SetValue(iter, iconCol, icon)
	if err != nil {
		return gtkErr("FilesTreeStore.addEntry", "ts.SetValue(iconCol)", err)
	}
	err = ts.treestore.SetValue(iter, textCol, text)
	if err != nil {
		return gtkErr("FilesTreeStore.addEntry", "ts.SetValue(textCol)", err)
	}
	return nil
}

func (s *FilesTreeStore) AddBehaviourFile(filename string, graph freesp.SignalGraph) error {
	s.behaviour = append(s.behaviour, behaviourFile{filename, graph})
	ts := s.treestore
	iter := ts.Append(nil)
	err := s.addEntry(iter, imageBehaviour, filename, graph)
	if err != nil {
		return gtkErr("FilesTreeStore.AddBehaviourFile", "ts.addEntry(filename)", err)
	}
	i := ts.Append(iter)
	err = s.addEntry(i, imagePortType, "Signal Types", nil)
	if err != nil {
		return gtkErr("FilesTreeStore.AddBehaviourFile", "ts.addEntry(SignalTypes dummy)", err)
	}
	for _, st := range graph.ItsType().SignalTypes() {
		j := ts.Append(i)
		err := s.addEntry(j, imagePortType, st.TypeName(), st)
		if err != nil {
			return gtkErr("FilesTreeStore.AddBehaviourFile", "ts.addEntry(SignalTypes)", err)
		}
	}
	err = s.addNodes(iter, graph.ItsType().InputNodes(), imageInputNode)
	if err != nil {
		return err
	}
	err = s.addNodes(iter, graph.ItsType().OutputNodes(), imageOutputNode)
	if err != nil {
		return err
	}
	err = s.addNodes(iter, graph.ItsType().ProcessingNodes(), imageProcessingNode)
	if err != nil {
		return err
	}
	return nil
}

func (s *FilesTreeStore) addNodes(iter *gtk.TreeIter, nodes []freesp.Node, kind *gdk.Pixbuf) error {
	ts := s.treestore
	for _, n := range nodes {
		i := ts.Append(iter)
		err := s.addEntry(i, kind, n.NodeName(), n)
		if err != nil {
			return gtkErr("FilesTreeStore.addNodes", "ts.addEntry(0)", err)
		}
		j := ts.Append(i)
		t := n.ItsType()
		err = s.addEntry(j, imageNodeType, n.ItsType().TypeName(), t)
		if err != nil {
			return gtkErr("FilesTreeStore.addNodes", "ts.addEntry(2)", err)
		}
		if len(n.InPorts()) > 0 {
			err = s.addPorts(i, n.InPorts(), imageInputPort)
			if err != nil {
				return err
			}
		}
		if len(n.OutPorts()) > 0 {
			err = s.addPorts(i, n.OutPorts(), imageOutputPort)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	err := s.addEntry(i, imagePortType, port.ItsType().TypeName(), t)
	if err != nil {
		return gtkErr("FilesTreeStore.addPortDetails", "s.addEntry()", err)
	}
	for _, c := range port.Connections() {
		con := freesp.Connection{port, c}
		i = ts.Append(iter)
		err := s.addEntry(i, imageConnected, c.Node().NodeName(), con)
		if err != nil {
			return gtkErr("FilesTreeStore.addPortDetails", "s.addEntry", err)
		}
	}
	return nil
}
