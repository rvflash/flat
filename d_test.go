// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"github.com/rvflash/flat"

	"gopkg.in/yaml.v3"
)

const (
	jsonStr = `
{
  "array": [
	1,
	2,
	3
  ],
  "boolean": true,
  "null": null,
  "number": 123,
  "object": {
	"a": "b",
	"c": "d",
	"e": "f"
  },
  "string": "Hello World"
}`
	xmlStr = `
<root xmlns:hyp="hyp">
  <array>1|2|3</array>
  <boolean>true</boolean>
  <null></null>
  <hyp:number>123</hyp:number>
  <object>
	<a>b</a>
	<c>d</c>
	<e>f</e>
  </object>
  <string>Hello World</string>
</root>
`
	yamlStr = `array:
- 1
- 2
- 3
boolean: true
'null':
number: 123
object:
  a: b
  c: d
  e: f
string: Hello World`
)

func TestD_Flatten(t *testing.T) {
	var (
		are = is.New(t)
		d   = map[string]interface{}{
			"array":      []interface{}{float64(1), float64(2), float64(3)},
			"boolean":    true,
			"null":       nil,
			"hyp:number": float64(123),
			"object": map[string]interface{}{
				"a": "b",
				"c": "d",
				"e": "f",
			},
			"string": "Hello World",
		}
		dt = map[string]struct {
			in  *flat.D
			not [][]string
			out map[string]interface{}
		}{
			"Default": {in: &flat.D{}},
			"With private fields": {
				in:  flat.New(d),
				not: [][]string{{"object", "c"}, {"hyp:number"}},
				out: map[string]interface{}{
					"array":    []interface{}{float64(1), float64(2), float64(3)},
					"boolean":  true,
					"null":     nil,
					"object_a": "b",
					"object_e": "f",
					"string":   "Hello World",
				},
			},
			"OK": {
				in: flat.New(d),
				out: map[string]interface{}{
					"array":      []interface{}{float64(1), float64(2), float64(3)},
					"boolean":    true,
					"null":       nil,
					"hyp_number": float64(123),
					"object_a":   "b",
					"object_c":   "d",
					"object_e":   "f",
					"string":     "Hello World",
				},
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.Flatten(tt.not...)
			are.Equal("", cmp.Diff(tt.out, out))
		})
	}
}

func TestD_Lookup(t *testing.T) {
	var (
		d = map[string]interface{}{
			"object": map[string]interface{}{
				"a": "b",
			},
		}
		are = is.New(t)
		dt  = map[string]struct {
			in   *flat.D
			keys []string
			out  interface{}
			err  error
		}{
			"Default":       {err: flat.ErrNotFound},
			"Blank":         {in: &flat.D{}, err: flat.ErrNotFound},
			"Unknown group": {in: flat.New(d), keys: []string{"object", "a", "b"}, err: flat.ErrNotFound},
			"Unknown value": {in: flat.New(d), keys: []string{"object", "b"}, err: flat.ErrNotFound},
			"OK":            {in: flat.New(d), keys: []string{"object", "a"}, out: "b"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.Lookup(tt.keys...)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch data
		})
	}
}

func TestD_JSONEncode(t *testing.T) {
	var (
		are = is.New(t)
		buf = bytes.Buffer{}
		err = flat.New(nil).JSONEncode(&buf)
	)
	are.NoErr(err)                    // unexpected error
	are.Equal("null\n", buf.String()) // mismatch value
}

func TestD_MarshalJSON(t *testing.T) {
	var (
		are    = is.New(t)
		d      = flat.New(nil)
		b, err = json.Marshal(d)
	)
	are.NoErr(err)               // unexpected error
	are.Equal("null", string(b)) // mismatch value
}

func TestD_UnmarshalJSON(t *testing.T) {
	var (
		d   = flat.D{}
		are = is.New(t)
		buf = []byte(jsonStr)
		err = json.Unmarshal(buf, &d)
	)
	are.NoErr(err)
	are.Equal("", cmp.Diff(d.Flatten(), map[string]interface{}{
		"array":    []interface{}{json.Number("1"), json.Number("2"), json.Number("3")},
		"boolean":  true,
		"null":     nil,
		"number":   json.Number("123"),
		"object_a": "b",
		"object_c": "d",
		"object_e": "f",
		"string":   "Hello World",
	}))
}

func TestD_UnmarshalJSON2(t *testing.T) {
	var (
		are = is.New(t)
		d   = flat.D{}
		err = d.UnmarshalJSON(nil)
	)
	are.NoErr(err)              // unexpected error
	are.Equal(nil, d.Flatten()) // mismatch value
}

