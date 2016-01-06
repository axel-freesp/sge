package freesp

import (
	"testing"
)

func TestGraph(t *testing.T) {
	case1 := []struct {
		library, graph     string
		nodes, connections int
	}{
		{`<library xmlns="http://www.freesp.de/xml/freeSP" version="1.0">
   <signal-type name="s1" scope="" mode="" c-type="" message-id=""></signal-type>
   <node-type name="Test">
      <intype port="" type="s1"></intype>
      <outtype port="" type="s1"></outtype>
   </node-type>
</library>
`, `<?xml version="1.0" encoding="UTF-8"?>
<signal-graph xmlns="http://www.freesp.de/xml/freeSP" version="1.0">
    <nodes>
        <input name="sensor">
            <outtype type="s1"/>
        </input>
        <output name="actuator">
            <intype type="s1"/>
        </output>
        <processing-node name="test" type="Test"></processing-node>
    </nodes>
    <connections>
        <connect from="sensor" to="test"/>
        <connect from="test" to="actuator"/>
    </connections>
</signal-graph>
`, 3, 2},
	}

	for i, c := range case1 {
		Init()
		var l Library = LibraryNew("test.alml", nil)
		buf := copyBuf(c.library)
		err := l.Read(buf)
		if err != nil {
			t.Errorf("Testcase %d: Failed to read from buffer: %v", i, err)
			return
		}
		var sg SignalGraph = SignalGraphNew("test.sml", nil)
		buf = copyBuf(c.graph)
		err = sg.Read(buf)
		if err != nil {
			t.Errorf("Testcase %d: Failed to read from buffer: %v", i, err)
			return
		}
		var st SignalGraphType = sg.ItsType()
		if len(st.Nodes()) != c.nodes {
			t.Errorf("Testcase %d: Node count mismatch", i)
			return
		}
	}
}
