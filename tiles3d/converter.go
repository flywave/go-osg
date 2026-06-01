package tiles3d

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/flywave/go-draco"
	"github.com/flywave/go-meshopt"
	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
	"gonum.org/v1/gonum/mat"
)

type GeometryConverter struct {
	opts       *ConverterOptions
	coordTrans *CoordinateTransformer
	geoidConv  *GeoidConverter
}

func NewGeometryConverter(opts *ConverterOptions, coordTrans *CoordinateTransformer, geoidConv *GeoidConverter) *GeometryConverter {
	return &GeometryConverter{
		opts:       opts,
		coordTrans: coordTrans,
		geoidConv:  geoidConv,
	}
}

func (c *GeometryConverter) Convert(node interface{}) *TileContent {
	content := &TileContent{}

	switch n := node.(type) {
	case *model.Group:
		children := n.GetChildren()
		for _, child := range children {
			childContent := c.Convert(child)
			if childContent != nil {
				content.Merge(childContent)
			}
		}
	case *model.Lod:
		children := n.GetChildren()
		if len(children) > 0 {
			content = c.Convert(children[0])
		}
	case *model.PagedLod:
		children := n.GetChildren()
		if len(children) > 0 {
			content = c.Convert(children[0])
		}
	case *model.MatrixTransform:
		children := n.GetChildren()
		for _, child := range children {
			childContent := c.Convert(child)
			if childContent != nil {
				c.applyMatrixTransform(childContent, n.Matrix)
				content.Merge(childContent)
			}
		}
	case *model.PositionAttitudeTransform:
		children := n.GetChildren()
		for _, child := range children {
			childContent := c.Convert(child)
			if childContent != nil {
				content.Merge(childContent)
			}
		}
	case *model.Geode:
		children := n.GetChildren()
		for _, child := range children {
			if geometry, ok := child.(*model.Geometry); ok {
				geomContent := c.extractGeometry(geometry)
				if geomContent != nil {
					content.Merge(geomContent)
				}
			}
		}
	}

	return content
}

func (c *GeometryConverter) extractGeometry(geom *model.Geometry) *TileContent {
	content := &TileContent{}

	if geom.VertexArray == nil || geom.VertexArray.Data == nil {
		return nil
	}

	vertices := c.extractVertices(geom.VertexArray)
	if len(vertices) == 0 {
		return nil
	}

	bbox := c.calculateBoundingBox(vertices)
	content.BoundingBox = bbox
	content.Vertices = vertices

	if geom.NormalArray != nil && geom.NormalArray.Data != nil {
		content.Normals = c.extractNormals(geom.NormalArray)
	}

	if len(geom.TexCoordArrayList) > 0 && geom.TexCoordArrayList[0] != nil {
		content.TexCoords = c.extractTexCoords(geom.TexCoordArrayList[0])
	} else {
		content.TexCoords = c.generateDefaultTexCoords(vertices)
	}

	content.Mode = model.GLTRIANGLES
	if geom.Primitives != nil && len(geom.Primitives) > 0 {
		for _, prim := range geom.Primitives {
			mode, indices := c.extractIndicesWithMode(prim)
			if mode >= 0 {
				content.Mode = mode
				if len(indices) > 0 {
					if content.Indices == nil {
						content.Indices = indices
					}
					content.Primitives = append(content.Primitives, PrimitiveInfo{
						Indices: indices,
						Mode:    mode,
					})
				}
			}
		}
	}

	if len(content.Normals) == 0 && len(content.Indices) > 0 {
		content.Normals = generateNormals(content.Vertices, content.Indices)
	}

	content.BatchLength = 1

	if len(vertices) >= 12 {
		fmt.Printf("DEBUG extractGeometry first 3 vertices: (%.6f, %.6f, %.6f)\n", vertices[0], vertices[1], vertices[2])
		fmt.Printf("DEBUG extractGeometry next 3 vertices: (%.6f, %.6f, %.6f)\n", vertices[len(vertices)-9], vertices[len(vertices)-8], vertices[len(vertices)-7])
		fmt.Printf("DEBUG extractGeometry last 3 vertices: (%.6f, %.6f, %.6f)\n", vertices[len(vertices)-6], vertices[len(vertices)-5], vertices[len(vertices)-4])
	}

	stateset := geom.GetStates()
	if stateset != nil {
		if tex := c.extractTextureFromStateSet(stateset, "", 0); tex != nil {
			content.Textures = append(content.Textures, tex.Data)
		}
	}

	return content
}

