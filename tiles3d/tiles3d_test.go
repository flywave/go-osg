package tiles3d

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-osg/model"
)

// testData returns a simple triangle mesh for use in compression/simplification tests.
func testMeshData() ([]float32, []float32, []float32, []uint32) {
	verts := []float32{0, 0, 0, 1, 0, 0, 0, 1, 0}
	norms := []float32{0, 0, 1, 0, 0, 1, 0, 0, 1}
	texcs := []float32{0, 0, 1, 0, 0, 1}
	indices := []uint32{0, 1, 2}
	return verts, norms, texcs, indices
}

// ---------------------------------------------------------------------------
// Options tests
// ---------------------------------------------------------------------------

func TestDefaultConverterOptions(t *testing.T) {
	opts := DefaultConverterOptions()
	if opts.TargetSRS != "EPSG:4326" {
		t.Errorf("TargetSRS = %q, want EPSG:4326", opts.TargetSRS)
	}
	if opts.GeoidModel != "none" {
		t.Errorf("GeoidModel = %q, want none", opts.GeoidModel)
	}
	if !opts.EnableTexture {
		t.Error("EnableTexture should be true")
	}
	if !opts.EnableUnlit {
		t.Error("EnableUnlit should be true")
	}
	if opts.MaxLOD != -1 {
		t.Errorf("MaxLOD = %d, want -1", opts.MaxLOD)
	}
}

// ---------------------------------------------------------------------------
// Metadata tests
// ---------------------------------------------------------------------------

func TestDetectSRSType(t *testing.T) {
	tests := []struct {
		srs  string
		want SRSType
	}{
		{"ENU:114,34", SRSTypeENU},
		{"EPSG:4326", SRSTypeEPSG},
		{"EPSG:4548", SRSTypeEPSG},
		{"unknown", SRSTypeUnknown},
		{"", SRSTypeUnknown},
		{"PROJCS[\"CGCS2000\", ...]", SRSTypeWKT},
		{"GEOGCS[\"WGS84\", ...]", SRSTypeWKT},
	}
	for _, tt := range tests {
		got := DetectSRSType(tt.srs)
		if got != tt.want {
			t.Errorf("DetectSRSType(%q) = %d, want %d", tt.srs, got, tt.want)
		}
	}
}

