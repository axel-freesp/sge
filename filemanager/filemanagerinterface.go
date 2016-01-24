package filemanager

import (
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mod "github.com/axel-freesp/sge/interface/model"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/views"
	//pf "github.com/axel-freesp/sge/interface/platform"
)

/*
 *  File handling Model
 */

type FilemanagerContextIf interface {
	views.Context
	mod.ModelContextIf
	ShowAll()
	FTS() tr.TreeMgrIf
	FTV() tr.TreeViewIf
	GVC() views.GraphViewCollectionIf
	CleanupNodeTypesFromNodes([]bh.NodeIf)
	CleanupSignalTypesFromNodes([]bh.NodeIf)
	NodeTypeIsInUse(bh.NodeTypeIf) bool
	CleanupNodeType(bh.NodeTypeIf)
	SignalTypeIsInUse(bh.SignalTypeIf) bool
	CleanupSignalType(bh.SignalTypeIf)
}