func (c *GeometryConverter) applyMatrixTransform(content *TileContent, matrix [4][4]float32) {
	if len(content.Vertices) == 0 {
		return
	}

	for i := 0; i < len(content.Vertices); i += 3 {
		x := float64(content.Vertices[i])
		y := float64(content.Vertices[i+1])
		z := float64(content.Vertices[i+2])

		m00 := float64(matrix[0][0])
		m01 := float64(matrix[0][1])
		m02 := float64(matrix[0][2])
		m03 := float64(matrix[0][3])
		m10 := float64(matrix[1][0])
		m11 := float64(matrix[1][1])
		m12 := float64(matrix[1][2])
		m13 := float64(matrix[1][3])
		m20 := float64(matrix[2][0])
		m21 := float64(matrix[2][1])
		m22 := float64(matrix[2][2])
		m23 := float64(matrix[2][3])

		newX := m00*x + m01*y + m02*z + m03
		newY := m10*x + m11*y + m12*z + m13
		newZ := m20*x + m21*y + m22*z + m23

		content.Vertices[i] = float32(newX)
		content.Vertices[i+1] = float32(newY)
		content.Vertices[i+2] = float32(newZ)
	}

	if len(content.Normals) > 0 {
		for i := 0; i < len(content.Normals); i += 3 {
			x := float64(content.Normals[i])
			y := float64(content.Normals[i+1])
			z := float64(content.Normals[i+2])

			m00 := float64(matrix[0][0])
			m01 := float64(matrix[0][1])
			m02 := float64(matrix[0][2])
			m10 := float64(matrix[1][0])
			m11 := float64(matrix[1][1])
			m12 := float64(matrix[1][2])
			m20 := float64(matrix[2][0])
			m21 := float64(matrix[2][1])
			m22 := float64(matrix[2][2])

			newX := m00*x + m01*y + m02*z
			newY := m10*x + m11*y + m12*z
			newZ := m20*x + m21*y + m22*z

			content.Normals[i] = float32(newX)
			content.Normals[i+1] = float32(newY)
			content.Normals[i+2] = float32(newZ)
		}
	}

	content.BoundingBox = calculateBoundingBox(content.Vertices)
}

const (
	meterToLat = 0.000000157891
	meterToLon = 0.000000156785
)

func meterToLatDeg(m float64) float64 {
	return m * meterToLat
}

func meterToLonDeg(m float64, latRad float64) float64 {
	return m * meterToLon / math.Cos(latRad)
}

