package tiles3d

import (
	"math"
	"testing"
)

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
