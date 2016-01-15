package graph

import (
    "image"
	//"github.com/gotk3/gotk3/cairo"
	interfaces "github.com/axel-freesp/sge/interface"
)

type Container struct {
    NamedBoxObject
	children []ContainerChild
}

type ContainerChild interface {
	Layout() image.Rectangle
}

func ContainerInit(children []ContainerChild, config DrawConfig, namer interfaces.Namer) (c Container) {
	c = Container{NamedBoxObjectInit(image.Rectangle{}, config, namer), children}
	c.box = c.Layout()
	return
}

func ContainerNew(children []ContainerChild, config DrawConfig, namer interfaces.Namer) (c *Container) {
	c = &Container{NamedBoxObjectInit(image.Rectangle{}, config, namer), children}
	c.box = c.Layout()
	return
}

func (c Container) Layout() (box image.Rectangle) {
	for _, ch := range c.children {
		box = ContainerFit(box, ch.Layout())
	}
	return
}

func ContainerFit(outer, inner image.Rectangle) image.Rectangle {
	borderTop := image.Point{-18, -30}
	borderBottom := image.Point{18, 18}
	test := image.Rectangle{inner.Min.Add(borderTop), inner.Max.Add(borderBottom)}
	if outer.Size().X == 0 {
		return test
	}
	return outer.Union(test)
}


