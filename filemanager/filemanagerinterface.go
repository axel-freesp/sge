package filemanager

import (
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	//"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
)

/*
 *  File handling Model
 */

type FilemanagerContextIf interface {
	interfaces.Context
	freesp.ContextIf
	ShowAll()
	FTS() freesp.TreeMgrIf
	FTV() freesp.TreeViewIf
	GVC() views.GraphViewCollection
	CleanupNodeTypesFromNodes([]freesp.NodeIf)
	CleanupSignalTypesFromNodes([]freesp.NodeIf)
	NodeTypeIsInUse(freesp.NodeTypeIf) bool
	CleanupNodeType(freesp.NodeTypeIf)
	SignalTypeIsInUse(freesp.SignalType) bool
	CleanupSignalType(freesp.SignalType)
}
