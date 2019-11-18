package model

const (
	USERDATACONTAINER_T        string = "osg::UserDataContainer"
	DEFAULTUSERDATACONTAINER_T string = "osg::DefaultUserDataContainer"
)

type UserDataContainer struct {
	DataVariance    int
	UserData        interface{}
	DescriptionList []string
	ObjectList      []interface{}
	Name            string
	Type            string
	Propertys       map[string]string
}

func NewUserDataContainer() UserDataContainer {
	return UserDataContainer{Type: USERDATACONTAINER_T}
}

type DefaultUserDataContainer struct {
	UserDataContainer
}

func NewDefaultUserDataContainer() DefaultUserDataContainer {
	obj := NewUserDataContainer()
	obj.Type = DEFAULTUSERDATACONTAINER_T
	return DefaultUserDataContainer{UserDataContainer: obj}
}
