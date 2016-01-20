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

func Init() (err error) {
	iconPath := os.Getenv(envIconPathKey)
	if len(iconPath) == 0 {
		err = fmt.Errorf("Missing Environment variable %s", iconPath)
		return
	}
	err = symbolInit(iconPath)
	if err != nil {
		fmt.Println("Warning: ", err)
	}

	return
}

type Element struct {
	prop freesp.Property
	elem freesp.TreeElement
}

type FilesTreeStore struct {
	treestore        *gtk.TreeStore
	lookup           map[string]Element
	libs             map[string]freesp.Library
	currentSelection *gtk.TreeIter
}

var _ freesp.Tree = (*FilesTreeStore)(nil)

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
	ret = &FilesTreeStore{ts, make(map[string]Element),
		make(map[string]freesp.Library),
		nil}
	err = nil
	return
}

func (s *FilesTreeStore) GetCurrentId() (id string) {
	id = ""
	if s.currentSelection == nil {
		return
	}
	p, err := s.treestore.GetPath(s.currentSelection)
	if err != nil {
		return
	}
	id = p.String()
	return
}

func (s *FilesTreeStore) GetObjectById(id string) (ret freesp.TreeElement, err error) {
	if len(id) == 0 {
		err = fmt.Errorf("FilesTreeStore.GetObjectById warning: empty id.\n")
		return
	}
	_, err = s.treestore.GetIterFromString(id)
	if err != nil {
		err = gtkErr("FilesTreeStore.GetObjectById", "GetIterFromString()", err)
		return
	}
	_, ok := s.lookup[id]
	if !ok {
		err = fmt.Errorf("FilesTreeStore.GetObjectById: no such id %s\n", id)
		return
	}
	ret = s.lookup[id].elem
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
		err = gtkErr("FilesTreeStore.SetValueById", "GetIterFromString()", err)
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
func (s *FilesTreeStore) GetObject(iter *gtk.TreeIter) (ret freesp.TreeElement, err error) {
	var ok bool
	ret, ok = s.getObject(iter)
	if !ok {
		err = fmt.Errorf("FilesTreeStore.GetObject: not found.")
		return
	}
	s.currentSelection = iter
	return
}

func (s *FilesTreeStore) AddSignalGraphFile(graph freesp.SignalGraph) (newId string, err error) {
	cursor := s.Append(rootCursor)
	graph.AddToTree(s, cursor)
	newId = cursor.Path
	return
}

func (s *FilesTreeStore) AddLibraryFile(lib freesp.Library) (cursor freesp.Cursor, err error) {
	cursor = s.Append(rootCursor)
	lib.AddToTree(s, cursor)
	return
}

func (s *FilesTreeStore) AddPlatformFile(plat freesp.Platform) (newId string, err error) {
	cursor := s.Append(rootCursor)
	plat.AddToTree(s, cursor)
	newId = cursor.Path
	return
}

func (s *FilesTreeStore) AddMappingFile(m freesp.Mapping) (cursor freesp.Cursor, err error) {
	cursor = s.Append(rootCursor)
	m.AddToTree(s, cursor)
	return
}

func (tree *FilesTreeStore) AddNewObject(parentId string, position int, obj freesp.TreeElement) (newId string, err error) {
	parent, _, err := tree.getObjAndIterById(parentId)
	if err != nil {
		return
	}
	cursor := freesp.Cursor{parentId, position}
	cursor, err = parent.AddNewObject(tree, cursor, obj)
	newId = cursor.Path
	return
}

// Root elements cannot be deleted. They need to be removed (e.g file->close)
func (tree *FilesTreeStore) DeleteObject(id string) (deleted []freesp.IdWithObject, err error) {
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

	cursor := freesp.Cursor{id, freesp.AppendCursor}
	deleted = parent.RemoveObject(tree, cursor)
	return
}

func (tree *FilesTreeStore) RemoveToplevel(id string) (deleted []freesp.IdWithObject, err error) {
	if strings.Contains(id, ":") {
		err = fmt.Errorf("FilesTreeStore.RemoveToplevel: not removing toplevel object (abort)")
		return
	}
	var obj freesp.TreeElement
	obj, err = tree.GetObjectById(id)
	if err != nil {
		return
	}
	obj.(freesp.Remover).Remove(tree)
	return
}

/*
 *      Tree Interface
 */

var rootCursor = freesp.Cursor{"", freesp.AppendCursor}

func (s *FilesTreeStore) Current() freesp.Cursor {
	return freesp.Cursor{s.GetCurrentId(), freesp.AppendCursor}
}

func (s *FilesTreeStore) Append(c freesp.Cursor) freesp.Cursor {
	var iter *gtk.TreeIter
	var err error
	if len(c.Path) > 0 {
		iter, err = s.treestore.GetIterFromString(c.Path)
		if err != nil {
			_, parentErr := s.treestore.GetIterFromString(s.Parent(c).Path)
			log.Printf("FilesTreeStore.Append: (%v) parent error: %s\n", s.Parent(c), parentErr)
			log.Printf("FilesTreeStore.Append: object = %T: %v\n", s.Object(c), s.Object(c))
			log.Fatalf("FilesTreeStore.Append: GetIterFromString(%v) failed: %s\n", c.Path, err)
		}
	}
	iter = s.treestore.Append(iter)
	if iter == nil {
		log.Fatal("FilesTreeStore.Append: zero child")
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		log.Fatal("FilesTreeStore.Append: no path")
	}
	return freesp.Cursor{path.String(), freesp.AppendCursor}
}

func (s *FilesTreeStore) Insert(c freesp.Cursor) freesp.Cursor {
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
	var iter *gtk.TreeIter
	var err error
	if len(c.Path) > 0 {
		iter, err = s.treestore.GetIterFromString(c.Path)
		if err != nil {
			log.Fatal("FilesTreeStore.Insert: zero iter")
		}
	}
	iter = s.treestore.Insert(iter, c.Position)
	if iter == nil {
		log.Fatal("FilesTreeStore.Insert: zero child")
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		log.Fatal("FilesTreeStore.Insert: no path")
	}
	return freesp.Cursor{path.String(), freesp.AppendCursor}
}

func (s *FilesTreeStore) Remove(c freesp.Cursor) (prefix string, index int) {
	s.deleteSubtree(c.Path)
	var suffix string
	if strings.Contains(c.Path, ":") {
		prefix = c.Path[:strings.LastIndex(c.Path, ":")]
		suffix = c.Path[strings.LastIndex(c.Path, ":")+1:]
	} else {
		suffix = c.Path
	}
	fmt.Sscanf(suffix, "%d", &index)
	i := index
	for ok := true; ok; i++ {
		ok = s.shiftLookup(prefix, fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", i))
	}
	_, ok := s.lookup[c.Path]
	if !ok {
		index--
		suffix = fmt.Sprintf("%d", index)
		if len(prefix) > 0 {
			c.Path = fmt.Sprintf("%s:%d", prefix, index)
		} else {
			c.Path = fmt.Sprintf("%d", index)
		}
		_, ok = s.lookup[c.Path]
		if !ok {
			s.currentSelection = nil
		} else {
			_, s.currentSelection, _ = s.getObjAndIterById(c.Path)
		}
	}
	return
}

func (s *FilesTreeStore) Parent(c freesp.Cursor) freesp.Cursor {
	prefix := c.Path[:strings.LastIndex(c.Path, ":")]
	suffix := c.Path[strings.LastIndex(c.Path, ":")+1:]
	var index int
	fmt.Sscanf(suffix, "%d", &index)
	return freesp.Cursor{prefix, index}
}

func (s *FilesTreeStore) Object(c freesp.Cursor) (obj freesp.TreeElement) {
	return s.lookup[c.Path].elem
}

// Toplevel search
func (s *FilesTreeStore) getIterAndPathFromObject(obj freesp.TreeElement) (iter *gtk.TreeIter, path string, err error) {
	var o Element
	ok := false
	for i := 0; !ok; i++ {
		path = fmt.Sprintf("%d", i)
		o, ok = s.lookup[path]
		if !ok {
			break
		} else if o.elem != obj {
			path, ok = s.getIdFromObjectRecursive(path, obj)
		}
	}
	if !ok {
		err = fmt.Errorf("getIterAndPathFromObject error: object not found.")
		for i := 0; !ok; i++ {
			path = fmt.Sprintf("%d", i)
			_, ok = s.lookup[path]
			if !ok {
				break
			} else {
				err = fmt.Errorf("%s, path:%s", err, path)
				ok = false
			}
		}
		return
	}
	iter, err = s.treestore.GetIterFromString(path)
	return
}

// Toplevel search
func (s *FilesTreeStore) Cursor(obj freesp.TreeElement) (cursor freesp.Cursor) {
	_, path, err := s.getIterAndPathFromObject(obj)
	if err != nil {
		log.Printf("%s\n", err)
		log.Panicf("FilesTreeStore.Cursor: obj %T: %v not found.\n", obj, obj)
	}
	return freesp.Cursor{path, freesp.AppendCursor}
}

// Subtree search
func (s *FilesTreeStore) CursorAt(start freesp.Cursor, obj freesp.TreeElement) (cursor freesp.Cursor) {
	path, ok := s.getIdFromObjectRecursive(start.Path, obj)
	if !ok {
		log.Println("FilesTreeStore.CursorAt: start =", start)
		log.Panicf("FilesTreeStore.CursorAt: obj %T: %v not found.\n", obj, obj)
	}
	return freesp.Cursor{path, freesp.AppendCursor}
}

func (s *FilesTreeStore) AddEntry(c freesp.Cursor, sym freesp.Symbol, text string, data freesp.TreeElement, prop freesp.Property) (err error) {
	var icon *gdk.Pixbuf
	if prop.IsReadOnly() {
		icon = readonlyPixbuf(sym)
	} else {
		icon = normalPixbuf(sym)
	}
	iter, err := s.treestore.GetIterFromString(c.Path)
	if err != nil {
		log.Fatal("FilesTreeStore.AddEntry: gtk.GetIterFromString failed: %s", err)
	}
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "GetPath", err)
		return
	}
	s.lookup[c.Path] = Element{prop, data}
	err = s.treestore.SetValue(iter, iconCol, icon)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "SetValue(iconCol)", err)
		return
	}
	err = s.treestore.SetValue(iter, textCol, text)
	if err != nil {
		err = gtkErr("FilesTreeStore.addEntry", "SetValue(textCol)", err)
		return
	}
	return nil
}