func (c *GeometryConverter) extractVertices(arr *model.Array) []float32 {
	if arr == nil || arr.Data == nil {
		return nil
	}

	fmt.Printf("DEBUG extractVertices: Array properties:\n")
	fmt.Printf("  Type: %d\n", arr.Type)
	fmt.Printf("  DataType: %d\n", arr.DataType)
	fmt.Printf("  DataSize: %d\n", arr.DataSize)
	fmt.Printf("  Binding: %d\n", arr.Binding)
	fmt.Printf("  Normalize: %v\n", arr.Normalize)
	fmt.Printf("  PreserveDataType: %v\n", arr.PreserveDataType)
	fmt.Printf("  Data type: %T\n", arr.Data)

	var vertices []float32

	switch data := arr.Data.(type) {
	case []float32:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		copy(vertices, data)
	case [][3]float32:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data)*3)
		for i, v := range data {
			vertices[i*3] = v[0]
			vertices[i*3+1] = v[1]
			vertices[i*3+2] = v[2]
		}
	case []float64:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	case []int16:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	case []uint16:
		if len(data) == 0 {
			return nil
		}
		vertices = make([]float32, len(data))
		for i, v := range data {
			vertices[i] = float32(v)
		}
	default:
		return nil
	}

	if len(vertices) < 3 {
		return nil
	}

	if len(vertices)%3 != 0 {
		vertices = vertices[:len(vertices)/3*3]
	}

	fmt.Printf("DEBUG extractVertices: len(vertices) = %d\n", len(vertices))

	// 打印原始顶点数据（前3个顶点）
	if len(vertices) >= 9 {
		fmt.Printf("DEBUG extractVertices: raw vertices (first 3):\n")
		for i := 0; i < 9; i += 3 {
			fmt.Printf("  [%d] x=%.6f, y=%.6f, z=%.6f\n", i/3, vertices[i], vertices[i+1], vertices[i+2])
		}
	} else {
		fmt.Printf("DEBUG extractVertices: len(vertices) < 9, cannot print raw vertices\n")
	}

	originOffset := c.coordTrans.GetOriginOffset()
	hasGeoRef := c.coordTrans.HasGeoReference()

	fmt.Printf("DEBUG extractVertices: originOffset=(%.6f, %.6f, %.6f)\n", originOffset[0], originOffset[1], originOffset[2])
	fmt.Printf("DEBUG extractVertices: hasGeoRef=%v\n", hasGeoRef)

	if !hasGeoRef {
		result := make([]float32, len(vertices))
		for i := 0; i < len(vertices); i += 3 {
			result[i] = vertices[i] + float32(originOffset[0])
			result[i+1] = vertices[i+1] + float32(originOffset[1])
			result[i+2] = vertices[i+2] + float32(originOffset[2])
		}
		return result
	}

	// Use least squares method for coordinate transformation (matching C++ 3dtiles approach)
	// Step 1: Calculate bounding box
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	minZ, maxZ := math.MaxFloat64, -math.MaxFloat64

	fmt.Printf("DEBUG extractVertices: calculating bounding box\n")
	for i := 0; i < len(vertices); i += 3 {
		x := float64(vertices[i]) + originOffset[0]
		y := float64(vertices[i+1]) + originOffset[1]
		z := float64(vertices[i+2]) + originOffset[2]

		if i < 9 {
			fmt.Printf("DEBUG extractVertices: vertex[%d] after offset: x=%.6f, y=%.6f, z=%.6f\n", i/3, x, y, z)
		}

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
	}

	// Step 2: Transform 8 corner points
	originalCorners := [8][3]float64{
		{minX, minY, minZ},
		{maxX, minY, minZ},
		{minX, maxY, minZ},
		{minX, minY, maxZ},
		{maxX, maxY, minZ},
		{minX, maxY, maxZ},
		{maxX, minY, maxZ},
		{maxX, maxY, maxZ},
	}

	correctedCorners := make([][3]float64, 8)
	for i, corner := range originalCorners {
		correctedCorners[i] = c.coordTrans.ToLocalENUFromSource(corner)
	}

	// Step 3: Calculate transformation matrix using least squares
	// Solve A * X = B where A is 8x4 (original corners), B is 8x4 (corrected corners)
	A := mat.NewDense(8, 4, nil)
	B := mat.NewDense(8, 4, nil)

	for i := 0; i < 8; i++ {
		A.Set(i, 0, originalCorners[i][0])
		A.Set(i, 1, originalCorners[i][1])
		A.Set(i, 2, originalCorners[i][2])
		A.Set(i, 3, 1.0)
		B.Set(i, 0, correctedCorners[i][0])
		B.Set(i, 1, correctedCorners[i][1])
		B.Set(i, 2, correctedCorners[i][2])
		B.Set(i, 3, 1.0)
	}

	var X mat.Dense
	svdOk := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("DEBUG: SVD panicked: %v, using direct transform\n", r)
			}
		}()
		var svd mat.SVD
		if ok := svd.Factorize(A, mat.SVDFull); !ok {
			return
		}
		_ = svd.SolveTo(&X, B, 4)
		svdOk = true
	}()
	if !svdOk {
		return c.transformVerticesDirectly(vertices, originOffset)
	}

	fmt.Printf("DEBUG: Transform matrix calculated:\n")
	for i := 0; i < 4; i++ {
		fmt.Printf("  [%.6f, %.6f, %.6f, %.6f]\n", X.At(i, 0), X.At(i, 1), X.At(i, 2), X.At(i, 3))
	}

	// Step 4: Apply transformation to all vertices
	result := make([]float32, len(vertices))
	for i := 0; i < len(vertices); i += 3 {
		x := float64(vertices[i]) + originOffset[0]
		y := float64(vertices[i+1]) + originOffset[1]
		z := float64(vertices[i+2]) + originOffset[2]

		// Apply 4x4 transformation matrix (column-major: A*X = B where X[i,j] maps component i of input to component j of output)
		newX := X.At(0, 0)*x + X.At(1, 0)*y + X.At(2, 0)*z + X.At(3, 0)
		newY := X.At(0, 1)*x + X.At(1, 1)*y + X.At(2, 1)*z + X.At(3, 1)
		newZ := X.At(0, 2)*x + X.At(1, 2)*y + X.At(2, 2)*z + X.At(3, 2)

		result[i] = float32(newX)
		result[i+1] = float32(newY)
		result[i+2] = float32(newZ)

		if i == 0 {
			fmt.Printf("DEBUG: First vertex: original=(%f, %f, %f) -> transformed=(%f, %f, %f)\n",
				x, y, z, newX, newY, newZ)
		}
	}

	return result
}

