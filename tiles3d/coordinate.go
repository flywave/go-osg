package tiles3d

import (
	"math"

	"github.com/flywave/go-proj"
)

type CoordinateTransformer struct {
	sourceSRS  string
	targetSRS  string
	center     [3]float64
	sourceProj *proj.Proj
	targetProj *proj.Proj
}

func NewCoordinateTransformer(sourceSRS, targetSRS string) *CoordinateTransformer {
	t := &CoordinateTransformer{
		sourceSRS: sourceSRS,
		targetSRS: targetSRS,
	}

	if sourceSRS != "" && targetSRS != "" {
		src, err := proj.NewProj(sourceSRS)
		if err == nil {
			t.sourceProj = src
		}
		dst, err := proj.NewProj(targetSRS)
		if err == nil {
			t.targetProj = dst
		}
	}

	return t
}

func (t *CoordinateTransformer) SetCenter(lon, lat, height float64) {
	t.center = [3]float64{lon, lat, height}
}

func (t *CoordinateTransformer) Transform(point [3]float64) [3]float64 {
	if t.sourceProj != nil && t.targetProj != nil {
		x, y, z, err := proj.Transform3(t.sourceProj, t.targetProj, point[0], point[1], point[2])
		if err == nil {
			return [3]float64{x, y, z}
		}
	}
	return point
}

func (t *CoordinateTransformer) ToWGS84(point [3]float64) [3]float64 {
	if t.targetSRS != "" {
		return t.Transform(point)
	}
	return point
}

func (t *CoordinateTransformer) ToLocalENU(point [3]float64) [3]float64 {
	lon := t.center[0]
	lat := t.center[1]
	height := t.center[2]

	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	sinLat := math.Sin(latRad)
	cosLat := math.Cos(latRad)
	sinLon := math.Sin(lonRad)
	cosLon := math.Cos(lonRad)

	xEnu := -sinLon*(point[0]-lon) + cosLon*(point[1]-lat)
	yEnu := -sinLat*cosLon*(point[0]-lon) - sinLat*sinLon*(point[1]-lat) + cosLat*(point[2]-height)
	zEnu := cosLat*cosLon*(point[0]-lon) + cosLat*sinLon*(point[1]-lat) + sinLat*(point[2]-height)

	return [3]float64{xEnu, yEnu, zEnu}
}

func (t *CoordinateTransformer) ToECEF(point [3]float64) [3]float64 {
	lat := point[1] * math.Pi / 180.0
	lon := point[0] * math.Pi / 180.0
	height := point[2]

	a := 6378137.0
	e2 := 0.00669437999014

	N := a / math.Sqrt(1-e2*math.Sin(lat)*math.Sin(lat))

	x := (N + height) * math.Cos(lat) * math.Cos(lon)
	y := (N + height) * math.Cos(lat) * math.Sin(lon)
	z := (N*(1-e2) + height) * math.Sin(lat)

	return [3]float64{x, y, z}
}

func (t *CoordinateTransformer) ToECEFFromLatLon(latRad, lonRad, height float64) [3]float64 {
	a := 6378137.0
	e2 := 0.00669437999014

	N := a / math.Sqrt(1-e2*math.Sin(latRad)*math.Sin(latRad))

	x := (N + height) * math.Cos(latRad) * math.Cos(lonRad)
	y := (N + height) * math.Cos(latRad) * math.Sin(lonRad)
	z := (N*(1-e2) + height) * math.Sin(latRad)

	return [3]float64{x, y, z}
}

func (t *CoordinateTransformer) GetCenter() [3]float64 {
	return t.center
}

func (t *CoordinateTransformer) HasProjection() bool {
	return t.sourceProj != nil && t.targetProj != nil
}

func (t *CoordinateTransformer) IsGeographicOutput() bool {
	if t.targetSRS != "" {
		if t.targetSRS == "EPSG:4326" || t.targetSRS == "EPSG:4978" {
			return true
		}
	}
	return false
}

func (t *CoordinateTransformer) IsECEFOutput() bool {
	if t.targetSRS != "" {
		return t.targetSRS == "EPSG:4978"
	}
	return false
}

func (t *CoordinateTransformer) GetTargetSRS() string {
	return t.targetSRS
}
