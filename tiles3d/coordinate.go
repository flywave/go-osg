package tiles3d

import (
	"fmt"
	"math"

	"github.com/flywave/go-proj"
)

type CoordinateTransformer struct {
	sourceSRS    string
	targetSRS    string
	srsType      SRSType
	center       [3]float64
	originOffset [3]float64
	sourceProj   *proj.Proj
	targetProj   *proj.Proj
	geoidConv    *GeoidConverter
}

func NewCoordinateTransformer(sourceSRS, targetSRS string) *CoordinateTransformer {
	fmt.Printf("DEBUG NewCoordinateTransformer: sourceSRS=%q, targetSRS=%q\n", sourceSRS, targetSRS)
	t := &CoordinateTransformer{
		sourceSRS: sourceSRS,
		targetSRS: targetSRS,
		srsType:   DetectSRSType(sourceSRS),
	}

	if sourceSRS != "" && t.srsType != SRSTypeENU {
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

func (t *CoordinateTransformer) SetGeoidConverter(geoidConv *GeoidConverter) {
	t.geoidConv = geoidConv
}

func (t *CoordinateTransformer) SetCenter(x, y, height float64) {
	t.SetGeoReference(x, y, height)
}

func (t *CoordinateTransformer) SetGeoReference(lon, lat, height float64) {
	fmt.Printf("DEBUG SetGeoReference: srsType=%d, lon=%f, lat=%f, h=%f\n", t.srsType, lon, lat, height)
	t.center = [3]float64{lon, lat, height}
}

func (t *CoordinateTransformer) SetGeoReferenceFromMetadata(metadata *ModelMetadata, projLib string) error {
	srsType := DetectSRSType(metadata.SRS)
	t.srsType = srsType

	originX, originY, originZ, err := ParseSRSOrigin(metadata.SRSOrigin)
	if err != nil {
		return fmt.Errorf("failed to parse SRSOrigin: %w", err)
	}

	t.originOffset = [3]float64{originX, originY, originZ}
	fmt.Printf("DEBUG SetGeoReferenceFromMetadata: srsType=%d, origin=(%f, %f, %f)\n", srsType, originX, originY, originZ)

	switch srsType {
	case SRSTypeENU:
		return t.setupENUMode(metadata.SRS)
	case SRSTypeEPSG:
		return t.setupEPSGMode(metadata.SRS, originX, originY, originZ)
	case SRSTypeWKT:
		return t.setupWKTMode(metadata.SRS, originX, originY, originZ)
	case SRSTypeUnknown:
		return t.setupUnknownMode(originX, originY, originZ)
	default:
		return fmt.Errorf("unknown SRS type: %s", metadata.SRS)
	}
}

func (t *CoordinateTransformer) setupENUMode(srs string) error {
	lon, lat, err := ParseENUOrigin(srs)
	if err != nil {
		return fmt.Errorf("failed to parse ENU origin: %w", err)
	}

	t.center = [3]float64{lon, lat, t.originOffset[2]}
	fmt.Printf("DEBUG setupENUMode: center=(%f, %f, %f)\n", t.center[0], t.center[1], t.center[2])

	if t.geoidConv != nil {
		correctedHeight := t.geoidConv.ConvertOrthometricToEllipsoidal(t.center[1], t.center[0], t.center[2])
		t.center[2] = correctedHeight
		fmt.Printf("DEBUG setupENUMode: geoid corrected height = %f\n", t.center[2])
	}

	return nil
}

func (t *CoordinateTransformer) setupEPSGMode(srs string, originX, originY, originZ float64) error {
	code, err := ParseEPSGCode(srs)
	if err != nil {
		return fmt.Errorf("failed to parse EPSG code: %w", err)
	}

	epsgStr := fmt.Sprintf("EPSG:%d", code)
	src, err := proj.NewProj(epsgStr)
	if err != nil {
		return fmt.Errorf("failed to create EPSG:%d projection: %w", code, err)
	}
	t.sourceProj = src

	wgs84, err := proj.NewProj("EPSG:4326")
	if err != nil {
		return fmt.Errorf("failed to create WGS84 projection: %w", err)
	}

	lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
	if err != nil {
		fmt.Printf("DEBUG setupEPSGMode: proj.Transform3 error, using raw values: %v\n", err)
		lon = originX
		lat = originY
		z = originZ
	} else {
		// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
		// We need to swap to match C++ behavior
		lon, lat = lat, lon
		lon = lon * 180.0 / math.Pi
		lat = lat * 180.0 / math.Pi
	}

	t.center = [3]float64{lon, lat, z}

	if t.geoidConv != nil {
		correctedHeight := t.geoidConv.ConvertOrthometricToEllipsoidal(t.center[1], t.center[0], t.center[2])
		t.center[2] = correctedHeight
		fmt.Printf("DEBUG setupEPSGMode: geoid corrected height = %f\n", t.center[2])
	}

	fmt.Printf("DEBUG setupEPSGMode: center=(%f, %f, %f)\n", t.center[0], t.center[1], t.center[2])

	return nil
}

func (t *CoordinateTransformer) setupWKTMode(srs string, originX, originY, originZ float64) error {
	src, err := proj.NewProj(srs)
	if err != nil {
		return fmt.Errorf("failed to create projection from WKT: %w", err)
	}
	t.sourceProj = src

	wgs84, err := proj.NewProj("EPSG:4326")
	if err != nil {
		return fmt.Errorf("failed to create WGS84 projection: %w", err)
	}

	lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
	if err != nil {
		fmt.Printf("DEBUG setupWKTMode: proj.Transform3 error, using raw values: %v\n", err)
		lon = originX
		lat = originY
		z = originZ
	} else {
		// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
		// We need to swap to match C++ behavior
		lon, lat = lat, lon
		lon = lon * 180.0 / math.Pi
		lat = lat * 180.0 / math.Pi
	}

	t.center = [3]float64{lon, lat, z}

	fmt.Printf("DEBUG setupWKTMode: center=(%f, %f, %f)\n", t.center[0], t.center[1], t.center[2])

	return nil
}

func (t *CoordinateTransformer) setupUnknownMode(originX, originY, originZ float64) error {
	t.originOffset = [3]float64{originX, originY, originZ}

	commonProjections := []string{
		"EPSG:4548",
		"EPSG:4547",
		"EPSG:4549",
		"EPSG:4490",
		"EPSG:4326",
	}

	for _, projStr := range commonProjections {
		src, err := proj.NewProj(projStr)
		if err == nil {
			wgs84, err := proj.NewProj("EPSG:4326")
			if err == nil {
				lon, lat, z, err := proj.Transform3(src, wgs84, originX, originY, originZ)
				if err == nil {
					// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
					// We need to swap to match C++ behavior
					lon, lat = lat, lon
					t.center = [3]float64{lon * 180.0 / math.Pi, lat * 180.0 / math.Pi, z}
					t.sourceProj = src
					fmt.Printf("DEBUG setupUnknownMode: detected projection %s, center=(%f, %f, %f)\n",
						projStr, t.center[0], t.center[1], t.center[2])
					return nil
				}
			}
		}
	}

	t.center = [3]float64{originX, originY, originZ}
	fmt.Printf("DEBUG setupUnknownMode: center=(%f, %f, %f)\n", t.center[0], t.center[1], t.center[2])
	return nil
}

func (t *CoordinateTransformer) GetSRSType() SRSType {
	return t.srsType
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
				// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
				// We need to swap to match C++ behavior
				return [3]float64{y, x, z}
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
	switch t.srsType {
	case SRSTypeENU:
		return t.toLocalENUFromENU(point)
	case SRSTypeEPSG, SRSTypeWKT:
		return t.toLocalENUFromProjected(point)
	case SRSTypeUnknown:
		return t.toLocalENUFromUnknown(point)
	default:
		if t.sourceProj == nil {
			return t.ToLocalENU(point)
		}
		ecefPoint := t.ToECEFFromSource(point)
		return t.ToLocalENUFromECEF(ecefPoint)
	}
}

func (t *CoordinateTransformer) toLocalENUFromENU(point [3]float64) [3]float64 {
	absX := point[0] + t.originOffset[0]
	absY := point[1] + t.originOffset[1]
	absZ := point[2] + t.originOffset[2]

	ecef := t.ToECEFFromLatLon(absY*math.Pi/180.0, absX*math.Pi/180.0, absZ)

	return t.ToLocalENUFromECEF(ecef)
}

func (t *CoordinateTransformer) toLocalENUFromProjected(point [3]float64) [3]float64 {
	absX := point[0]
	absY := point[1]
	absZ := point[2]

	fmt.Printf("DEBUG toLocalENUFromProjected: point=(%f, %f, %f), offset=(%f, %f, %f), center=(%f, %f, %f)\n",
		point[0], point[1], point[2], t.originOffset[0], t.originOffset[1], t.originOffset[2],
		t.center[0], t.center[1], t.center[2])

	if t.sourceProj == nil {
		return [3]float64{absX - t.originOffset[0], absY - t.originOffset[1], absZ}
	}

	wgs84, err := proj.NewProj("EPSG:4326")
	if err != nil {
		return [3]float64{absX - t.originOffset[0], absY - t.originOffset[1], absZ}
	}

	lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
	fmt.Printf("DEBUG toLocalENUFromProjected: proj.Transform3 result (raw): lon=%f(rad), lat=%f(rad), z=%f\n", lon, lat, z)
	if err != nil {
		fmt.Printf("DEBUG toLocalENUFromProjected: transform error: %v\n", err)
		return [3]float64{absX, absY, absZ}
	}

	// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
	// We need to swap to match C++ behavior
	lon, lat = lat, lon

	lonDeg := lon * 180.0 / math.Pi
	latDeg := lat * 180.0 / math.Pi
	fmt.Printf("DEBUG toLocalENUFromProjected: WGS84 (after swap): lon=%f(deg), lat=%f(deg)\n", lonDeg, latDeg)

	if t.geoidConv != nil {
		z = t.geoidConv.ConvertOrthometricToEllipsoidal(latDeg, lonDeg, z)
	}

	ecef := t.ToECEFFromLatLon(lat, lon, z)
	fmt.Printf("DEBUG toLocalENUFromProjected: ECEF=(%f, %f, %f)\n", ecef[0], ecef[1], ecef[2])

	enu := t.ToLocalENUFromECEF(ecef)
	fmt.Printf("DEBUG toLocalENUFromProjected: ENU=(%f, %f, %f)\n", enu[0], enu[1], enu[2])
	return enu
}

func (t *CoordinateTransformer) toLocalENUFromUnknown(point [3]float64) [3]float64 {
	absX := point[0] + t.originOffset[0]
	absY := point[1] + t.originOffset[1]
	absZ := point[2] + t.originOffset[2]

	fmt.Printf("DEBUG toLocalENUFromUnknown: abs=(%f, %f, %f), center=(%f, %f, %f)\n",
		absX, absY, absZ, t.center[0], t.center[1], t.center[2])

	if t.sourceProj != nil {
		wgs84, err := proj.NewProj("EPSG:4326")
		if err == nil {
			lon, lat, z, err := proj.Transform3(t.sourceProj, wgs84, absX, absY, absZ)
			fmt.Printf("DEBUG toLocalENUFromUnknown: proj.Transform3 result (raw): lon=%f(rad), lat=%f(rad), z=%f\n", lon, lat, z)
			if err == nil {
				// FIX: go-proj may return (lat, lon) in OGC order, while C++ OGR uses (lon, lat) in GIS order
				// We need to swap to match C++ behavior
				lon, lat = lat, lon

				lonDeg := lon * 180.0 / math.Pi
				latDeg := lat * 180.0 / math.Pi
				fmt.Printf("DEBUG toLocalENUFromUnknown: WGS84 (after swap): lon=%f(deg), lat=%f(deg)\n", lonDeg, latDeg)

				if t.geoidConv != nil {
					z = t.geoidConv.ConvertOrthometricToEllipsoidal(latDeg, lonDeg, z)
				}

				ecef := t.ToECEFFromLatLon(lat, lon, z)
				fmt.Printf("DEBUG toLocalENUFromUnknown: ECEF=(%f, %f, %f)\n", ecef[0], ecef[1], ecef[2])

				enu := t.ToLocalENUFromECEF(ecef)
				fmt.Printf("DEBUG toLocalENUFromUnknown: ENU=(%f, %f, %f)\n", enu[0], enu[1], enu[2])
				return enu
			}
		}
	}

	return [3]float64{absX - t.originOffset[0], absY - t.originOffset[1], absZ - t.originOffset[2]}
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

func (t *CoordinateTransformer) HasGeoReference() bool {
	return t.center[0] != 0 || t.center[1] != 0
}
