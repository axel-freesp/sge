package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
)

func SignalGraphNew(filename string) *signalGraph {
	return &signalGraph{filename, SignalGraphTypeNew()}
}

type signalGraph struct {
	filename string
	itsType  SignalGraphType
}

var _ SignalGraph = (*signalGraph)(nil)

func (s *signalGraph) Filename() string {
	return s.filename
}

func (s *signalGraph) ItsType() SignalGraphType {
	return s.itsType
}

func (s *signalGraph) Read(data []byte) error {
	g := backend.XmlSignalGraphNew()
	err := g.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.Read: %v", err))
	}
	fmt.Println("signalGraph.Read: call createSignalGraphTypeFromXml")
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename,
		func(_ string, _ PortDirection) *namedPortType { return nil })
	return err
}

func (s *signalGraph) ReadFile(filepath string) error {
	g := backend.XmlSignalGraphNew()
	err := g.ReadFile(filepath)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.ReadFile: %v", err))
	}
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename,
		func(_ string, _ PortDirection) *namedPortType { return nil })
	return err
}

func (s *signalGraph) Write() (data []byte, err error) {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	data, err = xmlsignalgraph.Write()
	return
}

func (s *signalGraph) WriteFile(filepath string) error {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	return xmlsignalgraph.WriteFile(filepath)
}

func (s *signalGraph) SetFilename(filename string) {
	s.filename = filename
}

//------------------------------

type signalGraphError struct {
	reason string
}

func (e *signalGraphError) Error() string {
	return fmt.Sprintf("signal graph error: %s", e.reason)
}

func newSignalGraphError(reason string) *signalGraphError {
	return &signalGraphError{reason}
}
