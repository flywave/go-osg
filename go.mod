module github.com/flywave/go-osg

require (
	github.com/flywave/gltf v0.20.4-0.20250703101259-1edf290bd287
	github.com/flywave/go-3dtile v0.0.0-00010101000000-000000000000
	github.com/flywave/go-draco v0.0.0-00010101000000-000000000000
	github.com/flywave/go-geoid v0.0.0-00010101000000-000000000000
	github.com/flywave/go-meshopt v0.0.0-00010101000000-000000000000
	github.com/flywave/go-proj v0.0.0-00010101000000-000000000000
	github.com/flywave/go3d v0.0.0-20250816053852-aed5d825659f
)

go 1.24

replace github.com/flywave/go-proj => ../go-proj

replace github.com/flywave/go-draco => ../go-draco

replace github.com/flywave/go-meshopt => ../go-meshopt

replace github.com/flywave/gltf => ../gltf

replace github.com/flywave/go-geoid => ../go-geoid

replace github.com/flywave/go-3dtile => ../go-3dtile
