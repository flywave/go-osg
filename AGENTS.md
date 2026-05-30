# go-osg — OpenSceneGraph format reader/writer in Go

## Project structure

- **root package** (`osg`) — core OSG format R/W: `.osgb` (binary), `.osgt` (ASCII), `.osg2` (extended).
- **`model/`** — data types: Node, Group, Geometry, Array, PagedLOD, StateSet, etc.
- **`tiles3d/`** — OSGB → 3D Tiles converter (B3DM / GLB / tileset.json). Depends on `go-proj`, `go-draco`, `go-meshopt`, `gltf`, `go-geoid`, `go-3dtile` via `replace` directives pointing to sibling `../` directories.
- **`serialization_test/`** — external test package (imports `osg` and `model`).
- **`cmd/check_vertices/`** — diagnostic tool for inspecting binary vertex data.

## Commands

| Action | Command |
|--------|---------|
| Run all tests | `go test ./...` |
| Run root-package tests only | `go test` |
| Run serialization tests | `go test ./serialization_test/...` |
| Build | `go build ./...` |

Tests currently pass in ~0.25s (91 tests, though some crash on missing `test_data/OSGB/` files).

## Critical gotchas

- **`go.mod` `replace` directives**: `go-proj`, `go-draco`, `go-meshopt`, `gltf`, `go-geoid`, `go-3dtile` all point to `../<repo>`. Full build requires those sibling repos checked out. For isolated work on the core OSG reader, only `model/` and root tests are safe.
- **Missing test data**: Some tests reference `test_data/OSGB/Data/*` and `test_data/0131/Data/*` which don't exist. Focus on tests that use `test_data/cessna.osgb`, `test_data/simpleroom.osgt`, `test_data/skydome.osgt`, and `test_data/Tile_+003_+003_L18_000.osgb`.
- **`tiles3d` coordinate transform** needs `PROJ_DATA` / `PROJ_LIB` env vars set to a proj-data directory. The example in `tiles3d_test/main.go` has a hardcoded macOS path — adjust for the current machine.

## Format version quirks (OSG binary)

File version thresholds baked into the reader:

| Version | Behaviour |
|---------|-----------|
| ≥ 112   | Arrays read as objects (via `ReadObject`, has wrapper/classname) |
| > 96    | PrimitiveSet includes `numInstances` field |
| > 94    | Image includes `ClassName` property |
| > 148   | Binary bracket block size uses 8-byte `uint64`; ≤ 148 uses 4-byte `int32` |

See `input_stream.go` and `binary_input.go` for implementation. The reference C++ is in `OpenSceneGraph-master/src/osgPlugins/osg/`.

## Entry point

```go
rw := osg.NewReadWrite()
res := rw.ReadNode("file.osgb", nil)
node := res.GetNode() // model.NodeInterface
```

## Key convention

All model types defined in `model/`. Serializers and reader/writer logic live in the root `osg` package. Registration of object wrappers happens in `init()` functions across many root-package files.
