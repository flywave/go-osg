package io

import (
	"strings"

	"github.com/flywave/go-osg/model"
)

type ObjectWrapperManager struct {
	Wraps map[string]*ObjectWrapper
}

var manager *ObjectWrapperManager

func GetObjectWrapperManager() *ObjectWrapperManager {
	return manager
}

func (man *ObjectWrapperManager) AddWrap(wrap *ObjectWrapper) {
	if wrap == nil {
		return
	}
	w, ok := manager.Wraps[wrap.Name]
	manager.Wraps[strings.ToLower(wrap.Name)] = wrap
}

func init() {
	manager = &ObjectWrapperManager{Wraps: make(map[string]*ObjectWrapper)}
}

type ObjectWrapperAssociate struct {
	FirstVersion int
	LastVersion  int
	Name         string
}

type CreateInstanceFuncType func() *model.Object

type ObjectWrapper struct {
	CreateInstanceFunc                   CreateInstanceFuncType
	Domain                               string
	Name                                 string
	Associates                           []*ObjectWrapperAssociate
	Serializers                          []*BaseSerializer
	BackupSerializers                    []*BaseSerializer
	TypeList                             []SerType
	Version                              int
	IsAssociatesRevisionsInheritanceDone bool
}

func NewObjectWrapper() ObjectWrapper {
	return ObjectWrapper{}
}

func (wp *ObjectWrapper) CreateInstance() *model.Object {
	return wp.CreateInstanceFunc()
}

func (wp *ObjectWrapper) AddSerializer(s *BaseSerializer, t SerType) {
	s.FirstVersion = wp.Version
	wp.Serializers = append(wp.Serializers, s)
	wp.TypeList = append(wp.TypeList, t)
}

func (wp *ObjectWrapper) MarkSerializerAsRemoved(name string) {
	for _, s := range wp.Serializers {
		ser := Serializer(s)
		if ser.GetSerializerName() == name {
			s.LastVersion = wp.Version - 1
		}
	}
}

func (wp *ObjectWrapper) GetSerializer(name string) *BaseSerializer {
	for _, s := range wp.Serializers {
		ser := Serializer(s)
		if ser.GetSerializerName() == name {
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager()[as.Name]
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := Serializer(s)
			if ser.GetSerializerName() == name {
				return s
			}
		}
	}
	return nil
}

func (wp *ObjectWrapper) GetSerializerAndType(name string, ty *SerType) *BaseSerializer {
	for i, s := range wp.Serializers {
		ser := Serializer(s)
		if ser.GetSerializerName() == name {
			*ty = wp.TypeList[i]
			return s
		}
	}

	for _, as := range wp.Associates {
		w := GetObjectWrapperManager()[as.Name]
		if w == nil {
			continue
		}
		for _, s := range w.Serializers {
			ser := Serializer(s)
			if ser.GetSerializerName() == name {
				*ty = w.TypeList[0]
				return s
			}
		}
	}
	*ty = RW_UNDEFINED
	return nil
}

func (wp *ObjectWrapper) Read(is *OsgIstream, obj *model.Object) {
	inputVersion := is.GetFileVersion(wp.Domain)
	for _, ser := range wp.Serializers {
		if ser.FirstVersion <= inputVersion &&
			inputVersion <= ser.LastVersion && ser.SupportsGetSet() {
			s := Serializer(ser)
			s.Read(is, obj)
		}
	}
}

func (wp *ObjectWrapper) Write(os *OsgOstream, obj *model.Object) {
	inputVersion := os.GetFileVersion(wp.Domain)
	for _, ser := range wp.Serializers {
		if ser.FirstVersion <= inputVersion &&
			inputVersion <= ser.LastVersion && ser.SupportsGetSet() {
			s := Serializer(ser)
			s.Write(os, obj)
		}
	}
}

func (wp *ObjectWrapper) ReadSchema(properties []string, types []SerType) bool {
	if len(wp.BackupSerializers) != 0 {
		wp.BackupSerializers = wp.Serializers
	}
	wp.Serializers = wp.Serializers[0:0]
	size := len(properties)
	serializersSize := len(wp.BackupSerializers)

	for i := 0; i < size; i++ {
		if serializersSize < i {
			break
		}
		prop := properties[i]
		if prop == wp.BackupSerializers[i].GetSerializerName() {
			wp.Serializers = append(wp.Serializers, wp.BackupSerializers[i])
		} else {
			for _, ser := range wp.Serializers {
				if prop != ser.GetSerializerName() {
					continue
				}
				wp.Serializers = append(wp.Serializers, ser)
			}
		}
	}
	return size == len(wp.Serializers)
}

func (wp *ObjectWrapper) WriteSchema(properties []string, types []SerType) {
	ssize := len(wp.Serializers)
	tsize := len(wp.TypeList)
	i := 0
	for {
		if ssize-1 <= i || tsize-1 <= i {
			break
		}
		ser := wp.Serializers[i]
		t := wp.TypeList[i]
		if ser.SupportsGetSet() {
			properties = append(properties, ser.GetSerializerName())
			types = append(types, t)
		}
		i++
	}
}