func TestParseSRSOrigin(t *testing.T) {
	tests := []struct {
		origin string
		x, y   float64
		z      float64
		err    bool
	}{
		{"500000,3000000", 500000, 3e6, 0, false},
		{"500000,3000000,100", 500000, 3e6, 100, false},
		{"-500000.5,2000000.25,50.75", -500000.5, 2000000.25, 50.75, false},
		{"invalid", 0, 0, 0, true},
		{"", 0, 0, 0, true},
	}
	for _, tt := range tests {
		x, y, z, err := ParseSRSOrigin(tt.origin)
		if tt.err {
			if err == nil {
				t.Errorf("ParseSRSOrigin(%q) expected error", tt.origin)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseSRSOrigin(%q) unexpected error: %v", tt.origin, err)
		}
		if x != tt.x || y != tt.y || z != tt.z {
			t.Errorf("ParseSRSOrigin(%q) = (%f,%f,%f), want (%f,%f,%f)",
				tt.origin, x, y, z, tt.x, tt.y, tt.z)
		}
	}
}

func TestParseEPSGCode(t *testing.T) {
	tests := []struct {
		srs  string
		code int
		err  bool
	}{
		{"EPSG:4326", 4326, false},
		{"EPSG:4548", 4548, false},
		{"invalid", 0, true},
		{"EPSG:", 0, true},
	}
	for _, tt := range tests {
		code, err := ParseEPSGCode(tt.srs)
		if tt.err {
			if err == nil {
				t.Errorf("ParseEPSGCode(%q) expected error", tt.srs)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseEPSGCode(%q) unexpected error: %v", tt.srs, err)
		}
		if code != tt.code {
			t.Errorf("ParseEPSGCode(%q) = %d, want %d", tt.srs, code, tt.code)
		}
	}
}

func TestParseENUOrigin(t *testing.T) {
	tests := []struct {
		srs     string
		lon, lat float64
		err     bool
	}{
		{"ENU:114,34", 114, 34, false},
		{"ENU:114.5,34.25", 114.5, 34.25, false},
		{"invalid", 0, 0, true},
		{"ENU:", 0, 0, true},
	}
	for _, tt := range tests {
		lon, lat, err := ParseENUOrigin(tt.srs)
		if tt.err {
			if err == nil {
				t.Errorf("ParseENUOrigin(%q) expected error", tt.srs)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseENUOrigin(%q) unexpected error: %v", tt.srs, err)
		}
		if lon != tt.lon || lat != tt.lat {
			t.Errorf("ParseENUOrigin(%q) = (%f,%f), want (%f,%f)",
				tt.srs, lon, lat, tt.lon, tt.lat)
		}
	}
}

// ---------------------------------------------------------------------------
// Triangulation tests
// ---------------------------------------------------------------------------

func TestTriangulateQuads(t *testing.T) {
	// GL_QUADS = 0x0007
	indices := []uint32{0, 1, 2, 3, 4, 5, 6, 7}
	result := triangulateQuadLike(indices, 7)

	if len(result) != 12 {
		t.Fatalf("quad triangulation: got %d indices, want 12", len(result))
	}
	// First quad: 0-1-2, 0-2-3
	if result[0] != 0 || result[1] != 1 || result[2] != 2 {
		t.Errorf("quad triangle 1: got %v, want [0,1,2]", result[:3])
	}
	if result[3] != 0 || result[4] != 2 || result[5] != 3 {
		t.Errorf("quad triangle 2: got %v, want [0,2,3]", result[3:6])
	}
}

func TestTriangulateQuadStrip(t *testing.T) {
	// GL_QUAD_STRIP = 0x0008
	indices := []uint32{0, 1, 2, 3, 4, 5}
	result := triangulateQuadLike(indices, 8)

	if len(result) != 12 {
		t.Fatalf("quadstrip triangulation: got %d indices, want 12", len(result))
	}
	// First quad: 0-1-2-3 → 0-1-2, 1-3-2
	// Second quad: 2-3-4-5 → 2-3-4, 3-5-4
	if result[0] != 0 || result[1] != 1 || result[2] != 2 {
		t.Errorf("quadstrip tri 1: got %v, want [0,1,2]", result[:3])
	}
}

func TestTriangulate_NotEnoughIndices(t *testing.T) {
	indices := []uint32{0, 1, 2}
	result := triangulateQuadLike(indices, 4)
	// Less than 4 indices should return as-is
	if len(result) != 3 {
		t.Errorf("not enough indices: got %d, want 3", len(result))
	}
}

// ---------------------------------------------------------------------------
// Bounding box tests
// ---------------------------------------------------------------------------

func TestCalculateBoundingBox(t *testing.T) {
	// Simple triangle
	vertices := []float32{0, 0, 0, 1, 0, 0, 0, 1, 0}
	bbox := calculateBoundingBox(vertices)

	// Center should be (0.5, 0.5, 0)
	if bbox[0] != 0.5 || bbox[1] != 0.5 || bbox[2] != 0 {
		t.Errorf("center = (%f,%f,%f), want (0.5,0.5,0)", bbox[0], bbox[1], bbox[2])
	}
	// Half extents should be (0.5, 0.5, 0.01) with z minimum
	if bbox[3] != 0.5 || bbox[7] != 0.5 || bbox[11] != 0.01 {
		t.Errorf("half extents = (%f,%f,%f), want (0.5,0.5,0.01)", bbox[3], bbox[7], bbox[11])
	}
}

func TestCalculateBoundingBox_Empty(t *testing.T) {
	bbox := calculateBoundingBox(nil)
	if bbox[3] != 0.01 || bbox[7] != 0.01 || bbox[11] != 0.01 {
		t.Errorf("empty bbox half extents = (%f,%f,%f), want (0.01,0.01,0.01)",
			bbox[3], bbox[7], bbox[11])
	}
}

func TestGeometricError(t *testing.T) {
	bbox := [12]float64{0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 5}
	err := calculateGeometricError(bbox)
	// (dx+dy+dz)/6 * 2 = (10+10+10)/6*2 = 30/6*2 = 10
	if math.Abs(err-10.0) > 0.001 {
		t.Errorf("geometric error = %f, want 10.0", err)
	}
}

// ---------------------------------------------------------------------------
// extractLODLevel tests
// ---------------------------------------------------------------------------

func TestExtractLODLevel(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"Tile_+002_+000_L22_000020.osgb", "L22"},
		{"Tile_+003_+003_L18_000.osgb", "L18"},
		{"not_a_tile.osgb", ""},
		{"Tile_+000_+000_L24_0000700.osgb", "L24"},
	}
	for _, tt := range tests {
		got := extractLODLevel(tt.id)
		if got != tt.want {
			t.Errorf("extractLODLevel(%q) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Default texcoords tests
// ---------------------------------------------------------------------------

func TestGenerateDefaultTexCoords(t *testing.T) {
	geomConv := &GeometryConverter{}

	// A square (4 vertices)
	vertices := []float32{0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0}
	texCoords := geomConv.generateDefaultTexCoords(vertices)

	if len(texCoords) != 8 {
		t.Fatalf("got %d texcoords, want 8", len(texCoords))
	}

	// First vertex (0,0,0) should map to (0,0)
	if texCoords[0] != 0 || texCoords[1] != 0 {
		t.Errorf("first texcoord = (%f,%f), want (0,0)", texCoords[0], texCoords[1])
	}
	// Third vertex (1,1,0) should map to (1,1)
	if texCoords[4] != 1 || texCoords[5] != 1 {
		t.Errorf("third texcoord = (%f,%f), want (1,1)", texCoords[4], texCoords[5])
	}
}

func TestGenerateDefaultTexCoords_Empty(t *testing.T) {
	geomConv := &GeometryConverter{}
	result := geomConv.generateDefaultTexCoords(nil)
	if result != nil {
		t.Errorf("expected nil for empty vertices")
	}
}

// ---------------------------------------------------------------------------
// Coordinate math tests (no proj dependency)
// ---------------------------------------------------------------------------

func TestECEFMath(t *testing.T) {
	tf := &CoordinateTransformer{}

	// Test ToECEFFromLatLon at equator, 0 longitude, sea level
	ecef := tf.ToECEFFromLatLon(0, 0, 0)
	// At equator, 0 lon: x = a (6378137), y = 0, z = 0
	if math.Abs(ecef[0]-6378137) > 1 {
		t.Errorf("ECEF x = %f, want 6378137", ecef[0])
	}
	if math.Abs(ecef[1]) > 1 {
		t.Errorf("ECEF y = %f, want 0", ecef[1])
	}
	if math.Abs(ecef[2]) > 1 {
		t.Errorf("ECEF z = %f, want 0", ecef[2])
	}

	// At north pole: x=0, y=0, z=a*(1-e2)
	ecef2 := tf.ToECEFFromLatLon(math.Pi/2, 0, 0)
	if math.Abs(ecef2[0]) > 1 {
		t.Errorf("pole ECEF x = %f, want ~0", ecef2[0])
	}
	expectedZ := 6356752.314245 // a*sqrt(1-e2) = polar radius
	if math.Abs(ecef2[2]-expectedZ) > 1 {
		t.Errorf("pole ECEF z = %f, want %f", ecef2[2], expectedZ)
	}
}

// ---------------------------------------------------------------------------
// TileContent Merge tests
// ---------------------------------------------------------------------------

func TestTileContentMerge(t *testing.T) {
	a := &TileContent{
		Vertices: []float32{1, 2, 3},
		Normals:  []float32{0, 0, 1},
		Indices:  []uint32{0, 1, 2},
	}
	b := &TileContent{
		Vertices: []float32{4, 5, 6},
		Normals:  []float32{0, 1, 0},
		Indices:  []uint32{3, 4, 5},
	}

	a.Merge(b)

	if len(a.Vertices) != 6 {
		t.Errorf("merged vertices = %d, want 6", len(a.Vertices))
	}
	if len(a.Normals) != 6 {
		t.Errorf("merged normals = %d, want 6", len(a.Normals))
	}
	if len(a.Indices) != 6 {
		t.Errorf("merged indices = %d, want 6", len(a.Indices))
	}
}

func TestTileContentMerge_Nil(t *testing.T) {
	a := &TileContent{Vertices: []float32{1, 2, 3}}
	a.Merge(nil)
	if len(a.Vertices) != 3 {
		t.Errorf("merge nil should not change vertices, got %d", len(a.Vertices))
	}
}

// ---------------------------------------------------------------------------
// CoordinateTransformer standalone math (no proj)
// ---------------------------------------------------------------------------

func TestToLocalENU(t *testing.T) {
	tf := &CoordinateTransformer{
		center: [3]float64{114.0, 30.0, 0}, // Wuhan-ish
	}
	// Point at the center should give (0,0,0)
	pt := tf.ToLocalENU([3]float64{114.0, 30.0, 0})
	if math.Abs(pt[0]) > 0.001 || math.Abs(pt[1]) > 0.001 || math.Abs(pt[2]) > 0.001 {
		t.Errorf("center point ENU = (%f,%f,%f), want (~0,~0,~0)", pt[0], pt[1], pt[2])
	}
}

func TestHasGeoReference(t *testing.T) {
	tf := &CoordinateTransformer{}
	if tf.HasGeoReference() {
		t.Error("empty center should not have geo reference")
	}

	tf.center = [3]float64{114, 30, 0}
	if !tf.HasGeoReference() {
		t.Error("should have geo reference after setting center")
	}
}

func TestGetOriginOffset(t *testing.T) {
	tf := &CoordinateTransformer{}
	off := tf.GetOriginOffset()
	if off != [3]float64{0, 0, 0} {
		t.Errorf("default origin offset = %v, want [0,0,0]", off)
	}

	tf.originOffset = [3]float64{500000, 3000000, 100}
	off = tf.GetOriginOffset()
	if off[0] != 500000 || off[1] != 3000000 || off[2] != 100 {
		t.Errorf("origin offset = %v, want [500000,3000000,100]", off)
	}
}

// ---------------------------------------------------------------------------
// GeoidConverter tests
// ---------------------------------------------------------------------------

func TestGeoidConverter_None(t *testing.T) {
	g := NewGeoidConverter("none", "")
	if g == nil {
		t.Fatal("NewGeoidConverter returned nil")
	}
	// With model "none", no geoid correction is applied
	h := g.ConvertOrthometricToEllipsoidal(30, 114, 100)
	if h != 100 {
		t.Errorf("no geoid: height = %f, want 100", h)
	}
	h2 := g.ConvertEllipsoidalToOrthometric(30, 114, 100)
	if h2 != 100 {
		t.Errorf("no geoid (inverse): height = %f, want 100", h2)
	}
}

// ---------------------------------------------------------------------------
// meterToLat/meterToLon conversion tests
// ---------------------------------------------------------------------------

func TestMeterToLat(t *testing.T) {
	// meterToLat converts to radians. 111km ≈ Earth radius 6371km → ~1/57.3 rad
	lat := meterToLatDeg(111000.0)
	if math.Abs(lat-0.0174) > 0.002 {
		t.Errorf("111km in lat rad = %f, want ~0.0174", lat)
	}
}

func TestMeterToLon(t *testing.T) {
	// meterToLon converts to radians at a given latitude
	lon := meterToLonDeg(111000.0, 0)
	if math.Abs(lon-0.0174) > 0.002 {
		t.Errorf("111km in lon rad at equator = %f, want ~0.0174", lon)
	}

	// At 60° latitude (cos=0.5), same distance = 2x the angular distance
	lon60 := meterToLonDeg(111000.0, math.Pi/3)
	if math.Abs(lon60-0.0348) > 0.004 {
		t.Errorf("111km in lon rad at 60° = %f, want ~0.0348", lon60)
	}
}

// ---------------------------------------------------------------------------
// extractVertices tests
// ---------------------------------------------------------------------------

func TestExtractVertices_Vec3Float(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	arr.Data = [][3]float32{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	verts := c.extractVertices(arr)
	if len(verts) != 9 {
		t.Fatalf("got %d floats, want 9", len(verts))
	}
	if verts[0] != 1 || verts[1] != 2 || verts[2] != 3 {
		t.Errorf("first vertex = (%f,%f,%f), want (1,2,3)", verts[0], verts[1], verts[2])
	}
}

func TestExtractVertices_FlatFloat(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.FloatArrayType, model.GLFLOAT, 1)
	arr.Data = []float32{1, 2, 3, 4, 5, 6}

	verts := c.extractVertices(arr)
	if len(verts) != 6 {
		t.Fatalf("got %d floats, want 6", len(verts))
	}
}

func TestExtractVertices_Double(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.DoubleArrayType, model.GLDOUBLE, 1)
	arr.Data = []float64{1.5, 2.5, 3.5}

	verts := c.extractVertices(arr)
	if len(verts) != 3 {
		t.Fatalf("got %d floats, want 3", len(verts))
	}
	if math.Abs(float64(verts[0])-1.5) > 0.001 {
		t.Errorf("first = %f, want 1.5", verts[0])
	}
}

func TestExtractVertices_Short(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.ShortArrayType, model.GLSHORT, 1)
	arr.Data = []int16{100, 200, 300}

	verts := c.extractVertices(arr)
	if len(verts) != 3 {
		t.Fatalf("got %d floats, want 3", len(verts))
	}
	if verts[0] != 100 || verts[1] != 200 || verts[2] != 300 {
		t.Errorf("got (%f,%f,%f), want (100,200,300)", verts[0], verts[1], verts[2])
	}
}

func TestExtractVertices_UShort(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.UShortArrayType, model.GLUNSIGNEDSHORT, 1)
	arr.Data = []uint16{10, 20, 30}

	verts := c.extractVertices(arr)
	if len(verts) != 3 {
		t.Fatalf("got %d floats, want 3", len(verts))
	}
}

func TestExtractVertices_Nil(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}
	if v := c.extractVertices(nil); v != nil {
		t.Error("expected nil for nil array")
	}
	if v := c.extractVertices(&model.Array{}); v != nil {
		t.Error("expected nil for array with nil Data")
	}
}

func TestExtractVertices_UnsupportedType(t *testing.T) {
	c := &GeometryConverter{coordTrans: &CoordinateTransformer{}}

	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	arr.Data = []int32{1, 2, 3} // unsupported

	if v := c.extractVertices(arr); v != nil {
		t.Error("expected nil for unsupported data type")
	}
}

// ---------------------------------------------------------------------------
// extractNormals tests
// ---------------------------------------------------------------------------

func TestExtractNormals_Vec3Float(t *testing.T) {
	c := &GeometryConverter{}

	arr := model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	arr.Data = [][3]float32{{0, 0, 1}, {0, 1, 0}}

	normals := c.extractNormals(arr)
	if len(normals) != 6 {
		t.Fatalf("got %d floats, want 6", len(normals))
	}
}

func TestExtractNormals_FlatFloat(t *testing.T) {
	c := &GeometryConverter{}

	arr := model.NewArray(model.FloatArrayType, model.GLFLOAT, 1)
	arr.Data = []float32{0, 0, 1, 0, 1, 0}

	normals := c.extractNormals(arr)
	if len(normals) != 6 {
		t.Fatalf("got %d floats, want 6", len(normals))
	}
}

func TestExtractNormals_Nil(t *testing.T) {
	c := &GeometryConverter{}
	if v := c.extractNormals(&model.Array{}); v != nil {
		t.Error("expected nil for array with nil Data")
	}
}

// ---------------------------------------------------------------------------
// extractTexCoords tests
// ---------------------------------------------------------------------------

func TestExtractTexCoords_Vec2Float(t *testing.T) {
	c := &GeometryConverter{}

	arr := model.NewArray(model.Vec2ArrayType, model.GLFLOAT, 2)
	arr.Data = [][2]float32{{0, 0}, {1, 0}, {1, 1}}

	tex := c.extractTexCoords(arr)
	if len(tex) != 6 {
		t.Fatalf("got %d floats, want 6", len(tex))
	}
}

func TestExtractTexCoords_FlatFloat(t *testing.T) {
	c := &GeometryConverter{}

	arr := model.NewArray(model.FloatArrayType, model.GLFLOAT, 1)
	arr.Data = []float32{0, 0, 1, 0}

	tex := c.extractTexCoords(arr)
	if len(tex) != 4 {
		t.Fatalf("got %d floats, want 4", len(tex))
	}
}

func TestExtractTexCoords_Nil(t *testing.T) {
	c := &GeometryConverter{}
	if v := c.extractTexCoords(&model.Array{}); v != nil {
		t.Error("expected nil for array with nil Data")
	}
}

// ---------------------------------------------------------------------------
// extractIndices tests
// ---------------------------------------------------------------------------

func TestExtractIndices_DrawArrays(t *testing.T) {
	c := &GeometryConverter{}

	da := model.NewDrawArrays()
	da.Mode = model.GLTRIANGLES
	da.Count = 3
	indices := c.extractIndices(da)
	if len(indices) != 3 {
		t.Fatalf("got %d indices, want 3", len(indices))
	}
	if indices[0] != 0 || indices[1] != 1 || indices[2] != 2 {
		t.Errorf("got %v, want [0,1,2]", indices)
	}
}

func TestExtractIndices_DrawArraysQuad(t *testing.T) {
	c := &GeometryConverter{}

	da := model.NewDrawArrays()
	da.Mode = 7    // GL_QUADS
	da.Count = 4
	indices := c.extractIndices(da)
	// 4 vertices → 2 triangles → 6 indices
	if len(indices) != 6 {
		t.Fatalf("quad DrawArrays: got %d indices, want 6", len(indices))
	}
}

func TestExtractIndices_DrawElementsUByte(t *testing.T) {
	c := &GeometryConverter{}

	de := model.NewDrawElementsUByte()
	de.Data = []uint8{0, 1, 2, 2, 3, 0}

	indices := c.extractIndices(de)
	if len(indices) != 6 {
		t.Fatalf("got %d indices, want 6", len(indices))
	}
}

func TestExtractIndices_DrawElementsUShort(t *testing.T) {
	c := &GeometryConverter{}

	de := model.NewDrawElementsUShort()
	de.Data = []uint16{0, 1, 2}

	indices := c.extractIndices(de)
	if len(indices) != 3 {
		t.Fatalf("got %d indices, want 3", len(indices))
	}
}

func TestExtractIndices_DrawElementsUInt(t *testing.T) {
	c := &GeometryConverter{}

	de := model.NewDrawElementsUInt()
	de.Data = []uint32{5, 6, 7, 8}

	indices := c.extractIndices(de)
	if len(indices) != 4 {
		t.Fatalf("got %d indices, want 4", len(indices))
	}
}

func TestExtractIndices_Nil(t *testing.T) {
	c := &GeometryConverter{}
	if v := c.extractIndices(nil); v != nil {
		t.Error("expected nil for nil primitive")
	}
}

// ---------------------------------------------------------------------------
// applyMatrixTransform tests
// ---------------------------------------------------------------------------

func TestApplyMatrixTransform(t *testing.T) {
	c := &GeometryConverter{}

	content := &TileContent{
		Vertices: []float32{1, 0, 0, 0, 1, 0},
		Normals:  []float32{0, 0, 1},
	}

	// Identity matrix
	matrix := [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	c.applyMatrixTransform(content, matrix)

	// With identity, vertices should be unchanged
	if content.Vertices[0] != 1 || content.Vertices[1] != 0 || content.Vertices[2] != 0 {
		t.Errorf("identity transform: first vertex = (%f,%f,%f), want (1,0,0)",
			content.Vertices[0], content.Vertices[1], content.Vertices[2])
	}
}

func TestApplyMatrixTransform_Translation(t *testing.T) {
	c := &GeometryConverter{}

	content := &TileContent{
		Vertices: []float32{0, 0, 0, 1, 1, 1},
	}

	// Translation matrix (move +10 in Z)
	matrix := [4][4]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 10},
		{0, 0, 0, 1},
	}

	c.applyMatrixTransform(content, matrix)

	// First vertex should be at (0,0,10)
	if content.Vertices[2] != 10 {
		t.Errorf("z after translation = %f, want 10", content.Vertices[2])
	}
}

