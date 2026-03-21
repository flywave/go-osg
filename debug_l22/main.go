package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/flywave/go-osg/tiles3d"
	_ "github.com/flywave/go-proj"
)

func init() {
	os.Setenv("PROJ_DATA", "/Users/xuning/Work/go-proj/proj_data")
	os.Setenv("PROJ_LIB", "/Users/xuning/Work/go-proj/proj_data")
}

func testLoad(path string) {
	fmt.Printf("\n=== Testing: %s ===\n", path)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\n", r)
			debug.PrintStack()
		}
	}()

	node, err := tiles3d.LoadOSGB(path)
	if err != nil {
		fmt.Printf("Error loading: %v\n", err)
		return
	}

	fmt.Printf("Loaded successfully! Node type: %T\n", node)
}

func main() {
	// Test a working file
	testLoad("/Users/xuning/Downloads/data/3/OSGB/Data/Tile_+002_+000/Tile_+002_+000_L21_00000.osgb")

	// Test the problematic file
	testLoad("/Users/xuning/Downloads/data/3/OSGB/Data/Tile_+002_+000/Tile_+002_+000_L22_000020.osgb")
}