func (c *GeometryConverter) transformVerticesDirectly(vertices []float32, originOffset [3]float64) []float32 {
	result := make([]float32, 0, len(vertices))
	invalidCount := 0

	for i := 0; i < len(vertices); i += 3 {
		absX := float64(vertices[i]) + originOffset[0]
		absY := float64(vertices[i+1]) + originOffset[1]
		absZ := float64(vertices[i+2]) + originOffset[2]

		const maxValidZ = 10000.0
		const maxValidXY = 1e12

		if absZ > maxValidZ || absZ < -maxValidZ {
			invalidCount++
			continue
		}
		if absX > maxValidXY || absX < -maxValidXY || absY > maxValidXY || absY < -maxValidXY {
			invalidCount++
			continue
		}

		localCoords := [3]float64{absX, absY, absZ}
		enuCoords := c.coordTrans.ToLocalENUFromSource(localCoords)

		if math.IsNaN(enuCoords[0]) || math.IsNaN(enuCoords[1]) || math.IsNaN(enuCoords[2]) ||
			math.IsInf(enuCoords[0], 0) || math.IsInf(enuCoords[1], 0) || math.IsInf(enuCoords[2], 0) {
			invalidCount++
			continue
		}

		result = append(result, float32(enuCoords[0]), float32(enuCoords[1]), float32(enuCoords[2]))
	}

	if invalidCount > 0 {
		fmt.Printf("DEBUG: Filtered %d invalid vertices out of %d total\n", invalidCount, len(vertices)/3)
	}

	return result
}

func normalizeVec3(x, y, z float32) (float32, float32, float32) {
	l := math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y) + float64(z)*float64(z))
	if l < 1e-20 {
		return 0.0, 0.0, 1.0
	}
	invLen := float32(1.0 / l)
	return x * invLen, y * invLen, z * invLen
}

func (c *GeometryConverter) extractNormals(arr *model.Array) []float32 {
	if arr.Data == nil {
		return nil
	}

	switch data := arr.Data.(type) {
	case []float32:
		normals := make([]float32, len(data))
		copy(normals, data)
		for i := 0; i < len(normals); i += 3 {
			nx, ny, nz := normalizeVec3(normals[i], normals[i+1], normals[i+2])
			normals[i], normals[i+1], normals[i+2] = nx, ny, nz
		}
		return normals
	case [][3]float32:
		normals := make([]float32, len(data)*3)
		for i, v := range data {
			nx, ny, nz := normalizeVec3(v[0], v[1], v[2])
			normals[i*3] = nx
			normals[i*3+1] = ny
			normals[i*3+2] = nz
		}
		return normals
	}

	return nil
}

func (c *GeometryConverter) extractTexCoords(arr *model.Array) []float32 {
	if arr.Data == nil {
		return nil
	}

	switch data := arr.Data.(type) {
	case []float32:
		return data
	case [][2]float32:
		texCoords := make([]float32, len(data)*2)
		for i, v := range data {
			texCoords[i*2] = v[0]
			texCoords[i*2+1] = v[1]
		}
		return texCoords
	}

	return nil
}

func (c *GeometryConverter) generateDefaultTexCoords(vertices []float32) []float32 {
	if len(vertices) == 0 {
		return nil
	}

	vertexCount := len(vertices) / 3

	minX, maxX := vertices[0], vertices[0]
	minY, maxY := vertices[1], vertices[1]

	for i := 0; i < len(vertices); i += 3 {
		if vertices[i] < minX {
			minX = vertices[i]
		}
		if vertices[i] > maxX {
			maxX = vertices[i]
		}
		if vertices[i+1] < minY {
			minY = vertices[i+1]
		}
		if vertices[i+1] > maxY {
			maxY = vertices[i+1]
		}
	}

	width := maxX - minX
	height := maxY - minY

	if width == 0 {
		width = 1.0
	}
	if height == 0 {
		height = 1.0
	}

	texCoords := make([]float32, vertexCount*2)
	for i := 0; i < vertexCount; i++ {
		x := vertices[i*3]
		y := vertices[i*3+1]
		texCoords[i*2] = (x - minX) / width
		texCoords[i*2+1] = (y - minY) / height
	}

	return texCoords
}

func (c *GeometryConverter) calculateBoundingBox(vertices []float32) [12]float64 {
	if len(vertices) == 0 {
		return [12]float64{0, 0, 0, 0.01, 0, 0, 0, 0.01, 0, 0, 0.01}
	}

	minX, maxX := float64(vertices[0]), float64(vertices[0])
	minY, maxY := float64(vertices[1]), float64(vertices[1])
	minZ, maxZ := float64(vertices[2]), float64(vertices[2])

	for i := 0; i < len(vertices); i += 3 {
		if float64(vertices[i]) < minX {
			minX = float64(vertices[i])
		}
		if float64(vertices[i]) > maxX {
			maxX = float64(vertices[i])
		}
		if float64(vertices[i+1]) < minY {
			minY = float64(vertices[i+1])
		}
		if float64(vertices[i+1]) > maxY {
			maxY = float64(vertices[i+1])
		}
		if float64(vertices[i+2]) < minZ {
			minZ = float64(vertices[i+2])
		}
		if float64(vertices[i+2]) > maxZ {
			maxZ = float64(vertices[i+2])
		}
	}

	centerX := (maxX + minX) / 2
	centerY := (maxY + minY) / 2
	centerZ := (maxZ + minZ) / 2

	xHalf := (maxX - minX) / 2
	yHalf := (maxY - minY) / 2
	zHalf := (maxZ - minZ) / 2

	if xHalf < 0.01 {
		xHalf = 0.01
	}
	if yHalf < 0.01 {
		yHalf = 0.01
	}
	if zHalf < 0.01 {
		zHalf = 0.01
	}

	return [12]float64{
		centerX, centerY, centerZ,
		xHalf, 0, 0,
		0, yHalf, 0,
		0, 0, zHalf,
	}
}

