package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlChannel struct {
	IOType string `xml:"io-type,attr"`
}

type XmlInChannel struct {
	XmlChannel
	XMLName xml.Name `xml:"input-channel"`
	Source  string   `xml:"source,attr"`
}

type XmlOutChannel struct {
	XmlChannel
	XMLName xml.Name `xml:"output-channel"`
	Dest    string   `xml:"dest,attr"`
}

func XmlInChannelNew(name, ioType, source string) *XmlInChannel {
	return &XmlInChannel{XmlChannel{ioType}, xml.Name{freespNamespace, "input-channel"}, source}
}

func XmlOutChannelNew(name, ioType, dest string) *XmlOutChannel {
	return &XmlOutChannel{XmlChannel{ioType}, xml.Name{freespNamespace, "output-channel"}, dest}
}

func (c *XmlInChannel) Read(data []byte) (err error) {
	err = xml.Unmarshal(data, c)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	return
}

func (c *XmlInChannel) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(c, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

func (c *XmlOutChannel) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, c)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (c *XmlOutChannel) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(c, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}
