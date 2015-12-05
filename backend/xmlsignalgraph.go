package backend

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/sge/tool"
)

type XmlSignalGraph struct {
	XMLName         xml.Name        `xml:"http://www.freesp.de/xml/freeSP signal-graph"`
	Version         string          `xml:"version,attr"`
	Libraries       []XmlLibraryRef `xml:"library"`
	SignalTypes     []XmlSignalType `xml:"signal-type"`
	InputNodes      []XmlInputNode  `xml:"nodes>input"`
	OutputNodes     []XmlOutputNode `xml:"nodes>output"`
	ProcessingNodes []XmlProcessingNode `xml:"nodes>processing-node"`
	Connections     []XmlConnect    `xml:"connections>connect"`
}

func XmlSignalGraphNew() *XmlSignalGraph {
    return &XmlSignalGraph{xml.Name{"http://www.freesp.de/xml/freeSP", "signal-graph"}, "1.0", nil, nil, nil, nil, nil, nil}
}

func (g *XmlSignalGraph) Read(data []byte) error {
	err := xml.Unmarshal(data, g)
	if err != nil {
		fmt.Printf("XmlSignalGraph.Read error: %v", err)
	}
	return err
}

func (g *XmlSignalGraph) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(g, "", "   ")
    if err != nil {
		err = fmt.Errorf("XmlSignalGraph.Write error: %v", err)
    }
    return
}

func (g *XmlSignalGraph) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("signalgraph.ReadFile error: Failed to read file %s", filepath)
	}
	err = g.Read(data)
	if err != nil {
		return fmt.Errorf("signalgraph.ReadFile error: %v", err)
	}
	return err
}

func (g *XmlSignalGraph) WriteFile(filepath string) error {
	// TODO
	return fmt.Errorf("WriteFile() interface not implemented")
}

func NewXmlSignalGraph() *XmlSignalGraph {
	return &XmlSignalGraph{xml.Name{"", ""}, "", nil, nil, nil, nil, nil, nil}
}

