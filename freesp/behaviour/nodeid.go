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
		return &nodeId{[]string{id}, ""}
	}
	return NodeIdFromString(fmt.Sprintf("%s/%s", parentId, id))
}

func NodeIdFromString(idString string) *nodeId {
	return &nodeId{strings.Split(idString, "/"), ""}
}

func (n nodeId) String() string {
	return strings.Join(n.path, "/")
}

func (n nodeId) Parent() bh.NodeIdIf {
	return NodeIdFromString(strings.Join(n.path[:len(n.path)-1], "/"))
}

func (n nodeId) IsAncestor(id bh.NodeIdIf) bool {
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
