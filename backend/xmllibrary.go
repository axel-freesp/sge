package backend

import (
	"encoding/xml"
	"fmt"
	"github.com/axel-freesp/sge/tool"
)

type XmlLibrary struct {
	XMLName     xml.Name        `xml:"http://www.freesp.de/xml/freeSP library"`
	Version     string          `xml:"version,attr"`
	Libraries   []XmlLibraryRef `xml:"library"`
	SignalTypes []XmlSignalType `xml:"signal-type"`
	NodeTypes   []XmlNodeType   `xml:"node-type"`
}

func XmlLibraryNew() *XmlLibrary {
	return &XmlLibrary{xml.Name{freespNamespace, "signal-graph"}, "1.0", nil, nil, nil}
}

func (g *XmlLibrary) Read(data []byte) error {
	err := xml.Unmarshal(data, g)
	if err != nil {
		fmt.Printf("XmlLibrary.Read error: %v", err)
	}
	return err
}

func (g *XmlLibrary) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(g, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlLibrary.Write error: %v", err)
	}
	return
}

func (g *XmlLibrary) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("XmlLibrary.ReadFile error: Failed to read file %s", filepath)
	}
	err = g.Read(data)
	if err != nil {
		return fmt.Errorf("XmlLibrary.ReadFile error: %v", err)
	}
	return err
}

func (g *XmlLibrary) WriteFile(filepath string) error {
	// TODO
	return fmt.Errorf("XmlLibrary.WriteFile() interface not implemented")
}