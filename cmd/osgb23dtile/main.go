package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flywave/go-osg/tiles3d"
)

func main() {
	input := flag.String("i", "", "input file or directory")
	output := flag.String("o", "", "output file or directory")
	format := flag.String("f", "osgb", "input format (osgb, gltf, b3dm)")
	srs := flag.String("srs", "", "source SRS (e.g. EPSG:4548, ENU:114,34)")
	maxLvl := flag.Int("max-lvl", -1, "max LOD level")
	enableDraco := flag.Bool("enable-draco", false, "enable Draco mesh compression")
	enableSimplify := flag.Bool("enable-simplify", false, "enable mesh simplification")
	enableUnlit := flag.Bool("enable-unlit", true, "enable KHR_materials_unlit")
	lon := flag.Float64("lon", 0, "center longitude")
	lat := flag.Float64("lat", 0, "center latitude")
	alt := flag.Float64("alt", 0, "center altitude")
	geoid := flag.String("geoid", "none", "geoid model (none, egm84, egm96, egm2008)")
	geoidPath := flag.String("geoid-path", "", "path to geoid data files")
	verbose := flag.Bool("v", false, "verbose output")
	projData := flag.String("proj-data", "", "PROJ_DATA path (overrides PROJ_DATA env)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -i <input> -o <output> -f <format> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFormats:\n")
		fmt.Fprintf(os.Stderr, "  osgb   OSGB directory/file -> 3D Tiles (tileset.json + b3dm)\n")
		fmt.Fprintf(os.Stderr, "  gltf   OSGB file -> GLB\n")
		fmt.Fprintf(os.Stderr, "  b3dm   B3DM file -> GLB\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*input); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: input %q does not exist\n", *input)
		os.Exit(1)
	}

	if *projData != "" {
		os.Setenv("PROJ_DATA", *projData)
		os.Setenv("PROJ_LIB", *projData)
	}

	if *verbose {
		fmt.Printf("Input: %s\n", *input)
		fmt.Printf("Output: %s\n", *output)
		fmt.Printf("Format: %s\n", *format)
		fmt.Printf("SRS: %s\n", *srs)
		fmt.Printf("MaxLOD: %d\n", *maxLvl)
		fmt.Printf("Draco: %v\n", *enableDraco)
		fmt.Printf("Simplify: %v\n", *enableSimplify)
		fmt.Printf("Unlit: %v\n", *enableUnlit)
		fmt.Printf("Center: lon=%f, lat=%f, alt=%f\n", *lon, *lat, *alt)
		fmt.Printf("Geoid: %s\n", *geoid)
	}

	switch *format {
	case "osgb":
		convertOSGB(*input, *output, *srs, *maxLvl, *enableSimplify, *enableDraco, *enableUnlit, *lon, *lat, *alt, *geoid, *geoidPath, *verbose)
	case "gltf":
		convertGLTF(*input, *output, *srs, *enableSimplify, *enableDraco, *enableUnlit, *lon, *lat, *alt, *geoid, *geoidPath, *verbose)
	case "b3dm":
		convertB3DM(*input, *output, *verbose)
	default:
		fmt.Fprintf(os.Stderr, "error: unsupported format %q\n", *format)
		os.Exit(1)
	}
}

func findMetadataInDir(dir string) string {
	candidates := []string{
		filepath.Join(dir, "metadata.xml"),
		filepath.Join(dir, "OSGB", "metadata.xml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func findRootOSGBInDir(dir string) string {
	var found, fallback string
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".osgb") {
			return nil
		}
		if fallback == "" {
			fallback = path
		}
		if name == "main.osgb" {
			found = path
			return filepath.SkipAll
		}
		if !strings.Contains(name, "_L") && found == "" {
			found = path
		}
		return nil
	})
	if found != "" {
		return found
	}
	return fallback
}

