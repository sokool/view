package view

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type Decorator func(v reflect.Value) (reflect.Value, error)

// Decoder structure to hold the encode methods
type Decoder struct{ decorators []Decorator }

func NewDecoder(r ...Decorator) *Decoder {
	return &Decoder{r}
}

func Decode(v any, r ...Decorator) ([]byte, error) {
	return NewDecoder(r...).Decode(v)
}

func (d *Decoder) Decode(v any) ([]byte, error) {
	// Call the helper function to handle the reflection-based encoding
	s, err := d.toValue(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}

	// Return the JSON as a byte slice
	return []byte(s), nil
}

func (d *Decoder) toValue(v reflect.Value) (string, error) {
	//o := v.Interface()
	//fmt.Printf("%T:%v \n", o, o)
	if v.IsValid() {
		for _, r := range d.decorators {
			w, err := r(v)
			if err != nil {
				return "", err
			}
			if w.IsValid() {
				v = w
			}
		}
	}
	// Handle different types using reflection
	switch v.Kind() {
	case reflect.Struct:
		return d.toStruct(v)
	case reflect.Map:
		return d.toMap(v)
	case reflect.Slice, reflect.Array:
		return d.toSlice(v)
	case reflect.String:
		return strconv.Quote(v.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Invalid:
		return "null", nil
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "null", nil
		}
		return d.toValue(v.Elem())
	default:
		return "", fmt.Errorf("unsupported %T type for JSON encoding", v.String())
	}
}

func (d *Decoder) toStruct(rv reflect.Value) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")
	t := rv.Type()
	first := true // Used to handle commas correctly

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := t.Field(i)
		tag := ft.Tag.Get("json")

		// Check for field visibility (must be exported)
		if !fv.CanInterface() {
			continue
		}

		// Parse JSON tag (e.g. `json:"name,omitempty"`)
		tagParts := strings.Split(tag, ",")
		jsonKey := tagParts[0]
		if jsonKey == "" {
			jsonKey = ft.Name
		}

		// Handle `json:"-"`, which means ignore this field
		if jsonKey == "-" {
			continue
		}

		// Handle omitempty: skip the field if it is zero-valued
		if len(tagParts) > 1 && tagParts[1] == "omitempty" && fv.IsZero() {
			continue
		}

		// Encode the field value recursively
		v, err := d.toValue(fv)
		if err != nil {
			return "", err
		}

		// Add a comma if needed (not for the first field)
		if !first {
			sb.WriteString(",")
		}
		first = false

		// Append the field name and value in JSON format
		sb.WriteString(`"` + jsonKey + `":` + v)
	}

	sb.WriteString("}")
	return sb.String(), nil
}

func (d *Decoder) toMap(rv reflect.Value) (string, error) {
	var sb strings.Builder
	sb.WriteString("{")

	// Get map keys
	keys := rv.MapKeys()

	// Create a slice to hold the string keys and a map to associate the string keys with their values
	strKeys := make([]string, len(keys))
	keyValueMap := make(map[string]reflect.Value, len(keys))

	// Populate strKeys and keyValueMap
	for i, key := range keys {
		k, err := d.toKey(key)
		if err != nil {
			return "", err
		}
		strKeys[i] = k
		keyValueMap[k] = rv.MapIndex(key)
	}

	// Sort the string keys alphabetically
	sort.Strings(strKeys)

	// Write sorted keys and their values to the string builder
	for i, strKey := range strKeys {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strKey + ":")

		// Retrieve the value from the map directly
		v, err := d.toValue(keyValueMap[strKey])
		if err != nil {
			return "", err
		}
		sb.WriteString(v)
	}

	sb.WriteString("}")
	return sb.String(), nil
}

func (d *Decoder) toSlice(v reflect.Value) (string, error) {
	var sb strings.Builder
	sb.WriteString("[")

	for i := 0; i < v.Len(); i++ {
		encodedValue, err := d.toValue(v.Index(i))
		if err != nil {
			return "", err
		}

		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(encodedValue)
	}
	sb.WriteString("]")
	return sb.String(), nil
}

func (d *Decoder) toKey(v reflect.Value) (string, error) {
	return strconv.Quote(v.String()), nil
}