func (s *FilesTreeStore) Property(c freesp.Cursor) freesp.Property {
	return s.lookup[c.Path].prop
}

/*
 *      Local functions
 */

func (s *FilesTreeStore) getObjAndIterById(id string) (obj freesp.TreeElement, iter *gtk.TreeIter, err error) {
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
	var fromId, toId string
	if len(parent) > 0 {
		fromId, toId = fmt.Sprintf("%s:%s", parent, from), fmt.Sprintf("%s:%s", parent, to)
	} else {
		fromId, toId = fmt.Sprintf("%s", from), fmt.Sprintf("%s", to)
	}
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

func (s *FilesTreeStore) getIterFromObject(obj freesp.TreeElement) (iter *gtk.TreeIter, err error) {
	var id string
	var o Element
	ok := false
	for i := 0; !ok; i++ {
		id = fmt.Sprintf("%d", i)
		o, ok = s.lookup[id]
		if !ok {
			break
		} else if o.elem != obj {
			id, ok = s.getIdFromObjectRecursive(id, obj)
		}
	}
	if !ok {
		err = fmt.Errorf("GetIterFromObject error: object not found.")
		return
	}
	iter, err = s.treestore.GetIterFromString(id)
	return
}

func (s *FilesTreeStore) getIdFromObject(obj freesp.TreeElement) (id string, err error) {
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

func (s *FilesTreeStore) getIdFromObjectRecursive(parent string, obj freesp.TreeElement) (id string, ok bool) {
	ok = false
	var o Element
	for i := 0; !ok; i++ {
		id = fmt.Sprintf("%s:%d", parent, i)
		o, ok = s.lookup[id]
		if !ok {
			break
		} else if o.elem != obj {
			id, ok = s.getIdFromObjectRecursive(id, obj)
		}
	}
	return
}

func (s *FilesTreeStore) getObject(iter *gtk.TreeIter) (obj freesp.TreeElement, ok bool) {
	ok = false
	if iter == nil {
		fmt.Println("FilesTreeStore.getObject: zero iter")
		return
	}
	path, err := s.treestore.GetPath(iter)
	if err != nil {
		return
	}
	var o Element
	o, ok = s.lookup[path.String()]
	obj = o.elem
	return
}
