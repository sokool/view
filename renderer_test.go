package view_test

import (
	"fmt"
	"testing"

	"github.com/sokool/view"
)

func TestJSON(t *testing.T) {
	//b, err := view.JSON([]foo{"two", "one", "six"}, "a", "b", "c")
	//assert.NoError(t, err)
	//assert.Equal(t, b, view.Bytes(`[{"bar":"-bar:hi-","baz":[{"a":1},{"b":"two"}],"luz":{"one":1,"two":2},"name":"-foo:two-"},{"bar":"-bar:hi-","baz":[{"a":1},{"b":"one"}],"luz":{"one":1,"two":2},"name":"-foo:one-"},{"bar":"-bar:hi-","baz":[{"a":1},{"b":"six"}],"luz":{"one":1,"two":2},"name":"-foo:six-"}]`))
	//
	//b, err = view.JSON(map[string]foo{"high": ("five"), "fourty": ("eight"), "ninthy": ("six")}, "1", "2", "3")
	//assert.NoError(t, err)
	//assert.Equal(t, b, view.Bytes(`{"fourty":{"bar":"-bar:hi-","baz":[{"a":1},{"b":"eight"}],"luz":{"one":1,"two":2},"name":"-foo:eight-"},"high":{"bar":"-bar:hi-","baz":[{"a":1},{"b":"five"}],"luz":{"one":1,"two":2},"name":"-foo:five-"},"ninthy":{"bar":"-bar:hi-","baz":[{"a":1},{"b":"six"}],"luz":{"one":1,"two":2},"name":"-foo:six-"}}`))
	//
	//b, err = view.JSON(foo("hi there"), "!", "@", "3")
	//assert.NoError(t, err)
	//assert.Equal(t, b, view.Bytes(`{"bar":"-bar:hi-","baz":[{"a":1},{"b":"hi there"}],"luz":{"one":1,"two":2},"name":"-foo:hi there-"}`))
}

type bar string

func (b bar) Render(...string) (any, error) {
	return fmt.Sprintf("-bar:%v-", b), nil
}

type foo string

func (f foo) Render(types ...string) (any, error) {
	return view.Writer{
		"bar":  bar("hi"),
		"name": fmt.Sprintf("-foo:%v-", f),
		"baz": view.Writers{
			{"a": 1},
			{"b": string(f)},
		},
		"luz": map[string]int{
			"one": 1,
			"two": 2,
		},
	}, nil
}
