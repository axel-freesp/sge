package tree

import (
	"github.com/axel-freesp/sge/interface/filedata"
	"github.com/axel-freesp/sge/interface/graph"
)

type TreeViewIf interface {
	SelectId(id string) error
}

type TreeMgrIf interface {
	AddToplevel(ToplevelTreeElementIf) (newId string, err error)
	RemoveToplevel(id string) (deleted []IdWithObject, err error)
	SetValueById(id, value string) error
	GetToplevelId(ToplevelTreeElementIf) (id string, err error)
	GetObjectById(string) (TreeElementIf, error)
}

type TreeIf interface {
	Current() Cursor
	Append(c Cursor) Cursor
	Insert(c Cursor) Cursor
	Remove(c Cursor) (prefix string, index int)
	Parent(c Cursor) Cursor
	Object(c Cursor) (obj TreeElementIf)
	Cursor(obj TreeElementIf) (cursor Cursor)
	CursorAt(start Cursor, obj TreeElementIf) (cursor Cursor)
	AddEntry(c Cursor, sym Symbol, text string, obj TreeElementIf, prop Property) (err error)
	Property(c Cursor) Property
}

type NamedTreeElementIf interface {
	TreeElementIf
	graph.Namer
}
type TreeElementIf interface {
	graph.XmlCreator
	AddToTree(tree TreeIf, cursor Cursor)
	AddNewObject(tree TreeIf, cursor Cursor, obj TreeElementIf) (newCursor Cursor, err error)
	RemoveObject(tree TreeIf, cursor Cursor) (removed []IdWithObject)
}

type ToplevelTreeElementIf interface {
	TreeElementIf
	filedata.FileDataIf
	RemoveFromTree(TreeIf)
}

type Cursor struct {
	Path     string
	Position int
}

const AppendCursor = -1

type IdWithObject struct {
	ParentId string
	Position int
	Object   TreeElementIf
}

type Symbol int

const (
	SymbolInputPort Symbol = iota
	SymbolOutputPort
	SymbolSignalType
	SymbolConnection
	SymbolImplElement
	SymbolImplGraph
	SymbolLibrary
	SymbolInputPortType
	SymbolOutputPortType
	SymbolInputNode
	SymbolOutputNode
	SymbolProcessingNode
	SymbolNodeType
	SymbolSignalGraph
	SymbolPlatform
	SymbolArch
	SymbolProcess
	SymbolIOType
	SymbolInChannel
	SymbolOutChannel
	SymbolMappings
	SymbolMapped
	SymbolUnmapped
)

type Property interface {
	IsReadOnly() bool
	MayAddObject() bool
	MayEdit() bool
	MayRemove() bool
}
