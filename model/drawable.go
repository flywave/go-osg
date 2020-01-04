package model

const (
	DRAWABLET string = "osg::Drawable"
)

type DrawCallback struct {
	Object
}

type ComputeBoundingBoxCallback struct {
	Object
}

type DrawableInterface interface {
	GetBoundingBox() *[2][3]float32
	SetBoundingBox(*[2][3]float32)
	GetShape() *Shape
	SetShape(*Shape)
	GetSupportsDisplayList() *bool
	SetSupportsDisplayList(bool)
	GetUseDisplayList() *bool
	SetUseDisplayList(bool)
	GetUseVertexBufferObjects() *bool
	SetUseVertexBufferObjects(bool)

	GetBoxCallback() *ComputeBoundingBoxCallback
	SetBoxCallback(*ComputeBoundingBoxCallback)

	GetDwCallback() *DrawCallback
	SetDwCallback(*DrawCallback)
}

type Drawable struct {
	Node
	BoundingBox            *[2][3]float32
	Shape                  *Shape
	SupportsDisplayList    bool
	UseDisplayList         bool
	UseVertexBufferObjects bool

	BoxCallback *ComputeBoundingBoxCallback
	DwCallback  *DrawCallback
}

func (d *Drawable) GetBoundingBox() *[2][3]float32 {
	return d.BoundingBox
}
func (d *Drawable) SetBoundingBox(bx *[2][3]float32) {
	d.BoundingBox = bx
}
func (d *Drawable) GetShape() *Shape {
	return d.Shape
}
func (d *Drawable) SetShape(sp *Shape) {
	d.Shape = sp
}
func (d *Drawable) GetSupportsDisplayList() *bool {
	return &d.SupportsDisplayList
}
func (d *Drawable) SetSupportsDisplayList(b bool) {
	d.SupportsDisplayList = b
}
func (d *Drawable) GetUseDisplayList() *bool {
	return &d.UseDisplayList
}
func (d *Drawable) SetUseDisplayList(b bool) {
	d.UseDisplayList = b
}
func (d *Drawable) GetUseVertexBufferObjects() *bool {
	return &d.UseVertexBufferObjects
}
func (d *Drawable) SetUseVertexBufferObjects(b bool) {
	d.UseVertexBufferObjects = b
}

func (d *Drawable) GetBoxCallback() *ComputeBoundingBoxCallback {
	return d.BoxCallback
}
func (d *Drawable) SetBoxCallback(cb *ComputeBoundingBoxCallback) {
	d.BoxCallback = cb
}

func (d *Drawable) GetDwCallback() *DrawCallback {
	return d.DwCallback
}
func (d *Drawable) SetDwCallback(cb *DrawCallback) {
	d.DwCallback = cb
}

func NewDrawable() *Drawable {
	n := NewNode()
	n.Type = DRAWABLET
	return &Drawable{Node: *n}
}
