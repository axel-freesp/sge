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
	treestore *gtk.TreeStore
	//mappings [] freesp.Mapping
	lookup map[string]interface{}
	graphs map[string]freesp.SignalGraphType
	libs   map[string]freesp.Library
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
		make(map[string]freesp.SignalGraphType),
		make(map[string]freesp.Library)}
	err = nil
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
	return
}

// Returns the object stored with the row
func (s *FilesTreeStore) GetObject(iter *gtk.TreeIter) (ret interface{}, err error) {
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetObject", "iter.GetPath", err)
		return
	}
	ret = s.lookup[path.String()]
	return
}

func (s *FilesTreeStore) AddSignalGraphFile(filename string, graph freesp.SignalGraph) error {
	ts := s.treestore
	iter := ts.Append(nil)
	err := s.addEntry(iter, imageSignalGraph, filename, graph)
	if err != nil {
		return gtkErr("FilesTreeStore.AddSignalGraphFile", "ts.addEntry(filename)", err)
	}
	return s.addSignalGraph(iter, graph)
}

func (s *FilesTreeStore) AddLibraryFile(filename string, lib freesp.Library) error {
	return s.addLibrary(s.treestore.Append(nil), lib)
}

/*
 *      Local functions
 */

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

func (s *FilesTreeStore) addSignalGraph(iter *gtk.TreeIter, graph freesp.SignalGraph) error {
	return s.addSignalGraphType(iter, graph.ItsType())
}

func (s *FilesTreeStore) addSignalGraphType(iter *gtk.TreeIter, gtype freesp.SignalGraphType) error {
	ts := s.treestore
	var err error
	if s.graphs[gtype.Name()] == nil {
		s.graphs[gtype.Name()] = gtype
		/*
		   for _, l := range gtype.Libraries() {
		       i := ts.Append(iter)
		       err = s.addLibrary(i, l)
		       if err != nil {
		           return gtkErr("FilesTreeStore.addSignalGraph", "s.addLibrary()", err)
		       }
		   }
		*/
		if len(gtype.SignalTypes()) > 0 {
			i := ts.Append(iter)
			err = s.addEntry(i, imageSignalType, "Signal Types", nil)
			if err != nil {
				return gtkErr("FilesTreeStore.addSignalGraph", "ts.addEntry(SignalTypes dummy)", err)
			}
			for _, st := range gtype.SignalTypes() {
				j := ts.Append(i)
				err := s.addEntry(j, imageSignalType, st.TypeName(), st)
				if err != nil {
					return gtkErr("FilesTreeStore.addSignalGraph", "ts.addEntry(SignalTypes)", err)
				}
			}
		}
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
			err := s.addEntry(j, imageSignalType, st.TypeName(), st)
			if err != nil {
				return gtkErr("FilesTreeStore.addLibrary", "ts.addEntry(SignalTypes)", err)
			}
		}
		err = s.addNodeTypes(iter, lib.NodeTypes())
		if err != nil {
			return err
		}
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
		return gtkErr("FilesTreeStore.addNodeTypes", "ts.addEntry(0)", err)
	}
	impl := ntype.Implementation()
	if impl != nil {
		switch impl.ImplementationType() {
		case freesp.NodeTypeElement:
			j := ts.Append(iter)
			err = s.addEntry(j, imageNodeTypeElem, impl.ElementName(), impl)
		default:
			j := ts.Append(iter)
			err := s.addEntry(j, imageSignalGraph, "implementation", impl)
			if err != nil {
				return gtkErr("FilesTreeStore.addNodeTypes", "ts.addEntry(implementation)", err)
			}
			err = s.addSignalGraphType(j, impl.Graph())
			if err != nil {
				return gtkErr("FilesTreeStore.addNodeTypes", "ts.addSignalGraph()", err)
			}
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

func (s *FilesTreeStore) addNamedPortTypes(iter *gtk.TreeIter, ports []freesp.NamedPortType, kind *gdk.Pixbuf) error {
	ts := s.treestore
	for _, p := range ports {
		i := ts.Append(iter)
		err := s.addEntry(i, kind, p.Name(), p)
		if err != nil {
			return gtkErr("FilesTreeStore.addNamedPortTypes", "s.addEntry()", err)
		}
		j := ts.Append(i)
		err = s.addEntry(j, imageSignalType, p.TypeName(), p.SignalType())
		if err != nil {
			return gtkErr("FilesTreeStore.addNamedPortTypes", "s.addEntry()", err)
		}
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
		err = s.addNodeType(j, t)
		if err != nil {
			return gtkErr("FilesTreeStore.addNodeType", "ts.addEntry(2)", err)
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
	err := s.addEntry(i, imageSignalType, port.ItsType().TypeName(), t)
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
