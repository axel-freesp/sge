package main

import (
	//"fmt"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	pf "github.com/axel-freesp/sge/interface/platform"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type NewElementDialog struct {
	EditMenuDialog
	selector        *gtk.ComboBoxText
	selectorChoices []elementType
}

var _ EditMenuDialogIf = (*NewElementDialog)(nil)

func NewElementDialogNew(fts *models.FilesTreeStore) (dialog *NewElementDialog, err error) {
	d, err := gtk.DialogNew()
	if err != nil {
		return
	}
	dialog = &NewElementDialog{EditMenuDialogInit(d, fts), nil, nil}
	var selector *gtk.Widget
	selector, err = dialog.init()
	if err != nil {
		return
	}
	err = dialog.CreateDialog(selector)
	if err != nil {
		return
	}
	if len(dialog.selectorChoices) > 0 {
		dialog.stack.SetVisibleChildName(string(dialog.selectorChoices[0]))
	}
	return
}

var choiceMap = map[elementType][]elementType{
	eSignalGraph:    {eNode, eInputNode, eOutputNode},
	eNode:           {eNode, eInputNode, eOutputNode},
	eNodeType:       {eNodeType, ePortType, eImplementation, eSignalType},
	ePort:           {eConnection},
	ePortType:       {ePortType},
	eConnection:     {eConnection},
	eSignalType:     {eSignalType, eNodeType},
	eLibrary:        {eSignalType, eNodeType},
	eImplementation: {eImplementation, eNode},
	ePlatform:       {eArch},
	eArch:           {eArch, eIOType, eProcess},
	eIOType:         {eIOType},
	eProcess:        {eProcess, eChannel},
	eChannel:        {eChannel},
}

// Lookup current selection in fts, choose which pageset to show.
func getSelectorChoices(fts *models.FilesTreeStore) []elementType {
	var activeElem elementType
	object, err := fts.GetObjectById(fts.GetCurrentId())
	if err != nil {
		return []elementType{}
	}
	switch object.(type) {
	case bh.SignalGraphIf, bh.SignalGraphTypeIf:
		activeElem = eSignalGraph
	case bh.NodeIf:
		activeElem = eNode
	case bh.NodeTypeIf:
		activeElem = eNodeType
	case bh.PortIf:
		activeElem = ePort
	case bh.PortTypeIf:
		activeElem = ePortType
	case bh.ConnectionIf:
		activeElem = eConnection
	case bh.SignalTypeIf:
		activeElem = eSignalType
	case bh.LibraryIf:
		activeElem = eLibrary
	case bh.ImplementationIf:
		activeElem = eImplementation
	case pf.PlatformIf:
		activeElem = ePlatform
	case pf.ArchIf:
		activeElem = eArch
	case pf.IOTypeIf:
		activeElem = eIOType
	case pf.ProcessIf:
		activeElem = eProcess
	case pf.ChannelIf:
		activeElem = eChannel
	default:
		return []elementType{}
	}
	return choiceMap[activeElem]
}

// Toplevel selector: control stack.
func comboSelectionChangedCB(dialog *NewElementDialog) {
	dialog.stack.SetVisibleChildName(dialog.selector.GetActiveText())
}

func (dialog *NewElementDialog) init() (selector *gtk.Widget, err error) {
	dialog.selector, err = gtk.ComboBoxTextNew()
	if err != nil {
		return
	}
	dialog.selectorChoices = getSelectorChoices(dialog.fts)
	for _, s := range dialog.selectorChoices {
		dialog.selector.AppendText(string(s))
	}
	dialog.selector.SetActive(0)
	dialog.selector.Connect("changed", func(t *gtk.ComboBoxText) {
		comboSelectionChangedCB(dialog)
	})
	var comboLayoutBox *gtk.Box
	comboLayoutBox, err = createLabeledRow("Select:", &dialog.selector.Widget)
	selector = &comboLayoutBox.Widget
	return
}

func (dialog *NewElementDialog) CollectJob(context string) *EditorJob {
	log.Println("Selector =", dialog.selector.GetActiveText())
	job := NewElementJobNew(context, elementType(dialog.selector.GetActiveText()))
	for _, i := range inputElementMap[job.elemType] {
		job.input[i] = dialog.readOut(i)
	}
	return EditorJobNew(JobNewElement, job)
}

func (dialog *NewElementDialog) Run(fts *models.FilesTreeStore) (job *EditorJob, ok bool) {
	context := fts.GetCurrentId()
	if context == "" {
		log.Fatal("NewElementDialog.Run error: Could not get context")
		return
	}
	ok = (gtk.ResponseType(dialog.Dialog().Run()) == gtk.RESPONSE_OK)
	if ok {
		job = dialog.CollectJob(context)
		//log.Println("editNew terminated: OK", job)
	}
	dialog.Dialog().Destroy()
	return
}
