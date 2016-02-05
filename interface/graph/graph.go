package graph

import (
	"fmt"
	"github.com/axel-freesp/sge/tool"
	"image"
	//"log"
	"strings"
)

var ValidModes = []PositionMode{
	PositionModeNormal,
	PositionModeMapping,
	PositionModeExpanded,
}

var StringFromMode = map[PositionMode]string{
	PositionModeNormal:   "normal",
	PositionModeMapping:  "mapping",
	PositionModeExpanded: "expanded",
}

var ModeFromString = map[string]PositionMode{
	"normal":   PositionModeNormal,
	"mapping":  PositionModeMapping,
	"expanded": PositionModeExpanded,
}

//
//	ModePositioner Proxy
//

type ModePositionerProxy struct {
	ref        ModePositioner
	activeMode PositionMode
}

var _ Positioner = (*ModePositionerProxy)(nil)

func ModePositionerProxyNew(ref ModePositioner, activeMode PositionMode) *ModePositionerProxy {
	return &ModePositionerProxy{ref, activeMode}
}

func (m *ModePositionerProxy) Position() image.Point {
	return m.ref.ModePosition(m.activeMode)
}

func (m *ModePositionerProxy) SetPosition(pos image.Point) {
	m.ref.SetModePosition(m.activeMode, pos)
}

//
//	PathModePositioner standard implementation
//

type PathModePositionerObject struct {
	ModePositionerObject
	pathlist   tool.StringList
	activePath string
}

func PathModePositionerObjectNew() (p *PathModePositionerObject) {
	p = &PathModePositionerObject{*ModePositionerObjectNew(), tool.StringListInit(), ""}
	p.pathlist.Append("")
	return
}

func (p *PathModePositionerObject) Position() image.Point {
	return p.ModePosition(p.activeMode)
}

func (p *PathModePositionerObject) SetPosition(pos image.Point) {
	p.SetModePosition(p.activeMode, pos)
}

func (p *PathModePositionerObject) ModePosition(mode PositionMode) (pos image.Point) {
	return p.PathModePosition(p.ActivePath(), mode)
}

func (p *PathModePositionerObject) SetModePosition(mode PositionMode, pos image.Point) {
	p.SetPathModePosition(p.ActivePath(), mode, pos)
}

func (p *PathModePositionerObject) PathModePosition(path string, mode PositionMode) (pos image.Point) {
	_, ok := p.pathlist.Find(path)
	if !ok {
		return
	}
	pos = p.position[CreatePathMode(path, mode)]
	return
}

func (p *PathModePositionerObject) SetPathModePosition(path string, mode PositionMode, pos image.Point) {
	_, ok := p.pathlist.Find(path)
	if !ok {
		p.pathlist.Append(path)
	}
	p.position[CreatePathMode(path, mode)] = pos
}

func (p *PathModePositionerObject) PathList() []string {
	return p.pathlist.Strings()
}

func (p *PathModePositionerObject) SetActivePath(path string) {
	_, ok := p.pathlist.Find(path)
	if !ok {
		p.pathlist.Append(path)
	}
	p.activePath = path
}

func (p *PathModePositionerObject) ActivePath() string {
	return p.activePath
}

//
//	ModePositioner standard implementation
//

type ModePositionerObject struct {
	position   map[PositionMode]image.Point
	activeMode PositionMode
}

func ModePositionerObjectNew() *ModePositionerObject {
	return &ModePositionerObject{make(map[PositionMode]image.Point), PositionMode("normal")}
}

func (m *ModePositionerObject) Position() image.Point {
	return m.ModePosition(m.activeMode)
}

func (m *ModePositionerObject) SetPosition(pos image.Point) {
	m.SetModePosition(m.activeMode, pos)
}

func (m *ModePositionerObject) ModePosition(mode PositionMode) (p image.Point) {
	return m.position[mode]
}

func (m *ModePositionerObject) SetModePosition(mode PositionMode, p image.Point) {
	m.position[mode] = p
}

func (m *ModePositionerObject) SetActiveMode(mode PositionMode) {
	m.activeMode = mode
}

func (m *ModePositionerObject) ActiveMode() PositionMode {
	return m.activeMode
}

//
//	Positioner standard implementation
//

type PositionerObject struct {
	position image.Point
}

func PositionerObjectNew() *PositionerObject {
	return &PositionerObject{}
}

func (p *PositionerObject) Position() image.Point {
	return p.position
}

func (p *PositionerObject) SetPosition(pos image.Point) {
	p.position = pos
}

//
//	PortDirection
//

var _ fmt.Stringer = PortDirection(false)

func (d PortDirection) String() (s string) {
	if d == InPort {
		s = "Input"
	} else {
		s = "Output"
	}
	return
}

func CreatePathMode(path string, mode PositionMode) (m PositionMode) {
	m = PositionMode(fmt.Sprintf("%s/%s", path, string(mode)))
	return
}

func SeparatePathMode(m string) (path string, mode PositionMode) {
	path = tool.Dirname(m)
	mode = PositionMode(tool.Basename(m))
	return
	if !strings.Contains(m, "/") {
		mode = PositionMode(m)
		return
	}
	p := strings.Split(m, "/")
	if len(p) < 2 {
		mode = PositionMode(m)
		return
	}
	l := len(p) - 1
	for i, s := range p {
		switch {
		case i == 0:
			path = s
		case i < l:
			path = fmt.Sprintf("%s/%s", path, s)
		}
	}
	mode = PositionMode(p[l])
	return
}
