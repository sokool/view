package view

import (
	_json "encoding/json"

	"github.com/bhmj/jsonslice"
)

type Reader struct {
	xpath map[string]any
}

func Read(xpath string, to any) *Reader {
	return &Reader{map[string]any{xpath: to}}
}

func (r *Reader) Read(xpath string, to any) *Reader { r.xpath[xpath] = to; return r }

func (r *Reader) UnmarshalJSON(from []byte) error {
	for q := range r.xpath {
		b, err := jsonslice.Get(from, q)
		if err != nil {
			return err
		}
		if b == nil {
			continue
		}
		if err = _json.Unmarshal(b, r.xpath[q]); err != nil {
			return err
		}
	}
	return nil
}
