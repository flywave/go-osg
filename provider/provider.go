package provider

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/flywave/go-osg"
	"github.com/flywave/go-osg/model"
)

// OsgFileInfo holds metadata about a discovered OSGB file.
type OsgFileInfo struct {
	Path      string
	FileName  string
	SRS       string
	SRSOrigin [3]float64
	BBox      [12]float64
	HasGeoRef bool
}

// TiledProvider discovers OSGB files and provides metadata for tiling.
type TiledProvider struct {
	rootDir  string
	fileList []OsgFileInfo
	srs      string
	origin   [3]float64
	bbox     [12]float64
}

// NewTiledProvider creates a provider that discovers OSGB files under rootDir.
func NewTiledProvider(rootDir string) *TiledProvider {
	p := &TiledProvider{
		rootDir:  rootDir,
		fileList: make([]OsgFileInfo, 0),
	}
	p.scanFiles()
	return p
}

// ComputeBBox individually loads each OSGB file to compute bounding boxes.
// This is expensive for large datasets; call only when needed.
func (p *TiledProvider) ComputeBBox() {
	p.computeGlobalBBox()
}

func (p *TiledProvider) scanFiles() {
	filepath.WalkDir(p.rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".osgb") {
			return nil
		}
		p.fileList = append(p.fileList, OsgFileInfo{Path: path, FileName: d.Name()})
		return nil
	})
}

func (p *TiledProvider) computeGlobalBBox() {
	if len(p.fileList) == 0 {
		return
	}
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	minZ, maxZ := math.MaxFloat64, -math.MaxFloat64
	hasValid := false

	for i := range p.fileList {
		bbox := p.loadFileBBox(&p.fileList[i])
		if bbox == nil {
			continue
		}
		hasValid = true
		p.fileList[i].BBox = *bbox
		cx, cy, cz := bbox[0], bbox[1], bbox[2]
		hx, hy, hz := bbox[3], bbox[7], bbox[11]
		if cx-hx < minX {
			minX = cx - hx
		}
		if cx+hx > maxX {
			maxX = cx + hx
		}
		if cy-hy < minY {
			minY = cy - hy
		}
		if cy+hy > maxY {
			maxY = cy + hy
		}
		if cz-hz < minZ {
			minZ = cz - hz
		}
		if cz+hz > maxZ {
			maxZ = cz + hz
		}
	}

	if !hasValid {
		return
	}

	centerX := (maxX + minX) / 2
	centerY := (maxY + minY) / 2
	centerZ := (maxZ + minZ) / 2
	p.bbox = [12]float64{
		centerX, centerY, centerZ,
		(maxX - minX) / 2, 0, 0,
		0, (maxY - minY) / 2, 0,
		0, 0, (maxZ - minZ) / 2,
	}
}

func (p *TiledProvider) loadFileBBox(info *OsgFileInfo) *[12]float64 {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(info.Path, nil)
	if res == nil {
		return nil
	}
	verts := extractRawVertices(res.GetNode())
	if len(verts) == 0 {
		return nil
	}
	bbox := calculateBoundingBox(verts)
	return &bbox
}

// SetSRSOrigin sets the SRS and origin for all files.
func (p *TiledProvider) SetSRSOrigin(srs string, origin [3]float64) {
	p.srs = srs
	p.origin = origin
	for i := range p.fileList {
		p.fileList[i].SRS = srs
		p.fileList[i].SRSOrigin = origin
		p.fileList[i].HasGeoRef = true
	}
}

// GetFileList returns discovered OSGB files.
func (p *TiledProvider) GetFileList() []OsgFileInfo {
	return p.fileList
}

// GetSRS returns the coordinate reference system string.
func (p *TiledProvider) GetSRS() string {
	return p.srs
}

// GetOrigin returns the SRS origin offset.
func (p *TiledProvider) GetOrigin() [3]float64 {
	return p.origin
}

// GetBBox returns the global bounding box.
func (p *TiledProvider) GetBBox() [12]float64 {
	return p.bbox
}

// GetFileCount returns the number of discovered files.
func (p *TiledProvider) GetFileCount() int {
	return len(p.fileList)
}

// SortFileListByLOD sorts files by LOD level extracted from filename.
func SortFileListByLOD(files []OsgFileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return extractLODLevel(filepath.Base(files[i].Path)) <
			extractLODLevel(filepath.Base(files[j].Path))
	})
}

