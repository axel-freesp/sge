package main

import (
	//"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type EditDialog struct {
	EditMenuDialog
}

var _ EditMenuDialogIf = (*EditDialog)(nil)

func EditDialogNew(fts *models.FilesTreeStore) (dialog *EditDialog, err error) {
	d, err := gtk.DialogNew()
	if err != nil {
		return
	}
	dialog = &EditDialog{EditMenuDialogInit(d, fts)}
	err = dialog.init()
	if err != nil {
		return
	}
	err = dialog.CreateDialog(nil)
	if err != nil {
		return
	}
	context, e := dialog.getActiveElementType()
	dialog.setCurrentValues(context)
	dialog.stack.SetVisibleChildName(string(e))
	return
}

func (dialog *EditDialog) setCurrentValues(context string) {
	obj, err := dialog.fts.GetObjectById(context)
	if err != nil {
		log.Fatalf("editdialog.go: getActiveElementType error: failed to get context: %s\n", err)
	}
	var i int
	var t string
	switch obj.(type) {
	case freesp.Node:
		// TODO: check auto-generated types
		if len(obj.(freesp.Node).InPorts()) > 0 {
			if len(obj.(freesp.Node).OutPorts()) > 0 {
				dialog.nodeNameEntry.SetText(obj.(freesp.Node).Name())
				for i, t = range freesp.GetRegisteredNodeTypes() {
					if obj.(freesp.Node).ItsType().TypeName() == t {
						break
					}
				}
				dialog.nodeTypeSelector.SetActive(i)
				dialog.nodeTypeSelector.SetSensitive(false)
			} else {
				// assume one input port
				dialog.outputNodeNameEntry.SetText(obj.(freesp.Node).Name())
				for i, t = range freesp.GetRegisteredSignalTypes() {
					if obj.(freesp.Node).InPorts()[0].SignalType().TypeName() == t {
						break
					}
				}
				dialog.outputTypeSelector.SetActive(i)
				dialog.outputTypeSelector.SetSensitive(false)
			}
		} else {
			// assume one output port
			dialog.inputNodeNameEntry.SetText(obj.(freesp.Node).Name())
			for i, t = range freesp.GetRegisteredSignalTypes() {
				if obj.(freesp.Node).OutPorts()[0].SignalType().TypeName() == t {
					break
				}
			}
			dialog.inputTypeSelector.SetActive(i)
			dialog.inputTypeSelector.SetSensitive(false)
		}
	case freesp.NodeType:
		dialog.typeNameEntry.SetText(obj.(freesp.NodeType).TypeName())
	case freesp.PortType:
		dialog.portNameEntry.SetText(obj.(freesp.PortType).Name())
		if obj.(freesp.PortType).Direction() == freesp.OutPort {
			dialog.directionSelector.SetActive(1)
		}
		for i, t = range freesp.GetRegisteredSignalTypes() {
			if obj.(freesp.PortType).SignalType().TypeName() == t {
				break
			}
		}
		dialog.signalTypeSelector.SetActive(i)
		dialog.directionSelector.SetSensitive(false)
	case freesp.SignalType:
		st := obj.(freesp.SignalType)
		dialog.signalTypeNameEntry.SetText(st.TypeName())
		dialog.cTypeEntry.SetText(st.CType())
		dialog.channelIdEntry.SetText(st.ChannelId())
		dialog.scopeSelector.SetActive(int(st.Scope()))
		dialog.modeSelector.SetActive(int(st.Mode()))
	case freesp.Implementation:
		dialog.implNameEntry.SetText(obj.(freesp.Implementation).ElementName())
	default:
		log.Fatalf("editdialog.go: getActiveElementType error: invalid active object type %T\n", obj)
	}
	return
}

func (dialog *EditDialog) getActiveElementType() (context string, e elementType) {
	context = dialog.fts.GetCurrentId()
	obj, err := dialog.fts.GetObjectById(context)
	if err != nil {
		log.Fatalf("editdialog.go: getActiveElementType error: failed to get context: %s\n", err)
	}
	switch obj.(type) {
	case freesp.SignalGraph:
		e = eSignalGraph
		log.Fatalf("editdialog.go: getActiveElementType error: SignalGraph is read-only\n")
	case freesp.Node:
		// TODO: check auto-generated types
		if len(obj.(freesp.Node).InPorts()) > 0 {
			if len(obj.(freesp.Node).OutPorts()) > 0 {
				e = eNode
			} else {
				// assume one input port
				e = eOutputNode
			}
		} else {
			// assume one output port
			e = eInputNode
		}
	case freesp.NodeType:
		e = eNodeType
	case freesp.Port:
		e = ePort
	case freesp.PortType:
		e = ePortType
	case freesp.Connection:
		e = eConnection
		log.Fatalf("editdialog.go: getActiveElementType error: Connection is read-only\n")
	case freesp.SignalType:
		e = eSignalType
	case freesp.Library:
		e = eLibrary
		log.Fatalf("editdialog.go: getActiveElementType error: Library is read-only\n")
	case freesp.Implementation:
		e = eImplementation
		if obj.(freesp.Implementation).ImplementationType() == freesp.NodeTypeGraph {
			log.Fatalf("editdialog.go: getActiveElementType error: Implementation/graph is read-only\n")
		}
	default:
		log.Fatalf("editdialog.go: getActiveElementType error: invalid active object type %T\n", obj)
	}
	return
}

func (dialog *EditDialog) init() (err error) {
	return
}

func (dialog *EditDialog) CollectJob() *EditorJob {
	context, e := dialog.getActiveElementType()
	job := EditJobNew(context, e)
	for _, i := range inputElementMap[job.elemType] {
		job.detail[i] = dialog.readOut(i)
	}
	return EditorJobNew(JobEdit, job)
}

func (dialog *EditDialog) Run(fts *models.FilesTreeStore) (job *EditorJob, ok bool) {
	ok = (gtk.ResponseType(dialog.Dialog().Run()) == gtk.RESPONSE_OK)
	if ok {
		job = dialog.CollectJob()
		//log.Println("editEdit terminated: OK", job)
	}
	dialog.Dialog().Destroy()
	return
}
