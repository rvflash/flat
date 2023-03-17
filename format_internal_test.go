// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat

import (
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/matryer/is"
)

func TestFmtString(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			// inputs
			in  interface{}
			sep string
			// outputs
			out string
		}{
			"Default":       {},
			"False":         {in: false, out: "false"},
			"True":          {in: true, out: "true"},
			"String":        {in: "string", out: "string"},
			"Pi":            {in: float64(3.14), out: "3.14"},
			"JSON number":   {in: json.Number("-42"), out: "-42"},
			"Not supported": {in: int64(-42), out: ""},
			"Slice":         {in: []interface{}{"4", "2"}, sep: DefaultXMLArraySep, out: "4|2"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out := fmtString(tt.in, tt.sep)
			are.Equal(tt.out, out)
		})
	}
}

func TestToBool(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  interface{}
			out bool
			err error
		}{
			"Default": {err: ErrOutOfRange},
			"Invalid": {in: "", out: false, err: strconv.ErrSyntax},
			"String":  {in: "true", out: true},
			"OK":      {in: true, out: true},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := toBool(tt.in)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch result
		})
	}
}

func TestToFloat64(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  interface{}
			out float64
			err error
		}{
			"Default": {err: ErrOutOfRange},
			"Invalid": {in: "", out: 0, err: strconv.ErrSyntax},
			"Number":  {in: json.Number("3.14"), out: 3.14},
			"String":  {in: "3.14", out: 3.14},
			"OK":      {in: float64(3.14), out: 3.14},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := toFloat64(tt.in)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch result
		})
	}
}

func TestToInt64(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  interface{}
			out int64
			err error
		}{
			"Default": {err: ErrOutOfRange},
			"Invalid": {in: "", out: 0, err: strconv.ErrSyntax},
			"Number":  {in: json.Number("-42"), out: -42},
			"String":  {in: "-42", out: -42},
			"OK":      {in: float64(-42), out: -42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := toInt64(tt.in)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch result
		})
	}
}

func TestToString(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  interface{}
			out string
			err error
		}{
			"Default": {err: ErrOutOfRange},
			"Bool":    {in: true, out: "", err: ErrOutOfRange},
			"Number":  {in: json.Number("-42"), out: "-42"},
			"OK":      {in: "oops", out: "oops"},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := toString(tt.in)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch result
		})
	}
}

func TestToUint64(t *testing.T) {
	var (
		are = is.New(t)
		dt  = map[string]struct {
			in  interface{}
			out uint64
			err error
		}{
			"Default": {err: ErrOutOfRange},
			"Invalid": {in: "", out: 0, err: strconv.ErrSyntax},
			"Number":  {in: json.Number("42"), out: 42},
			"String":  {in: "42", out: 42},
			"OK":      {in: float64(42), out: 42},
		}
	)
	for name, tt := range dt {
		tt := tt
		t.Run(name, func(t *testing.T) {
			out, err := toUint64(tt.in)
			are.True(errors.Is(err, tt.err)) // unexpected error
			are.Equal(tt.out, out)           // mismatch result
		})
	}
}