func extractLODLevel(tileID string) string {
	parts := strings.Split(tileID, "_")
	for _, part := range parts {
		if len(part) >= 2 && part[0] == 'L' {
			return part
		}
	}
	return ""
}

// --- MeshProvider ---

// MeshProvider loads geometry from a single OSGB file.
type MeshProvider struct {
	path     string
	node     interface{}
	geometry *model.Geometry
	verts    []float32
	indices  []uint32
	normals  []float32
	texcoord []float32
	textures [][]byte
}

// NewMeshProvider creates a MeshProvider from a single OSGB file.
func NewMeshProvider(path string) *MeshProvider {
	p := &MeshProvider{path: path}
	p.load()
	return p
}

func (p *MeshProvider) load() {
	rw := osg.NewReadWrite()
	res := rw.ReadNode(p.path, nil)
	if res == nil {
		return
	}
	p.node = res.GetNode()

	var extract func(n interface{})
	extract = func(n interface{}) {
		switch v := n.(type) {
		case *model.Geode:
			for _, c := range v.GetChildren() {
				if g, ok := c.(*model.Geometry); ok {
					p.extractGeometry(g)
				}
			}
		case *model.Group:
			for _, c := range v.GetChildren() {
				extract(c)
			}
		case *model.PagedLod:
			for _, c := range v.GetChildren() {
				extract(c)
			}
		case *model.PositionAttitudeTransform:
			for _, c := range v.GetChildren() {
				extract(c)
			}
		case *model.MatrixTransform:
			for _, c := range v.GetChildren() {
				extract(c)
			}
		}
	}
	extract(p.node)
}

func (p *MeshProvider) extractGeometry(g *model.Geometry) {
	p.geometry = g

	if g.VertexArray != nil && g.VertexArray.Data != nil {
		switch data := g.VertexArray.Data.(type) {
		case [][3]float32:
			p.verts = make([]float32, len(data)*3)
			for i, v := range data {
				p.verts[i*3] = v[0]
				p.verts[i*3+1] = v[1]
				p.verts[i*3+2] = v[2]
			}
		case []float32:
			p.verts = data
		}
	}

	if g.NormalArray != nil && g.NormalArray.Data != nil {
		switch data := g.NormalArray.Data.(type) {
		case [][3]float32:
			p.normals = make([]float32, len(data)*3)
			for i, v := range data {
				nx, ny, nz := normalizeVec3(v[0], v[1], v[2])
				p.normals[i*3] = nx
				p.normals[i*3+1] = ny
				p.normals[i*3+2] = nz
			}
		case []float32:
			p.normals = make([]float32, len(data))
			copy(p.normals, data)
		}
	}

	if len(g.TexCoordArrayList) > 0 && g.TexCoordArrayList[0] != nil && g.TexCoordArrayList[0].Data != nil {
		switch data := g.TexCoordArrayList[0].Data.(type) {
		case [][2]float32:
			p.texcoord = make([]float32, len(data)*2)
			for i, v := range data {
				p.texcoord[i*2] = v[0]
				p.texcoord[i*2+1] = v[1]
			}
		case []float32:
			p.texcoord = data
		}
	}

	if g.Primitives != nil {
		for _, prim := range g.Primitives {
			switch pr := prim.(type) {
			case *model.DrawElementsUInt:
				if pr.Data != nil {
					p.indices = make([]uint32, len(pr.Data))
					copy(p.indices, pr.Data)
				}
			case *model.DrawElementsUShort:
				if pr.Data != nil {
					p.indices = make([]uint32, len(pr.Data))
					for i, idx := range pr.Data {
						p.indices[i] = uint32(idx)
					}
				}
			case *model.DrawElementsUByte:
				if pr.Data != nil {
					p.indices = make([]uint32, len(pr.Data))
					for i, idx := range pr.Data {
						p.indices[i] = uint32(idx)
					}
				}
			case *model.DrawArrays:
				count := int(pr.Count)
				first := int(pr.First)
				p.indices = make([]uint32, count)
				for i := 0; i < count; i++ {
					p.indices[i] = uint32(first + i)
				}
			}
		}
	}

	if len(p.normals) == 0 && len(p.indices) > 0 {
		p.normals = generateNormals(p.verts, p.indices)
	}

	if ss := g.GetStates(); ss != nil {
		for _, attrList := range ss.TextureAttributeList {
			for _, pair := range attrList {
				if pair == nil {
					continue
				}
				if tex, ok := pair.First.(*model.Texture); ok && tex.Image != nil {
					data := extractImageData(tex.Image)
					if data != nil {
						p.textures = append(p.textures, data)
					}
				}
			}
		}
	}
}

