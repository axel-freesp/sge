package behaviour

import (
	"fmt"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	"log"
	"strings"
)

type nodeId struct {
	path     []string
	filename string
}

var EmptyNodeId bh.NodeIdIf = (*nodeId)(nil)

func NodeIdNew(parentId bh.NodeIdIf, id string) *nodeId {
	if parentId == EmptyNodeId {
		log.Panic("NodeIdNew FIXME: parent = nil")
	}
	return NodeIdFromString(fmt.Sprintf("%s/%s", parentId, id), parentId.Filename())
}

func NodeIdFromString(idString, filename string) *nodeId {
	return &nodeId{strings.Split(idString, "/"), filename}
}

func (n nodeId) String() string {
	return strings.Join(n.path, "/")
}

func (n nodeId) Parent() bh.NodeIdIf {
	return NodeIdFromString(strings.Join(n.path[:len(n.path)-1], "/"), n.filename)
}

func (n nodeId) IsAncestor(id bh.NodeIdIf) bool {
	if n.filename != id.Filename() {
		return false
	}
	if id == EmptyNodeId {
		return false
	}
	nid := id.(*nodeId)
	if len(n.path) == 0 {
		return true
	}
	if len(nid.path) == 0 {
		return false
	}
	if nid.path[0] != n.path[0] {
		return false
	}
	return (&nodeId{n.path[1:], n.filename}).IsAncestor(&nodeId{nid.path[1:], nid.filename})
}

func (n *nodeId) SetFilename(filename string) {
	n.filename = filename
}

func (n nodeId) Filename() (filename string) {
	return n.filename
}

func (n nodeId) First() string {
	if len(n.path) == 0 {
		return ""
	}
	return n.path[0]
}

func (n *nodeId) TruncFirst() {
	if len(n.path) == 0 {
		log.Panicf("nodeId.TruncFirst FIXME: empty path\n")
	}
	n.path = n.path[1:]
}

/*
 *      NodeIdList
 *
 */

type NodeIdList struct {
	nodeIds []bh.NodeIdIf
}

func NodeIdListInit() NodeIdList {
	return NodeIdList{nil}
}

func (l *NodeIdList) Append(st bh.NodeIdIf) {
	l.nodeIds = append(l.nodeIds, st)
}

func (l *NodeIdList) Remove(st bh.NodeIdIf) {
	var i int
	for i = range l.nodeIds {
		if st == l.nodeIds[i] {
			break
		}
	}
	if i >= len(l.nodeIds) {
		for _, v := range l.nodeIds {
			log.Printf("NodeIdList.RemoveNodeType have bh.NodeId %v\n", v)
		}
		log.Fatalf("NodeIdList.RemoveNodeType error: bh.NodeId %v not in this list\n", st)
	}
	for i++; i < len(l.nodeIds); i++ {
		l.nodeIds[i-1] = l.nodeIds[i]
	}
	l.nodeIds = l.nodeIds[:len(l.nodeIds)-1]
}

func (l *NodeIdList) NodeIds() []bh.NodeIdIf {
	return l.nodeIds
}
