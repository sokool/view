package view

import (
	"bytes"
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
	b := r.bytes
	switch value.(type) {
	case *int, *int8, *int16, *int32, *int64, *uint, *uint8, *uint16, *uint32, *uint64,
		*float32, *float64, *complex64, *complex128:
		// Check if the first and last characters are apostrophes
		if b = bytes.TrimSpace(r.bytes); len(b) > 1 && b[0] == '"' && b[len(b)-1] == '"' {
			b = b[1 : len(b)-1] // Remove the apostrophes
		}
	default:

	}

	return _json.Unmarshal(b, value)
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
	return len(r.bytes) == 0 || bytes.Equal(r.bytes, []byte("null"))
}

func (r *Reader) String() string {
	return string(r.bytes)
}

func (r *Reader) Sprintf(msg string, jsonPath ...string) string {
	var args = make([]any, len(jsonPath))
	for i, p := range jsonPath {
		args[i] = r.Text(p)
	}
	return fmt.Sprintf(msg, args...)
}

func (r *Reader) report(err error) {
	r.err = err
}
