package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
	"image"
	"log"
	"strings"
)

type NewElementJob struct {
	parentId, newId string
	elemType        elementType
	input           map[inputElement]string
	extra           string // used for pasting nodes and connections
}

func NewElementJobNew(context string, t elementType) *NewElementJob {
	return &NewElementJob{context, "", t, make(map[inputElement]string), ""}
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
	case eNode, eInputNode, eOutputNode:
		var context freesp.SignalGraphTypeIf
		switch parentObject.(type) {
		case freesp.NodeIf:
			context = parentObject.(freesp.NodeIf).Context()
			j.parentId = getParentId(j.parentId)
		case freesp.SignalGraphIf:
			context = parentObject.(freesp.SignalGraphIf).ItsType()
		case freesp.SignalGraphTypeIf:
			context = parentObject.(freesp.SignalGraphTypeIf)
		case freesp.ImplementationIf:
			if parentObject.(freesp.ImplementationIf).ImplementationType() == freesp.NodeTypeGraph {
				context = parentObject.(freesp.ImplementationIf).Graph()
			} else {
				log.Fatal("NewElementJob.CreateObject(eNode) error: parent implementation is no graph...")
			}
		default:
			log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject wrong type...")
		}
		if j.elemType == eNode {
			ntype, ok := freesp.GetNodeTypeByName(j.input[iNodeTypeSelect])
			if !ok {
				log.Fatal("NewElementJob.CreateObject(eNode) error: referenced parentObject type wrong...")
			}
			ret, err = freesp.NodeNew(j.input[iNodeName], ntype, context)
		} else if j.elemType == eInputNode {
			ret, err = freesp.InputNodeNew(j.input[iInputNodeName], j.input[iInputTypeSelect], context)
		} else {
			ret, err = freesp.OutputNodeNew(j.input[iOutputNodeName], j.input[iOutputTypeSelect], context)
		}
		if len(j.extra) > 0 {
			coords := strings.Split(j.extra, "|")
			var x, y int
			fmt.Sscanf(coords[0], "%d", &x)
			fmt.Sscanf(coords[1], "%d", &y)
			pos := image.Point{x, y}
			//log.Printf("NewElementJob.CreateObject(eNode) setting position %s: %v\n", j.extra, pos)
			ret.(freesp.NodeIf).SetPosition(pos)
		}

	case eNodeType:
		var context string
		switch parentObject.(type) {
		case freesp.NodeTypeIf:
			context = parentObject.(freesp.NodeTypeIf).DefinedAt()
			j.parentId = getParentId(j.parentId)
		case freesp.SignalType:
			j.parentId = getParentId(j.parentId)
			parentObject, err = fts.GetObjectById(j.parentId)
			context = parentObject.(freesp.LibraryIf).Filename()
		case freesp.LibraryIf:
			context = parentObject.(freesp.LibraryIf).Filename()
		default:
			log.Fatal("NewElementJob.CreateObject(eNodeType) error: referenced parentObject wrong type...")
		}
		ret = freesp.NodeTypeNew(j.input[iTypeName], context)

	case eConnection:
		switch parentObject.(type) {
		case freesp.Port:
		case freesp.SignalGraphTypeIf:
			fromTo := strings.Split(j.extra, "/")
			var n freesp.NodeIf
			for _, n = range parentObject.(freesp.SignalGraphTypeIf).Nodes() {
				if n.Name() == fromTo[0] {
					for _, parentObject = range n.OutPorts() {
						if parentObject.(freesp.Port).Name() == fromTo[1] {
							break
						}
					}
					break
				}
			}
			if parentObject == nil {
				log.Fatalf("NewElementJob.CreateObject(eNodeType) error: no valid FROM port for edge job %v\n", j)
			}
			_ = parentObject.(freesp.Port)
		case freesp.ImplementationIf:
			fromTo := strings.Split(j.extra, "/")
			var n freesp.NodeIf
			for _, n = range parentObject.(freesp.ImplementationIf).Graph().Nodes() {
				if n.Name() == fromTo[0] {
					for _, parentObject = range n.OutPorts() {
						if parentObject.(freesp.Port).Name() == fromTo[1] {
							break
						}
					}
					break
				}
			}
			if parentObject == nil {
				log.Fatalf("NewElementJob.CreateObject(eNodeType) error: no valid FROM port for edge job %v\n", j)
			}
			_ = parentObject.(freesp.Port)
		default:
			log.Fatalf("NewElementJob.CreateObject(eConnection) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ports := getMatchingPorts(fts, parentObject)
		for _, p := range ports {
			s := fmt.Sprintf("%s/%s", p.Node().Name(), p.Name())
			if j.input[iPortSelect] == s {
				var from, to freesp.Port
				if p.Direction() == interfaces.InPort {
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
		case freesp.NodeTypeIf:
		default:
			log.Fatal("NewElementJob.CreateObject(ePortType) error: referenced parentObject wrong type...")
		}
		_, ok := freesp.GetSignalTypeByName(j.input[iSignalTypeSelect])
		if !ok {
			err = fmt.Errorf("NewElementJob.CreateObject(ePortType) error: referenced signal type wrong...")
			return
		}
		ret = freesp.PortTypeNew(j.input[iPortName], j.input[iSignalTypeSelect], string2direction[j.input[iDirection]])

	case eSignalType:
		switch parentObject.(type) {
		case freesp.SignalType:
			j.parentId = getParentId(j.parentId)
		case freesp.NodeTypeIf:
			j.parentId = getParentId(j.parentId)
		case freesp.LibraryIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		name := j.input[iSignalTypeName]
		cType := j.input[iCType]
		channelId := j.input[iChannelId]
		scope := string2scope[j.input[iScope]]
		mode := string2mode[j.input[iSignalMode]]
		ret, err = freesp.SignalTypeNew(name, cType, channelId, scope, mode)
		if err != nil {
			log.Printf("NewElementJob.CreateObject(eSignalType) error: SignalTypeNew failed: %s\n", err)
			return
		}

	case eImplementation:
		switch parentObject.(type) {
		case freesp.ImplementationIf:
			j.parentId = getParentId(j.parentId)
		case freesp.NodeTypeIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		implType := string2implType[j.input[iImplementationType]]
		ret = freesp.ImplementationNew(j.input[iImplName], implType, &global)

	case eArch:
		switch parentObject.(type) {
		case freesp.ArchIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(freesp.ArchIf).Platform()
		case freesp.PlatformIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eArch) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret = freesp.ArchNew(j.input[iArchName], parentObject.(freesp.PlatformIf))

	case eIOType:
		var p freesp.PlatformIf
		switch parentObject.(type) {
		case freesp.IOTypeIf:
			j.parentId = getParentId(j.parentId)
			p = parentObject.(freesp.IOTypeIf).Platform()
		case freesp.ArchIf:
			p = parentObject.(freesp.ArchIf).Platform()
		default:
			log.Fatalf("NewElementJob.CreateObject(eIOType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret, err = freesp.IOTypeNew(j.input[iIOTypeName], interfaces.IOMode(j.input[iIOModeSelect]), p)

	case eProcess:
		switch parentObject.(type) {
		case freesp.ProcessIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(freesp.ProcessIf).Arch()
		case freesp.ArchIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eProcess) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret = freesp.ProcessNew(j.input[iProcessName], parentObject.(freesp.ArchIf))

	case eChannel:
		switch parentObject.(type) {
		case freesp.ChannelIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(freesp.ChannelIf).Process()
		case freesp.ProcessIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eChannel) error: referenced parentObject wrong type %T\n", parentObject)
		}
		processes := getOtherProcesses(fts, parentObject)
		var p freesp.ProcessIf
		for _, p = range processes {
			s := fmt.Sprintf("%s/%s", p.Arch().Name(), p.Name())
			if s == j.input[iChannelLinkSelect] {
				break
			}
		}
		if p == nil {
			log.Fatalf("NewElementJob.CreateObject(eChannel) error: can't find chosen process\n", j.input[iChannelLinkSelect])
		}
		ioType, ok := freesp.GetIOTypeByName(j.input[iIOTypeSelect])
		if !ok {
			log.Fatalf("NewElementJob.CreateObject(eChannel) error: can't find chosen ioType\n", j.input[iIOTypeSelect])
		}
		ret = freesp.ChannelNew(string2direction[j.input[iChannelDirection]], ioType, parentObject.(freesp.ProcessIf), j.input[iChannelLinkSelect])

	default:
		log.Fatal("NewElementJob.CreateObject error: invalid elemType ", j.elemType)
	}
	return
}
