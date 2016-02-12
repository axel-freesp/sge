package model

import (
	"github.com/axel-freesp/sge/interface/tree"
)

type ModelContextIf interface {
	SignalGraphMgr() FileManagerIf
	LibraryMgr() FileManagerIf
	PlatformMgr() FileManagerIf
	MappingMgr() FileManagerMappingIf
}

type FileManagerIf interface {
	New() (tree.ToplevelTreeElementIf, error)
	Access(name string) (tree.ToplevelTreeElementIf, error)
	Remove(name string)
	Rename(oldName, newName string) error
	Store(name string) error
}

type FileManagerMappingIf interface {
	FileManagerIf
	SetGraphForNew(g interface{})
	SetPlatformForNew(p interface{})
}
