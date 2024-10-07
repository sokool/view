package view

import (
	"reflect"
)

type Renderer interface {
	Render(string) any
}

//func String(v any) string {
//	b, err := Decode(v)
//	if err != nil {
//		return fmt.Sprintf("view: %T to string failed due %s", v, err)
//	}
//	var w bytes.Buffer
//	if err = json.Indent(&w, b, "", "\t"); err != nil {
//		return fmt.Sprintf("view: %T to string failed due %s", v, err)
//	}
//	return w.String()
//}

func RendererDecorator(v reflect.Value) (reflect.Value, error) {
	if w, ok := v.Interface().(Renderer); ok {
		return reflect.ValueOf(w.Render("n")), nil
	}
	return reflect.Value{}, nil
}
