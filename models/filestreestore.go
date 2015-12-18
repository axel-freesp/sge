package models

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"image"
	_ "image/png"
	"log"
	"os"
	"strings"
)

type IdWithObject struct {
	ParentId string
	Position int
	Object   interface{}
}

type TreeElement interface {
	AddToTree(tree *FilesTreeStore, cursor Cursor)
	AddNewObject(tree *FilesTreeStore, cursor Cursor, obj interface{}) (newCursor Cursor)
	RemoveObject(tree *FilesTreeStore, cursor Cursor) (removed []IdWithObject)
}

const envIconPathKey = "SGE_ICON_PATH"
const (
	iconCol = 0
	textCol = 1
)

var (
	imageSignalGraphAlt *gdk.Pixbuf = nil
)

var modulesInit = []func(string) error{init_port,
	init_signaltype,
	init_connection,
	init_node,
	init_nodetype,
	init_implementation,
	init_namedporttype,
	init_signalgraphtype,
	init_library,
	init_signalgraph,
}

func Init() (err error) {
	iconPath := os.Getenv(envIconPathKey)
	if len(iconPath) == 0 {
		err = fmt.Errorf("Missing Environment variable %s", iconPath)
		return
	}

	for _, init := range modulesInit {
		err = init(iconPath)
		if err != nil {
			return
		}
	}
	return

	file, err := os.Open(fmt.Sprintf("%s/test1.png", iconPath))
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.Init error: image.Decode failed: %s", err)
		return
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	data := make([]byte, width*height*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		j := width * 4 * y
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[j+4*x+0] = byte(r)
			data[j+4*x+1] = byte(g)
			data[j+4*x+2] = byte(b)
			data[j+4*x+3] = byte(a)
		}
	}

	imageSignalGraphAlt, err = gdk.PixbufNewFromBytes(data, gdk.COLORSPACE_RGB, true, 8, width, height, width*4)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.Init error: PixbufNewFromXpmData failed: %s", err)
	}
	fmt.Println("models.Init: all icons intitialized")
	return
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
	cursor := s.Append(rootCursor)
	SignalGraph{graph}.AddToTree(s, cursor)
	newId = cursor.Path
	return
	/*
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
	*/
}

func (s *FilesTreeStore) AddLibraryFile(filename string, lib freesp.Library) (newId string, err error) {
	cursor := s.Append(rootCursor)
	Library{lib}.AddToTree(s, cursor)
	newId = cursor.Path
	return

	/*
		iter := s.treestore.Append(nil)
		err = s.addLibrary(iter, lib)
		if err != nil {
			return
		}
		newId, err = s.getIdFromIter(iter)
		return
	*/
}

func (tree *FilesTreeStore) AddNewObject(parentId string, position int, obj interface{}) (newId string, err error) {
	//log.Printf("FilesTreeStore.AddNewObject %T, %v\n", obj, obj)
	parent, parentIter, err := tree.getObjAndIterById(parentId)
	if err != nil {
		return
	}
	//log.Printf("FilesTreeStore.AddNewObject parent = %T, %v\n", parent, parent)
	cursor := Cursor{parentIter, parentId, position}
	switch parent.(type) {
	case freesp.NodeType:
		cursor = NodeType{parent.(freesp.NodeType)}.AddNewObject(tree, cursor, obj)

	case freesp.SignalGraphType:
		cursor = SignalGraphType{parent.(freesp.SignalGraphType)}.AddNewObject(tree, cursor, obj)

	case freesp.SignalGraph:
		cursor = SignalGraph{parent.(freesp.SignalGraph)}.AddNewObject(tree, cursor, obj)

	case freesp.Port:
		cursor = Port{parent.(freesp.Port)}.AddNewObject(tree, cursor, obj)

	case freesp.Library:
		cursor = Library{parent.(freesp.Library)}.AddNewObject(tree, cursor, obj)

	case freesp.Implementation:
		cursor = Implementation{parent.(freesp.Implementation)}.AddNewObject(tree, cursor, obj)

	default:
		err = fmt.Errorf("FilesTreeStore.AddNewObject: invalid parent type %T\n", parent)
		return
	}
	newId = cursor.Path

	return
}