func (c *TileContent) Merge(other *TileContent) {
	if other == nil {
		return
	}

	c.Vertices = append(c.Vertices, other.Vertices...)
	c.Normals = append(c.Normals, other.Normals...)
	c.TexCoords = append(c.TexCoords, other.TexCoords...)
	c.Indices = append(c.Indices, other.Indices...)
	c.Textures = append(c.Textures, other.Textures...)
	c.Primitives = append(c.Primitives, other.Primitives...)

	if c.BatchLength == 0 {
		c.BoundingBox = other.BoundingBox
		c.BatchLength = other.BatchLength
	} else {
		c.BatchLength += other.BatchLength
	}
}

func LoadOSGB(path string) (interface{}, error) {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(path, nil)
	if res == nil {
		return nil, fmt.Errorf("failed to read %s", path)
	}
	return res.GetNode(), nil
}

func (c *GeometryConverter) extractMaterial(geode *model.Geode) []*Material {
	materials := make([]*Material, 0)
	children := geode.GetChildren()
	for _, child := range children {
		if drawable, ok := child.(model.DrawableInterface); ok {
			geom, ok := drawable.(*model.Geometry)
			if !ok {
				continue
			}
			if stateset := geom.GetStates(); stateset != nil {
				mat := c.extractMaterialFromStateSet(stateset)
				if mat != nil {
					materials = append(materials, mat)
				}
			}
		}
	}
	return materials
}

func (c *GeometryConverter) extractMaterialFromStateSet(stateset *model.StateSet) *Material {
	mat := &Material{
		BaseColorFactor: [4]float32{1.0, 1.0, 1.0, 1.0},
		MetallicFactor:  0.0,
		RoughnessFactor: 1.0,
		DoubleSided:     false,
		Unlit:           true,
	}

	if stateset.AttributeList == nil {
		return mat
	}

	for _, attrPair := range stateset.AttributeList {
		if attrPair == nil {
			continue
		}
		if material, ok := attrPair.First.(*model.Material); ok {
			diffuse := material.GetDiffuse(model.GLFRONTANDBACK)
			if diffuse != nil {
				mat.BaseColorFactor = *diffuse
			}
			specular := material.GetSpecular(model.GLFRONTANDBACK)
			if specular != nil && (specular[0] > 0 || specular[1] > 0 || specular[2] > 0) {
				mat.MetallicFactor = 0.9
				mat.RoughnessFactor = 0.1
			}
			emission := material.GetEmission(model.GLFRONTANDBACK)
			if emission != nil && (emission[0] > 0 || emission[1] > 0 || emission[2] > 0) {
				mat.BaseColorFactor[0] = emission[0]
				mat.BaseColorFactor[1] = emission[1]
				mat.BaseColorFactor[2] = emission[2]
			}
		}
	}
	return mat
}

func (c *GeometryConverter) extractTexture(geode *model.Geode, basePath string) []*Texture {
	textures := make([]*Texture, 0)
	children := geode.GetChildren()
	texIndex := 0
	for _, child := range children {
		if drawable, ok := child.(model.DrawableInterface); ok {
			geom, ok := drawable.(*model.Geometry)
			if !ok {
				continue
			}
			if stateset := geom.GetStates(); stateset != nil {
				if tex := c.extractTextureFromStateSet(stateset, basePath, texIndex); tex != nil {
					textures = append(textures, tex)
					texIndex++
				}
			}
		}
	}
	return textures
}

func (c *GeometryConverter) extractTextureFromStateSet(stateset *model.StateSet, basePath string, index int) *Texture {
	if stateset.TextureAttributeList == nil || len(stateset.TextureAttributeList) == 0 {
		return nil
	}

	for unit, attrList := range stateset.TextureAttributeList {
		if unit > 0 {
			break
		}
		for _, attrPair := range attrList {
			if attrPair == nil {
				continue
			}
			if tex, ok := attrPair.First.(*model.Texture); ok {
				if tex.Image != nil {
					texData := c.extractImageData(tex.Image, basePath)
					if texData != nil {
						return &Texture{
							ID:   fmt.Sprintf("texture_%d", index),
							Data: texData,
							Mime: getImageMimeType(tex.Image),
						}
					}
				}
			}
		}
	}
	return nil
}

