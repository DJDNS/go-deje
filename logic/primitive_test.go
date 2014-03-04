package logic

/*
import (
    "testing"
    "reflect"
)

func TestTraverse_Map(t *testing.T) {
    container := map[string]interface{}{
        "hello": "world",
        "deep": map[string]interface{}{
            "area": "contents",
        },
    }

    retrieved, err := traverse_path_item(
        "hello",
        reflect.ValueOf(container),
    )
    if err != nil {
        t.Fatal(err)
    }

    path := []interface{}{"deep", "area"}
    retrieved, err = traverse_path(path, reflect.ValueOf(container))
    if err != nil {
        t.Fatal(err)
    }
    t.Fatal(retrieved.Type().Kind())
}
*/
