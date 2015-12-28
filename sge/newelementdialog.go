package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

type inputElement string

const (
	iNodeName           inputElement = "NodeName"
	iTypeName                        = "TypeName"
	iPortName                        = "PortName"
	iImplName                        = "ImplName"
	iSignalTypeName                  = "SignalTypeName"
	iNodeTypeSelect                  = "NodeTypeSelect"
	iSignalTypeSelect                = "SignalTypeSelect"
	iImplementationType              = "ImplementationType"
	iPortSelect                      = "PortSelect"
	iCType                           = "CType"
	iChannelId                       = "ChannelId"
	iScope                           = "Scope"
	iSignalMode                      = "SignalMode"
	iDirection                       = "Direction"
)

type elementType string

const (
	eSignalGraph    elementType = "SignalGraph"
	eNode                       = "Node"
	eNodeType                   = "NodeType"
	ePort                       = "Port"
	ePortType                   = "NamedPortType"
	eConnection                 = "Connection"
	eSignalType                 = "SignalType"
	eLibrary                    = "Library"
	eImplementation             = "Implementation"
)

var choiceMap = map[elementType][]elementType{
	eSignalGraph:    {eNode},
	eNode:           {eNode},
	eNodeType:       {eNodeType, ePortType, eImplementation},
	ePort:           {eConnection},
	ePortType:       {ePortType},
	eConnection:     {eConnection},
	eSignalType:     {eSignalType},
	eLibrary:        {eSignalType, eNodeType},
	eImplementation: {eImplementation, eNode},
}

var inputElementMap = map[elementType][]inputElement{
	eNode:           {iNodeName, iNodeTypeSelect},
	eNodeType:       {iTypeName},
	ePortType:       {iPortName, iSignalTypeSelect, iDirection},
	eConnection:     {iPortSelect},
	eSignalType:     {iSignalTypeName, iCType, iChannelId, iScope, iSignalMode},
	eImplementation: {iImplName, iImplementationType},
}

type NewElementDialog struct {
	dialog   *gtk.Dialog
	fts      *models.FilesTreeStore
	selector *gtk.ComboBoxText
	stack    *gtk.Stack

	nodeTypeSelector       *gtk.ComboBoxText
	signalTypeSelector     *gtk.ComboBoxText
	scopeSelector          *gtk.ComboBoxText
	modeSelector           *gtk.ComboBoxText
	directionSelector      *gtk.ComboBoxText
	implementationSelector *gtk.ComboBoxText
	portSelector           *gtk.ComboBoxText

	nodeNameEntry       *gtk.Entry
	typeNameEntry       *gtk.Entry
	portNameEntry       *gtk.Entry
	implNameEntry       *gtk.Entry
	signalTypeNameEntry *gtk.Entry
	cTypeEntry          *gtk.Entry
	channelIdEntry      *gtk.Entry
}

