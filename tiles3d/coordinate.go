package tiles3d

import (
	"fmt"
	"math"

	"github.com/flywave/go-proj"
)

type CoordinateTransformer struct {
	sourceSRS    string
	targetSRS    string
	center       [3]float64
	originOffset [3]float64
	sourceProj   *proj.Proj
	targetProj   *proj.Proj
}

func NewCoordinateTransformer(sourceSRS, targetSRS string) *CoordinateTransformer {
	fmt.Printf("DEBUG NewCoordinateTransformer: sourceSRS=%q, targetSRS=%q\n", sourceSRS, targetSRS)
	t := &CoordinateTransformer{
		sourceSRS: sourceSRS,
		targetSRS: targetSRS,
	}

	if sourceSRS != "" {
		src, err := proj.NewProj(sourceSRS)
		fmt.Printf("DEBUG NewCoordinateTransformer: sourceProj created, err=%v\n", err)
		if err == nil {
			t.sourceProj = src
		}
	}
	if targetSRS != "" {
		dst, err := proj.NewProj(targetSRS)
		if err == nil {
			t.targetProj = dst
		}
	}

	return t
}

func (t *CoordinateTransformer) SetCenter(x, y, height float64) {
	// If sourceProj is set, assume input is in source SRS (e.g., EPSG:4548)
	// If no sourceProj, assume input is WGS84 lon/lat in degrees
	fmt.Printf("DEBUG SetCenter: sourceProj=%v, x=%f, y=%f, h=%f\n", t.sourceProj != nil, x, y, height)
	if t.sourceProj != nil {
		// Input is in source SRS (e.g., EPSG:4548), store as offset and convert to WGS84
		t.originOffset = [3]float64{x, y, height}
		wgs84, err := proj.NewProj("EPSG:4326")
		fmt.Printf("DEBUG SetCenter: wgs84 created, err=%v\n", err)
		if err == nil {
			lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, x, y, height)
			fmt.Printf("DEBUG SetCenter: transform result: lon=%f, lat=%f, z=%f, err=%v\n", lon, lat, z, err)
			if err == nil {
				// PROJ returns radians, convert to degrees
				t.center = [3]float64{lon * 180.0 / math.Pi, lat * 180.0 / math.Pi, z}
				fmt.Printf("DEBUG SetCenter: center converted to WGS84: (%f, %f, %f)\n", t.center[0], t.center[1], t.center[2])
				return
			}
		}
		t.center = [3]float64{x, y, height}
	} else {
		// No sourceProj - assume input is already WGS84 lon/lat in degrees
		t.originOffset = [3]float64{0, 0, 0}
		t.center = [3]float64{x, y, height}
	}
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
	if t.sourceProj != nil {
		wgs84, err := proj.NewProj("EPSG:4326")
		if err == nil {
			x, y, z, err := proj.Transform3(t.sourceProj, wgs84, point[0], point[1], point[2])
			if err == nil {
				return [3]float64{x, y, z}
			}
		}
	}
	return point
}

func (t *CoordinateTransformer) ToECEFFromSource(point [3]float64) [3]float64 {
	if t.sourceProj != nil {
		ecef, err := proj.NewProj("EPSG:4978")
		if err == nil {
			x, y, z, err := proj.Transform3(t.sourceProj, ecef, point[0], point[1], point[2])
			if err == nil {
				return [3]float64{x, y, z}
			}
		}
	}
	return point
}

func (t *CoordinateTransformer) ToLocalENUFromECEF(ecefPoint [3]float64) [3]float64 {
	lon := t.center[0]
	lat := t.center[1]
	height := t.center[2]

	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	sinLat := math.Sin(latRad)
	cosLat := math.Cos(latRad)
	sinLon := math.Sin(lonRad)
	cosLon := math.Cos(lonRad)

	originECEF := t.ToECEFFromLatLon(latRad, lonRad, height)

	dx := ecefPoint[0] - originECEF[0]
	dy := ecefPoint[1] - originECEF[1]
	dz := ecefPoint[2] - originECEF[2]

	xEnu := -sinLon*dx + cosLon*dy
	yEnu := -sinLat*cosLon*dx - sinLat*sinLon*dy + cosLat*dz
	zEnu := cosLat*cosLon*dx + cosLat*sinLon*dy + sinLat*dz

	return [3]float64{xEnu, yEnu, zEnu}
}

func (t *CoordinateTransformer) ToLocalENUFromSource(point [3]float64) [3]float64 {
	// If no sourceProj, assume input is WGS84 lon/lat in degrees
	if t.sourceProj == nil {
		// point is WGS84 lon/lat/height in degrees
		return t.ToLocalENU(point)
	}
	// Otherwise convert from source SRS to ECEF then to ENU
	ecefPoint := t.ToECEFFromSource(point)
	enuPoint := t.ToLocalENUFromECEF(ecefPoint)
	return enuPoint
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

func (t *CoordinateTransformer) GetOriginOffset() [3]float64 {
	return t.originOffset
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
