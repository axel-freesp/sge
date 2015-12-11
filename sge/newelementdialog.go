package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
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
	eNamedPortType              = "NamedPortType"
	eConnection                 = "Connection"
	eSignalType                 = "SignalType"
	eLibrary                    = "Library"
	eImplementation             = "Implementation"
)

var choiceMap = map[elementType][]elementType{
	eSignalGraph:    {eNode},
	eNode:           {eNode},
	eNodeType:       {eNodeType, eNamedPortType, eImplementation},
	ePort:           {eConnection},
	eNamedPortType:  {eNamedPortType},
	eConnection:     {eConnection},
	eSignalType:     {eSignalType},
	eLibrary:        {eSignalType, eNodeType},
	eImplementation: {eImplementation},
}

var inputElementMap = map[elementType][]inputElement{
	eNode:           {iNodeName, iNodeTypeSelect},
	eNodeType:       {iTypeName},
	eNamedPortType:  {iPortName, iSignalTypeSelect, iDirection},
	eConnection:     {iPortSelect},
	eSignalType:     {iSignalTypeName, iCType, iChannelId, iScope, iSignalMode},
	eImplementation: {iImplName, iImplementationType},
}

type NewElementDialog struct {
	dialog   *gtk.Dialog
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
	dialog = &NewElementDialog{d, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
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
	case freesp.NamedPortType:
		activeElem = eNamedPortType
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
func getMatchingPorts(fts *models.FilesTreeStore) (ret []freesp.Port) {
	object, err := fts.GetObjectById(fts.GetCurrentId())
	if err != nil {
		return []freesp.Port{}
	}
	var thisPort freesp.Port
	switch object.(type) {
	case freesp.Port:
		thisPort = object.(freesp.Port)
	case freesp.Connection:
		conn := object.(freesp.Connection)
		thisPort = conn.From
	default:
		return []freesp.Port{}
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
			if p.ItsType().SignalType() == thisPort.ItsType().SignalType() {
				ret = append(ret, p)
			}
		}
	}
	return
}

// Construct input element (some depending on fts).
func (dialog *NewElementDialog) createInputElement(i inputElement, fts *models.FilesTreeStore) (widget *gtk.Widget, err error) {
	var label string
	var obj *gtk.Widget
	switch i {
	case iNodeName:
		label = "Name:"
		dialog.nodeNameEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.nodeNameEntry.Widget
	case iTypeName:
		label = "Name:"
		dialog.typeNameEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.typeNameEntry.Widget
	case iPortName:
		label = "Name:"
		dialog.portNameEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.portNameEntry.Widget
	case iImplName:
		label = "Name:"
		dialog.implNameEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.implNameEntry.Widget
	case iSignalTypeName:
		label = "Name:"
		dialog.signalTypeNameEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.signalTypeNameEntry.Widget
	case iNodeTypeSelect:
		label = "Select node type:"
		registeredTypes := freesp.GetRegisteredNodeTypes()
		dialog.nodeTypeSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range registeredTypes {
			dialog.nodeTypeSelector.AppendText(s)
		}
		dialog.nodeTypeSelector.SetActive(0)
		obj = &dialog.nodeTypeSelector.Widget
	case iSignalTypeSelect:
		label = "Select signal type:"
		registeredTypes := freesp.GetRegisteredSignalTypes()
		dialog.signalTypeSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range registeredTypes {
			dialog.signalTypeSelector.AppendText(s)
		}
		dialog.signalTypeSelector.SetActive(0)
		obj = &dialog.signalTypeSelector.Widget
	case iImplementationType:
		label = "Implementation type:"
		dialog.implementationSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range []string{"Elementary Type", "Signal Graph"} {
			dialog.implementationSelector.AppendText(s)
		}
		dialog.implementationSelector.SetActive(0)
		obj = &dialog.implementationSelector.Widget
	case iPortSelect:
		label = "Select port to connect:"
		ports := getMatchingPorts(fts)
		dialog.portSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, p := range ports {
			s := fmt.Sprintf("%s/%s", p.Node().NodeName(), p.PortName())
			dialog.portSelector.AppendText(s)
		}
		dialog.portSelector.SetActive(0)
		obj = &dialog.portSelector.Widget
	case iCType:
		label = "C type:"
		dialog.cTypeEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.cTypeEntry.Widget
	case iChannelId:
		label = "Channel id:"
		dialog.channelIdEntry, err = gtk.EntryNew()
		if err != nil {
			return
		}
		obj = &dialog.channelIdEntry.Widget
	case iScope:
		label = "Scope:"
		dialog.scopeSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range []string{"Global", "Local"} {
			dialog.scopeSelector.AppendText(s)
		}
		dialog.scopeSelector.SetActive(0)
		obj = &dialog.scopeSelector.Widget
	case iSignalMode:
		label = "Transfer mode:"
		dialog.modeSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range []string{"Asynchronous", "Isochronous"} {
			dialog.modeSelector.AppendText(s)
		}
		dialog.modeSelector.SetActive(0)
		obj = &dialog.modeSelector.Widget
	case iDirection:
		label = "Direction:"
		dialog.directionSelector, err = gtk.ComboBoxTextNew()
		if err != nil {
			return
		}
		for _, s := range []string{"InPort", "OutPort"} {
			dialog.directionSelector.AppendText(s)
		}
		dialog.directionSelector.SetActive(0)
		obj = &dialog.directionSelector.Widget
	default:
		err = fmt.Errorf("createInputElement invalid input element")
		return
	}
	row, err := createLabeledRow(label, obj)
	if err != nil {
		return
	}
	widget = &row.Widget
	return
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
	eNamedPortType,
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
	d.SetDefaultResponse(gtk.RESPONSE_CANCEL)
	d.ShowAll()
	if len(selectorChoices) > 0 {
		dialog.stack.SetVisibleChildName(string(selectorChoices[0]))
	}
	return
}

func (dialog *NewElementDialog) readOut(i inputElement) string {
	switch i {
	case iNodeName:
		text, _ := dialog.nodeNameEntry.GetText()
		return text
	case iTypeName:
		text, _ := dialog.typeNameEntry.GetText()
		return text
	case iPortName:
		text, _ := dialog.portNameEntry.GetText()
		return text
	case iImplName:
		text, _ := dialog.implNameEntry.GetText()
		return text
	case iSignalTypeName:
		text, _ := dialog.signalTypeNameEntry.GetText()
		return text
	case iNodeTypeSelect:
		return dialog.nodeTypeSelector.GetActiveText()
	case iSignalTypeSelect:
		return dialog.signalTypeSelector.GetActiveText()
	case iImplementationType:
		return dialog.implementationSelector.GetActiveText()
	case iPortSelect:
		return dialog.portSelector.GetActiveText()
	case iCType:
		text, _ := dialog.cTypeEntry.GetText()
		return text
	case iChannelId:
		text, _ := dialog.channelIdEntry.GetText()
		return text
	case iScope:
		return dialog.scopeSelector.GetActiveText()
	case iSignalMode:
		return dialog.modeSelector.GetActiveText()
	case iDirection:
		return dialog.directionSelector.GetActiveText()
	}
	return ""
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

// TODO: should go to separate module
type NewElementJob struct {
	parentId, newId string
	elemType        elementType
	input           map[inputElement]string
}

func NewElementJobNew(context string, t elementType) *NewElementJob {
	return &NewElementJob{context, "", t, make(map[inputElement]string)}
}

func (j *NewElementJob) String() string {
	ret := fmt.Sprintf("%s (context=%s, newPath=%s)", j.elemType, j.parentId, j.newId)
	for _, i := range inputElementMap[j.elemType] {
		ret = fmt.Sprintf("%s, %s=%s", ret, i, j.input[i])
	}
	return ret
}

func getParentId(id string) string {
	split := strings.Split(id, ":")
	return strings.Join(split[:len(split)-1], ":")
}

func (j *NewElementJob) CreateObject(fts *models.FilesTreeStore) interface{} {
	parentObject, err := fts.GetObjectById(j.parentId)
	if err != nil {
		log.Fatal("NewElementJob.CreateObject error: referenced parentObject run away...")
	}
	switch j.elemType {
	case eNode:
		var context freesp.SignalGraphType
		switch parentObject.(type) {
		case freesp.Node:
			context = parentObject.(freesp.Node).Context()
			j.parentId = getParentId(j.parentId)
		case freesp.SignalGraph:
			context = parentObject.(freesp.SignalGraph).ItsType()
		case freesp.SignalGraphType:
			context = parentObject.(freesp.SignalGraphType)
		default:
			log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject wrong type...")
		}
		ntype, ok := freesp.GetNodeTypeByName(j.input[iNodeTypeSelect])
		if !ok {
			log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject type wrong...")
		}
		return freesp.NodeNew(j.input[iNodeName], ntype, context)

	case eNodeType:
		var context string
		switch parentObject.(type) {
		case freesp.NodeType:
			context = parentObject.(freesp.NodeType).DefinedAt()
			j.parentId = getParentId(j.parentId)
		case freesp.Library:
			context = parentObject.(freesp.Library).Filename()
		default:
			log.Fatal("NewElementJob.CreateObject(eNodeType) error: referenced parentObject wrong type...")
		}
		return freesp.NodeTypeNew(j.input[iTypeName], context)

	case eConnection:
		switch parentObject.(type) {
		case freesp.Port:
		default:
			log.Fatal("NewElementJob.CreateObject(eConnection) error: referenced parentObject wrong type...")
		}
		ports := getMatchingPorts(fts)
		for _, p := range ports {
			s := fmt.Sprintf("%s/%s", p.Node().NodeName(), p.PortName())
			if j.input[iPortSelect] == s {
				return p
			}
		}

	case eNamedPortType:
		switch parentObject.(type) {
		case freesp.NamedPortType:
			j.parentId = getParentId(j.parentId)
		case freesp.NodeType:
		default:
			log.Fatal("NewElementJob.CreateObject(eNamedPortType) error: referenced parentObject wrong type...")
		}
		_, ok := freesp.GetSignalTypeByName(j.input[iSignalTypeSelect])
		if !ok {
			log.Fatal("NewElementJob.CreateObject(eNamedPortType) error: referenced signal type wrong...")
		}
		var dir freesp.PortDirection
		if j.input[iDirection] == "InPort" {
			dir = freesp.InPort
		} else {
			dir = freesp.OutPort
		}
		return freesp.NamedPortTypeNew(j.input[iPortName], j.input[iSignalTypeSelect], dir)

	case eSignalType:
		name := j.input[iSignalTypeName]
		cType := j.input[iCType]
		channelId := j.input[iChannelId]
		var scope freesp.Scope
		if j.input[iScope] == "Local" {
			scope = freesp.Local
		} else {
			scope = freesp.Global
		}
		var mode freesp.Mode
		if j.input[iSignalMode] == "Asynchronous" {
			mode = freesp.Asynchronous
		} else {
			mode = freesp.Synchronous
		}
		return freesp.SignalTypeNew(name, cType, channelId, scope, mode)

	case eImplementation:
		var implType freesp.ImplementationType
		if j.input[iImplementationType] == "Elementary Type" {
			implType = freesp.NodeTypeElement
		} else {
			implType = freesp.NodeTypeGraph
		}
		return freesp.ImplementationNew(j.input[iImplName], implType)

	default:
		log.Fatal("NewElementJob.CreateObject error: invalid elemType ", j.elemType)
	}
	return nil
}