func (c *GeometryConverter) extractImageData(img *model.Image, basePath string) []byte {
	if img == nil {
		return nil
	}

	if len(img.FileName) > 0 {
		texPath := img.FileName
		if !filepath.IsAbs(texPath) && len(basePath) > 0 {
			texPath = filepath.Join(basePath, texPath)
		}
		data, err := os.ReadFile(texPath)
		if err == nil {
			return data
		}
	}

	if img.Data != nil && len(img.Data) > 0 {
		data := img.Data
		w := int(img.S)
		h := int(img.T)
		if w > 0 && h > 0 {
			rgbaSize := w * h * 4
			rgbSize := w * h * 3
			if len(data) == rgbaSize {
				return encodeRGBAAsJPEG(data, w, h)
			}
			if len(data) == rgbSize {
				return encodeRGBAsJPEG(data, w, h)
			}
		}
		return data
	}

	return nil
}

func encodeRGBAAsJPEG(data []byte, w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 4
			img.Set(x, y, color.RGBA{data[i], data[i+1], data[i+2], data[i+3]})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}

func encodeRGBAsJPEG(data []byte, w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 3
			img.Set(x, y, color.RGBA{data[i], data[i+1], data[i+2], 255})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}

func getImageMimeType(img *model.Image) string {
	return "image/jpeg"
}

func (c *GeometryConverter) extractIndices(prim interface{}) []uint32 {
	if prim == nil {
		return nil
	}

	switch p := prim.(type) {
	case *model.DrawElementsUByte:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, uint32(idx))
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawElementsUShort:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, uint32(idx))
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawElementsUInt:
		indices := make([]uint32, 0)
		if p.Data != nil {
			for _, idx := range p.Data {
				indices = append(indices, idx)
			}
		}
		mode := p.Mode
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	case *model.DrawArrays:
		mode := p.Mode
		count := int(p.Count)
		first := int(p.First)
		indices := make([]uint32, count)
		for i := 0; i < count; i++ {
			indices[i] = uint32(first + i)
		}
		if mode == model.QUADS || mode == model.QUADSTRIP {
			return triangulateQuadLike(indices, mode)
		}
		return indices
	}
	return nil
}

func fixupPrimitiveMode(mode int32, indices []uint32) int32 {
	if mode == model.GLPOINTS && len(indices) >= 3 && len(indices)%3 == 0 {
		return model.GLTRIANGLES
	}
	if mode == model.GLQUADS && len(indices) >= 4 && len(indices)%4 == 0 {
		return model.GLTRIANGLES
	}
	if mode == model.GLQUADSTRIP && len(indices) >= 4 {
		return model.GLTRIANGLES
	}
	return mode
}

func (c *GeometryConverter) extractIndicesWithMode(prim interface{}) (int32, []uint32) {
	if prim == nil {
		return -1, nil
	}

	switch p := prim.(type) {
	case *model.DrawElementsUByte:
		if p.Data == nil {
			return -1, nil
		}
		indices := make([]uint32, len(p.Data))
		for i, idx := range p.Data {
			indices[i] = uint32(idx)
		}
		mode := fixupPrimitiveMode(p.Mode, indices)
		if mode == model.GLQUADS {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADS)
		}
		if mode == model.GLQUADSTRIP {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADSTRIP)
		}
		return mode, indices
	case *model.DrawElementsUShort:
		if p.Data == nil {
			return -1, nil
		}
		indices := make([]uint32, len(p.Data))
		for i, idx := range p.Data {
			indices[i] = uint32(idx)
		}
		mode := fixupPrimitiveMode(p.Mode, indices)
		if mode == model.GLQUADS {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADS)
		}
		if mode == model.GLQUADSTRIP {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADSTRIP)
		}
		return mode, indices
	case *model.DrawElementsUInt:
		if p.Data == nil {
			return -1, nil
		}
		indices := make([]uint32, len(p.Data))
		copy(indices, p.Data)
		mode := fixupPrimitiveMode(p.Mode, indices)
		if mode == model.GLQUADS {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADS)
		}
		if mode == model.GLQUADSTRIP {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADSTRIP)
		}
		return mode, indices
	case *model.DrawArrays:
		count := int(p.Count)
		first := int(p.First)
		indices := make([]uint32, count)
		for i := 0; i < count; i++ {
			indices[i] = uint32(first + i)
		}
		mode := fixupPrimitiveMode(p.Mode, indices)
		if mode == model.GLQUADS {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADS)
		}
		if mode == model.GLQUADSTRIP {
			return model.GLTRIANGLES, triangulateQuadLike(indices, model.GLQUADSTRIP)
		}
		return mode, indices
	}
	return -1, nil
}

