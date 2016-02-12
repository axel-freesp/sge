package main

import (
	"fmt"
	"github.com/axel-freesp/sge/freesp"
	"github.com/axel-freesp/sge/freesp/behaviour"
	"github.com/axel-freesp/sge/freesp/platform"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	gr "github.com/axel-freesp/sge/interface/graph"
	pf "github.com/axel-freesp/sge/interface/platform"
	tr "github.com/axel-freesp/sge/interface/tree"
	"github.com/axel-freesp/sge/models"
	//"image"
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

func (j *NewElementJob) CreateObject(fts *models.FilesTreeStore) (ret tr.TreeElementIf, err error) {
	var parentObject tr.TreeElementIf
	parentObject, err = fts.GetObjectById(j.parentId)
	if err != nil {
		log.Fatal("NewElementJob.CreateObject error: referenced parentObject run away...")
	}
	switch j.elemType {
	case eNode, eInputNode, eOutputNode:
		var context bh.SignalGraphTypeIf
		switch parentObject.(type) {
		case bh.NodeIf:
			context = parentObject.(bh.NodeIf).Context()
			j.parentId = getParentId(j.parentId)
		case bh.SignalGraphIf:
			context = parentObject.(bh.SignalGraphIf).ItsType()
		case bh.SignalGraphTypeIf:
			context = parentObject.(bh.SignalGraphTypeIf)
		case bh.ImplementationIf:
			if parentObject.(bh.ImplementationIf).ImplementationType() == bh.NodeTypeGraph {
				context = parentObject.(bh.ImplementationIf).Graph()
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
			ret, err = behaviour.NodeNew(j.input[iNodeName], ntype, context)
		} else if j.elemType == eInputNode {
			ret, err = behaviour.InputNodeNew(j.input[iInputNodeName], j.input[iInputTypeSelect], context)
		} else {
			ret, err = behaviour.OutputNodeNew(j.input[iOutputNodeName], j.input[iOutputTypeSelect], context)
		}
		if len(j.extra) > 0 {
			coords := strings.Split(j.extra, "|")
			var x, y int
			fmt.Sscanf(coords[0], "%d", &x)
			fmt.Sscanf(coords[1], "%d", &y)
			//pos := image.Point{x, y}
			//log.Printf("NewElementJob.CreateObject(eNode) setting position %s: %v\n", j.extra, pos)
			//ret.(bh.NodeIf).SetPosition(pos)
			// TODO: SetModePosition...
		}

	case eNodeType:
		var context string
		switch parentObject.(type) {
		case bh.NodeTypeIf:
			context = parentObject.(bh.NodeTypeIf).DefinedAt()
			j.parentId = getParentId(j.parentId)
		case bh.SignalTypeIf:
			j.parentId = getParentId(j.parentId)
			parentObject, err = fts.GetObjectById(j.parentId)
			context = parentObject.(bh.LibraryIf).Filename()
		case bh.LibraryIf:
			context = parentObject.(bh.LibraryIf).Filename()
		default:
			log.Fatal("NewElementJob.CreateObject(eNodeType) error: referenced parentObject wrong type...")
		}
		ret = behaviour.NodeTypeNew(j.input[iTypeName], context)

	case eConnection:
		switch parentObject.(type) {
		case bh.PortIf:
		case bh.SignalGraphTypeIf:
			fromTo := strings.Split(j.extra, "/")
			var n bh.NodeIf
			for _, n = range parentObject.(bh.SignalGraphTypeIf).Nodes() {
				if n.Name() == fromTo[0] {
					for _, parentObject = range n.OutPorts() {
						if parentObject.(bh.PortIf).Name() == fromTo[1] {
							break
						}
					}
					break
				}
			}
			if parentObject == nil {
				log.Fatalf("NewElementJob.CreateObject(eNodeType) error: no valid FROM port for edge job %v\n", j)
			}
			_ = parentObject.(bh.PortIf)
		case bh.ImplementationIf:
			fromTo := strings.Split(j.extra, "/")
			var n bh.NodeIf
			for _, n = range parentObject.(bh.ImplementationIf).Graph().Nodes() {
				if n.Name() == fromTo[0] {
					for _, parentObject = range n.OutPorts() {
						if parentObject.(bh.PortIf).Name() == fromTo[1] {
							break
						}
					}
					break
				}
			}
			if parentObject == nil {
				log.Fatalf("NewElementJob.CreateObject(eNodeType) error: no valid FROM port for edge job %v\n", j)
			}
			_ = parentObject.(bh.PortIf)
		default:
			log.Fatalf("NewElementJob.CreateObject(eConnection) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ports := getMatchingPorts(fts, parentObject)
		for _, p := range ports {
			s := fmt.Sprintf("%s/%s", p.Node().Name(), p.Name())
			if j.input[iPortSelect] == s {
				var from, to bh.PortIf
				if p.Direction() == gr.InPort {
					from = parentObject.(bh.PortIf)
					to = p
				} else {
					from = p
					to = parentObject.(bh.PortIf)
				}
				ret = behaviour.ConnectionNew(from, to)
				break
			}
		}

	case ePortType:
		switch parentObject.(type) {
		case bh.PortTypeIf:
			j.parentId = getParentId(j.parentId)
		case bh.NodeTypeIf:
		default:
			log.Fatal("NewElementJob.CreateObject(ePortType) error: referenced parentObject wrong type...")
		}
		_, ok := freesp.GetSignalTypeByName(j.input[iSignalTypeSelect])
		if !ok {
			err = fmt.Errorf("NewElementJob.CreateObject(ePortType) error: referenced signal type wrong...")
			return
		}
		ret = behaviour.PortTypeNew(j.input[iPortName], j.input[iSignalTypeSelect], string2direction[j.input[iDirection]])

	case eSignalType:
		switch parentObject.(type) {
		case bh.SignalTypeIf:
			j.parentId = getParentId(j.parentId)
		case bh.NodeTypeIf:
			j.parentId = getParentId(j.parentId)
		case bh.LibraryIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		parentObject, err = fts.GetObjectById(j.parentId)
		if err != nil {
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: parent library object not found\n")
		}
		name := j.input[iSignalTypeName]
		cType := j.input[iCType]
		channelId := j.input[iChannelId]
		scope := string2scope[j.input[iScope]]
		mode := string2mode[j.input[iSignalMode]]
		ret, err = behaviour.SignalTypeNew(name, cType, channelId, scope, mode, parentObject.(bh.LibraryIf).Filename())
		if err != nil {
			log.Printf("NewElementJob.CreateObject(eSignalType) error: SignalTypeNew failed: %s\n", err)
			return
		}

	case eImplementation:
		switch parentObject.(type) {
		case bh.ImplementationIf:
			j.parentId = getParentId(j.parentId)
		case bh.NodeTypeIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eSignalType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		implType := string2implType[j.input[iImplementationType]]
		ret = behaviour.ImplementationNew(j.input[iImplName], implType, &global)

	case eArch:
		switch parentObject.(type) {
		case pf.ArchIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(pf.ArchIf).Platform()
		case pf.PlatformIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eArch) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret = platform.ArchNew(j.input[iArchName], parentObject.(pf.PlatformIf))

	case eIOType:
		var p pf.PlatformIf
		switch parentObject.(type) {
		case pf.IOTypeIf:
			j.parentId = getParentId(j.parentId)
			p = parentObject.(pf.IOTypeIf).Platform()
		case pf.ArchIf:
			p = parentObject.(pf.ArchIf).Platform()
		default:
			log.Fatalf("NewElementJob.CreateObject(eIOType) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret, err = platform.IOTypeNew(j.input[iIOTypeName], gr.IOMode(j.input[iIOModeSelect]), p)

	case eProcess:
		switch parentObject.(type) {
		case pf.ProcessIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(pf.ProcessIf).Arch()
		case pf.ArchIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eProcess) error: referenced parentObject wrong type %T\n", parentObject)
		}
		ret = platform.ProcessNew(j.input[iProcessName], parentObject.(pf.ArchIf))

	case eChannel:
		switch parentObject.(type) {
		case pf.ChannelIf:
			j.parentId = getParentId(j.parentId)
			parentObject = parentObject.(pf.ChannelIf).Process()
		case pf.ProcessIf:
		default:
			log.Fatalf("NewElementJob.CreateObject(eChannel) error: referenced parentObject wrong type %T\n", parentObject)
		}
		processes := getOtherProcesses(fts, parentObject)
		var p pf.ProcessIf
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
		ret = platform.ChannelNew(string2direction[j.input[iChannelDirection]], ioType, parentObject.(pf.ProcessIf), j.input[iChannelLinkSelect])

	default:
		log.Fatal("NewElementJob.CreateObject error: invalid elemType ", j.elemType)
	}
	return
}
