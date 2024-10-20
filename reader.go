package view

import (
	_json "encoding/json"
	"fmt"
	"strings"

	"github.com/bhmj/jsonslice"
)

type Reader struct {
	bytes []byte
	err   error
}

func JSONReader(json []byte) *Reader {
	return &Reader{json, nil}
}

// Number returns the float64 at the given JSON path.
func (r *Reader) Number(path string) (f float64) {
	if err := r.Select(path).To(&f); err != nil {
		r.report(err)
	}
	return f
}

// Text returns the string at the given JSON path.
func (r *Reader) Text(path string) (s string) {
	if err := r.Select(path).To(&s); err != nil {
		r.report(err)
	}
	return s
}

// Bool returns the bool at the given JSON path.
func (r *Reader) Bool(path string) (b bool) {
	if err := r.Select(path).To(&b); err != nil {
		r.report(err)
	}
	return b
}

// Select returns the Writer at the given JSON path.
func (r *Reader) Select(path string, args ...any) *Reader {
	if r.err != nil {
		return r
	}
	path = fmt.Sprintf(path, args...)
	if path = strings.TrimSpace(path); len(path) >= 1 && path[0] != '$' {
		path = "$." + path
	}
	b, err := jsonslice.Get(r.bytes, path)
	if err != nil {
		return &Reader{nil, err}
	}
	return JSONReader(b)
}

// Read reads the value at the given JSON path into the given value.
func (r *Reader) Read(path string, to any) *Reader {
	if err := r.Select(path).To(to); err != nil {
		r.report(err)
		return r
	}
	return r
}

// To transform a Writer into the given value.
func (r *Reader) To(value any) error {
	if r.err != nil {
		return r.err
	}
	if r.IsEmpty() {
		return nil
	}
	return _json.Unmarshal(r.bytes, value)
}

func (r *Reader) Each(fn func(*Reader) bool) {
	for i := 0; ; i++ {
		if n := r.Select("[%d]", i); n.IsEmpty() || !fn(n) {
			return
		}
	}
}

func (r *Reader) Error() error {
	return r.err
}

func (r *Reader) IsEmpty() bool {
	return r.bytes == nil
}

func (r *Reader) String() string {
	return string(r.bytes)
}

func (r *Reader) report(err error) {
	r.err = err
}
