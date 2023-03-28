package view_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sokool/view"
)

func TestWriter_MarshalJSON(t *testing.T) {
	var j = []byte(`{"intro":"Hello word","number":3.14}`)
	var w = view.Write("intro", "Hello word").Write("number", 3.14)
	var b []byte
	var err error

	if b, err = w.MarshalJSON(); !bytes.Equal(b, j) || err != nil {
		t.Fatal()
	}
	if b, err = json.Marshal(w); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, j) {
		t.Fatal()
	}
}
