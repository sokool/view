package view_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	. "github.com/sokool/view"
)

func TestRenderer(t *testing.T) {
	type (
		scenario struct {
			payload any
			expects string
		}
		object = map[string]any
	)
	cases := map[string]scenario{
		"nil->null": {
			payload: nil,
			expects: `null`,
		},
		"int->number": {
			payload: 56,
			expects: `56`,
		},
		"float->number": {
			payload: 8.929349,
			expects: `8.929349`,
		},
		"slice->array": {
			payload: []any{1, 2, 3},
			expects: `[1,2,3]`,
		},
		"map->object": {
			payload: object{"one": uint64(1), "two": int8(2)},
			expects: `{"one":1,"two":2}`,
		},
		"type with Render": {
			payload: Name("John"),
			expects: `"Mr. John"`,
		},
		"type int with marshals text": {
			payload: Quantity(34),
			expects: `"34qty"`,
		},
		"type with unexported fields and marshals json": {
			payload: NewEmail("john", "gmail.com"),
			expects: `"john@gmail.com"`,
		},
		//"writer renderer": {
		//	payload: Writer{
		//		"internal": Writer{
		//			"one": 1,
		//		},
		//		"name": Name("John"),
		//	},
		//	expects: `{"internal":{"one":1},"name":"Mr. John"}`,
		//},
		"x": {
			payload: foo{One: "John", Mail: NewEmail("john", "gmail.com")},
			expects: `{"firstname":"Mr. John","Mail":"john@gmail.com"}`,
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			b, err := Decode(c.payload, RendererDecorator)
			if err != nil {
				t.Fatalf("no error expected, got %v", err)
			}
			if bytes.Compare(b, []byte(c.expects)) != 0 {
				t.Fatalf("expected %v, got %s", c.expects, b)
			}
		})
	}

}

type Quantity int

func (q Quantity) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(int(q)) + "qty"), nil
}

type Email struct{ name, host string }

func NewEmail(name, host string) Email {
	return Email{name, host}
}

func (e Email) MarshalJSON() ([]byte, error) {
	if e.host == "" {
		return []byte(`null`), nil
	}
	return []byte(fmt.Sprintf(`"%s@%s"`, e.name, e.host)), nil
}

type Name string

func (n Name) Render(s string) any {
	if n == "" {
		return ""
	}
	return fmt.Sprintf("Mr. %s", n)
}

type foo struct {
	One  Name `json:"firstname"`
	Two  Name `json:"lastname,omitempty"`
	Mail Email
	Test int `json:"-"`
}
