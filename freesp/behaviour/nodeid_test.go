package behaviour

import (
	//"fmt"
	"testing"
)

func TestNodeId(t *testing.T) {
	case1 := []struct {
		id, parent string
	}{
		{"node1", ""},
		{"node1/child1", "node1"},
		{"some/node/with/deep/path", "some/node/with/deep"},
	}
	for i, c := range case1 {
		id := NodeIdFromString(c.id)
		p := id.Parent()
		if p.String() != c.parent {
			t.Errorf("Testcase %d failed: %s is parent of %s\n", i, p, id)
		}
	}
}