// Root elements cannot be deleted. They need to be removed (e.g file->close)
func (tree *FilesTreeStore) DeleteObject(id string) (deleted []IdWithObject, err error) {
	_, iter, err := tree.getObjAndIterById(id)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.DeleteObject error: object not found, id:", id)
		return
	}
	li := strings.LastIndex(id, ":")
	if li < 0 {
		err = fmt.Errorf("FilesTreeStore.DeleteObject error: Try to delete root element")
		return
	}
	parentId := id[:li]
	parent, err := tree.GetObjectById(parentId)
	if err != nil {
		err = fmt.Errorf("FilesTreeStore.DeleteObject error: parent object not found")
		return
	}
	var position int
	fmt.Sscanf(id[li+1:], "%d", &position)

	cursor := Cursor{iter, id, AppendCursor}
	log.Println("FilesTreeStore.DeleteObject cursor =", cursor)

	switch parent.(type) {
	case freesp.NodeType:
		deleted = NodeType{parent.(freesp.NodeType)}.RemoveObject(tree, cursor)

	case freesp.SignalGraphType:
		deleted = SignalGraphType{parent.(freesp.SignalGraphType)}.RemoveObject(tree, cursor)

	case freesp.SignalGraph:
		deleted = SignalGraph{parent.(freesp.SignalGraph)}.RemoveObject(tree, cursor)

	case freesp.Port:
		deleted = Port{parent.(freesp.Port)}.RemoveObject(tree, cursor)

	case freesp.Library:
		deleted = Library{parent.(freesp.Library)}.RemoveObject(tree, cursor)

	case freesp.Implementation:
		deleted = Implementation{parent.(freesp.Implementation)}.RemoveObject(tree, cursor)

	default:
		err = fmt.Errorf("FilesTreeStore.DeleteObject: invalid parent type %T\n", parent)
		return
	}
	return
}

/*
 *      Tree Interface
 */

type Cursor struct {
	Iter     *gtk.TreeIter
	Path     string
	Position int
}

const AppendCursor = -1

var rootCursor = Cursor{nil, "", AppendCursor}

func (s *FilesTreeStore) Append(c Cursor) Cursor {
	iter := s.treestore.Append(c.Iter)
	if iter == nil {
		log.Fatal("FilesTreeStore.Append: zero iter")
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		log.Fatal("FilesTreeStore.Append: no path")
	}
	return Cursor{iter, path.String(), AppendCursor}
}

func (s *FilesTreeStore) Insert(c Cursor) Cursor {
	if c.Iter == nil {
		log.Fatal("FilesTreeStore.Insert FIXME: zero c.Iiter")
	}
	if c.Position >= 0 {
		index := 0
		for ok := true; ok; index++ {
			_, ok = s.lookup[fmt.Sprintf("%s:%d", c.Path, index+1)]
		}
		//fmt.Printf("FilesTreeStore.insertElement(%d): last element: %s:%d\n", position, prefix, index)
		for ; index >= c.Position; index-- {
			s.shiftLookup(c.Path, fmt.Sprintf("%d", index), fmt.Sprintf("%d", index+1))
		}
	}
	iter := s.treestore.Insert(c.Iter, c.Position)
	if iter == nil {
		log.Fatal("FilesTreeStore.Insert: zero iter")
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		log.Fatal("FilesTreeStore.Insert: no path")
	}
	return Cursor{iter, path.String(), AppendCursor}
}

