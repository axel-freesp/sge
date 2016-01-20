package filemanager

import (
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
	"github.com/axel-freesp/sge/views"
)

/*
 *  File handling Model
 */

type FilemanagerContextIf interface {
	interfaces.Context
	freesp.Context
	ShowAll()
	FTS() *models.FilesTreeStore
	FTV() *views.FilesTreeView
	GVC() views.GraphViewCollection
	CleanupNodeTypesFromNodes([]freesp.Node)
	CleanupSignalTypesFromNodes([]freesp.Node)
	NodeTypeIsInUse(freesp.NodeType) bool
	CleanupNodeType(freesp.NodeType)
	SignalTypeIsInUse(freesp.SignalType) bool
	CleanupSignalType(freesp.SignalType)
}
