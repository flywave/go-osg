package model

const (
	USERDATACONTAINERT        string = "osg::UserDataContainer"
	DEFAULTUSERDATACONTAINERT string = "osg::DefaultUserDataContainer"
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

func NewUserDataContainer() *UserDataContainer {
	return &UserDataContainer{Type: USERDATACONTAINERT}
}

type DefaultUserDataContainer struct {
	UserDataContainer
}

func NewDefaultUserDataContainer() *DefaultUserDataContainer {
	obj := NewUserDataContainer()
	obj.Type = DEFAULTUSERDATACONTAINERT
	return &DefaultUserDataContainer{UserDataContainer: *obj}
}
