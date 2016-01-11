package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// Menu control wants to see this:
type EditMenuDialogIf interface {
	Run(fts *models.FilesTreeStore) (job *EditorJob, ok bool)
}

type inputElement string

const (
	iNodeName           inputElement = "NodeName"
	iInputNodeName                   = "InputNodeName"
	iOutputNodeName                  = "OutputNodeName"
	iTypeName                        = "TypeName"
	iPortName                        = "PortName"
	iImplName                        = "ImplName"
	iSignalTypeName                  = "SignalTypeName"
	iNodeTypeSelect                  = "NodeTypeSelect"
	iSignalTypeSelect                = "SignalTypeSelect"
	iInputTypeSelect                 = "InputTypeSelect"
	iOutputTypeSelect                = "OutputTypeSelect"
	iImplementationType              = "ImplementationType"
	iPortSelect                      = "PortSelect"
	iCType                           = "CType"
	iChannelId                       = "ChannelId"
	iScope                           = "Scope"
	iSignalMode                      = "SignalMode"
	iDirection                       = "Direction"
	iChannelDirection                = "ChannelDirection"
	iIOTypeSelect                    = "IOTypeSelect"
	iChannelLinkSelect               = "ChannelLinkSelect"
	iIOTypeName                      = "IOTypeName"
	iIOModeSelect                    = "IOModeSelect"
	iProcessName                     = "ProcessName"
	iArchName                        = "ArchName"
)

type elementType string

const (
	eSignalGraph    elementType = "SignalGraph"
	eNode                       = "Node"
	eInputNode                  = "InputNode"
	eOutputNode                 = "OutputNode"
	eNodeType                   = "NodeType"
	ePort                       = "Port"
	ePortType                   = "PortType"
	eConnection                 = "Connection"
	eSignalType                 = "SignalType"
	eLibrary                    = "Library"
	eImplementation             = "Implementation"
	ePlatform                   = "Platform"
	eArch                       = "Arch"
	eProcess                    = "Process"
	eIOType                     = "IOType"
	eChannel                    = "Channel"
)

var inputElementMap = map[elementType][]inputElement{
	eNode:           {iNodeName, iNodeTypeSelect},
	eInputNode:      {iInputNodeName, iInputTypeSelect},
	eOutputNode:     {iOutputNodeName, iOutputTypeSelect},
	eNodeType:       {iTypeName},
	ePortType:       {iPortName, iSignalTypeSelect, iDirection},
	eConnection:     {iPortSelect},
	eSignalType:     {iSignalTypeName, iCType, iChannelId, iScope, iSignalMode},
	eImplementation: {iImplName, iImplementationType},
	eChannel:        {iChannelDirection, iIOTypeSelect, iChannelLinkSelect},
	eIOType:         {iIOTypeName, iIOModeSelect},
	eProcess:        {iProcessName},
	eArch:           {iArchName},
}

var scopeStrings = []string{"Local", "Global"}
var modeStrings = []string{"Isochronous", "Asynchronous"}
var directionStrings = []string{"InPort", "OutPort"}
var implTypeStrings = []string{"Elementary Type", "Signal Graph"}
var ioModeStrings = []string{
	string(interfaces.IOModeShmem),
	string(interfaces.IOModeAsync),
	string(interfaces.IOModeSync),
}

var scope2string = map[freesp.Scope]string{
	freesp.Local:  scopeStrings[freesp.Local],
	freesp.Global: scopeStrings[freesp.Global],
}

var mode2string = map[freesp.Mode]string{
	freesp.Synchronous:  modeStrings[freesp.Synchronous],
	freesp.Asynchronous: modeStrings[freesp.Asynchronous],
}

var direction2string = map[interfaces.PortDirection]string{
	interfaces.InPort:  directionStrings[0],
	interfaces.OutPort: directionStrings[1],
}

var implType2string = map[freesp.ImplementationType]string{
	freesp.NodeTypeElement: implTypeStrings[freesp.NodeTypeElement],
	freesp.NodeTypeGraph:   implTypeStrings[freesp.NodeTypeGraph],
}

var string2scope = map[string]freesp.Scope{
	"Local":  freesp.Local,
	"Global": freesp.Global,
}

var string2mode = map[string]freesp.Mode{
	"Isochronous":  freesp.Synchronous,
	"Asynchronous": freesp.Asynchronous,
}

var string2direction = map[string]interfaces.PortDirection{
	"InPort":  interfaces.InPort,
	"OutPort": interfaces.OutPort,
}

var string2implType = map[string]freesp.ImplementationType{
	implTypeStrings[freesp.NodeTypeElement]: freesp.NodeTypeElement,
	implTypeStrings[freesp.NodeTypeGraph]:   freesp.NodeTypeGraph,
}

