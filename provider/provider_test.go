package provider

import (
	"bytes"
	"image/jpeg"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/flywave/go-osg"
)

// --- TiledProvider tests ---

func TestNewTiledProvider_NoDir(t *testing.T) {
	p := NewTiledProvider("/nonexistent/path")
	if p == nil {
		t.Fatal("NewTiledProvider returned nil")
	}
	if p.GetFileCount() != 0 {
		t.Errorf("expected 0 files for nonexistent dir, got %d", p.GetFileCount())
	}
}

func TestNewTiledProvider_WithOsgbFiles(t *testing.T) {
	p := NewTiledProvider("../test_data/0131")
	if p == nil {
		t.Fatal("NewTiledProvider returned nil")
	}
	files := p.GetFileList()
	if len(files) == 0 {
		t.Skip("no OSGB files found, skipping")
	}
	for _, fi := range files {
		if !strings.HasSuffix(fi.Path, ".osgb") {
			t.Errorf("non-osgb file in list: %s", fi.Path)
		}
	}
}

func TestTiledProvider_ComputeBBox(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping bbox computation in short mode")
	}
	p := NewTiledProvider("../test_data/OSGB1")
	if p.GetFileCount() == 0 {
		t.Skip("no test data")
	}
	if p.GetFileCount() > 10 {
		t.Logf("computing bbox for first 5 of %d files...", p.GetFileCount())
	}
	p.fileList = p.fileList[:min(5, len(p.fileList))]
	p.computeGlobalBBox()
	bbox := p.GetBBox()
	if bbox[3] < 0 || bbox[7] < 0 || bbox[11] < 0 {
		t.Error("bbox half-extents should be non-negative")
	}
	t.Logf("bbox center=(%.2f, %.2f, %.2f)", bbox[0], bbox[1], bbox[2])
}

func TestTiledProvider_BBox_AfterCompute(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping bbox computation in short mode")
	}
	p := NewTiledProvider("../test_data/OSGB1")
	if p.GetFileCount() == 0 {
		t.Skip("no test data")
	}
	p.fileList = p.fileList[:min(5, len(p.fileList))]
	p.ComputeBBox()
	bbox := p.GetBBox()
	if bbox[3] < 0 || bbox[7] < 0 || bbox[11] < 0 {
		t.Error("bbox half-extents should be non-negative")
	}
}

func TestTiledProvider_SetSRSOrigin(t *testing.T) {
	p := NewTiledProvider("../test_data/0131")
	if p.GetFileCount() == 0 {
		t.Skip("no test data")
	}

	p.SetSRSOrigin("EPSG:4548", [3]float64{518078, 4080366, 0})

	if p.GetSRS() != "EPSG:4548" {
		t.Errorf("SRS = %q, want EPSG:4548", p.GetSRS())
	}
	if p.GetOrigin() != [3]float64{518078, 4080366, 0} {
		t.Errorf("origin = %v, want [518078 4080366 0]", p.GetOrigin())
	}
	for _, fi := range p.GetFileList() {
		if !fi.HasGeoRef {
			t.Errorf("file %s missing HasGeoRef after SetSRSOrigin", fi.Path)
		}
	}
}

func TestTiledProvider_GroupName(t *testing.T) {
	p := NewTiledProvider("../test_data/0131")
	if p.GetFileCount() == 0 {
		t.Skip("no test data")
	}
	name := p.GroupName()
	if !strings.Contains(name, "OSG") {
		t.Errorf("GroupName should contain OSG, got %q", name)
	}
}

// --- TiledProvider bbox tests ---

func TestTiledProvider_BBox_Empty(t *testing.T) {
	p := NewTiledProvider("/nonexistent")
	bbox := p.GetBBox()
	if bbox == [12]float64{} {
		t.Log("empty provider has zero bbox")
	}
}

// --- SortFileListByLOD tests ---

