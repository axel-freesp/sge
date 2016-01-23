package backend

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/sge/tool"
)

type XmlSignalGraph struct {
	XMLName         xml.Name            `xml:"http://www.freesp.de/xml/freeSP signal-graph"`
	Version         string              `xml:"version,attr"`
	Libraries       []XmlLibraryRef     `xml:"library"`
	InputNodes      []XmlInputNode      `xml:"nodes>input"`
	OutputNodes     []XmlOutputNode     `xml:"nodes>output"`
	ProcessingNodes []XmlProcessingNode `xml:"nodes>processing-node"`
	Connections     []XmlConnect        `xml:"connections>connect"`
}

func XmlSignalGraphNew() *XmlSignalGraph {
	return &XmlSignalGraph{xml.Name{freespNamespace, "signal-graph"}, "1.0", nil, nil, nil, nil, nil}
}

func (g *XmlSignalGraph) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, g)
	cnt = len(data)
	return
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
		return fmt.Errorf("XmlSignalGraph.ReadFile error: Failed to read file %s", filepath)
	}
	_, err = g.Read(data)
	if err != nil {
		return fmt.Errorf("XmlSignalGraph.ReadFile error: %v", err)
	}
	return err
}

func (g *XmlSignalGraph) WriteFile(filepath string) error {
	data, err := g.Write()
	if err != nil {
		return err
	}
	buf := make([]byte, len(data)+len(xmlHeader))
	for i := 0; i < len(xmlHeader); i++ {
		buf[i] = xmlHeader[i]
	}
	for i := 0; i < len(data); i++ {
		buf[i+len(xmlHeader)] = data[i]
	}
	return tool.WriteFile(filepath, buf)
}
