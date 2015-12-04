package freesp

import (
	"testing"
	//"fmt"
)

func TestGraph(t *testing.T) {
	case1 := []struct {
		xml                string
		nodes, connections int
	}{
		{`<?xml version="1.0" encoding="UTF-8"?>
<signal-graph version="1.0">
    <nodes>
        <input name="sensor">
            <outtype c-type="VideoLine_t" mode="sync" message-id="MSG_TSR_EDGEIMAGE" sizeof="48" cycle-time-us="90" type="TSR-EdgeImage"/>
        </input>
        <output name="actuator1">
            <intype scope="local" c-type="uint16" message-id="MSG_TSR_EDGE_LEVEL" type="TSR-Edgelevel"/>
        </output>
        <output name="actuator2">
            <intype scope="local" c-type="Tsr_Contourlist_t" message-id="MSG_TSR_CONTOURLIST" type="TSR-Contourlist"/>
        </output>
        <processing-node name="process1" type="filter">
            <outtype scope="local" c-type="uint16" message-id="MSG_TSR_EDGE_LEVEL" port="edge-level" type="TSR-Edgelevel"/>
            <outtype scope="local" c-type="Tsr_Contourlist_t" message-id="MSG_TSR_CONTOURLIST" port="result" type="TSR-Contourlist"/>
            <intype c-type="VideoLine_t" mode="sync" message-id="MSG_TSR_EDGEIMAGE" sizeof="48" cycle-time-us="90" type="TSR-EdgeImage"/>
            <implementation name="HAGL_Framework" stylesheet="HAGL-Framework/TSR.xsl">
                <interface header="TSR-Contour.h"/>
            </implementation>
        </processing-node>
    </nodes>
    <connections>
        <connect from="sensor" to="process1"/>
        <connect from="process1" to="actuator1" from-port="edge-level"/>
        <connect from="process1" to="actuator2" from-port="result"/>
    </connections>
</signal-graph>
`, 4, 3},
	}

	for i, c := range case1 {
		SignalGraphInit()
		var sg SignalGraph = SignalGraphNew()
		buf := make([]byte, len(c.xml))
		for i := 0; i < len(c.xml); i++ {
			buf[i] = c.xml[i]
		}
		err := sg.Read(buf)
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