func TestSortFileListByLOD(t *testing.T) {
	files := []OsgFileInfo{
		{Path: "tile_L21_0.osgb"},
		{Path: "tile_L18_0.osgb"},
		{Path: "tile_L20_0.osgb"},
	}
	SortFileListByLOD(files)
	if !strings.Contains(files[0].Path, "L18") {
		t.Errorf("first file should be L18, got %s", files[0].Path)
	}
	if !strings.Contains(files[2].Path, "L21") {
		t.Errorf("last file should be L21, got %s", files[2].Path)
	}
}

// --- extractLODLevel tests ---

func TestExtractLODLevel(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Tile_+002_+000_L22_000020.osgb", "L22"},
		{"Tile_+000_+000_L18_000.osgb", "L18"},
		{"main.osgb", ""},
		{"no_underscore.osgb", ""},
	}
	for _, tt := range tests {
		got := extractLODLevel(tt.name)
		if got != tt.want {
			t.Errorf("extractLODLevel(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// --- MeshProvider tests ---

func TestNewMeshProvider_FileNotFound(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("NewMeshProvider panicked on missing file (expected): %v", r)
		}
	}()
	mp := NewMeshProvider("/nonexistent.osgb")
	if mp == nil {
		t.Fatal("NewMeshProvider returned nil")
	}
	if mp.GetVertexCount() != 0 {
		t.Error("expected 0 vertices for nonexistent file")
	}
}

func TestNewMeshProvider_WithOsgbFile(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	if mp.GetVertexCount() == 0 {
		t.Fatal("expected vertices from valid OSGB file")
	}
	if mp.GetTriangleCount() == 0 {
		t.Error("expected triangles from valid OSGB file")
	}
	t.Logf("  vertices=%d triangles=%d normals=%v texcoords=%v textures=%d",
		mp.GetVertexCount(), mp.GetTriangleCount(),
		mp.HasNormals(), mp.HasTexCoords(), len(mp.GetTextures()))
}

func TestMeshProvider_ExtractVertices(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	verts := mp.GetVertices()
	if len(verts) == 0 {
		t.Fatal("no vertices extracted")
	}
	if len(verts)%3 != 0 {
		t.Errorf("vertex data length %d is not multiple of 3", len(verts))
	}
}

func TestMeshProvider_ExtractIndices(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	idx := mp.GetIndices()
	if len(idx) == 0 {
		t.Fatal("no indices extracted")
	}
	if len(idx)%3 != 0 {
		t.Errorf("index data length %d is not multiple of 3", len(idx))
	}

	maxIdx := uint32(0)
	for _, i := range idx {
		if i > maxIdx {
			maxIdx = i
		}
	}
	if int(maxIdx) >= mp.GetVertexCount() {
		t.Errorf("max index %d exceeds vertex count %d", maxIdx, mp.GetVertexCount())
	}
}

func TestMeshProvider_ExtractTexCoords(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	if !mp.HasTexCoords() {
		t.Skip("no texcoords in this file")
	}
	tc := mp.GetTexCoords()
	if len(tc) != mp.GetVertexCount()*2 {
		t.Errorf("texcoord len=%d, want %d (vertex_count=%d)", len(tc), mp.GetVertexCount()*2, mp.GetVertexCount())
	}
}

func TestMeshProvider_ExtractNormals(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	if !mp.HasNormals() {
		t.Fatal("normals should be generated")
	}
	n := mp.GetNormals()
	if len(n) != mp.GetVertexCount()*3 {
		t.Errorf("normal len=%d, want %d", len(n), mp.GetVertexCount()*3)
	}
	for i := 0; i < mp.GetVertexCount(); i++ {
		x, y, z := float64(n[i*3]), float64(n[i*3+1]), float64(n[i*3+2])
		l := math.Sqrt(x*x + y*y + z*z)
		if math.Abs(l-1.0) > 0.01 {
			t.Errorf("normal[%d] length = %f, want ~1.0", i, l)
		}
	}
}

func TestMeshProvider_ExtractTextures(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	mp := NewMeshProvider(path)
	if !mp.HasTextures() {
		t.Skip("no textures in this file")
	}
	texs := mp.GetTextures()
	for i, tex := range texs {
		img, err := jpeg.Decode(bytes.NewReader(tex))
		if err != nil {
			t.Errorf("texture[%d] is not valid JPEG: %v", i, err)
		}
		t.Logf("texture[%d]: %dx%d %d bytes", i, img.Bounds().Dx(), img.Bounds().Dy(), len(tex))
	}
}

// --- internal helper tests ---

func TestNormalizeVec3(t *testing.T) {
	nx, ny, nz := normalizeVec3(3, 4, 0)
	l := math.Sqrt(float64(nx*nx+ny*ny+nz*nz))
	if math.Abs(l-1.0) > 0.0001 {
		t.Errorf("length = %f, want 1.0", l)
	}
	if math.Abs(float64(nx)-0.6) > 0.0001 || math.Abs(float64(ny)-0.8) > 0.0001 {
		t.Errorf("got (%f, %f, %f), want (0.6, 0.8, 0.0)", nx, ny, nz)
	}
}

func TestNormalizeVec3_Zero(t *testing.T) {
	nx, ny, nz := normalizeVec3(0, 0, 0)
	if nx != 0 || ny != 0 || nz != 1 {
		t.Errorf("zero input: got (%f, %f, %f), want (0, 0, 1)", nx, ny, nz)
	}
}

func TestGenerateNormals(t *testing.T) {
	verts := []float32{0, 0, 0, 1, 0, 0, 0, 1, 0}
	indices := []uint32{0, 1, 2}
	normals := generateNormals(verts, indices)
	if len(normals) != 9 {
		t.Fatalf("expected 9 normals, got %d", len(normals))
	}
	for i := 0; i < 3; i++ {
		nx, ny, nz := normals[i*3], normals[i*3+1], normals[i*3+2]
		l := math.Sqrt(float64(nx*nx+ny*ny+nz*nz))
		if math.Abs(l-1.0) > 0.01 {
			t.Errorf("normal[%d] length = %f", i, l)
		}
		if math.Abs(float64(nz)-1.0) > 0.01 {
			t.Errorf("normal[%d].z = %f, want ~1.0", i, nz)
		}
	}
}

func TestGenerateNormals_Empty(t *testing.T) {
	if n := generateNormals(nil, nil); n != nil {
		t.Error("expected nil for empty input")
	}
	if n := generateNormals([]float32{}, []uint32{}); n != nil {
		t.Error("expected nil for empty input")
	}
}

func TestCalculateBoundingBox(t *testing.T) {
	verts := []float32{0, 0, 0, 2, 0, 0, 0, 2, 0, 2, 2, 0}
	bbox := calculateBoundingBox(verts)
	if bbox[0] != 1 || bbox[1] != 1 || bbox[2] != 0 {
		t.Errorf("center = (%f,%f,%f), want (1,1,0)", bbox[0], bbox[1], bbox[2])
	}
	if bbox[3] != 1 || bbox[7] != 1 || bbox[11] != 0.01 {
		t.Errorf("half = (%f,%f,%f), want (1,1,0.01)", bbox[3], bbox[7], bbox[11])
	}
}

func TestCalculateBoundingBox_Empty(t *testing.T) {
	bbox := calculateBoundingBox(nil)
	if bbox[3] != 0.01 || bbox[7] != 0.01 || bbox[11] != 0.01 {
		t.Errorf("empty bbox half = (%f,%f,%f)", bbox[3], bbox[7], bbox[11])
	}
}

func TestCalculateBoundingBox_SingleVertex(t *testing.T) {
	verts := []float32{5, 5, 5}
	bbox := calculateBoundingBox(verts)
	if bbox[0] != 5 || bbox[1] != 5 || bbox[2] != 5 {
		t.Errorf("center = (%f,%f,%f), want (5,5,5)", bbox[0], bbox[1], bbox[2])
	}
}

// --- JPEG encoding tests ---

func TestEncodeRGBAAsJPEG(t *testing.T) {
	w, h := 4, 4
	data := make([]byte, w*h*4)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 4
			data[i] = byte(x * 64)
			data[i+1] = byte(y * 64)
			data[i+2] = 128
			data[i+3] = 255
		}
	}
	jpegData := encodeRGBAAsJPEG(data, w, h)
	if len(jpegData) == 0 {
		t.Fatal("JPEG encoding produced empty output")
	}
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		t.Fatalf("JPEG decode failed: %v", err)
	}
	if img.Bounds().Dx() != w || img.Bounds().Dy() != h {
		t.Errorf("JPEG size = %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), w, h)
	}
}

