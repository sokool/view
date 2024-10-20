package view_test

import (
	"fmt"
	"testing"

	"github.com/sokool/view"
)

func TestReader_JSON(t *testing.T) {
	j := []byte(`{
	"id": 1,
	"name": "John",
	"age": 30.6,
	"hired": true,
	"location": {
		"point": {
			"lat": 40.7128,
			"lon": -74.0060
		},
		"address": "New York Hudson 60"
	},
	"jobs": [
		{"title": "developer", "salary": "100$"},
		{"title": "manager", "salary": "200$"},
		{"title": "ceo", "salary": "300$"}
	],
	"skills": ["go", "js", "python"]
}
`)

	r := view.JSONReader(j)

	if r.Text("name") != "John" {
		t.Fatalf("expected John, got %s", r.Text("name"))
	}
	if r.Number("location.point.lat") != 40.7128 {
		t.Fatalf("expected 40.7128, got %f", r.Number("location.point.lat"))
	}
	if r.Bool("hired") != true {
		t.Fatalf("expected true, got %t", r.Bool("hired"))
	}
	var s string
	var f float64
	var p map[string]any
	if err := r.Read("name", &s).Read("age", &f).Read("location.point", &p).Error(); err != nil {
		t.Fatalf("expected nil, got %s err", err)
	}
	if s != "John" {
		t.Fatalf("expected John, got %s", s)
	}
	if f != 30.6 {
		t.Fatalf("expected 30.6, got %f", f)
	}
	s = ""
	if err := r.Select("abc").To(&s); err != nil || s != "" { // abc attribute not exists
		t.Fatalf("expected empty string withou error, got %v error and %s value", err, s)
	}
	if n := r.Select("a").Select("b"); !n.IsEmpty() {
		t.Fatalf("expected empty meta")
	}
	if f = r.Select("location.point").Number("lat"); f != 40.7128 {
		t.Fatalf("expected 40.7128, got %f", f)
	}
	if s = r.Select("jobs[1]").Text("title"); s != "manager" {
		t.Fatalf("expected manager, got %s", s)
	}
	if s = r.Select("jobs[?(@.title == 'manager')]").Text("[0].salary"); s != "200$" {
		t.Fatalf("expected 200$, got %s", s)
	}
	if err := r.Select("jobs[1.title").Error(); err == nil {
		t.Fatalf("expected error, got nil")
	}
	var ss []string
	if err := r.Select("jobs[*].title").To(&ss); err != nil || fmt.Sprintf("%v", ss) != "[developer manager ceo]" {
		t.Fatalf("expected nil, got %s", err)
	}
}

func TestMeta_Each(t *testing.T) {
	slice := []byte(`[{"name": "Tom", "age": 32},{"name": "Jerry", "age": 30},
			    {"name": "Spike", "age": 40},{"name": "Tyke", "age": 20}]`)

	var m = view.JSONReader(slice)
	var ss []string
	if err := m.Select("[*].name").To(&ss); err != nil || fmt.Sprintf("%v", ss) != "[Tom Jerry Spike Tyke]" {
		t.Fatalf("expected nil, got %s", err)
	}

	var s string
	for n := range m.Each {
		s += n.Text("name")
	}
	if s != "TomJerrySpikeTyke" {
		t.Fatalf("expected TomJerrySpikeTyke, got %s", s)
	}
}
