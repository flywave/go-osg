package tiles3d

import (
	"github.com/flywave/go-geoid"
)

type GeoidConverter struct {
	model string
	geoid *geoid.Geoid
}

func NewGeoidConverter(model, dataPath string) *GeoidConverter {
	g := &GeoidConverter{
		model: model,
	}

	if model != "none" && model != "" {
		verticalDatum := geoid.VerticalDatumFromString(model)
		g.geoid = geoid.NewGeoid(verticalDatum, false)
		if dataPath != "" {
			geoid.SetGeoidPath(dataPath)
		}
	}

	return g
}

func (g *GeoidConverter) ConvertOrthometricToEllipsoidal(lat, lon, orthometricHeight float64) float64 {
	if g.geoid == nil {
		return orthometricHeight
	}

	geoidHeight := g.geoid.GetHeight(lat, lon)
	return orthometricHeight + geoidHeight
}

func (g *GeoidConverter) ConvertEllipsoidalToOrthometric(lat, lon, ellipsoidalHeight float64) float64 {
	if g.geoid == nil {
		return ellipsoidalHeight
	}

	geoidHeight := g.geoid.GetHeight(lat, lon)
	return ellipsoidalHeight - geoidHeight
}