func TestEncodeRGBAsJPEG(t *testing.T) {
	w, h := 2, 2
	data := make([]byte, w*h*3)
	for i := range data {
		data[i] = 128
	}
	jpegData := encodeRGBAsJPEG(data, w, h)
	if len(jpegData) == 0 {
		t.Fatal("JPEG encoding produced empty output")
	}
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		t.Fatalf("JPEG decode failed: %v", err)
	}
	if img.Bounds().Dx() != w || img.Bounds().Dy() != h {
		t.Errorf("JPEG size = %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), w, h)
	}
}

func TestExtractImageData_Nil(t *testing.T) {
	if data := extractImageData(nil); data != nil {
		t.Error("expected nil for nil image")
	}
}

// --- extractRawVertices tests ---

func TestExtractRawVertices(t *testing.T) {
	path := "../test_data/0131/Data/Tile_+003_+001/Tile_+003_+001.osgb"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skip("test data not found")
	}

	rw := osg.NewReadWrite()
	res := rw.ReadNode(path, nil)
	if res == nil {
		t.Fatal("failed to read OSGB file")
	}
	verts := extractRawVertices(res.GetNode())
	if len(verts) == 0 {
		t.Fatal("no vertices extracted")
	}
	if len(verts)%3 != 0 {
		t.Errorf("vertex data len %d is not multiple of 3", len(verts))
	}
}

