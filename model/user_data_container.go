package model

const (
	USERDATACONTAINER_T        string = "osg::UserDataContainer"
	DEFAULTUSERDATACONTAINER_T string = "osg::DefaultUserDataContainer"
)

type UserDataContainer struct {
	Object
	User_data       interface{}
	DescriptionList []string
	ObjectList      []*Object
}

func NewUserDataContainer() UserDataContainer {
	obj := NewObject()
	obj.Type = USERDATACONTAINER_T
	return UserDataContainer{Object: obj}
}

type DefaultUserDataContainer struct {
	UserDataContainer
}

func NewDefaultUserDataContainer() DefaultUserDataContainer {
	obj := NewUserDataContainer()
	obj.Type = DEFAULTUSERDATACONTAINER_T
	return DefaultUserDataContainer{UserDataContainer: obj}
}