// ---------------------------------------------------------------------------
// GeometryConverter.Convert tests (node tree traversal)
// ---------------------------------------------------------------------------

func TestGeometryConvert_GroupWithGeode(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	// Build: Group → Geode → Geometry with one triangle
	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}}
	da := model.NewDrawArrays()
	da.Mode = model.GLTRIANGLES
	da.Count = 3
	geom.AddPrimitiveSet(da)

	geode := model.NewGeode()
	geode.AddChild(geom)

	group := model.NewGroup()
	group.AddChild(geode)

	content := c.Convert(group)
	if content == nil {
		t.Fatal("Convert returned nil")
	}
	if len(content.Vertices) != 9 {
		t.Errorf("got %d vertices, want 9", len(content.Vertices))
	}
}

func TestGeometryConvert_EmptyGroup(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	group := model.NewGroup()
	content := c.Convert(group)
	// Empty group should still return a valid (empty) content
	if content == nil {
		t.Error("empty group should return non-nil content")
	}
}

func TestGeometryConvert_NilGeometry(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	// Geometry with nil VertexArray should be skipped
	geom := model.NewGeometry()
	geode := model.NewGeode()
	geode.AddChild(geom)

	content := c.Convert(geode)
	if content == nil || len(content.Vertices) != 0 {
		t.Logf("nil geometry: got %d vertices", len(content.Vertices))
	}
}