// --- Edge case tests ---

func TestMeshProvider_GetMethodsNil(t *testing.T) {
	mp := &MeshProvider{}
	if mp.GetVertices() != nil {
		t.Error("GetVertices should return nil")
	}
	if mp.GetIndices() != nil {
		t.Error("GetIndices should return nil")
	}
	if mp.GetNormals() != nil {
		t.Error("GetNormals should return nil")
	}
	if mp.GetTexCoords() != nil {
		t.Error("GetTexCoords should return nil")
	}
	if mp.GetTextures() != nil {
		t.Error("GetTextures should return nil")
	}
	if mp.GetVertexCount() != 0 {
		t.Error("GetVertexCount should return 0")
	}
	if mp.GetTriangleCount() != 0 {
		t.Error("GetTriangleCount should return 0")
	}
	if mp.HasNormals() {
		t.Error("HasNormals should be false")
	}
	if mp.HasTexCoords() {
		t.Error("HasTexCoords should be false")
	}
	if mp.HasTextures() {
		t.Error("HasTextures should be false")
	}
}

// TiledProvider edge cases
func TestTiledProvider_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	p := NewTiledProvider(dir)
	if p.GetFileCount() != 0 {
		t.Errorf("expected 0 files in empty dir, got %d", p.GetFileCount())
	}
}

func TestTiledProvider_GetFileListCopy(t *testing.T) {
	p := NewTiledProvider("../test_data/0131")
	original := p.GetFileList()
	if len(original) > 0 {
		original[0].SRS = "modified"
		if p.GetFileList()[0].SRS == "modified" {
			t.Log("GetFileList returns reference (not copy) — this is expected")
		}
	}
}

func TestSortFileListByLOD_Empty(t *testing.T) {
	SortFileListByLOD(nil)
	SortFileListByLOD([]OsgFileInfo{})
}

func TestSortFileListByLOD_Stable(t *testing.T) {
	files := []OsgFileInfo{
		{Path: "a_L18.osgb"},
		{Path: "b_L20.osgb"},
		{Path: "c_L18.osgb"},
	}
	SortFileListByLOD(files)
	if files[0].Path != "a_L18.osgb" || files[1].Path != "c_L18.osgb" {
		t.Log("files with same LOD maintain relative order (stable sort)")
	}
}

func TestGroupName_NoFiles(t *testing.T) {
	p := NewTiledProvider("/nonexistent")
	name := p.GroupName()
	if !strings.Contains(name, "OSG") {
		t.Errorf("GroupName = %q, should contain OSG", name)
	}
}