func triangulateQuadLike(indices []uint32, mode int32) []uint32 {
	if len(indices) < 4 {
		return indices
	}

	var result []uint32

	if mode == model.QUADS {
		if len(indices)%4 != 0 {
			indices = indices[:len(indices)/4*4]
		}
		quadCount := len(indices) / 4
		result = make([]uint32, 0, quadCount*6)
		for q := 0; q < quadCount; q++ {
			base := q * 4
			result = append(result, indices[base], indices[base+1], indices[base+2])
			result = append(result, indices[base], indices[base+2], indices[base+3])
		}
		return result
	}

	if mode == model.QUADSTRIP {
		if len(indices)%2 != 0 {
			indices = indices[:len(indices)/2*2]
		}
		pairCount := len(indices) / 2
		if pairCount < 2 {
			return indices
		}
		result = make([]uint32, 0, (pairCount-1)*6)
		for i := 0; i < pairCount-1; i++ {
			base := i * 2
			if base+3 >= len(indices) {
				break
			}
			a := indices[base]
			b := indices[base+1]
			c := indices[base+2]
			d := indices[base+3]
			result = append(result, a, b, c, b, d, c)
		}
		return result
	}

	return indices
}

func (c *GeometryConverter) simplifyMesh(vertices, normals, texcoords []float32, indices []uint32) ([]float32, []float32, []float32, []uint32) {
	if c.opts == nil || !c.opts.EnableSimplify || c.opts.SimplifyRatio >= 1.0 {
		return vertices, normals, texcoords, indices
	}

	mesh := &meshopt.Mesh{}
	vertexCount := len(vertices) / 3
	indexCount := len(indices)

	if indexCount == 0 || vertexCount == 0 {
		return vertices, normals, texcoords, indices
	}

	mesh.Indices = make([][3]uint32, indexCount/3)
	for i := 0; i < indexCount/3; i++ {
		mesh.Indices[i] = [3]uint32{indices[i*3], indices[i*3+1], indices[i*3+2]}
	}

	mesh.Nodes = append(mesh.Nodes, meshopt.Node{
		Type: meshopt.AttributeTypePosition,
		Data: c.float32ToBytes(vertices),
	})

	if len(normals) > 0 {
		mesh.Nodes = append(mesh.Nodes, meshopt.Node{
			Type: meshopt.AttributeTypeNormal,
			Data: c.float32ToBytes(normals),
		})
	}

	if len(texcoords) > 0 {
		mesh.Nodes = append(mesh.Nodes, meshopt.Node{
			Type: meshopt.AttributeTypeTexCoord,
			Data: c.float32ToBytes(texcoords),
		})
	}

	settings := &meshopt.Settings{
		SimplifyThreshold:  c.opts.SimplifyRatio,
		SimplifyAggressive: false,
		CompressMore:       false,
	}

	err := meshopt.ProcessMesh(mesh, settings)
	if err != nil {
		return vertices, normals, texcoords, indices
	}

	newIndices := make([]uint32, len(mesh.Indices)*3)
	for i, tri := range mesh.Indices {
		newIndices[i*3] = tri[0]
		newIndices[i*3+1] = tri[1]
		newIndices[i*3+2] = tri[2]
	}

	return vertices, normals, texcoords, newIndices
}

func (c *GeometryConverter) float32ToBytes(data []float32) []byte {
	result := make([]byte, len(data)*4)
	for i, v := range data {
		bits := float32ToBits(v)
		result[i*4] = byte(bits)
		result[i*4+1] = byte(bits >> 8)
		result[i*4+2] = byte(bits >> 16)
		result[i*4+3] = byte(bits >> 24)
	}
	return result
}

func float32ToBits(f float32) uint32 {
	bits := uint32(0)
	copy((*[4]byte)(unsafe.Pointer(&bits))[:], (*[4]byte)(unsafe.Pointer(&f))[:])
	return bits
}

