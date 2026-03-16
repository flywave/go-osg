package osg

import (
	"testing"

	"github.com/flywave/go-osg/model"
)

func TestUserSerializer_Read(t *testing.T) {
	lk := NewIntLookup()
	lk.Add("TEST", 1)

	ser := NewUserSerializer("TestUser", nil, func(is *OsgIstream, obj interface{}) {
	}, func(os *OsgOstream, obj interface{}) {
	})

	if ser.GetSerializerName() != "TestUser" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestUser")
	}
}

func TestPropByValSerializer(t *testing.T) {
	ser := NewPropByValSerializer("TestProp", false, func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestProp" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestProp")
	}
}

func TestMatrixSerializer(t *testing.T) {
	ser := NewMatrixSerializer("TestMatrix", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestMatrix" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestMatrix")
	}
}

func TestGlenumSerializer(t *testing.T) {
	ser := NewGlenumSerializer("TestGlenum", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestGlenum" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestGlenum")
	}
}

func TestStringSerializer(t *testing.T) {
	ser := NewStringSerializer("TestString", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestString" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestString")
	}
}

func TestObjectSerializer(t *testing.T) {
	ser := NewObjectSerializer("TestObject", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestObject" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestObject")
	}
}

func TestImageSerializer(t *testing.T) {
	ser := NewImageSerializer("TestImage", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestImage" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestImage")
	}
}

func TestEnumSerializer(t *testing.T) {
	ser := NewEnumSerializer("TestEnum", func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestEnum" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestEnum")
	}

	ser.Add("VALUE1", 1)
	ser.Add("VALUE2", 2)

	if ser.LookUp.GetValue("VALUE1") != 1 {
		t.Error("EnumSerializer.Add() failed")
	}
}

func TestVectorSerializer(t *testing.T) {
	ser := NewVectorSerializer("TestVector", RWVECTOR, &model.Array{}, func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestVector" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestVector")
	}
}

func TestIsAVectorSerializer(t *testing.T) {
	ser := NewIsAVectorSerializer("TestIsAVector", RWINT, 4, func(obj interface{}) interface{} {
		return nil
	}, func(obj interface{}, val interface{}) {
	})

	if ser.GetSerializerName() != "TestIsAVector" {
		t.Errorf("GetSerializerName() = %q, want %q", ser.GetSerializerName(), "TestIsAVector")
	}
}

func TestSerTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    SerType
		expected SerType
	}{
		{"RWUNDEFINED", RWUNDEFINED, 0},
		{"RWUSER", RWUSER, 1},
		{"RWOBJECT", RWOBJECT, 2},
		{"RWIMAGE", RWIMAGE, 3},
		{"RWLIST", RWLIST, 4},
		{"RWBOOL", RWBOOL, 5},
		{"RWCHAR", RWCHAR, 6},
		{"RWUCHAR", RWUCHAR, 7},
		{"RWSHORT", RWSHORT, 8},
		{"RWUSHORT", RWUSHORT, 9},
		{"RWINT", RWINT, 10},
		{"RWUINT", RWUINT, 11},
		{"RWFLOAT", RWFLOAT, 12},
		{"RWDOUBLE", RWDOUBLE, 13},
		{"RWVEC2F", RWVEC2F, 14},
		{"RWVEC2D", RWVEC2D, 15},
		{"RWVEC3F", RWVEC3F, 16},
		{"RWVEC3D", RWVEC3D, 17},
		{"RWVEC4F", RWVEC4F, 18},
		{"RWVEC4D", RWVEC4D, 19},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestUsageConstants(t *testing.T) {
	if READWRITEPROPERTY != 1 {
		t.Errorf("READWRITEPROPERTY = %d, want 1", READWRITEPROPERTY)
	}
	if GETPROPERTY != 2 {
		t.Errorf("GETPROPERTY = %d, want 2", GETPROPERTY)
	}
	if SETPROPERTY != 4 {
		t.Errorf("SETPROPERTY = %d, want 4", SETPROPERTY)
	}
	if GETSETPROPERTY != (GETPROPERTY | SETPROPERTY) {
		t.Errorf("GETSETPROPERTY = %d, want %d", GETSETPROPERTY, GETPROPERTY|SETPROPERTY)
	}
}

func TestBaseSerializer_SupportsReadWrite(t *testing.T) {
	ser := NewBaseSerializer(READWRITEPROPERTY)
	if !ser.SupportsReadWrite() {
		t.Error("SupportsReadWrite() should be true for READWRITEPROPERTY")
	}

	ser2 := NewBaseSerializer(GETPROPERTY)
	if ser2.SupportsReadWrite() {
		t.Error("SupportsReadWrite() should be false for GETPROPERTY")
	}
}

func TestBaseSerializer_SupportsGetSet(t *testing.T) {
	ser := NewBaseSerializer(GETSETPROPERTY)
	if !ser.SupportsGetSet() {
		t.Error("SupportsGetSet() should be true for GETSETPROPERTY")
	}

	ser2 := NewBaseSerializer(READWRITEPROPERTY)
	if ser2.SupportsGetSet() {
		t.Error("SupportsGetSet() should be false for READWRITEPROPERTY")
	}
}
