package model

const (
	GEODET string = "osg::Geode"
)

type Geode struct {
	Group
}

func NewGeode() *Geode {
	g := NewGroup()
	g.Type = GEODET
	return &Geode{Group: *g}
}

func (g *Geode) Accept(nv *NodeVisitor) {
	if nv.ValidNodeMask(g) {
		nv.PushOntoNodePath(g)
		nv.Apply(g)
		nv.PopFromNodePath()
	}
}

func (g *Geode) AddDrawable(d *Drawable) {
	g.AddChild(d)
}

func (g *Geode) RemoveDrawable(d *Drawable) error {
	return g.RemoveChild(d)
}

func (g *Geode) RemoveDrawableCount(pos int, count int) error {
	return g.RemoveChild2(pos, count)
}

func (g *Geode) ReplaceDrawable(o *Drawable, n *Drawable) error {
	return g.ReplaceChild(o, n)
}

func (g *Geode) SetDrawable(pos int, d *Drawable) error {
	return g.SetChild(pos, d)
}
