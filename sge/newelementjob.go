package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/models"
	"log"
	"strings"
)

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

func (j *NewElementJob) CreateObject(fts *models.FilesTreeStore) (ret freesp.TreeElement, err error) {
	var parentObject freesp.TreeElement
	parentObject, err = fts.GetObjectById(j.parentId)
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
		case freesp.Implementation:
			if parentObject.(freesp.Implementation).ImplementationType() == freesp.NodeTypeGraph {
				context = parentObject.(freesp.Implementation).Graph()
			} else {
				log.Fatal("NewElementJob.CreateObject(eNode) error: parent implementation is no graph...")
			}
		default:
			log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject wrong type...")
		}
		ntype, ok := freesp.GetNodeTypeByName(j.input[iNodeTypeSelect])
		if !ok {
			log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject type wrong...")
		}
		ret, err = freesp.NodeNew(j.input[iNodeName], ntype, context)

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
		ret = freesp.NodeTypeNew(j.input[iTypeName], context)

	case eConnection:
		switch parentObject.(type) {
		case freesp.Port:
		default:
			log.Fatal("NewElementJob.CreateObject(eConnection) error: referenced parentObject wrong type...")
		}
		ports := getMatchingPorts(fts, parentObject)
		for _, p := range ports {
			s := fmt.Sprintf("%s/%s", p.Node().Name(), p.Name())
			if j.input[iPortSelect] == s {
				var from, to freesp.Port
				if p.Direction() == freesp.InPort {
					from = parentObject.(freesp.Port)
					to = p
				} else {
					from = p
					to = parentObject.(freesp.Port)
				}
				ret = freesp.ConnectionNew(from, to)
				break
			}
		}

	case ePortType:
		switch parentObject.(type) {
		case freesp.PortType:
			j.parentId = getParentId(j.parentId)
		case freesp.NodeType:
		default:
			log.Fatal("NewElementJob.CreateObject(ePortType) error: referenced parentObject wrong type...")
		}
		_, ok := freesp.GetSignalTypeByName(j.input[iSignalTypeSelect])
		if !ok {
			err = fmt.Errorf("NewElementJob.CreateObject(ePortType) error: referenced signal type wrong...")
			return
		}
		var dir freesp.PortDirection
		if j.input[iDirection] == "InPort" {
			dir = freesp.InPort
		} else {
			dir = freesp.OutPort
		}
		ret = freesp.PortTypeNew(j.input[iPortName], j.input[iSignalTypeSelect], dir)

	case eSignalType:
		switch parentObject.(type) {
		case freesp.SignalType:
			j.parentId = getParentId(j.parentId)
		case freesp.Library:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
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
		ret, err = freesp.SignalTypeNew(name, cType, channelId, scope, mode)
		if err != nil {
			log.Printf("NewElementJob.CreateObject(eSignalType) error: SignalTypeNew failed: %s\n", err)
			return
		}

	case eImplementation:
		switch parentObject.(type) {
		case freesp.Implementation:
			j.parentId = getParentId(j.parentId)
		case freesp.NodeType:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		var implType freesp.ImplementationType
		if j.input[iImplementationType] == "Elementary Type" {
			implType = freesp.NodeTypeElement
		} else {
			implType = freesp.NodeTypeGraph
		}
		ret = freesp.ImplementationNew(j.input[iImplName], implType, &global)

	default:
		log.Fatal("NewElementJob.CreateObject error: invalid elemType ", j.elemType)
	}
	return
}
