// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
)

func TestSimplify(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  map[string]interface{}
			out map[string]interface{}
		}{
			"Default": {},
			"Short":   {in: map[string]interface{}{"key": "value"}, out: map[string]interface{}{"key": "value"}},
			"Common part but inside keys name": {
				in:  map[string]interface{}{"geek1": "value", "geek2": "value"},
				out: map[string]interface{}{"geek1": "value", "geek2": "value"},
			},
			"Only some keys have a common prefix": {
				in: map[string]interface{}{
					"array":    "value",
					"object_a": "value",
					"object_c": "value",
					"object_e": "value",
					"string":   "value",
				},
				out: map[string]interface{}{
					"array":    "value",
					"object_a": "value",
					"object_c": "value",
					"object_e": "value",
					"string":   "value",
				},
			},
			"OK": {
				in:  map[string]interface{}{"geek_name": "value", "geek_age": float64(42)},
				out: map[string]interface{}{"name": "value", "age": float64(42)},
			},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := simplify(tt.in)
			are.Equal("", cmp.Diff(tt.out, out)) // mismatch data
		})
	}
}