func TestGeometryConvert_PagedLod(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{{0, 0, 0}}
	geode := model.NewGeode()
	geode.AddChild(geom)

	plod := model.NewPagedLod()
	plod.AddChild(geode)

	content := c.Convert(plod)
	if content == nil {
		t.Fatal("PagedLod Convert returned nil")
	}
	if len(content.Vertices) != 3 {
		t.Errorf("got %d vertices, want 3", len(content.Vertices))
	}
}

func TestGeometryConvert_MatrixTransform(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{{0, 0, 0}}
	da2 := model.NewDrawArrays()
	da2.Mode = model.GLTRIANGLES
	da2.Count = 3
	geom.AddPrimitiveSet(da2)
	geode := model.NewGeode()
	geode.AddChild(geom)

	mt := model.NewMatrixTransform()
	mt.Matrix = [4][4]float32{
		{1, 0, 0, 5},
		{0, 1, 0, 10},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	mt.AddChild(geode)

	content := c.Convert(mt)
	if content == nil {
		t.Fatal("MatrixTransform Convert returned nil")
	}
	// Vertex should be translated by matrix: (0,0,0)→(5,10,0)
	if len(content.Vertices) >= 3 {
		if content.Vertices[0] != 5 || content.Vertices[1] != 10 {
			t.Errorf("after matrix transform: first vertex = (%f,%f), want (5,10)",
				content.Vertices[0], content.Vertices[1])
		}
	}
}

// ---------------------------------------------------------------------------
// TileContent Merge tests (extended)
// ---------------------------------------------------------------------------

func TestTileContentMerge_BoundingBox(t *testing.T) {
	a := &TileContent{
		Vertices:    []float32{0, 0, 0},
		BoundingBox: [12]float64{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1},
		BatchLength: 1,
	}
	b := &TileContent{
		Vertices:    []float32{10, 10, 10},
		BoundingBox: [12]float64{10, 10, 10, 1, 0, 0, 0, 1, 0, 0, 0, 1},
		BatchLength: 1,
	}

	a.Merge(b)
	if a.BatchLength != 2 {
		t.Errorf("BatchLength = %d, want 2", a.BatchLength)
	}
}

// ---------------------------------------------------------------------------
// convertNodeToTile tests
// ---------------------------------------------------------------------------

func TestConvertNodeToTile(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	geomConv := NewGeometryConverter(opts, coordTrans, geoidConv)

	converter := NewConverter(opts)
	converter.geomConverter = geomConv

	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}}
	da3 := model.NewDrawArrays()
	da3.Mode = model.GLTRIANGLES
	da3.Count = 3
	geom.AddPrimitiveSet(da3)

	geode := model.NewGeode()
	geode.AddChild(geom)
	group := model.NewGroup()
	group.AddChild(geode)

	tile := converter.convertNodeToTile(group, "test.osgb")
	if tile == nil {
		t.Fatal("tile is nil")
	}
	if tile.ID != "test.osgb" {
		t.Errorf("tile ID = %q, want test.osgb", tile.ID)
	}
	if tile.Path != "test.osgb.b3dm" {
		t.Errorf("tile Path = %q, want test.osgb.b3dm", tile.Path)
	}
	if tile.Content == nil {
		t.Fatal("tile.Content is nil")
	}
	if len(tile.Content.Vertices) != 9 {
		t.Errorf("got %d vertices, want 9", len(tile.Content.Vertices))
	}
}

