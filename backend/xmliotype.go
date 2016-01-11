package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlIOType struct {
	XMLName xml.Name  `xml:"io-type"`
	Name    string    `xml:"name,attr"`
	Mode    XmlIOMode `xml:"mode,attr"`
	Hint    XmlHint   `xml:"hint"`
}

func XmlIOTypeNew(name string, mode XmlIOMode, hint XmlHint) *XmlIOType {
	return &XmlIOType{xml.Name{freespNamespace, "io-type"}, name, mode, hint}
}

type XmlIOMode string

const (
	IOModeShmem XmlIOMode = "shmem"
	IOModeAsync XmlIOMode = "async"
	IOModeSync  XmlIOMode = "sync"
)

func (t *XmlIOType) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, t)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	return
}

func (t *XmlIOType) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(t, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}