func NewElementDialogNew(fts *models.FilesTreeStore) (dialog *NewElementDialog, err error) {
	d, err := gtk.DialogNew()
	if err != nil {
		return
	}
	dialog = &NewElementDialog{d, fts, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
	err = dialog.init(fts)
	return
}

// Lookup current selection in fts, choose which pageset to show.
func getSelectorChoices(fts *models.FilesTreeStore) []elementType {
	var activeElem elementType
	object, err := fts.GetObjectById(fts.GetCurrentId())
	if err != nil {
		return []elementType{}
	}
	switch object.(type) {
	case freesp.SignalGraph, freesp.SignalGraphType:
		activeElem = eSignalGraph
	case freesp.Node:
		activeElem = eNode
	case freesp.NodeType:
		activeElem = eNodeType
	case freesp.Port:
		activeElem = ePort
	case freesp.PortType:
		activeElem = ePortType
	case freesp.Connection:
		activeElem = eConnection
	case freesp.SignalType:
		activeElem = eSignalType
	case freesp.Library:
		activeElem = eLibrary
	case freesp.Implementation:
		activeElem = eImplementation
	default:
		return []elementType{}
	}
	return choiceMap[activeElem]
}

// Toplevel selector: control stack.
func comboSelectionChangedCB(dialog *NewElementDialog) {
	dialog.stack.SetVisibleChildName(dialog.selector.GetActiveText())
}

// For connection only: lookup matching ports to connect.
func getMatchingPorts(fts *models.FilesTreeStore, object freesp.TreeElement) (ret []freesp.Port) {
	var thisPort freesp.Port
	switch object.(type) {
	case freesp.Port:
		thisPort = object.(freesp.Port)
	case freesp.Connection:
		log.Fatal("getMatchingPorts error: expecting Port, not Connection")
	default:
		log.Fatal("getMatchingPorts error: expecting Port")
	}
	thisNode := thisPort.Node()
	graph := thisNode.Context()
	for _, n := range graph.Nodes() {
		var ports []freesp.Port
		if thisPort.Direction() == freesp.InPort {
			ports = n.OutPorts()
		} else {
			ports = n.InPorts()
		}
		for _, p := range ports {
			if p.SignalType().TypeName() == thisPort.SignalType().TypeName() {
				ret = append(ret, p)
			}
		}
	}
	return
}

type inputElementHandling struct {
	label         string
	readOut       func(dialog *NewElementDialog) string
	createElement func(dialog *NewElementDialog) (obj *gtk.Widget, err error)
}

func getText(entry *gtk.Entry) string {
	text, _ := entry.GetText()
	return text
}

func newEntry(entryP **gtk.Entry) (obj *gtk.Widget, err error) {
	*entryP, err = gtk.EntryNew()
	if err != nil {
		return
	}
	obj = &(*entryP).Widget
	return
}

func newComboBox(comboBoxP **gtk.ComboBoxText, choices []string) (obj *gtk.Widget, err error) {
	*comboBoxP, err = gtk.ComboBoxTextNew()
	if err != nil {
		return
	}
	for _, s := range choices {
		(*comboBoxP).AppendText(s)
	}
	(*comboBoxP).SetActive(0)
	obj = &(*comboBoxP).Widget
	return
}

var inputHandling = map[inputElement]inputElementHandling{
	iNodeName: {"Name:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.nodeNameEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.nodeNameEntry)
		},
	},
	iTypeName: {"Name:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.typeNameEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.typeNameEntry)
		},
	},
	iPortName: {"Name:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.portNameEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.portNameEntry)
		},
	},
	iImplName: {"Name:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.implNameEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.implNameEntry)
		},
	},
	iSignalTypeName: {"Name:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.signalTypeNameEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.signalTypeNameEntry)
		},
	},
	iNodeTypeSelect: {"Select node type:",
		func(dialog *NewElementDialog) string {
			return dialog.nodeTypeSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.nodeTypeSelector, freesp.GetRegisteredNodeTypes())
		},
	},
	iSignalTypeSelect: {"Select signal type:",
		func(dialog *NewElementDialog) string {
			return dialog.signalTypeSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.signalTypeSelector, freesp.GetRegisteredSignalTypes())
		},
	},
	iImplementationType: {"Implementation type:",
		func(dialog *NewElementDialog) string {
			return dialog.implementationSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.implementationSelector, []string{"Elementary Type", "Signal Graph"})
		},
	},
	iPortSelect: {"Select port to connect:",
		func(dialog *NewElementDialog) string {
			return dialog.portSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			var choices []string
			object, err := dialog.fts.GetObjectById(dialog.fts.GetCurrentId())
			if err != nil {
				log.Fatalf("Internal error: FileTreeStore.GetObjectById(GetCurrentId()) failed\n")
			}
			switch object.(type) {
			case freesp.Port:
			case freesp.Connection:
				object = dialog.fts.Object(dialog.fts.Parent(freesp.Cursor{dialog.fts.GetCurrentId(), -1}))
			default:
				return
			}
			for _, p := range getMatchingPorts(dialog.fts, object) {
				choices = append(choices, fmt.Sprintf("%s/%s", p.Node().Name(), p.Name()))
			}
			return newComboBox(&dialog.portSelector, choices)
		},
	},
	iCType: {"C type:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.cTypeEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.cTypeEntry)
		},
	},
	iChannelId: {"Channel id:",
		func(dialog *NewElementDialog) string {
			return getText(dialog.channelIdEntry)
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.channelIdEntry)
		},
	},
	iScope: {"Scope:",
		func(dialog *NewElementDialog) string {
			return dialog.scopeSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.scopeSelector, []string{"Global", "Local"})
		},
	},
	iSignalMode: {"Transfer mode:",
		func(dialog *NewElementDialog) string {
			return dialog.modeSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.modeSelector, []string{"Asynchronous", "Isochronous"})
		},
	},
	iDirection: {"Direction:",
		func(dialog *NewElementDialog) string {
			return dialog.directionSelector.GetActiveText()
		},
		func(dialog *NewElementDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.directionSelector, []string{"InPort", "OutPort"})
		},
	},
}

