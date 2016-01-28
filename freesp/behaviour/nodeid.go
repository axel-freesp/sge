package behaviour

import (
	"fmt"
	bh "github.com/axel-freesp/sge/interface/behaviour"
	"strings"
)

type nodeId struct {
	path []string
}

var EmptyNodeId bh.NodeIdIf = (*nodeId)(nil)

func NodeIdNew(parentId bh.NodeIdIf, id string) *nodeId {
	if parentId == EmptyNodeId {
		return &nodeId{[]string{id}}
	}
	return NodeIdFromString(fmt.Sprintf("%s/%s", parentId, id))
}

func NodeIdFromString(idString string) *nodeId {
	return &nodeId{strings.Split(idString, "/")}
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
	return (&nodeId{n.path[1:]}).IsAncestor(&nodeId{nid.path[1:]})
}