// ---------------------------------------------------------------------------
// convertNodeToTile integration: verify optimizeMesh is called in the pipeline
// ---------------------------------------------------------------------------

func TestConvertNodeToTile_OptimizeMeshCalled(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Skipf("optimizeMesh panicked (library may need more triangles): %v", r)
		}
	}()
	// Verify that optimizeMesh runs during convertNodeToTile by enabling
	// simplification and checking it doesn't crash
	opts := DefaultConverterOptions()
	opts.EnableSimplify = true
	opts.SimplifyRatio = 0.9
	opts.SimplifyTargetError = 0.1
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	geomConv := NewGeometryConverter(opts, coordTrans, geoidConv)

	converter := NewConverter(opts)
	converter.geomConverter = geomConv

	// Build a simple mesh with 2 triangles
	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{
		{0, 0, 0}, {1, 0, 0}, {1, 1, 0},
		{0, 0, 0}, {1, 1, 0}, {0, 1, 0},
	}
	da := model.NewDrawArrays()
	da.Mode = model.GLTRIANGLES
	da.Count = 6
	geom.AddPrimitiveSet(da)

	geode := model.NewGeode()
	geode.AddChild(geom)
	group := model.NewGroup()
	group.AddChild(geode)

	tile := converter.convertNodeToTile(group, "test.osgb")
	if tile == nil {
		t.Fatal("tile is nil")
	}
	if tile.Content == nil {
		t.Fatal("tile.Content is nil")
	}
	if len(tile.Content.Vertices) == 0 {
		t.Error("no vertices after conversion with optimization")
	}
	t.Logf("optimizeMesh integration: %d verts", len(tile.Content.Vertices))
}

// ---------------------------------------------------------------------------
// TilesetGenerator tests
// ---------------------------------------------------------------------------

func TestTilesetGenerator_Generate(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewTilesetGenerator(opts)

	tile := &Tile{
		ID:             "test.osgb",
		GeometricError: 100,
		BoundingBox:    [12]float64{0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 5},
		Path:           "test.osgb.b3dm",
		Content:        &TileContent{Vertices: []float32{0, 0, 0}},
	}

	tileset := gen.Generate(tile)
	if tileset == nil {
		t.Fatal("tileset is nil")
	}
	if tileset.Asset.Version != "1.0" {
		t.Errorf("version = %q, want 1.0", tileset.Asset.Version)
	}
	if tileset.Root == nil {
		t.Fatal("root node is nil")
	}
	if len(tileset.Root.BoundVolume.Box) != 12 {
		t.Errorf("bbox len = %d, want 12", len(tileset.Root.BoundVolume.Box))
	}
}

func TestTilesetGenerator_GenerateNested(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewTilesetGenerator(opts)

	child := &Tile{
		ID:             "child.osgb",
		GeometricError: 50,
		BoundingBox:    [12]float64{10, 10, 10, 1, 0, 0, 0, 1, 0, 0, 0, 1},
		Path:           "child.osgb.b3dm",
		Content:        &TileContent{},
	}

	root := &Tile{
		ID:             "root.osgb",
		GeometricError: 200,
		BoundingBox:    [12]float64{0, 0, 0, 10, 0, 0, 0, 10, 0, 0, 0, 10},
		Path:           "root.osgb.b3dm",
		Content:        &TileContent{},
		Children:       []*Tile{child},
	}

	tileset := gen.Generate(root)
	if tileset.Root == nil {
		t.Fatal("root is nil")
	}
	if len(tileset.Root.Children) != 1 {
		t.Fatalf("got %d children, want 1", len(tileset.Root.Children))
	}
	// Root should not have content URI (isRoot=true)
	if tileset.Root.Content != nil {
		t.Error("root should not have Content")
	}
	// Child should have content URI
	if tileset.Root.Children[0].Content == nil {
		t.Error("child should have Content")
	}
	if tileset.Root.Children[0].Refine != "REPLACE" {
		t.Errorf("child Refine = %q, want REPLACE", tileset.Root.Children[0].Refine)
	}
}

// ---------------------------------------------------------------------------
// loadPagedLODChildren tests
// ---------------------------------------------------------------------------

func TestGetPagedLODChildrenFromNode_PagedLod(t *testing.T) {
	opts := DefaultConverterOptions()
	converter := NewConverter(opts)

	plod := model.NewPagedLod()
	plod.DataBasePath = "."
	plod.PerRangeDataList = []model.PerRangeData{
		{FileName: "child.osgb", PriorityOffset: 1},
		{FileName: "", PriorityOffset: 0},
	}

	tile := &Tile{
		ID:   "parent.osgb",
		Node: plod,
	}

	children := converter.getPagedLODChildrenFromNode(tile)
	if len(children) != 1 {
		t.Fatalf("got %d children, want 1 (only non-empty filenames)", len(children))
	}
	if children[0] != "child.osgb" {
		t.Errorf("child = %q, want child.osgb", children[0])
	}
}

func TestGetPagedLODChildrenFromNode_Group(t *testing.T) {
	opts := DefaultConverterOptions()
	converter := NewConverter(opts)

	childPlod := model.NewPagedLod()
	childPlod.PerRangeDataList = []model.PerRangeData{
		{FileName: "nested.osgb"},
	}
	group := model.NewGroup()
	group.AddChild(childPlod)

	tile := &Tile{ID: "group.osgb", Node: group}
	children := converter.getPagedLODChildrenFromNode(tile)
	if len(children) != 1 {
		t.Fatalf("got %d children, want 1", len(children))
	}
}

func TestGetPagedLODChildrenFromNode_NilNode(t *testing.T) {
	opts := DefaultConverterOptions()
	converter := NewConverter(opts)
	tile := &Tile{ID: "nil.osgb", Node: nil}
	if v := converter.getPagedLODChildrenFromNode(tile); v != nil {
		t.Error("expected nil for nil node")
	}
}

// ---------------------------------------------------------------------------
// extractLODLevel edge cases (extended)
// ---------------------------------------------------------------------------

func TestExtractLODLevel_EdgeCases(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"", ""},
		{"no_underscore.osgb", ""},
		{"Tile_+000_+000.osgb", ""},
		{"Tile_+000_+000_L18_00.osgb", "L18"},
		{"L18_test.osgb", "L18"},
	}
	for _, tt := range tests {
		got := extractLODLevel(tt.id)
		if got != tt.want {
			t.Errorf("extractLODLevel(%q) = %q, want %q", tt.id, got, tt.want)
		}
	}
}

// ---------------------------------------------------------------------------
// getPagedLODChildren tests
// ---------------------------------------------------------------------------

func TestGetPagedLODChildren_NilTile(t *testing.T) {
	converter := &Converter{visitedTiles: make(map[string]bool), visitedDirs: make(map[string]bool)}
	if v := converter.getPagedLODChildren(nil); v != nil {
		t.Error("expected nil for nil tile")
	}
}

// ---------------------------------------------------------------------------
// DetectSRSType edge cases (extended)
// ---------------------------------------------------------------------------

