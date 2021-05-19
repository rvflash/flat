// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	base10    = 10
	bits64    = 64
	precision = -1
)

func fmtString(x interface{}, xmlArraySep string) string {
	switch d := x.(type) {
	case []interface{}:
		a := make([]string, len(d))
		for k, v := range d {
			a[k] = fmtString(v, xmlArraySep)
		}
		return strings.Join(a, xmlArraySep)
	case bool:
		return strconv.FormatBool(d)
	case float64:
		return strconv.FormatFloat(d, 'g', precision, bits64)
	case string:
		return d
	case json.Number:
		return d.String()
	default:
		return ""
	}
}

func toBool(m interface{}) (bool, error) {
	switch v := m.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		var x bool
		return x, newErrOutOfRange(x, v)
	}
}

func toFloat64(m interface{}) (float64, error) {
	switch v := m.(type) {
	case float64:
		return v, nil
	case json.Number:
		return v.Float64()
	case string:
		return strconv.ParseFloat(v, bits64)
	default:
		var x float64
		return x, newErrOutOfRange(x, v)
	}
}

func toInt64(m interface{}) (int64, error) {
	switch v := m.(type) {
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case string:
		return strconv.ParseInt(v, base10, bits64)
	default:
		var x int64
		return x, newErrOutOfRange(x, v)
	}
}

func toString(m interface{}) (string, error) {
	s, ok := m.(string)
	if !ok {
		return "", newErrOutOfRange("", m)
	}
	return s, nil
}

func toUint64(m interface{}) (uint64, error) {
	switch v := m.(type) {
	case float64:
		return uint64(v), nil
	case json.Number:
		return strconv.ParseUint(v.String(), base10, bits64)
	case string:
		return strconv.ParseUint(v, base10, bits64)
	default:
		var x uint64
		return x, newErrOutOfRange(x, v)
	}
}