// Construct input element (some depending on fts).
func (dialog *NewElementDialog) createInputElement(i inputElement, fts *models.FilesTreeStore) (widget *gtk.Widget, err error) {
	obj, err := inputHandling[i].createElement(dialog)
	row, err := createLabeledRow(inputHandling[i].label, obj)
	if err != nil {
		return
	}
	widget = &row.Widget
	return
}

// Read out value from input element
func (dialog *NewElementDialog) readOut(i inputElement) string {
	return inputHandling[i].readOut(dialog)
}

// Collect all input elements to be constructed.
func (dialog *NewElementDialog) fillBox(box *gtk.Box, e elementType, fts *models.FilesTreeStore) error {
	for _, i := range inputElementMap[e] {
		widget, err := dialog.createInputElement(i, fts)
		if err != nil {
			return err
		}
		box.Add(widget)
	}
	return nil
}

// List of keys into the named stack
var stackAlternatives = []elementType{
	eNode,
	eNodeType,
	ePortType,
	eConnection,
	eSignalType,
	eImplementation,
}

func (dialog *NewElementDialog) fillStack(fts *models.FilesTreeStore) error {
	for _, a := range stackAlternatives {
		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
		if err != nil {
			return err
		}
		err = dialog.fillBox(box, a, fts)
		if err != nil {
			return err
		}
		dialog.stack.AddNamed(box, string(a))
	}
	return nil
}

func (dialog *NewElementDialog) init(fts *models.FilesTreeStore) (err error) {
	d := dialog.dialog
	d.SetTitle("New Element")

	active, err := fts.GetValueById(fts.GetCurrentId())
	if err != nil {
		return
	}

	box, err := d.GetContentArea()
	if err != nil {
		return
	}

	activeText, err := gtk.LabelNew(fmt.Sprintf("%s", active))
	if err != nil {
		return
	}
	activeLayoutBox, err := createLabeledRow("Active element:", &activeText.Widget)
	if err != nil {
		return
	}
	box.PackStart(activeLayoutBox, false, false, 6)

	dialog.selector, err = gtk.ComboBoxTextNew()
	if err != nil {
		return
	}
	selectorChoices := getSelectorChoices(fts)
	for _, s := range selectorChoices {
		dialog.selector.AppendText(string(s))
	}
	dialog.selector.SetActive(0)
	dialog.selector.Connect("changed", func(t *gtk.ComboBoxText) {
		comboSelectionChangedCB(dialog)
	})
	comboLayoutBox, err := createLabeledRow("Select:", &dialog.selector.Widget)
	if err != nil {
		return
	}
	box.PackStart(comboLayoutBox, false, false, 6)

	dialog.stack, err = gtk.StackNew()
	if err != nil {
		return
	}
	err = dialog.fillStack(fts)
	if err != nil {
		return
	}
	box.PackStart(dialog.stack, true, true, 6)

	d.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	d.AddButton("OK", gtk.RESPONSE_OK)
	d.SetDefaultResponse(gtk.RESPONSE_OK)
	d.ShowAll()
	if len(selectorChoices) > 0 {
		dialog.stack.SetVisibleChildName(string(selectorChoices[0]))
	}
	return
}

func (dialog *NewElementDialog) collectJob(context string) *EditorJob {
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
	ok = (gtk.ResponseType(dialog.dialog.Run()) == gtk.RESPONSE_OK)
	if ok {
		job = dialog.collectJob(context)
		log.Println("editNew terminated: OK", job)
	}
	dialog.dialog.Destroy()
	return
}