type EditMenuDialog struct {
	dialog *gtk.Dialog
	fts    *models.FilesTreeStore
	stack  *gtk.Stack

	nodeTypeSelector         *gtk.ComboBoxText
	signalTypeSelector       *gtk.ComboBoxText
	inputTypeSelector        *gtk.ComboBoxText
	outputTypeSelector       *gtk.ComboBoxText
	scopeSelector            *gtk.ComboBoxText
	modeSelector             *gtk.ComboBoxText
	directionSelector        *gtk.ComboBoxText
	implementationSelector   *gtk.ComboBoxText
	portSelector             *gtk.ComboBoxText
	channelDirectionSelector *gtk.ComboBoxText
	ioTypeSelector           *gtk.ComboBoxText
	processSelector          *gtk.ComboBoxText
	ioModeSelector           *gtk.ComboBoxText

	nodeNameEntry       *gtk.Entry
	inputNodeNameEntry  *gtk.Entry
	outputNodeNameEntry *gtk.Entry
	typeNameEntry       *gtk.Entry
	portNameEntry       *gtk.Entry
	implNameEntry       *gtk.Entry
	signalTypeNameEntry *gtk.Entry
	cTypeEntry          *gtk.Entry
	channelIdEntry      *gtk.Entry
	ioTypeNameEntry     *gtk.Entry
	processNameEntry    *gtk.Entry
	archNameEntry       *gtk.Entry
}

func EditMenuDialogInit(d *gtk.Dialog, fts *models.FilesTreeStore) (ret EditMenuDialog) {
	ret = EditMenuDialog{}
	ret.dialog = d
	ret.fts = fts
	return
}

func (dialog *EditMenuDialog) Dialog() *gtk.Dialog {
	return dialog.dialog
}

