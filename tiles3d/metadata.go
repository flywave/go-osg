package tiles3d

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ModelMetadata struct {
	XMLName   xml.Name `xml:"ModelMetadata"`
	Version   string   `xml:"version,attr"`
	SRS       string   `xml:"SRS"`
	SRSOrigin string   `xml:"SRSOrigin"`
}

type SRSType int

const (
	SRSTypeUnknown SRSType = iota
	SRSTypeENU
	SRSTypeEPSG
	SRSTypeWKT
)

type GeoReference struct {
	Lon        float64
	Lat        float64
	Height     float64
	OffsetX    float64
	OffsetY    float64
	OffsetZ    float64
	SRSOriginX float64
	SRSOriginY float64
	SRSOriginZ float64
}

func ParseMetadataXML(metadataPath string) (*ModelMetadata, error) {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata.xml: %w", err)
	}

	var metadata ModelMetadata
	if err := xml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata.xml: %w", err)
	}

	return &metadata, nil
}

func DetectSRSType(srs string) SRSType {
	srs = strings.TrimSpace(srs)
	if strings.HasPrefix(srs, "ENU") {
		return SRSTypeENU
	}
	if strings.HasPrefix(srs, "EPSG") {
		return SRSTypeEPSG
	}
	if srs == "unknown" || srs == "" {
		return SRSTypeUnknown
	}
	return SRSTypeWKT
}

func ParseSRSOrigin(origin string) (x, y, z float64, err error) {
	parts := strings.Split(origin, ",")
	if len(parts) < 2 {
		return 0, 0, 0, fmt.Errorf("invalid SRSOrigin format: %s", origin)
	}

	x, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse SRSOrigin X: %w", err)
	}

	y, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse SRSOrigin Y: %w", err)
	}

	z = 0.0
	if len(parts) >= 3 {
		z, err = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse SRSOrigin Z: %w", err)
		}
	}

	return x, y, z, nil
}

func ParseENUOrigin(srs string) (lon, lat float64, err error) {
	parts := strings.Split(srs, ":")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid ENU format: %s", srs)
	}

	coords := strings.Split(parts[1], ",")
	if len(coords) < 2 {
		return 0, 0, fmt.Errorf("invalid ENU coordinates: %s", parts[1])
	}

	lon, err = strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse ENU lon: %w", err)
	}

	lat, err = strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse ENU lat: %w", err)
	}

	return lon, lat, nil
}

func ParseEPSGCode(srs string) (int, error) {
	parts := strings.Split(srs, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid EPSG format: %s", srs)
	}

	code, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, fmt.Errorf("failed to parse EPSG code: %w", err)
	}

	return code, nil
}

func FindMetadataFile(basePath string) (string, error) {
	paths := []string{
		filepath.Join(basePath, "metadata.xml"),
		filepath.Join(basePath, "OSGB", "metadata.xml"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("metadata.xml not found in %s", basePath)
}