func TestD_XMLEncode(t *testing.T) {
	var (
		are = is.New(t)
		buf = &bytes.Buffer{}
		err = flat.New(nil).XMLEncode(buf)
	)
	are.NoErr(err)              // unexpected error
	are.Equal("", buf.String()) // mismatch value
}

func TestD_MarshalXML(t *testing.T) {
	var (
		are    = is.New(t)
		d      = flat.New(nil)
		b, err = xml.Marshal(d)
	)
	are.NoErr(err)           // unexpected error
	are.Equal("", string(b)) // mismatch value
}

func TestD_UnmarshalXML(t *testing.T) {
	var (
		d   = flat.D{}
		are = is.New(t)
		buf = []byte(xmlStr)
		err = xml.Unmarshal(buf, &d)
	)
	are.NoErr(err)
	are.Equal("", cmp.Diff(d.Flatten(), map[string]interface{}{
		"array":      "1|2|3", // todo in the next release: []interface{}{"1","2","3"}
		"boolean":    "true",  // todo in the next release: true
		"null":       "\n  ",  // todo in the next release: nil
		"hyp_number": "123",
		"object_a":   "b",
		"object_c":   "d",
		"object_e":   "f",
		"string":     "Hello World",
	}))
}

func TestD_YAMLEncode(t *testing.T) {
	var (
		are = is.New(t)
		buf = bytes.Buffer{}
		err = flat.New(nil).YAMLEncode(&buf)
	)
	are.NoErr(err)                  // unexpected error
	are.Equal("{}\n", buf.String()) // mismatch value
}

func TestD_UnmarshalYAML(t *testing.T) {
	var (
		d   = flat.D{}
		are = is.New(t)
		buf = []byte(yamlStr)
		err = yaml.Unmarshal(buf, &d)
	)
	are.NoErr(err)
	are.Equal("", cmp.Diff(d.Flatten(), map[string]interface{}{
		"array":    []interface{}{1, 2, 3},
		"boolean":  true,
		"null":     nil,
		"number":   123,
		"object_a": "b",
		"object_c": "d",
		"object_e": "f",
		"string":   "Hello World",
	}))
}

func TestD_UnmarshalYAML2(t *testing.T) {
	var (
		are = is.New(t)
		d   = flat.D{}
		err = d.UnmarshalYAML(nil)
	)
	are.NoErr(err)              // unexpected error
	are.Equal(nil, d.Flatten()) // mismatch value
}

func TestD_Bool(t *testing.T) {
	var (
		d   = flat.New(map[string]interface{}{"bool": true})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out bool
			err error
		}{
			"Default": {err: flat.ErrNotFound},
			"Blank":   {keys: []string{"bool"}, err: flat.ErrNotFound},
			"Unknown": {in: d, err: flat.ErrNotFound, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"bool"}, out: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.Bool(tt.keys...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(tt.out, out)           // mismatch default value
		})
	}
}

func TestD_ShouldBool(t *testing.T) {
	var (
		d   = flat.New(map[string]interface{}{"bool": true})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out bool
		}{
			"Default": {},
			"Blank":   {keys: []string{"bool"}},
			"Unknown": {in: d, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"bool"}, out: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.ShouldBool(tt.keys...)
			are.Equal(tt.out, out)
		})
	}
}

func TestD_Float64(t *testing.T) {
	var (
		f   = float64(3.14)
		d   = flat.New(map[string]interface{}{"float64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out float64
			err error
		}{
			"Default": {err: flat.ErrNotFound},
			"Blank":   {keys: []string{"float64"}, err: flat.ErrNotFound},
			"Unknown": {in: d, err: flat.ErrNotFound, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"float64"}, out: f},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.Float64(tt.keys...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(tt.out, out)           // mismatch default value
		})
	}
}

func TestD_ShouldFloat64(t *testing.T) {
	var (
		f   = float64(3.14)
		d   = flat.New(map[string]interface{}{"float64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out float64
		}{
			"Default": {},
			"Blank":   {keys: []string{"float64"}},
			"Unknown": {in: d, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"float64"}, out: f},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.ShouldFloat64(tt.keys...)
			are.Equal(tt.out, out)
		})
	}
}

