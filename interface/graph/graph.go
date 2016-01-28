package graph

import (
	"fmt"
	"github.com/axel-freesp/sge/tool"
	"image"
	"strings"
)

var _ fmt.Stringer = PortDirection(false)

func (d PortDirection) String() (s string) {
	if d == InPort {
		s = "Input"
	} else {
		s = "Output"
	}
	return
}

func CreatePathMode(path string, mode PositionMode) (m string) {
	m = fmt.Sprintf("%s/%s", path, string(mode))
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

func CreateModePositioner(path string, pmp PathModePositioner) (mp ModePositioner) {
	return &ModePositionerConverter{path, pmp}
}

var _ ModePositioner = (*ModePositionerConverter)(nil)

type ModePositionerConverter struct {
	path string
	ref  PathModePositioner
}

func (c ModePositionerConverter) ModePosition(mode PositionMode) (pos image.Point) {
	return c.ref.PathModePosition(c.path, mode)
}

func (c *ModePositionerConverter) SetModePosition(mode PositionMode, pos image.Point) {
	c.ref.SetPathModePosition(c.path, mode, pos)
}
