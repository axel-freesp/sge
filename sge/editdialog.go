package main

import (
	//"fmt"
	"github.com/axel-freesp/sge/freesp"
	//tr "github.com/axel-freesp/sge/interface/tree"
	interfaces "github.com/axel-freesp/sge/interface"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	mp "github.com/axel-freesp/sge/interface/mapping"
	pf "github.com/axel-freesp/sge/interface/platform"
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
	case bh.NodeIf:
		// TODO: check auto-generated types
		if len(obj.(bh.NodeIf).InPorts()) > 0 {
			if len(obj.(bh.NodeIf).OutPorts()) > 0 {
				dialog.nodeNameEntry.SetText(obj.(bh.NodeIf).Name())
				for i, t = range freesp.GetRegisteredNodeTypes() {
					if obj.(bh.NodeIf).ItsType().TypeName() == t {
						break
					}
				}
				dialog.nodeTypeSelector.SetActive(i)
				dialog.nodeTypeSelector.SetSensitive(false)
			} else {
				// assume one input port
				dialog.outputNodeNameEntry.SetText(obj.(bh.NodeIf).Name())
				for i, t = range freesp.GetRegisteredSignalTypes() {
					if obj.(bh.NodeIf).InPorts()[0].SignalType().TypeName() == t {
						break
					}
				}
				dialog.outputTypeSelector.SetActive(i)
				dialog.outputTypeSelector.SetSensitive(false)
			}
		} else {
			// assume one output port
			dialog.inputNodeNameEntry.SetText(obj.(bh.NodeIf).Name())
			for i, t = range freesp.GetRegisteredSignalTypes() {
				if obj.(bh.NodeIf).OutPorts()[0].SignalType().TypeName() == t {
					break
				}
			}
			dialog.inputTypeSelector.SetActive(i)
			dialog.inputTypeSelector.SetSensitive(false)
		}
	case bh.NodeTypeIf:
		dialog.typeNameEntry.SetText(obj.(bh.NodeTypeIf).TypeName())
	case bh.PortType:
		dialog.portNameEntry.SetText(obj.(bh.PortType).Name())
		if obj.(bh.PortType).Direction() == interfaces.OutPort {
			dialog.directionSelector.SetActive(1)
		}
		for i, t = range freesp.GetRegisteredSignalTypes() {
			if obj.(bh.PortType).SignalType().TypeName() == t {
				break
			}
		}
		dialog.signalTypeSelector.SetActive(i)
		dialog.directionSelector.SetSensitive(false)
	case bh.SignalType:
		st := obj.(bh.SignalType)
		dialog.signalTypeNameEntry.SetText(st.TypeName())
		dialog.cTypeEntry.SetText(st.CType())
		dialog.channelIdEntry.SetText(st.ChannelId())
		dialog.scopeSelector.SetActive(int(st.Scope()))
		dialog.modeSelector.SetActive(int(st.Mode()))
	case bh.ImplementationIf:
		dialog.implNameEntry.SetText(obj.(bh.ImplementationIf).ElementName())
	case pf.ArchIf:
		dialog.archNameEntry.SetText(obj.(pf.ArchIf).Name())
	case pf.ProcessIf:
		dialog.processNameEntry.SetText(obj.(pf.ProcessIf).Name())
	case pf.IOTypeIf:
		dialog.ioTypeNameEntry.SetText(obj.(pf.IOTypeIf).Name())
		for i, t = range ioModeStrings {
			if string(obj.(pf.IOTypeIf).IOMode()) == t {
				break
			}
		}
		dialog.ioModeSelector.SetActive(i)
	case pf.ChannelIf:
		if obj.(pf.ChannelIf).Direction() == interfaces.OutPort {
			dialog.channelDirectionSelector.SetActive(1)
		}
		dialog.channelDirectionSelector.SetSensitive(false)
		for i, t = range freesp.GetRegisteredIOTypes() {
			if obj.(pf.ChannelIf).IOType().Name() == t {
				break
			}
		}
		dialog.ioTypeSelector.SetActive(i)
		pr := getOtherProcesses(dialog.fts, obj)
		var p pf.ProcessIf
		for i, p = range pr {
			if obj.(pf.ChannelIf).Link().Process() == p {
				break
			}
		}
		dialog.processSelector.SetActive(i)
		dialog.processSelector.SetSensitive(false)
	case mp.MappedElementIf:
		pr := obj.(mp.MappedElementIf).Process()
		i := 0
		for _, a := range obj.(mp.MappedElementIf).Mapping().Platform().Arch() {
			for _, p := range a.Processes() {
				if pr == p {
					dialog.processMapSelector.SetActive(i)
					return
				}
				i++
			}
		}
		dialog.processMapSelector.SetActive(i)
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
	case bh.SignalGraphIf:
		e = eSignalGraph
		log.Fatalf("editdialog.go: getActiveElementType error: SignalGraphIf is read-only\n")
	case bh.NodeIf:
		// TODO: check auto-generated types
		if len(obj.(bh.NodeIf).InPorts()) > 0 {
			if len(obj.(bh.NodeIf).OutPorts()) > 0 {
				e = eNode
			} else {
				// assume one input port
				e = eOutputNode
			}
		} else {
			// assume one output port
			e = eInputNode
		}
	case bh.NodeTypeIf:
		e = eNodeType
	case bh.Port:
		e = ePort
	case bh.PortType:
		e = ePortType
	case bh.Connection:
		e = eConnection
		log.Fatalf("editdialog.go: getActiveElementType error: Connection is read-only\n")
	case bh.SignalType:
		e = eSignalType
	case bh.LibraryIf:
		e = eLibrary
		log.Fatalf("editdialog.go: getActiveElementType error: LibraryIf is read-only\n")
	case bh.ImplementationIf:
		e = eImplementation
		if obj.(bh.ImplementationIf).ImplementationType() == bh.NodeTypeGraph {
			log.Fatalf("editdialog.go: getActiveElementType error: ImplementationIf/graph is read-only\n")
		}
	case pf.ArchIf:
		e = eArch
	case pf.ProcessIf:
		e = eProcess
	case pf.IOTypeIf:
		e = eIOType
	case pf.ChannelIf:
		e = eChannel
	case mp.MappedElementIf:
		e = eMapElement
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
