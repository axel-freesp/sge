package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlSignalType struct {
	XMLName xml.Name `xml:"signal-type"`
	Name    string   `xml:"name,attr"`
	Scope   string   `xml:"scope,attr"`
	Mode    string   `xml:"mode,attr"`
	Ctype   string   `xml:"c-type,attr"`
	Msgid   string   `xml:"message-id,attr"`
}

func XmlSignalTypeNew(name, scope, mode, ctype, msgid string) *XmlSignalType {
	return &XmlSignalType{xml.Name{freespNamespace, "signal-type"}, name, scope, mode, ctype, msgid}
}

func (t *XmlSignalType) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, t)
	if err != nil {
		fmt.Printf("XmlSignalType.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (t *XmlSignalType) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(t, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlSignalType.Write error: %v", err)
	}
	return
}
