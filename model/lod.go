package model

// _"github.com/ungerik/go3d/vec3"

type CenterMode uint32
type RangeMode uint32

const (
	USE_BOUNDING_SPHERE_CENTER                CenterMode = 0
	USER_DEFINED_CENTER                       CenterMode = 1
	UNION_OF_BOUNDING_SPHERE_AND_USER_DEFINED CenterMode = 2
	DISTANCE_FROM_EYE_POINT                   RangeMode  = 0
	PIXEL_SIZE_ON_SCREEN                      RangeMode  = 1
	LodType                                   string     = "osg::Object"
)

type MinMaxPair [2]float32

type RangeListType []MinMaxPair

type Lod struct {
	Group
	Cmode CenterMode
	// Center    vec3.T
	Radius    float
	Rmode     RangeMode
	RangeList RangeListType
}
