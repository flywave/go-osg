package osg

import (
	"testing"
)

func TestObjectWrapper_MarkSerializerAsAdded(t *testing.T) {
	wp := NewObjectWrapper("TestObject", nil, "Parent1")
	wp.Version = 5

	wp.MarkSerializerAsAdded("Parent1")

	for _, as := range wp.Associates {
		if as.Name == "Parent1" {
			if as.FirstVersion != 5 {
				t.Errorf("MarkSerializerAsAdded() FirstVersion = %d, want 5", as.FirstVersion)
			}
			return
		}
	}
	t.Error("Associate Parent1 not found")
}

func TestObjectWrapper_GetSerializer(t *testing.T) {
	wp := NewObjectWrapper("TestObject", nil, "")

	ser := NewBaseSerializer(READWRITEPROPERTY)
	ser2 := &UserSerializer{BaseSerializer: *ser, Name: "TestSerializer"}
	wp.AddSerializer(ser2, RWUSER)

	found := wp.GetSerializer("TestSerializer")
	if found == nil {
		t.Error("GetSerializer() returned nil")
	}
}

func TestObjectWrapper_GetSerializer_NotFound(t *testing.T) {
	wp := NewObjectWrapper("TestObject", nil, "")

	found := wp.GetSerializer("NonExistent")
	if found != nil {
		t.Error("GetSerializer() should return nil for non-existent serializer")
	}
}

func TestObjectWrapper_AddSerializer(t *testing.T) {
	wp := NewObjectWrapper("TestObject", nil, "")

	ser := NewBaseSerializer(READWRITEPROPERTY)
	us := &UserSerializer{BaseSerializer: *ser, Name: "TestSer"}

	initialCount := len(wp.Serializers)
	wp.AddSerializer(us, RWUSER)

	if len(wp.Serializers) != initialCount+1 {
		t.Error("AddSerializer() failed to add serializer")
	}

	if len(wp.TypeList) != initialCount+1 {
		t.Error("AddSerializer() failed to add type")
	}
}

func TestNewRegisterCustomWrapperProxy(t *testing.T) {
	NewRegisterCustomWrapperProxy(func() interface{} { return "test" }, "TestDomain", "TestWrapper", "Parent1 Parent2")

	wrap := GetObjectWrapperManager().FindWrap("osg::testwrapper")
	if wrap == nil {
		t.Error("NewRegisterCustomWrapperProxy() failed to register wrapper")
	}
}

func TestUpdateWrapperVersionProxy(t *testing.T) {
	wp := NewObjectWrapper("TestObject", nil, "")
	wp.Version = 10

	proxy := AddUpdateWrapperVersionProxy(wp, 20)

	if wp.Version != 20 {
		t.Errorf("Version = %d, want 20", wp.Version)
	}

	proxy.SetLastVersion()

	if wp.Version != 10 {
		t.Errorf("Version after SetLastVersion = %d, want 10", wp.Version)
	}
}

func TestObjectWrapperManager_RemoveWrap(t *testing.T) {
	manager := GetObjectWrapperManager()

	wp := NewObjectWrapper("TestRemoveWrap", nil, "")
	manager.AddWrap(wp)

	found := manager.FindWrap("osg::testremovewrap")
	if found == nil {
		t.Error("AddWrap() failed")
	}

	manager.RemoveWrap(wp)
	found = manager.FindWrap("osg::testremovewrap")
	if found != nil {
		t.Error("RemoveWrap() failed")
	}
}

func TestNewObjectWrapper2(t *testing.T) {
	wp := NewObjectWrapper2("TestName", "TestDomain", func() interface{} { return nil }, "Parent1")

	if wp.Name != "TestName" {
		t.Errorf("Name = %q, want %q", wp.Name, "TestName")
	}

	if wp.Domain != "TestDomain" {
		t.Errorf("Domain = %q, want %q", wp.Domain, "TestDomain")
	}
}
