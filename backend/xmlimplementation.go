package backend

import (
	"encoding/xml"
)

type XmlImplementation struct {
	XMLName     xml.Name         `xml:"implementation"`
	Name        string           `xml:"name,attr"`
	SignalGraph []XmlSignalGraph `xml:"signal-graph"`
}

func XmlImplementationNew(name string) *XmlImplementation {
	return &XmlImplementation{xml.Name{freespNamespace, "implementation"}, name, nil}
}