func (c *GeometryConverter) compressMeshDraco(vertices, normals, texcoords []float32, indices []uint32) ([]byte, map[string]int, error) {
	if c.opts == nil || !c.opts.EnableDraco {
		return nil, nil, nil
	}

	meshBuilder := draco.NewMeshBuilder()

	posQuantization := int32(c.opts.DracoPositionBits)
	if posQuantization == 0 {
		posQuantization = 11
	}
	normQuantization := int32(c.opts.DracoNormalBits)
	if normQuantization == 0 {
		normQuantization = 10
	}
	texQuantization := int32(c.opts.DracoTexCoordBits)
	if texQuantization == 0 {
		texQuantization = 12
	}

	vertexCount := len(vertices) / 3
	meshBuilder.Start(vertexCount)

	posData := make([]float32, len(vertices))
	copy(posData, vertices)
	meshBuilder.SetAttribute(vertexCount, posData, draco.GAT_POSITION)

	if len(normals) > 0 {
		normData := make([]float32, len(normals))
		copy(normData, normals)
		meshBuilder.SetAttribute(vertexCount, normData, draco.GAT_NORMAL)
	}

	if len(texcoords) > 0 {
		texData := make([]float32, len(texcoords))
		copy(texData, texcoords)
		meshBuilder.SetAttribute(vertexCount, texData, draco.GAT_TEX_COORD)
	}

	mesh := meshBuilder.GetMesh()

	encoder := draco.NewEncoder()
	encoder.SetAttributeQuantization(draco.GAT_POSITION, posQuantization)
	encoder.SetAttributeQuantization(draco.GAT_NORMAL, normQuantization)
	encoder.SetAttributeQuantization(draco.GAT_TEX_COORD, texQuantization)

	err, compressed := encoder.EncodeMesh(mesh)
	if err != nil {
		return nil, nil, err
	}

	attrIds := map[string]int{
		"POSITION": 0,
	}

	if len(normals) > 0 {
		attrIds["NORMAL"] = 1
	}
	if len(texcoords) > 0 {
		attrIds["TEXCOORD_0"] = 2
	}

	return compressed, attrIds, nil
}

func (c *GeometryConverter) optimizeMesh(vertices, normals, texcoords []float32, indices []uint32) ([]float32, []float32, []float32, []uint32) {
	if c.opts == nil {
		return vertices, normals, texcoords, indices
	}

	if c.opts.EnableSimplify && c.opts.SimplifyRatio < 1.0 {
		return c.simplifyMesh(vertices, normals, texcoords, indices)
	}

	return vertices, normals, texcoords, indices
}

func (c *GeometryConverter) GetDatabasePath(node interface{}) string {
	switch n := node.(type) {
	case *model.PagedLod:
		return n.DataBasePath
	}
	return ""
}

func (c *GeometryConverter) GetFileNames(node interface{}) []string {
	switch n := node.(type) {
	case *model.PagedLod:
		names := make([]string, 0)
		children := n.GetChildren()
		for i := 0; i < len(children); i++ {
			if i < len(n.PerRangeDataList) {
				if fname := n.PerRangeDataList[i].FileName; fname != "" {
					names = append(names, fname)
				}
			}
		}
		return names
	}
	return nil
}

func generateNormals(vertices []float32, indices []uint32) []float32 {
	vertexCount := len(vertices) / 3
	if vertexCount == 0 || len(indices) < 3 {
		return nil
	}

	normals := make([]float32, vertexCount*3)
	triangleCount := len(indices) / 3

	for t := 0; t < triangleCount; t++ {
		i0 := indices[t*3]
		i1 := indices[t*3+1]
		i2 := indices[t*3+2]

		if i0 >= uint32(vertexCount) || i1 >= uint32(vertexCount) || i2 >= uint32(vertexCount) {
			continue
		}

		ax := float64(vertices[i0*3])
		ay := float64(vertices[i0*3+1])
		az := float64(vertices[i0*3+2])
		bx := float64(vertices[i1*3])
		by := float64(vertices[i1*3+1])
		bz := float64(vertices[i1*3+2])
		cx := float64(vertices[i2*3])
		cy := float64(vertices[i2*3+1])
		cz := float64(vertices[i2*3+2])

		ux, uy, uz := bx-ax, by-ay, bz-az
		vx, vy, vz := cx-ax, cy-ay, cz-az

		nx := uy*vz - uz*vy
		ny := uz*vx - ux*vz
		nz := ux*vy - uy*vx

		l := math.Sqrt(nx*nx + ny*ny + nz*nz)
		if l < 1e-20 {
			continue
		}
		invLen := float32(1.0 / l)
		nx_f := float32(nx) * invLen
		ny_f := float32(ny) * invLen
		nz_f := float32(nz) * invLen

		normals[i0*3] += nx_f
		normals[i0*3+1] += ny_f
		normals[i0*3+2] += nz_f
		normals[i1*3] += nx_f
		normals[i1*3+1] += ny_f
		normals[i1*3+2] += nz_f
		normals[i2*3] += nx_f
		normals[i2*3+1] += ny_f
		normals[i2*3+2] += nz_f
	}

	for i := 0; i < vertexCount; i++ {
		nx, ny, nz := normalizeVec3(normals[i*3], normals[i*3+1], normals[i*3+2])
		normals[i*3], normals[i*3+1], normals[i*3+2] = nx, ny, nz
	}

	return normals
}
