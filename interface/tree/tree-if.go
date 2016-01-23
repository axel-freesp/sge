package tree

import (
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/interface/filedata"
)

type TreeViewIf interface {
	SelectId(id string) error
}

type TreeMgrIf interface {
	AddToplevel(ToplevelTreeElement) (newId string, err error)
	RemoveToplevel(id string) (deleted []IdWithObject, err error)
	SetValueById(id, value string) error
	GetToplevelId(ToplevelTreeElement) (id string, err error)
}

type TreeIf interface {
	Current() Cursor
	Append(c Cursor) Cursor
	Insert(c Cursor) Cursor
	Remove(c Cursor) (prefix string, index int)
	Parent(c Cursor) Cursor
	Object(c Cursor) (obj TreeElement)
	Cursor(obj TreeElement) (cursor Cursor)
	CursorAt(start Cursor, obj TreeElement) (cursor Cursor)
	AddEntry(c Cursor, sym Symbol, text string, obj TreeElement, prop Property) (err error)
	Property(c Cursor) Property
}

type NamedTreeElementIf interface {
	TreeElement
	interfaces.Namer
}
type TreeElement interface {
	AddToTree(tree TreeIf, cursor Cursor)
	AddNewObject(tree TreeIf, cursor Cursor, obj TreeElement) (newCursor Cursor, err error)
	RemoveObject(tree TreeIf, cursor Cursor) (removed []IdWithObject)
}

type ToplevelTreeElement interface {
	TreeElement
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
	Object   TreeElement
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