func TestD_Int64(t *testing.T) {
	var (
		f   = float64(-42)
		d   = flat.New(map[string]interface{}{"int64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out int64
			err error
		}{
			"Default": {err: flat.ErrNotFound},
			"Blank":   {keys: []string{"int64"}, err: flat.ErrNotFound},
			"Unknown": {in: d, err: flat.ErrNotFound, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"int64"}, out: -42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.Int64(tt.keys...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(tt.out, out)           // mismatch default value
		})
	}
}

func TestD_ShouldInt64(t *testing.T) {
	var (
		f   = float64(-42)
		d   = flat.New(map[string]interface{}{"int64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out int64
		}{
			"Default": {},
			"Blank":   {keys: []string{"int64"}},
			"Unknown": {in: d, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"int64"}, out: -42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.ShouldInt64(tt.keys...)
			are.Equal(tt.out, out)
		})
	}
}

func TestD_String(t *testing.T) {
	var (
		s   = "hi"
		i   = "42"
		d   = flat.New(map[string]interface{}{"string": s, "number": json.Number(i)})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out string
			err error
		}{
			"Default": {err: flat.ErrNotFound},
			"Blank":   {keys: []string{"string"}, err: flat.ErrNotFound},
			"Unknown": {in: d, err: flat.ErrNotFound, keys: []string{"oops"}},
			"Number":  {in: d, keys: []string{"number"}, out: i},
			"OK":      {in: d, keys: []string{"string"}, out: s},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.String(tt.keys...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(tt.out, out)           // mismatch default value
		})
	}
}

func TestD_ShouldString(t *testing.T) {
	var (
		s   = "hi"
		i   = "42"
		d   = flat.New(map[string]interface{}{"string": s, "number": json.Number(i)})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out string
		}{
			"Default": {},
			"Blank":   {keys: []string{"string"}},
			"Unknown": {in: d, keys: []string{"oops"}},
			"Number":  {in: d, keys: []string{"number"}, out: i},
			"OK":      {in: d, keys: []string{"string"}, out: s},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.ShouldString(tt.keys...)
			are.Equal(tt.out, out)
		})
	}
}

func TestD_Strings(t *testing.T) {
	var (
		are = is.New(t)
		d   = flat.New(map[string]interface{}{
			"numbers":  []interface{}{json.Number("1")},
			"booleans": []interface{}{true},
			"bool":     true,
			"strings":  []interface{}{"4", "2"},
		})
		dt = map[string]struct {
			keys []string
			out  []string
			err  error
		}{
			"Default":    {err: flat.ErrNotFound},
			"Unknown":    {keys: []string{"oops"}, err: flat.ErrNotFound},
			"Invalid":    {keys: []string{"bool"}, err: flat.ErrOutOfRange},
			"Wrong type": {keys: []string{"booleans"}, err: flat.ErrOutOfRange},
			"Number":     {keys: []string{"numbers"}, out: []string{"1"}},
			"OK":         {keys: []string{"strings"}, out: []string{"4", "2"}},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := d.Strings(tt.keys...)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch data
		})
	}
}

func TestD_Time(t *testing.T) {
	var (
		are = is.New(t)
		d   = flat.New(map[string]interface{}{
			"time": "08/1983",
			"bool": true,
		})
		x  = time.Date(1983, time.August, 1, 0, 0, 0, 0, time.UTC)
		dt = map[string]struct {
			layout string
			keys   []string
			out    time.Time
			err    error
		}{
			"Default":    {err: flat.ErrNotFound},
			"Unknown":    {keys: []string{"oops"}, err: flat.ErrNotFound},
			"Wrong type": {keys: []string{"bool"}, err: flat.ErrOutOfRange},
			"OK":         {keys: []string{"time"}, layout: "01/2006", out: x},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := d.Time(tt.layout, tt.keys...)
			are.True(errors.Is(err, tt.err))                       // unexpected error
			are.Equal(tt.out, out)                                 // mismatch data
			are.Equal(tt.out, d.ShouldTime(tt.layout, tt.keys...)) // mismatch should data
		})
	}
}

func TestD_Uint64(t *testing.T) {
	var (
		f   = float64(42)
		d   = flat.New(map[string]interface{}{"uint64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out uint64
			err error
		}{
			"Default": {err: flat.ErrNotFound},
			"Blank":   {keys: []string{"uint64"}, err: flat.ErrNotFound},
			"Unknown": {in: d, err: flat.ErrNotFound, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"uint64"}, out: 42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := tt.in.Uint64(tt.keys...)
			are.True(errors.Is(err, tt.err)) // mismatch error
			are.Equal(tt.out, out)           // mismatch default value
		})
	}
}

func TestD_ShouldUint64(t *testing.T) {
	var (
		f   = float64(42)
		d   = flat.New(map[string]interface{}{"uint64": f})
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in   *flat.D
			keys []string
			// outputs
			out uint64
		}{
			"Default": {},
			"Blank":   {keys: []string{"uint64"}},
			"Unknown": {in: d, keys: []string{"oops"}},
			"OK":      {in: d, keys: []string{"uint64"}, out: 42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := tt.in.ShouldUint64(tt.keys...)
			are.Equal(tt.out, out)
		})
	}
}