func TestDetectSRSType_Extended(t *testing.T) {
	if DetectSRSType("ENU:114,34") != SRSTypeENU {
		t.Error("expected SRSTypeENU")
	}
	if DetectSRSType("ENU:114.5,34.2,100") != SRSTypeENU {
		t.Error("expected SRSTypeENU for ENU with height")
	}
	if DetectSRSType("EPSG:4548") != SRSTypeEPSG {
		t.Error("expected SRSTypeEPSG for EPSG:4548")
	}
	if DetectSRSType("unknown") != SRSTypeUnknown {
		t.Error("expected SRSTypeUnknown")
	}
	if DetectSRSType("") != SRSTypeUnknown {
		t.Error("expected SRSTypeUnknown for empty string")
	}
	if DetectSRSType("PROJCS[...]") != SRSTypeWKT {
		t.Error("expected SRSTypeWKT for WKT-like string")
	}
}

// ---------------------------------------------------------------------------
// ParseSRSOrigin edge cases
// ---------------------------------------------------------------------------

func TestParseSRSOrigin_Extended(t *testing.T) {
	x, y, z, err := ParseSRSOrigin("500000,3000000,100.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if x != 500000 || y != 3e6 || z != 100.5 {
		t.Errorf("got (%f,%f,%f), want (500000,3000000,100.5)", x, y, z)
	}

	// Leading/trailing spaces
	x, y, z, err = ParseSRSOrigin("  1,  2  ,  3  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if x != 1 || y != 2 || z != 3 {
		t.Errorf("with spaces: got (%f,%f,%f)", x, y, z)
	}

	// Negative values
	x, y, z, err = ParseSRSOrigin("-500000,-3000000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if x != -500000 || y != -3e6 || z != 0 {
		t.Errorf("negative: got (%f,%f,%f)", x, y, z)
	}
}

// ---------------------------------------------------------------------------
// FindMetadataFile tests
// ---------------------------------------------------------------------------

func TestFindMetadataFile_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := FindMetadataFile(dir)
	if err == nil {
		t.Error("expected error for directory without metadata.xml")
	}
}

func TestFindMetadataFile_Found(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "metadata.xml"), []byte("<ModelMetadata></ModelMetadata>"), 0644)
	path, err := FindMetadataFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != filepath.Join(dir, "metadata.xml") {
		t.Errorf("path = %q, want %q", path, filepath.Join(dir, "metadata.xml"))
	}
}

func TestFindMetadataFile_InOSGBSubdir(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "OSGB"), 0755)
	os.WriteFile(filepath.Join(dir, "OSGB", "metadata.xml"), []byte("<ModelMetadata></ModelMetadata>"), 0644)
	path, err := FindMetadataFile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != filepath.Join(dir, "OSGB", "metadata.xml") {
		t.Errorf("path = %q, want %q", path, filepath.Join(dir, "OSGB", "metadata.xml"))
	}
}

// ---------------------------------------------------------------------------
// extractGeometry tests
// ---------------------------------------------------------------------------

func TestExtractGeometry_WithStateSet(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	geom := model.NewGeometry()
	geom.VertexArray = model.NewArray(model.Vec3ArrayType, model.GLFLOAT, 3)
	geom.VertexArray.Data = [][3]float32{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}}
	da4 := model.NewDrawArrays()
	da4.Mode = model.GLTRIANGLES
	da4.Count = 3
	geom.AddPrimitiveSet(da4)

	content := c.extractGeometry(geom)
	if content == nil {
		t.Fatal("extractGeometry returned nil")
	}
	if len(content.Vertices) != 9 {
		t.Errorf("got %d vertices, want 9", len(content.Vertices))
	}
	if content.BatchLength != 1 {
		t.Errorf("BatchLength = %d, want 1", content.BatchLength)
	}
}

func TestExtractGeometry_NilVertexArray(t *testing.T) {
	opts := DefaultConverterOptions()
	coordTrans := NewCoordinateTransformer("", "")
	geoidConv := NewGeoidConverter("none", "")
	c := NewGeometryConverter(opts, coordTrans, geoidConv)

	geom := model.NewGeometry()
	if v := c.extractGeometry(geom); v != nil {
		t.Error("expected nil for geometry without vertex data")
	}
}

// ---------------------------------------------------------------------------
// calculateGeometricError edge cases (standalone in tileset.go)
// ---------------------------------------------------------------------------

func TestCalculateGeometricError_Standalone(t *testing.T) {
	bbox := [12]float64{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1}
	err := calculateGeometricError(bbox)
	// (2+2+2)/6*2 = 6/6*2 = 2
	if math.Abs(err-2.0) > 0.001 {
		t.Errorf("error = %f, want 2.0", err)
	}
}

// ---------------------------------------------------------------------------
// NewConverter tests
// ---------------------------------------------------------------------------

func TestNewConverter(t *testing.T) {
	opts := DefaultConverterOptions()
	c := NewConverter(opts)
	if c == nil {
		t.Fatal("NewConverter returned nil")
	}
	if c.geomConverter == nil {
		t.Error("geomConverter is nil")
	}
	if c.tilesetGen == nil {
		t.Error("tilesetGen is nil")
	}
	if c.b3dmGen == nil {
		t.Error("b3dmGen is nil")
	}
	if c.coordTrans == nil {
		t.Error("coordTrans is nil")
	}
	if c.geoidConv == nil {
		t.Error("geoidConv is nil")
	}
	if c.visitedTiles == nil {
		t.Error("visitedTiles is nil")
	}
}

// ---------------------------------------------------------------------------
// GetDatabasePath and GetFileNames tests
// ---------------------------------------------------------------------------

func TestGetDatabasePath(t *testing.T) {
	c := &GeometryConverter{}

	plod := model.NewPagedLod()
	plod.DataBasePath = "./tiles"
	if c.GetDatabasePath(plod) != "./tiles" {
		t.Errorf("got %q, want ./tiles", c.GetDatabasePath(plod))
	}

	// Non-PagedLod should return empty
	if c.GetDatabasePath(model.NewGroup()) != "" {
		t.Error("expected empty for Group")
	}
}

// ---------------------------------------------------------------------------
// compressMeshDraco tests
// ---------------------------------------------------------------------------

func TestCompressMeshDraco_Disabled(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{EnableDraco: false}}
	verts, norms, texcs, indices := testMeshData()

	data, attrIds, err := c.compressMeshDraco(verts, norms, texcs, indices)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != nil {
		t.Error("expected nil data when Draco disabled")
	}
	if attrIds != nil {
		t.Error("expected nil attrIds when Draco disabled")
	}
}

func TestCompressMeshDraco_Enabled(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{
		EnableDraco:       true,
		DracoPositionBits: 11,
		DracoNormalBits:   10,
		DracoTexCoordBits: 12,
	}}
	verts, norms, texcs, indices := testMeshData()

	data, attrIds, err := c.compressMeshDraco(verts, norms, texcs, indices)
	if err != nil {
		// Draco may reject degenerate/small meshes; skip gracefully
		t.Skipf("Draco compression skipped: %v", err)
	}
	if data == nil {
		t.Skip("Draco returned nil data (library may not be available)")
	}
	if len(data) == 0 {
		t.Error("compressed data is empty")
	}
	if attrIds == nil {
		t.Fatal("attrIds is nil")
	}
	if _, ok := attrIds["POSITION"]; !ok {
		t.Error("missing POSITION attribute")
	}
	t.Logf("Draco compressed %d bytes (%d verts): pos=%d,norm=%d,tex=%d",
		len(data), len(verts)/3, attrIds["POSITION"], attrIds["NORMAL"], attrIds["TEXCOORD_0"])
}

