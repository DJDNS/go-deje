package state

type MapContainer struct{
    Parent *Container
    ParentKey interface{}
    Value map[string]*Container
}

/*
func MakeMapContainer(m map[string]interface{}) (*Container, error) {
    c := MapContainer{
        nil,
        nil,
        make(map[string]*Container),
    }
    for key, value := range m {
        child, err := MakeContainer(value)
        if err != nil { return nil, err }

        child.SetParentage(&c, key)
        c[key] = child
    }
    return &c, nil
}
*/
