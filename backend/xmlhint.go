package backend

import (
	"encoding/xml"
	"fmt"
)

type XmlNodeHint struct {
	XmlModeHint
	Expanded bool `xml:"expanded,attr"`
}

type XmlHint struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

func XmlHintNew(x, y int) *XmlHint {
	return &XmlHint{x, y}
}

func (h *XmlHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlShape struct {
	W int `xml:"w,attr"`
	H int `xml:"h,attr"`
}

func XmlShapeNew(w, h int) *XmlShape {
	return &XmlShape{w, h}
}

func (h *XmlShape) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlShape) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlRectangle struct {
	XmlHint
	XmlShape
}

func XmlRectangleNew(x, y, w, h int) *XmlRectangle {
	return &XmlRectangle{XmlHint{x, y}, XmlShape{w, h}}
}

func (h *XmlRectangle) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlRectangle) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlModeHint struct {
	Entry []XmlModeHintEntry `xml:"hint"`
}

func XmlModeHintNew() *XmlModeHint {
	return &XmlModeHint{}
}

func (h *XmlModeHint) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlModeHint) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlModeRectangle struct {
	Entry []XmlModeRectEntry `xml:"hint"`
}

func XmlModeRectNew() *XmlModeRectangle {
	return &XmlModeRectangle{}
}

func (h *XmlModeRectangle) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlModeRectangle) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlModeHintEntry struct {
	Mode string `xml:"mode,attr"`
	XmlHint
}

func XmlModeHintEntryNew(mode string, x, y int) *XmlModeHintEntry {
	return &XmlModeHintEntry{mode, XmlHint{x, y}}
}

func (h *XmlModeHintEntry) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlModeHintEntry) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}

type XmlModeRectEntry struct {
	Mode string `xml:"mode,attr"`
	XmlHint
	XmlShape
}

func XmlModeRectEntryNew(mode string, x, y, w, h int) *XmlModeRectEntry {
	return &XmlModeRectEntry{mode, XmlHint{x, y}, XmlShape{w, h}}
}

func (h *XmlModeRectEntry) Read(data []byte) (cnt int, err error) {
	err = xml.Unmarshal(data, h)
	if err != nil {
		err = fmt.Errorf("XmlConnect.Read error: %v", err)
	}
	cnt = len(data)
	return
}

func (h *XmlModeRectEntry) Write() (data []byte, err error) {
	data, err = xml.MarshalIndent(h, "", "   ")
	if err != nil {
		err = fmt.Errorf("XmlConnect.Write error: %v", err)
	}
	return
}
