package freesp

type Tree interface {
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

type TreeElement interface {
	AddToTree(tree Tree, cursor Cursor)
	AddNewObject(tree Tree, cursor Cursor, obj TreeElement) (newCursor Cursor)
	RemoveObject(tree Tree, cursor Cursor) (removed []IdWithObject)
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
)

type Property interface {
	IsReadOnly() bool
	MayAddObject() bool
	MayEdit() bool
	MayRemove() bool
}