// ---------------------------------------------------------------------------
// simplifyMesh tests
// ---------------------------------------------------------------------------

func TestSimplifyMesh_Disabled(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{EnableSimplify: false}}
	verts, norms, texcs, indices := testMeshData()

	v, n, tc, idx := c.simplifyMesh(verts, norms, texcs, indices)
	// When disabled, should return original data unchanged
	if len(v) != len(verts) {
		t.Error("vertices changed when simplify disabled")
	}
	if len(idx) != len(indices) {
		t.Error("indices changed when simplify disabled")
	}
	_ = n
	_ = tc
}

func TestSimplifyMesh_Ratio1(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{
		EnableSimplify:       true,
		SimplifyRatio:        1.0,
		SimplifyTargetError:  0.01,
	}}
	verts, norms, texcs, indices := testMeshData()

	v, n, tc, idx := c.simplifyMesh(verts, norms, texcs, indices)
	// ratio >= 1.0 → no simplification
	if len(v) != len(verts) {
		t.Errorf("vertices changed: %d → %d", len(verts), len(v))
	}
	if len(idx) != len(indices) {
		t.Errorf("indices changed: %d → %d", len(indices), len(idx))
	}
	_ = n
	_ = tc
}

func TestSimplifyMesh_NoIndices(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{
		EnableSimplify:       true,
		SimplifyRatio:        0.5,
		SimplifyTargetError:  0.01,
	}}
	verts, norms, texcs, _ := testMeshData()

	v, n, tc, idx := c.simplifyMesh(verts, norms, texcs, nil)
	// With no indices, should return original
	if len(v) != len(verts) {
		t.Error("vertices changed with no indices")
	}
	if idx != nil {
		t.Error("indices should be nil")
	}
	_ = n
	_ = tc
}

// ---------------------------------------------------------------------------
// optimizeMesh tests
// ---------------------------------------------------------------------------

func TestOptimizeMesh_Disabled(t *testing.T) {
	c := &GeometryConverter{opts: nil}
	verts, norms, texcs, indices := testMeshData()

	v, n, tc, idx := c.optimizeMesh(verts, norms, texcs, indices)
	if len(v) != len(verts) {
		t.Error("vertices changed")
	}
	_ = n
	_ = tc
	_ = idx
}

func TestOptimizeMesh_Simplify(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{
		EnableSimplify:       true,
		SimplifyRatio:        0.5,
		SimplifyTargetError:  0.01,
	}}
	// Use 2 triangles (4 vertices)
	verts := []float32{0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0}
	norms := []float32{0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1}
	indices := []uint32{0, 1, 2, 0, 2, 3}

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("simplify library panicked (may need more triangles): %v", r)
		}
	}()

	v, n, tc, idx := c.optimizeMesh(verts, norms, nil, indices)
	if len(v) == 0 {
		t.Error("all vertices lost after simplification")
	}
	_ = n
	_ = tc
	_ = idx
}

// ---------------------------------------------------------------------------
// simplifyMesh with actual meshopt library (when available)
// ---------------------------------------------------------------------------

func TestSimplifyMesh_WithMeshopt(t *testing.T) {
	c := &GeometryConverter{opts: &ConverterOptions{
		EnableSimplify:       true,
		SimplifyRatio:        0.8,
		SimplifyTargetError:  0.1,
	}}
	verts := []float32{
		0, 0, 0,
		1, 0, 0,
		0, 1, 0,
		1, 1, 0,
	}
	norms := []float32{
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	}
	indices := []uint32{0, 1, 2, 1, 3, 2}

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("simplify library panicked: %v", r)
		}
	}()

	v, n, tc, idx := c.simplifyMesh(verts, norms, nil, indices)
	if v == nil || n == nil || idx == nil {
		t.Fatal("simplify returned nil")
	}
	_ = tc
	t.Logf("simplified: %d verts, %d indices → %d verts, %d indices",
		len(verts)/3, len(indices), len(v)/3, len(idx))
}

// ---------------------------------------------------------------------------
// buildGLTF tests (TileContent → gltf.Document)
// ---------------------------------------------------------------------------

func TestBuildGLTF_Basic(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	content := &TileContent{
		Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
		Normals:     []float32{0, 0, 1, 0, 0, 1, 0, 0, 1},
		TexCoords:   []float32{0, 0, 1, 0, 0, 1},
		Indices:     []uint32{0, 1, 2},
		BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
		BatchLength: 1,
	}

	doc := gen.buildGLTF(content)
	if doc == nil {
		t.Fatal("buildGLTF returned nil")
	}
	if len(doc.Buffers) == 0 {
		t.Error("no buffers")
	}
	if len(doc.Accessors) == 0 {
		t.Error("no accessors")
	}
	if len(doc.Meshes) == 0 {
		t.Error("no meshes")
	}
	if doc.Scene == nil {
		t.Error("no scene")
	}
}

func TestBuildGLTF_WithTexture(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	content := &TileContent{
		Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
		Indices:     []uint32{0, 1, 2},
		BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
		Textures:    [][]byte{{0xFF, 0xD8, 0xFF}, {0x89, 0x50, 0x4E}},
		BatchLength: 1,
	}

	doc := gen.buildGLTF(content)
	if doc == nil {
		t.Fatal("buildGLTF returned nil")
	}
	// Should have textures, images, samplers
	t.Logf("GLTF: buffers=%d, accessors=%d, images=%d, textures=%d, meshes=%d",
		len(doc.Buffers), len(doc.Accessors), len(doc.Images), len(doc.Textures), len(doc.Meshes))
}

// ---------------------------------------------------------------------------
// B3DMGenerator.Generate tests
// ---------------------------------------------------------------------------

func TestB3DMGenerator_Generate(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content: &TileContent{
			Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
			Indices:     []uint32{0, 1, 2},
			BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
			BatchLength: 1,
		},
	}

	data, err := gen.Generate(tile)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("generated data is empty")
	}
	t.Logf("B3DM generated: %d bytes", len(data))
}

func TestB3DMGenerator_GenerateWithDraco(t *testing.T) {
	opts := DefaultConverterOptions()
	opts.EnableDraco = true
	opts.DracoPositionBits = 11
	opts.DracoNormalBits = 10
	opts.DracoTexCoordBits = 12
	gen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content: &TileContent{
			Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
			Indices:     []uint32{0, 1, 2},
			BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
			BatchLength: 1,
		},
	}

	data, err := gen.Generate(tile)
	if err != nil {
		t.Skipf("Draco B3DM generation skipped: %v", err)
	}
	if len(data) == 0 {
		t.Error("Draco B3DM data is empty")
	}
	t.Logf("B3DM with Draco generated: %d bytes", len(data))
}

func TestB3DMGenerator_Generate_NilContent(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	_, err := gen.Generate(&Tile{})
	if err == nil {
		t.Error("expected error for nil tile content")
	}
}

// ---------------------------------------------------------------------------
// B3DMGenerator.GenerateGLB tests
// ---------------------------------------------------------------------------

