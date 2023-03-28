package view

import (
	"bytes"
	"encoding/json"
)

var Null Bytes = []byte(`null`)

type Bytes []byte

func (b Bytes) Empty() bool {
	return b.Eq(Null) || len(b) == 0
}

func (b Bytes) Eq(c []byte) bool {
	return bytes.Equal(b, c)
}

func (b Bytes) IsJSON() bool {
	if b.Empty() {
		return true
	}
	switch b[0] {
	case '[':
		return true
	case '{':
		return true
	default:
		return false
	}
}

func (b Bytes) JSON() json.RawMessage {
	return json.RawMessage(b)
}

func (b Bytes) IsYAML() bool {
	return false
}

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) Buffer() *bytes.Buffer {
	return bytes.NewBuffer(b)
}
