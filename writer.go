package view

import "encoding/json"

func Write(name string, v any) Writer { return Writer{name: v} }

type Writer map[string]any

func (w Writer) Write(name string, v any) Writer {
	w[name] = v
	return w
}

func (w Writer) Render(types ...string) (any, error) {
	return map[string]any(w), nil
}

func (w Writer) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(w))
}

func (w Writer) String() string {
	return String(w)
}

type Writers []Writer

func (w Writers) Render(...string) (any, error) {
	return w, nil
}

func (w Writers) MarshalJSON() (_ []byte, err error) {
	s := len(w)
	if s == 0 {
		return Null, nil
	}
	var j = make([]json.RawMessage, s)
	for i := range w {
		if j[i], err = w[i].MarshalJSON(); err != nil {
			return nil, err
		}
	}
	return json.Marshal(j)
}

func (w Writers) YAML() ([]byte, error) {
	return YAML(w)
}