// GetVertices returns vertex positions [x0,y0,z0, x1,y1,z1, ...].
func (p *MeshProvider) GetVertices() []float32 {
	return p.verts
}

// GetIndices returns triangle indices [i0,i1,i2, i3,i4,i5, ...].
func (p *MeshProvider) GetIndices() []uint32 {
	return p.indices
}

// GetNormals returns vertex normals.
func (p *MeshProvider) GetNormals() []float32 {
	return p.normals
}

// GetTexCoords returns texture coordinates [u0,v0, u1,v1, ...].
func (p *MeshProvider) GetTexCoords() []float32 {
	return p.texcoord
}

// GetTextures returns embedded texture images (JPEG encoded).
func (p *MeshProvider) GetTextures() [][]byte {
	return p.textures
}

// GetVertexCount returns the number of vertices.
func (p *MeshProvider) GetVertexCount() int {
	if p.verts != nil {
		return len(p.verts) / 3
	}
	return 0
}

// GetTriangleCount returns the number of triangles.
func (p *MeshProvider) GetTriangleCount() int {
	if p.indices != nil {
		return len(p.indices) / 3
	}
	return 0
}

// HasNormals returns true if normals are available.
func (p *MeshProvider) HasNormals() bool {
	return len(p.normals) > 0
}

// HasTexCoords returns true if texture coordinates are available.
func (p *MeshProvider) HasTexCoords() bool {
	return len(p.texcoord) > 0
}

// HasTextures returns true if texture images are available.
func (p *MeshProvider) HasTextures() bool {
	return len(p.textures) > 0
}

// GroupName returns a human-readable summary.
func (p *TiledProvider) GroupName() string {
	return fmt.Sprintf("OSG[%s] %d files", p.srs, len(p.fileList))
}

// --- internal helpers ---

func normalizeVec3(x, y, z float32) (float32, float32, float32) {
	l := math.Sqrt(float64(x)*float64(x) + float64(y)*float64(y) + float64(z)*float64(z))
	if l < 1e-20 {
		return 0.0, 0.0, 1.0
	}
	invLen := float32(1.0 / l)
	return x * invLen, y * invLen, z * invLen
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

func calculateBoundingBox(vertices []float32) [12]float64 {
	if len(vertices) == 0 {
		return [12]float64{0, 0, 0, 0.01, 0, 0, 0, 0.01, 0, 0, 0, 0.01}
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

func extractRawVertices(node interface{}) []float32 {
	var verts []float32
	var visit func(n interface{})
	visit = func(n interface{}) {
		switch v := n.(type) {
		case *model.Geode:
			for _, c := range v.GetChildren() {
				if g, ok := c.(*model.Geometry); ok && g.VertexArray != nil && g.VertexArray.Data != nil {
					switch data := g.VertexArray.Data.(type) {
					case [][3]float32:
						for _, p := range data {
							verts = append(verts, p[0], p[1], p[2])
						}
					case []float32:
						verts = append(verts, data...)
					}
				}
			}
		case *model.Group:
			for _, c := range v.GetChildren() {
				visit(c)
			}
		case *model.PagedLod:
			for _, c := range v.GetChildren() {
				visit(c)
			}
		case *model.PositionAttitudeTransform:
			for _, c := range v.GetChildren() {
				visit(c)
			}
		case *model.MatrixTransform:
			for _, c := range v.GetChildren() {
				visit(c)
			}
		}
	}
	visit(node)
	return verts
}

func extractImageData(img *model.Image) []byte {
	if img == nil || img.Data == nil {
		return nil
	}
	w, h := int(img.S), int(img.T)
	if w <= 0 || h <= 0 {
		return nil
	}
	data := img.Data
	rgbaSize := w * h * 4
	rgbSize := w * h * 3
	if len(data) == rgbaSize {
		return encodeRGBAAsJPEG(data, w, h)
	}
	if len(data) == rgbSize {
		return encodeRGBAsJPEG(data, w, h)
	}
	return data
}

func encodeRGBAAsJPEG(data []byte, w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			srcY := h - 1 - y
			i := (srcY*w + x) * 4
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
			srcY := h - 1 - y
			i := (srcY*w + x) * 3
			img.Set(x, y, color.RGBA{data[i], data[i+1], data[i+2], 255})
		}
	}
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return buf.Bytes()
}
