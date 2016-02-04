package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlPlatformHint struct {
	XMLName xml.Name         `xml:"http://www.freesp.de/xml/freeSP hints"`
	Ref     string           `xml:"ref,attr"`
	Arch    []XmlArchPosHint `xml:"arch"`
}

func XmlPlatformHintNew(ref string) *XmlPlatformHint {
	return &XmlPlatformHint{xml.Name{freespNamespace, "hints"}, ref, nil}
}

func (h *XmlPlatformHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlHint.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlPlatformHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlHint.Write error: %v", err)
	}
	return
}

type XmlArchPosHint struct {
	Name      string               `xml:"name,attr"`
	ArchPorts []XmlArchPortPosHint `xml:"arch-port"`
	Processes []XmlProcessPosHint  `xml:"process"`
	XmlModeHint
}

func XmlArchPosHintNew(name string) *XmlArchPosHint {
	return &XmlArchPosHint{name, nil, nil, XmlModeHint{}}
}

type XmlArchPortPosHint struct {
	XmlModeHint
}

func XmlArchPortPosHintNew() *XmlArchPortPosHint {
	return &XmlArchPortPosHint{XmlModeHint{}}
}

type XmlProcessPosHint struct {
	Name        string              `xml:"name,attr"`
	InChannels  []XmlChannelPosHint `xml:"in-channel"`
	OutChannels []XmlChannelPosHint `xml:"out-channel"`
	XmlModeHint
}

func XmlProcessPosHintNew(name string) *XmlProcessPosHint {
	return &XmlProcessPosHint{name, nil, nil, XmlModeHint{}}
}

type XmlChannelPosHint struct {
	XmlModeHint
}

func XmlChannelPosHintNew() *XmlChannelPosHint {
	return &XmlChannelPosHint{XmlModeHint{}}
}
