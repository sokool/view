package view_test

import (
	"testing"

	"github.com/sokool/view"
)

func TestReader_JSON(t *testing.T) {
	j := []byte(`
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`)

	type name struct {
		First,
		Last string
	}
	type model struct {
		name    name
		friends []name
	}

	var m model
	if err := view.Read("$.name", &m.name).Read("$.friends", &m.friends).UnmarshalJSON(j); err != nil {
		t.Fatal(err)
	}
	if m.name.First != "Tom" || m.name.Last != "Anderson" {
		t.Fatal()
	}
	if len(m.friends) != 3 {
		t.Fatal()
	}
	if n := m.friends[0]; n.First != "Dale" {
		t.Fatal()
	}
	if n := m.friends[1]; n.First != "Roger" {
		t.Fatal()
	}
	if n := m.friends[2]; n.First != "Jane" {
		t.Fatal()
	}
	var n []name
	if err := view.Read("$.friends[?(@.age > 50)]", &n).UnmarshalJSON(j); err != nil {
		t.Fatal(err)
	}
	if len(n) != 1 || n[0].Last != "Craig" {
		t.Fatal()
	}
}