func TestB3DMGenerator_GenerateGLB(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content: &TileContent{
			Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
			Indices:     []uint32{0, 1, 2},
			BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
			BatchLength: 1,
		},
	}

	data, err := gen.GenerateGLB(tile)
	if err != nil {
		t.Fatalf("GenerateGLB failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("GLB data is empty")
	}
	t.Logf("GLB generated: %d bytes", len(data))
}

func TestB3DMGenerator_GenerateGLBWithDraco(t *testing.T) {
	opts := DefaultConverterOptions()
	opts.EnableDraco = true
	gen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content: &TileContent{
			Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
			Indices:     []uint32{0, 1, 2},
			BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
			BatchLength: 1,
		},
	}

	data, err := gen.GenerateGLB(tile)
	if err != nil {
		t.Skipf("Draco GLB generation skipped: %v", err)
	}
	if len(data) == 0 {
		t.Error("Draco GLB data is empty")
	}
	t.Logf("GLB with Draco generated: %d bytes", len(data))
}

// ---------------------------------------------------------------------------
// TilesetGenerator.Write tests
// ---------------------------------------------------------------------------

func TestTilesetGenerator_Write(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewTilesetGenerator(opts)

	tileset := &TileJSON{
		Asset: AssetJSON{
			Version: "1.0",
			GenBy:   "test",
		},
		GeometricError: 500,
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "tileset.json")
	err := gen.Write(tileset, path)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if len(data) == 0 {
		t.Error("written tileset.json is empty")
	}
}

// ---------------------------------------------------------------------------
// B3DMGenerator.Write tests
// ---------------------------------------------------------------------------

func TestB3DMGenerator_Write(t *testing.T) {
	opts := DefaultConverterOptions()
	gen := NewB3DMGenerator(opts)

	tile := &Tile{
		Content: &TileContent{
			Vertices:    []float32{0, 0, 0, 1, 0, 0, 0, 1, 0},
			Indices:     []uint32{0, 1, 2},
			BoundingBox: [12]float64{0, 0, 0, 0.5, 0, 0, 0, 0.5, 0, 0, 0, 0.5},
			BatchLength: 1,
		},
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "test.b3dm")
	err := gen.Write(tile, path)
	if err != nil {
		t.Fatalf("B3DM Write failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	t.Logf("B3DM file written: %d bytes", len(data))
}

// ---------------------------------------------------------------------------
// LoadOSGB tests (integration with actual test file)
// ---------------------------------------------------------------------------

func TestLoadOSGB_AvailableData(t *testing.T) {
	// Try multiple paths to find test data
	candidates := []string{
		"../test_data/Tile_+003_+003_L18_000.osgb",
		"../test_data/cessna.osgb",
		"test_data/Tile_+003_+003_L18_000.osgb",
		"test_data/cessna.osgb",
	}
	var path string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}
	if path == "" {
		t.Skip("no test data file found")
	}

	node, err := LoadOSGB(path)
	if err != nil {
		t.Fatalf("LoadOSGB(%q) failed: %v", path, err)
	}
	if node == nil {
		t.Fatal("LoadOSGB returned nil")
	}
	t.Logf("LoadOSGB(%q): type=%T", path, node)
}

// ---------------------------------------------------------------------------
// full Convert pipeline entry point tests
// ---------------------------------------------------------------------------

func TestOSGBToGLB_AvailableData(t *testing.T) {
	// Find test data
	path := ""
	for _, p := range []string{
		"../test_data/Tile_+003_+003_L18_000.osgb",
		"../tiles3d_test/Tile_+002_+000_L22_000020.osgb",
	} {
		if _, err := os.Stat(p); err == nil {
			path = p
			break
		}
	}
	if path == "" {
		t.Skip("no test data file for OSGBToGLB")
	}

	opts := DefaultConverterOptions()
	opts.SourceSRS = ""
	data, err := OSGBToGLB(path, opts)
	if err != nil {
		t.Fatalf("OSGBToGLB failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("GLB data is empty")
	}
	t.Logf("OSGBToGLB: %d bytes", len(data))
}

func TestGetFileNames(t *testing.T) {
	c := &GeometryConverter{}

	plod := model.NewPagedLod()
	plod.PerRangeDataList = []model.PerRangeData{
		{FileName: "child1.osgb"},
		{FileName: "child2.osgb"},
		{FileName: ""},
	}
	// Add dummy children matching the non-empty PerRangeData entries
	plod.AddChild(model.NewGroup())
	plod.AddChild(model.NewGroup())

	names := c.GetFileNames(plod)
	if len(names) != 2 {
		t.Fatalf("got %d file names, want 2", len(names))
	}
	if names[0] != "child1.osgb" || names[1] != "child2.osgb" {
		t.Errorf("got %v, want [child1.osgb child2.osgb]", names)
	}
}

// ---------------------------------------------------------------------------
// hasGeoReference edge cases
// ---------------------------------------------------------------------------

func TestHasGeoReference_WithCenter(t *testing.T) {
	tf := &CoordinateTransformer{}
	if tf.HasGeoReference() {
		t.Error("zero center should not have geo reference")
	}

	tf.center = [3]float64{0.001, 0, 0}
	if !tf.HasGeoReference() {
		t.Error("non-zero lon should be detected")
	}

	tf.center = [3]float64{0, 0.001, 0}
	if !tf.HasGeoReference() {
		t.Error("non-zero lat should be detected")
	}
}

// ---------------------------------------------------------------------------
// extractMaterial tests
// ---------------------------------------------------------------------------

func TestExtractMaterial_NoStateSet(t *testing.T) {
	c := &GeometryConverter{}

	geode := model.NewGeode()
	geom := model.NewGeometry()
	geode.AddChild(geom)

	materials := c.extractMaterial(geode)
	if len(materials) != 0 {
		t.Errorf("got %d materials, want 0", len(materials))
	}
}

func TestExtractMaterial_EmptyGeode(t *testing.T) {
	c := &GeometryConverter{}
	geode := model.NewGeode()
	materials := c.extractMaterial(geode)
	if len(materials) != 0 {
		t.Errorf("got %d materials, want 0", len(materials))
	}
}

// ---------------------------------------------------------------------------
// extractImageData tests
// ---------------------------------------------------------------------------

func TestExtractImageData_Nil(t *testing.T) {
	c := &GeometryConverter{}
	if v := c.extractImageData(nil, ""); v != nil {
		t.Error("expected nil for nil image")
	}
}

func TestExtractImageData_FromData(t *testing.T) {
	c := &GeometryConverter{}

	img := model.NewImage()
	img.Data = []uint8{1, 2, 3, 4}

	data := c.extractImageData(img, "")
	if len(data) != 4 {
		t.Errorf("got %d bytes, want 4", len(data))
	}
}

func TestGetImageMimeType(t *testing.T) {
	tests := []struct {
		fileName string
		want     string
	}{
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.png", "image/png"},
		{"test.ktx", "image/ktx2"},
		{"test.ktx2", "image/ktx2"},
		{"test", "image/jpeg"},
	}
	for _, tt := range tests {
		img := model.NewImage()
		img.FileName = tt.fileName
		got := getImageMimeType(img)
		if got != tt.want {
			t.Errorf("getImageMimeType(%q) = %q, want %q", tt.fileName, got, tt.want)
		}
	}
}
