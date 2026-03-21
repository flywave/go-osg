package main

import (
	"fmt"
	"os"

	"github.com/flywave/go-osg/tiles3d"
	_ "github.com/flywave/go-proj"
)

func init() {
	os.Setenv("PROJ_DATA", "/Users/xuning/Work/go-proj/proj_data")
	os.Setenv("PROJ_LIB", "/Users/xuning/Work/go-proj/proj_data")
}

func main() {
	inputPath := "/Users/xuning/Downloads/data/3/OSGB/Data/Tile_+002_+000/Tile_+002_+000.osgb"
	outputPath := "/tmp/test_output"

	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	fmt.Printf("Input: %s\n", inputPath)
	fmt.Printf("Output: %s\n", outputPath)

	opts := tiles3d.DefaultConverterOptions()
	opts.EnableTexture = true
	opts.EnableUnlit = true
	opts.MaxLOD = 10

	// Source SRS is EPSG:4548 (CGCS2000 / 3-degree Gauss-Kruger CM 117E)
	opts.SourceSRS = "EPSG:4548"
	opts.TargetSRS = ""

	// From metadata.xml: SRSOrigin="518078.000000,4080366.000000,0.000000" (EPSG:4548)
	// SetCenter will convert this to WGS84 center
	opts.CenterLongitude = 518078.0
	opts.CenterLatitude = 4080366.0
	opts.CenterHeight = 0.0

	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	fmt.Println("Starting conversion...")

	result, err := tiles3d.OSGBTo3DTiles(inputPath, outputPath, opts)
	if err != nil {
		fmt.Printf("Conversion failed: %v\n", err)
		return
	}

	fmt.Println("Conversion successful!")
	fmt.Printf("JSON length: %d\n", len(result.JSON))
	fmt.Printf("Bounding box: %v\n", result.BoundingBox)

	tilesetPath := outputPath + "/tileset.json"
	if _, err := os.Stat(tilesetPath); err == nil {
		fmt.Println("tileset.json created successfully")
	}

	files, _ := os.ReadDir(outputPath)
	fmt.Printf("Output files (%d):\n", len(files))
	for _, f := range files {
		fmt.Printf("  - %s\n", f.Name())
	}
}