func (s *FilesTreeStore) Remove(c Cursor) (prefix string, index int) {
	s.deleteSubtree(c.Path)
	prefix = c.Path[:strings.LastIndex(c.Path, ":")]
	suffix := c.Path[strings.LastIndex(c.Path, ":")+1:]
	fmt.Sscanf(suffix, "%d", &index)
	i := index
	for ok := true; ok; i++ {
		ok = s.shiftLookup(prefix, fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", i))
	}
	return
}

func (s *FilesTreeStore) Parent(c Cursor) Cursor {
	if c.Iter == nil {
		log.Fatal("FilesTreeStore.Parent: zero c.Iiter")
	}
	prefix := c.Path[:strings.LastIndex(c.Path, ":")]
	suffix := c.Path[strings.LastIndex(c.Path, ":")+1:]
	var index int
	fmt.Sscanf(suffix, "%d", &index)
	iter, err := s.treestore.GetIterFromString(prefix)
	if err != nil {
		log.Fatal("FilesTreeStore.Parent: no path")
	}
	return Cursor{iter, prefix, index}
}

func (s *FilesTreeStore) Object(c Cursor) (obj interface{}) {
	return s.lookup[c.Path]
}

// Toplevel search
func (s *FilesTreeStore) getIterAndPathFromObject(obj interface{}) (iter *gtk.TreeIter, path string, err error) {
	var o interface{}
	ok := false
	for i := 0; !ok; i++ {
		path = fmt.Sprintf("%d", i)
		o, ok = s.lookup[path]
		if !ok {
			break
		} else if o != obj {
			path, ok = s.getIdFromObjectRecursive(path, obj)
		} else {
			fmt.Println("match: ", ok, path)
		}
	}
	if !ok {
		err = fmt.Errorf("getIterAndPathFromObject error: object not found.")
		return
	}
	iter, err = s.treestore.GetIterFromString(path)
	return
}

// Toplevel search
func (s *FilesTreeStore) Cursor(obj interface{}) (cursor Cursor) {
	iter, path, err := s.getIterAndPathFromObject(obj)
	if err != nil {
		log.Fatal("FilesTreeStore.Cursor: obj %T: %v not found.", obj, obj)
	}
	return Cursor{iter, path, AppendCursor}
}

// Subtree search
func (s *FilesTreeStore) CursorAt(start Cursor, obj interface{}) (cursor Cursor) {
	path, ok := s.getIdFromObjectRecursive(start.Path, obj)
	if !ok {
		log.Println("FilesTreeStore.CursorAt: start =", start)
		log.Fatalf("FilesTreeStore.CursorAt: obj %T: %v not found.", obj, obj)
	}
	iter, err := s.treestore.GetIterFromString(path)
	if err != nil {
		log.Fatal("FilesTreeStore.CursorAt: gtk.GetIterFromString failed: %s", err)
	}
	return Cursor{iter, path, AppendCursor}
}

func (s *FilesTreeStore) AddEntry(c Cursor, icon *gdk.Pixbuf, text string, data interface{}) (err error) {
	if c.Iter == nil {
		err = fmt.Errorf("FilesTreeStore.addEntry: zero iter")
		return
	}
	path, err := s.treestore.GetPath(c.Iter)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "GetPath", err)
		return
	}
	s.lookup[path.String()] = data
	err = s.treestore.SetValue(c.Iter, iconCol, icon)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "SetValue(iconCol)", err)
		return
	}
	err = s.treestore.SetValue(c.Iter, textCol, text)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "SetValue(textCol)", err)
		return
	}
	return nil
}

/*
 *      Local functions
 */

func (s *FilesTreeStore) getObjAndIterById(id string) (obj interface{}, iter *gtk.TreeIter, err error) {
	obj, err = s.GetObjectById(id)
	if err != nil {
		return
	}
	iter, err = s.treestore.GetIterFromString(id)
	return
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
	if iter == nil {
		err = fmt.Errorf("FilesTreeStore.getIdFromObject: zero iter")
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
	if iter == nil {
		fmt.Println("FilesTreeStore.getObject: zero iter")
		return
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		return
	}
	obj, ok = s.lookup[path.String()]
	return
}
