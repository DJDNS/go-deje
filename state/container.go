package state

//import "reflect"

type Container interface {
    Remove() error
    RemoveChild(interface{}) error
    SetParentage(Container, interface{})

    Set(key, value interface{}) error
    Export() interface{}
}

/*
func MakeContainer(value interface{}) (*Container, error) {
    switch reflect.TypeOf(value).Kind() {
    case reflect.Map:
        return MakeMapContainer(value)
    }
}
*/
