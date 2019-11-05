package model

import "unsafe"

const (
	GeodeType string = "osg::Geode"
)

type Geode struct {
	Group
}

func (g *Geode) AddDrawable(d *Drawable) {
	uP := unsafe.Pointer(d)
	g.AddChild((*Node)(uP))
}

func (g *Geode) RemoveDrawable(d *Drawable) error {
	uP := unsafe.Pointer(d)
	return g.RemoveChild((*Node)(uP))
}

func (g *Geode) RemoveDrawableCount(pos int, count int) error {
	return g.RemoveChild2(pos, count)
}

func (g *Geode) ReplaceDrawable(o *Drawable, n *Drawable) error {
	uP := unsafe.Pointer(o)
	uP1 := unsafe.Pointer(n)
	return g.ReplaceChild((*Node)(uP), (*Node)(uP1))
}

func (g *Geode) SetDrawable(pos int, d *Drawable) error {
	uP := unsafe.Pointer(d)
	return g.SetChild(pos, (*Node)(uP))
}