func (dialog *EditMenuDialog) CreateDialog(extra *gtk.Widget) (err error) {
	d := dialog.dialog
	d.SetTitle("New Element")

	var active string
	active, err = dialog.fts.GetValueById(dialog.fts.GetCurrentId())
	if err != nil {
		return
	}

	var box *gtk.Box
	box, err = d.GetContentArea()
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

	if extra != nil {
		box.PackStart(extra, false, false, 6)
	}

	dialog.stack, err = gtk.StackNew()
	if err != nil {
		return
	}
	err = dialog.fillStack()
	if err != nil {
		return
	}
	box.PackStart(dialog.stack, true, true, 6)

	d.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	d.AddButton("OK", gtk.RESPONSE_OK)
	d.SetDefaultResponse(gtk.RESPONSE_OK)
	d.ShowAll()
	return
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
		if thisPort.Direction() == interfaces.InPort {
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

func getOtherProcesses(fts *models.FilesTreeStore, object freesp.TreeElement) (ret []freesp.Process) {
	var thisProcess freesp.Process
	switch object.(type) {
	case freesp.Channel:
		thisProcess = object.(freesp.Channel).Process()
	case freesp.Process:
		thisProcess = object.(freesp.Process)
	default:
		log.Fatalf("getOtherProcesses error: invalid type %T\n", object)
	}
	platform := thisProcess.Arch().Platform()
	for _, a := range platform.Arch() {
		for _, p := range a.Processes() {
			if p != thisProcess {
				ret = append(ret, p)
			}
		}
	}
	return
}

type inputElementHandling struct {
	label         string
	readOut       func(dialog *EditMenuDialog) string
	createElement func(dialog *EditMenuDialog) (obj *gtk.Widget, err error)
}

var inputHandling = map[inputElement]inputElementHandling{
	iNodeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.nodeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.nodeNameEntry)
		},
	},
	iInputNodeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.inputNodeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.inputNodeNameEntry)
		},
	},
	iOutputNodeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.outputNodeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.outputNodeNameEntry)
		},
	},
	iTypeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.typeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.typeNameEntry)
		},
	},
	iPortName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.portNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.portNameEntry)
		},
	},
	iImplName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.implNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.implNameEntry)
		},
	},
	iSignalTypeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.signalTypeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.signalTypeNameEntry)
		},
	},
	iIOTypeName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.ioTypeNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.ioTypeNameEntry)
		},
	},
	iProcessName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.processNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.processNameEntry)
		},
	},
	iArchName: {"Name:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.archNameEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.archNameEntry)
		},
	},
	iNodeTypeSelect: {"Select node type:",
		func(dialog *EditMenuDialog) string {
			return dialog.nodeTypeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.nodeTypeSelector, freesp.GetRegisteredNodeTypes())
		},
	},
	iSignalTypeSelect: {"Select signal type:",
		func(dialog *EditMenuDialog) string {
			return dialog.signalTypeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.signalTypeSelector, freesp.GetRegisteredSignalTypes())
		},
	},
	iInputTypeSelect: {"Select signal type:",
		func(dialog *EditMenuDialog) string {
			return dialog.inputTypeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.inputTypeSelector, freesp.GetRegisteredSignalTypes())
		},
	},
	iOutputTypeSelect: {"Select signal type:",
		func(dialog *EditMenuDialog) string {
			return dialog.outputTypeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.outputTypeSelector, freesp.GetRegisteredSignalTypes())
		},
	},
	iIOTypeSelect: {"Select IO type:",
		func(dialog *EditMenuDialog) string {
			return dialog.ioTypeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.ioTypeSelector, freesp.GetRegisteredIOTypes())
		},
	},
	iImplementationType: {"Implementation type:",
		func(dialog *EditMenuDialog) string {
			return dialog.implementationSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.implementationSelector, implTypeStrings)
		},
	},
	iIOModeSelect: {"Transmission mode:",
		func(dialog *EditMenuDialog) string {
			return dialog.ioModeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.ioModeSelector, ioModeStrings)
		},
	},
	iPortSelect: {"Select port to connect:",
		func(dialog *EditMenuDialog) string {
			return dialog.portSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			fts := dialog.fts
			var choices []string
			object := fts.Object(fts.Current())
			switch object.(type) {
			case freesp.Port:
			case freesp.Connection:
				object = fts.Object(fts.Parent(fts.Current()))
			default:
				return
			}
			for _, p := range getMatchingPorts(dialog.fts, object) {
				choices = append(choices, fmt.Sprintf("%s/%s", p.Node().Name(), p.Name()))
			}
			return newComboBox(&dialog.portSelector, choices)
		},
	},
	iChannelLinkSelect: {"Select arch/process to connect:",
		func(dialog *EditMenuDialog) string {
			return dialog.processSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			fts := dialog.fts
			var choices []string
			object := fts.Object(fts.Current())
			switch object.(type) {
			case freesp.Process:
			case freesp.Channel:
				object = fts.Object(fts.Parent(fts.Current()))
			default:
				return
			}
			for _, p := range getOtherProcesses(dialog.fts, object) {
				choices = append(choices, fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name()))
			}
			return newComboBox(&dialog.processSelector, choices)
		},
	},
	iCType: {"C type:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.cTypeEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.cTypeEntry)
		},
	},
	iChannelId: {"Channel id:",
		func(dialog *EditMenuDialog) string {
			return getText(dialog.channelIdEntry)
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newEntry(&dialog.channelIdEntry)
		},
	},
	iScope: {"Scope:",
		func(dialog *EditMenuDialog) string {
			return dialog.scopeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.scopeSelector, scopeStrings)
		},
	},
	iSignalMode: {"Transfer mode:",
		func(dialog *EditMenuDialog) string {
			return dialog.modeSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.modeSelector, modeStrings)
		},
	},
	iDirection: {"Direction:",
		func(dialog *EditMenuDialog) string {
			return dialog.directionSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.directionSelector, directionStrings)
		},
	},
	iChannelDirection: {"Direction:",
		func(dialog *EditMenuDialog) string {
			return dialog.channelDirectionSelector.GetActiveText()
		},
		func(dialog *EditMenuDialog) (obj *gtk.Widget, err error) {
			return newComboBox(&dialog.channelDirectionSelector, directionStrings)
		},
	},
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

// Construct input element (some depending on fts).
func (dialog *EditMenuDialog) createInputElement(i inputElement) (widget *gtk.Widget, err error) {
	obj, err := inputHandling[i].createElement(dialog)
	row, err := createLabeledRow(inputHandling[i].label, obj)
	if err != nil {
		return
	}
	widget = &row.Widget
	return
}

// Read out value from input element
func (dialog *EditMenuDialog) readOut(i inputElement) string {
	return inputHandling[i].readOut(dialog)
}

// Collect all input elements to be constructed.
func (dialog *EditMenuDialog) fillBox(box *gtk.Box, e elementType) error {
	for _, i := range inputElementMap[e] {
		widget, err := dialog.createInputElement(i)
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
	eInputNode,
	eOutputNode,
	eNodeType,
	ePortType,
	eConnection,
	eSignalType,
	eImplementation,
	eArch,
	eIOType,
	eProcess,
	eChannel,
}

func (dialog *EditMenuDialog) fillStack() error {
	for _, a := range stackAlternatives {
		box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
		if err != nil {
			return err
		}
		err = dialog.fillBox(box, a)
		if err != nil {
			return err
		}
		dialog.stack.AddNamed(box, string(a))
	}
	return nil
}