func convertOSGB(input, output, srs string, maxLvl int, enableSimplify, enableDraco, enableUnlit bool, lon, lat, alt float64, geoid, geoidPath string, verbose bool) {
	opts := tiles3d.DefaultConverterOptions()
	opts.SourceSRS = srs
	opts.EnableSimplify = enableSimplify
	opts.EnableDraco = enableDraco
	opts.EnableUnlit = enableUnlit
	opts.MaxLOD = maxLvl
	opts.GeoidModel = geoid
	opts.GeoidDataPath = geoidPath

	info, err := os.Stat(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: input %q: %v\n", input, err)
		os.Exit(1)
	}

	var inputPath string
	if info.IsDir() {
		osgbFile := findRootOSGBInDir(input)
		if osgbFile == "" {
			fmt.Fprintf(os.Stderr, "error: no .osgb file found in %q\n", input)
			os.Exit(1)
		}
		inputPath = osgbFile
		if verbose {
			fmt.Printf("Input is a directory, found root OSGB: %s\n", inputPath)
			if meta := findMetadataInDir(input); meta != "" {
				fmt.Printf("Found metadata.xml: %s\n", meta)
			}
		}
	} else {
		inputPath = input
	}

	if lon != 0 || lat != 0 {
		opts.CenterLongitude = lon
		opts.CenterLatitude = lat
		opts.CenterHeight = alt
	}

	if err := os.MkdirAll(output, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Input: %s\n", inputPath)
		fmt.Printf("Output: %s\n", output)
		fmt.Println("Starting OSGB to 3D Tiles conversion...")
	}

	_, err = tiles3d.OSGBTo3DTiles(inputPath, output, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: conversion failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Conversion complete: %s -> %s\n", inputPath, output)
}

func convertGLTF(input, output, srs string, enableSimplify, enableDraco, enableUnlit bool, lon, lat, alt float64, geoid, geoidPath string, verbose bool) {
	opts := tiles3d.DefaultConverterOptions()
	opts.SourceSRS = srs
	opts.EnableSimplify = enableSimplify
	opts.EnableDraco = enableDraco
	opts.EnableUnlit = enableUnlit
	opts.GeoidModel = geoid
	opts.GeoidDataPath = geoidPath

	if lon != 0 || lat != 0 {
		opts.CenterLongitude = lon
		opts.CenterLatitude = lat
		opts.CenterHeight = alt
	}

	if verbose {
		fmt.Printf("Converting OSGB to GLB: %s -> %s\n", input, output)
	}

	err := tiles3d.OSGBToGLBFile(input, output, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: GLB conversion failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("GLB conversion complete: %s\n", output)
}

func convertB3DM(input, output string, verbose bool) {
	if !strings.HasSuffix(output, ".glb") && !strings.HasSuffix(output, ".gltf") {
		output += ".glb"
	}

	data, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to read B3DM file: %v\n", err)
		os.Exit(1)
	}

	if len(data) < 28 {
		fmt.Fprintf(os.Stderr, "error: invalid B3DM file (too small)\n")
		os.Exit(1)
	}

	magic := string(data[0:4])
	if magic != "b3dm" {
		fmt.Fprintf(os.Stderr, "error: not a valid B3DM file (magic: %q)\n", magic)
		os.Exit(1)
	}

	featureJSONLen := int(data[12]) | int(data[13])<<8 | int(data[14])<<16 | int(data[15])<<24
	batchJSONLen := int(data[20]) | int(data[21])<<8 | int(data[22])<<16 | int(data[23])<<24

	glbOffset := 28 + featureJSONLen + batchJSONLen
	if glbOffset >= len(data) {
		fmt.Fprintf(os.Stderr, "error: invalid B3DM file (GLB offset out of range)\n")
		os.Exit(1)
	}

	glbData := data[glbOffset:]

	if err := os.WriteFile(output, glbData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write GLB file: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Extracted GLB from B3DM: %s (%d bytes)\n", output, len(glbData))
	} else {
		fmt.Printf("B3DM to GLB complete: %s\n", output)
	}
}
