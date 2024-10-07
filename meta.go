package view

func Write(name string, v any) Writer { return Writer{name: v} }

type Writer map[string]any

func (w Writer) Write(name string, v any) Writer {
	w[name] = v
	return w
}
