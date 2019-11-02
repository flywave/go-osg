package model

const (
	UserDataContainerType        string = "osg::UserDataContainer"
	DefaultUserDataContainerType string = "osg::DefaultUserDataContainer"
)

type UserDataContainer struct {
	Object
	User_data       interface{}
	DescriptionList []string
	ObjectList      []*Object
}

func NewUserDataContainer() UserDataContainer {
	obj := NewObject()
	obj.ObjectType = UserDataContainerType
	return UserDataContainer{Object: obj}
}

type DefaultUserDataContainer struct {
	UserDataContainer
}

func NewDefaultUserDataContainer() DefaultUserDataContainer {
	obj := NewUserDataContainer()
	obj.ObjectType = DefaultUserDataContainerType
	return UserDataContainer{UserDataContainer: obj}
}
