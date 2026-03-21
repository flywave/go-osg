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
	inputPath := "/Users/xuning/Work/go-osg/tiles3d_test/Tile_+002_+000_L22_000020.osgb"

	fmt.Printf("Loading: %s\n", inputPath)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\n", r)
		}
	}()

	node, err := tiles3d.LoadOSGB(inputPath)
	if err != nil {
		fmt.Printf("Error loading: %v\n", err)
		return
	}

	fmt.Printf("Loaded successfully! Node type: %T\n", node)
}
