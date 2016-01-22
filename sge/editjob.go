package main

import (
	"fmt"
	"log"
	//"strings"
	"github.com/axel-freesp/sge/freesp"
	interfaces "github.com/axel-freesp/sge/interface"
	"github.com/axel-freesp/sge/models"
)

type EditJob struct {
	objId       string
	elemType    elementType
	detail, old map[inputElement]string
}

type EditJobDirection int

const (
	EditJobForward EditJobDirection = iota
	EditJobRevert
)

func EditJobNew(context string, t elementType) *EditJob {
	return &EditJob{context, t, make(map[inputElement]string), make(map[inputElement]string)}
}

func (j *EditJob) String() string {
	ret := fmt.Sprintf("%s (objId=%s)\n\told: ", j.elemType, j.objId)
	for _, i := range inputElementMap[j.elemType] {
		ret = fmt.Sprintf("%s, %s=%s", ret, i, j.old[i])
	}
	ret = fmt.Sprintf("%s\n\tnew: ", ret)
	for _, i := range inputElementMap[j.elemType] {
		ret = fmt.Sprintf("%s, %s=%s", ret, i, j.detail[i])
	}
	return ret
}

func (j *EditJob) EditObject(fts *models.FilesTreeStore, direction EditJobDirection) (state string, err error) {
	var detail, old *map[inputElement]string
	if direction == EditJobForward {
		detail, old = &j.detail, &j.old
	} else {
		old, detail = &j.detail, &j.old
	}
	var obj freesp.TreeElement
	obj, err = fts.GetObjectById(j.objId)
	state = j.objId
	switch j.elemType {
	case eNode:
		n := obj.(freesp.Node)
		(*old)[iNodeName] = n.Name()
		n.SetName((*detail)[iNodeName])
		fts.SetValueById(j.objId, n.Name())
		for _, p := range n.InPorts() {
			updateConnections(p, fts)
		}
		for _, p := range n.OutPorts() {
			updateConnections(p, fts)
		}
	case eOutputNode:
		n := obj.(freesp.Node)
		(*old)[iOutputNodeName] = n.Name()
		n.SetName((*detail)[iOutputNodeName])
		fts.SetValueById(j.objId, n.Name())
		for _, p := range n.InPorts() {
			updateConnections(p, fts)
		}
	case eInputNode:
		n := obj.(freesp.Node)
		(*old)[iInputNodeName] = n.Name()
		n.SetName((*detail)[iInputNodeName])
		fts.SetValueById(j.objId, n.Name())
		for _, p := range n.OutPorts() {
			updateConnections(p, fts)
		}
	case eNodeType:
		nt := obj.(freesp.NodeType)
		if len(nt.Instances()) > 0 {
			log.Printf("jobApplier.Apply(JobEdit): WARNING: NodeType %s has instances.\n", nt.TypeName())
			log.Printf("jobApplier.Apply(JobEdit): Editing is not implemented in this case.\n")
			return
		}
		(*old)[iTypeName] = nt.TypeName()
		nt.SetTypeName((*detail)[iTypeName])
		fts.SetValueById(j.objId, nt.TypeName())
	case ePortType:
		pt := obj.(freesp.PortType)
		ptCursor := fts.Cursor(pt)
		ntCursor := fts.Parent(ptCursor)
		nt := fts.Object(ntCursor).(freesp.NodeType)
		if len(nt.Instances()) > 0 {
			log.Printf("jobApplier.Apply(JobEdit): WARNING: NodeType %s has instances.\n", nt.TypeName())
			log.Printf("jobApplier.Apply(JobEdit): Editing is not implemented in this case.\n")
			return
		}
		(*old)[iPortName] = pt.Name()
		(*old)[iSignalTypeSelect] = pt.SignalType().TypeName()
		(*old)[iDirection] = direction2string[pt.Direction()]
		fts.DeleteObject(ptCursor.Path)
		fts.AddNewObject(ntCursor.Path, ntCursor.Position,
			freesp.PortTypeNew((*detail)[iPortName],
				(*detail)[iSignalTypeSelect],
				string2direction[(*detail)[iDirection]]))
		state = ptCursor.Path
	case eSignalType:
		st := obj.(freesp.SignalType)
		if (*detail)[iSignalTypeName] != st.TypeName() {
			log.Printf("jobApplier.Apply(JobEdit): Renaming SignalType is not implemented.\n")
		}
		(*old)[iCType] = st.CType()
		st.SetCType((*detail)[iCType])
		(*old)[iChannelId] = st.ChannelId()
		st.SetChannelId((*detail)[iChannelId])
		(*old)[iScope] = scope2string[st.Scope()]
		st.SetScope(string2scope[(*detail)[iScope]])
		(*old)[iSignalMode] = mode2string[st.Mode()]
		st.SetMode(string2mode[(*detail)[iSignalMode]])
	case eImplementation:
		impl := obj.(freesp.Implementation)
		(*old)[iImplName] = impl.ElementName()
		impl.SetElemName((*detail)[iImplName])
		fts.SetValueById(j.objId, (*detail)[iImplName])
	case eArch:
		a := obj.(freesp.Arch)
		(*old)[iArchName] = a.Name()
		a.SetName((*detail)[iArchName])
		for _, p := range a.Processes() {
			for _, c := range p.InChannels() {
				link := c.Link()
				id := fts.Cursor(link)
				fts.SetValueById(id.Path, link.Name())
			}
			for _, c := range p.OutChannels() {
				link := c.Link()
				id := fts.Cursor(link)
				fts.SetValueById(id.Path, link.Name())
			}
		}
		fts.SetValueById(j.objId, a.Name())
	case eProcess:
		p := obj.(freesp.Process)
		(*old)[iProcessName] = p.Name()
		p.SetName((*detail)[iProcessName])
		for _, c := range p.InChannels() {
			link := c.Link()
			id := fts.Cursor(link)
			fts.SetValueById(id.Path, link.Name())
		}
		for _, c := range p.OutChannels() {
			link := c.Link()
			id := fts.Cursor(link)
			fts.SetValueById(id.Path, link.Name())
		}
		fts.SetValueById(j.objId, p.Name())
	case eIOType:
		t := obj.(freesp.IOType)
		(*old)[iIOTypeName] = t.Name()
		t.SetName((*detail)[iIOTypeName])
		for _, a := range t.Platform().Arch() {
			aCursor := fts.Cursor(a)
			for _, p := range a.Processes() {
				pCursor := fts.CursorAt(aCursor, p)
				for _, c := range p.InChannels() {
					if t == c.IOType() {
						cCursor := fts.CursorAt(pCursor, c)
						fts.SetValueById(cCursor.Path, c.Name())
					}
				}
				for _, c := range p.OutChannels() {
					if t == c.IOType() {
						cCursor := fts.CursorAt(pCursor, c)
						fts.SetValueById(cCursor.Path, c.Name())
					}
				}
			}
		}
		(*old)[iIOModeSelect] = string(t.IOMode())
		t.SetIOMode(interfaces.IOMode((*detail)[iIOModeSelect]))
		fts.SetValueById(j.objId, t.Name())
	case eChannel:
		c := obj.(freesp.Channel)
		//(*old)[iChannelDirection] = direction2string[c.Direction()]
		//c.SetDirection(string2direction[(*detail)[iChannelDirection]])
		(*old)[iIOTypeSelect] = c.IOType().Name()
		iot, ok := freesp.GetIOTypeByName((*detail)[iIOTypeSelect])
		if ok {
			c.SetIOType(iot)
		} else {
			log.Printf("jobApplier.Apply(JobEdit): ERROR: IOType %s not registered.\n", (*detail)[iIOTypeSelect])
		}
		fts.SetValueById(j.objId, c.Name())
		link := c.Link()
		id := fts.Cursor(link)
		fts.SetValueById(id.Path, link.Name())
	case eMapElement:
		c := obj.(freesp.MappedElement)
		pr := c.Process()
		if pr == nil {
			(*old)[iProcessSelect] = "<unmapped>"
		} else {
			(*old)[iProcessSelect] = fmt.Sprintf("%s/%s", pr.Arch().Name(), pr.Name())
		}
		for _, a := range c.Mapping().Platform().Arch() {
			for _, pr = range a.Processes() {
				if (*detail)[iProcessSelect] == fmt.Sprintf("%s/%s", a.Name(), pr.Name()) {
					c.SetProcess(pr)
					fts.SetValueById(j.objId, (*detail)[iProcessSelect])
					return
				}
			}
		}
		c.SetProcess(nil)
		fts.SetValueById(j.objId, (*detail)[iProcessSelect])
	default:
		log.Printf("jobApplier.Apply(JobEdit): error: invalid job description\n")
	}
	return
}

func updateConnections(p freesp.Port, fts *models.FilesTreeStore) {
	nodeCursor := fts.Cursor(p.Node())
	portCursor := fts.CursorAt(nodeCursor, p)
	for _, c := range p.Connections() {
		conn := p.Connection(c)
		connCursor := fts.CursorAt(portCursor, conn)
		otherNode := c.Node()
		otherNodeCursor := fts.Cursor(otherNode)
		otherPortCursor := fts.CursorAt(otherNodeCursor, c)
		otherConnCursor := fts.CursorAt(otherPortCursor, conn)
		connText := fmt.Sprintf("%s/%s -> %s/%s",
			conn.From().Node().Name(), conn.From().Name(),
			conn.To().Node().Name(), conn.To().Name())
		fts.SetValueById(connCursor.Path, connText)
		fts.SetValueById(otherConnCursor.Path, connText)
	}
}
