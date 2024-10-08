package view

import (
	"reflect"
)

type Renderer interface {
	Render(string) any
}

func RendererDecorator(v reflect.Value) (reflect.Value, error) {
	if w, ok := v.Interface().(Renderer); ok {
		return reflect.ValueOf(w.Render("n")), nil
	}
	return reflect.Value{}, nil
}
