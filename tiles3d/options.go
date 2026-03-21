package tiles3d

type ConverterOptions struct {
	InputPath  string
	OutputPath string

	SourceSRS string
	TargetSRS string

	GeoidModel    string
	GeoidDataPath string

	EnableSimplify      bool
	SimplifyRatio       float32
	SimplifyTargetError float32
	EnableDraco         bool
	DracoPositionBits   int
	DracoNormalBits     int
	DracoTexCoordBits   int

	EnableTexture   bool
	TextureCompress string
	EnableKTX2      bool

	EnableUnlit bool

	CenterLongitude float64
	CenterLatitude  float64
	CenterHeight    float64

	GeometricError float64

	MaxLOD int
}

type Tile struct {
	ID             string
	Path           string
	GeometricError float64
	BoundingBox    [12]float64
	Children       []*Tile
	Content        *TileContent
	Node           interface{}
}

type TileContent struct {
	Vertices    []float32
	Normals     []float32
	TexCoords   []float32
	Indices     []uint32
	Material    *Material
	Textures    [][]byte
	BatchLength int
	BoundingBox [12]float64
}

type Material struct {
	BaseColorFactor [4]float32
	MetallicFactor  float32
	RoughnessFactor float32
	AlphaCutoff     float32
	DoubleSided     bool
	Unlit           bool
}

type Texture struct {
	ID   string
	Data []byte
	Mime string
}

type BoundingVolume struct {
	Box []float64
}

type TileJSON struct {
	Asset          AssetJSON     `json:"asset"`
	GeometricError float64       `json:"geometricError"`
	Root           *TileJSONNode `json:"root"`
	ExtensionsUsed []string      `json:"extensionsUsed,omitempty"`
	Schema         *SchemaJSON   `json:"schema,omitempty"`
}

type AssetJSON struct {
	Version string `json:"version"`
	GenBy   string `json:"genBy,omitempty"`
}

type TileJSONNode struct {
	BoundVolume    BoundVolumeJSON `json:"boundingVolume"`
	GeometricError float64         `json:"geometricError"`
	Refine         string          `json:"refine"`
	Content        *ContentJSON    `json:"content,omitempty"`
	Children       []*TileJSONNode `json:"children,omitempty"`
}

type BoundVolumeJSON struct {
	Box    []float64 `json:"box,omitempty"`
	Sphere []float64 `json:"sphere,omitempty"`
	Region []float64 `json:"region,omitempty"`
}

type ContentJSON struct {
	URI string `json:"uri"`
}

type SchemaJSON struct {
	Classes map[string]ClassJSON `json:"classes,omitempty"`
}

type ClassJSON struct {
	Properties map[string]PropertyJSON `json:"properties,omitempty"`
}

type PropertyJSON struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable,omitempty"`
}

func DefaultConverterOptions() *ConverterOptions {
	return &ConverterOptions{
		TargetSRS:           "EPSG:4326",
		GeoidModel:          "none",
		EnableTexture:       true,
		TextureCompress:     "jpeg",
		SimplifyRatio:       1.0,
		SimplifyTargetError: 0.01,
		DracoPositionBits:   11,
		DracoNormalBits:     10,
		DracoTexCoordBits:   12,
		EnableUnlit:         true,
		GeometricError:      500.0,
		MaxLOD:              -1,
	}
}
