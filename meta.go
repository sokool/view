package view

import (
	"encoding/json"
	"fmt"
)

func Write(name string, v any) Writer { return Writer{name: v} }

type Writer map[string]any

func (w Writer) Write(name string, v any) Writer {
	w[name] = v
	return w
}

func (w Writer) GoString() string {
	return stringify(w)
}

func stringify(v any) string {
	b, _ := json.MarshalIndent(v, "", "\t")
	return fmt.Sprintf("%T\n%s\n", v, b)
}
