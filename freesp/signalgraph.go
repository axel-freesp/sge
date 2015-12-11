package freesp

import (
	"fmt"
	"github.com/axel-freesp/sge/backend"
	"github.com/axel-freesp/sge/tool"
)

func SignalGraphNew(filename string) *signalGraph {
	return &signalGraph{filename, SignalGraphTypeNew()}
}

type signalGraph struct {
	filename string
	itsType  SignalGraphType
}

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
	s.itsType, err = createSignalGraphTypeFromXml(g, s.filename,
		func(_ string, _ PortDirection) *namedPortType { return nil })
	return err
}

func (s *signalGraph) ReadFile(filepath string) error {
	data, err := tool.ReadFile(filepath)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalGraph.ReadFile: %v", err))
	}
	err = s.Read(data)
	if err != nil {
		return newSignalGraphError(fmt.Sprintf("signalgraph.ReadFile: %v", err))
	}
	return err
}

func (s *signalGraph) Write() (data []byte, err error) {
	xmlsignalgraph := CreateXmlSignalGraph(s)
	data, err = xmlsignalgraph.Write()
	return
}

func (s *signalGraph) WriteFile(filepath string) error {
	data, err := s.Write()
	if err != nil {
		return err
	}
	err = tool.WriteFile(filepath, data)
	return nil
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
