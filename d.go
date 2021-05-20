// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package flat provides methods to handle XML or JSON data as a map[string]interface.
package flat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/rvflash/naming"
)

// Settings allows to customize the data during the marshalling or unmarshalling processes.
type Settings func(*D)

// XMLArray defines the separator used to handle XML array.
func XMLArray(sep string) Settings {
	return func(d *D) {
		if sep != "" {
			d.xmlArraySep = sep
		}
	}
}

// XMLName allows to define the XML name of the data.
func XMLName(s string) Settings {
	return func(d *D) {
		if s != "" {
			d.xmlName = s
		}
	}
}

// XMLNS allows to define the XML namespace of the data.
func XMLNS(s string) Settings {
	return func(d *D) {
		if s != "" {
			d.xmlns = s
		}
	}
}

// XMLAttributes sets the given list of attributes on the XML root data.
func XMLAttributes(list []xml.Attr) Settings {
	return func(d *D) {
		d.xmlAttributes = list
	}
}

const (
	// DefaultXMLName is the default XML name of the data.
	DefaultXMLName = "d"
	// DefaultXMLArraySep is the default XML separator of each array values.
	DefaultXMLArraySep = "|"
)

// New creates a new instance of D based on the given data and options.
// Supported types are:
// bool, for booleans
// float64, for numbers
// json.Number, for numbers (float64, int64 or uint64).
// string, for string literals
// nil, for null
// []interface{}, for arrays
// map[string]interface{}, for objects.
func New(m map[string]interface{}, opts ...Settings) *D {
	d := &D{D: m}
	for _, opt := range append([]Settings{
		XMLArray(DefaultXMLArraySep),
		XMLName(DefaultXMLName),
	}, opts...) {
		opt(d)
	}
	return d
}

// D represents a data.
type D struct {
	D             map[string]interface{}
	xmlArraySep   string
	xmlAttributes []xml.Attr
	xmlName       string
	xmlns         string
}

const (
	levelSep = " "
	rootName = ""
	keySep   = '_'
)

// Flatten allows to export D in a single dimension.
// Any of its properties, absent from the list of ignored keys, are lifted to the first level.
// Each property has a new name, using the snake case, based on names of its hierarchy.
// Common prefix in keys name are omitted to limit the length of each ones.
func (d D) Flatten(ignoredKeys ...[]string) map[string]interface{} {
	if len(d.D) == 0 {
		return nil
	}
	not := make(map[string]struct{}, len(ignoredKeys))
	for _, v := range ignoredKeys {
		not[naming.SnakeCase(strings.Join(v, levelSep))] = struct{}{}
	}
	return simplify(flatten(d.D, not, rootName))
}

func flatten(in map[string]interface{}, not map[string]struct{}, root string) map[string]interface{} {
	var (
		out = make(map[string]interface{})
		fk  string
		ok  bool
	)
	for k, v := range in {
		fk = naming.SnakeCase(root + levelSep + k)
		if _, ok = not[fk]; ok {
			continue
		}
		switch d := v.(type) {
		case map[string]interface{}:
			for kf, vf := range flatten(d, not, fk) {
				out[kf] = vf
			}
		default:
			out[fk] = d
		}
	}
	return out
}

func simplify(in map[string]interface{}) map[string]interface{} {
	prefix := commonPrefix(in)
	if prefix == "" {
		return in
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[strings.TrimPrefix(k, prefix)] = v
	}
	return out
}

func commonPrefix(in map[string]interface{}) string {
	n := len(in)
	if n <= 1 {
		return ""
	}
	var (
		i   int
		x   = make([]string, n)
		min = func(a, b int) int {
			if a > b {
				return b
			}
			return a
		}
	)
	// Sorts keys.
	for k := range in {
		x[i] = k
		i++
	}
	sort.Strings(x)
	// Identifies the common prefix.
	r1, r2 := []rune(x[0]), []rune(x[n-1])
	c := min(len(r1), len(r2))
	i = 0
	for i < c && r1[i] == r2[i] {
		i++
	}
	if i == 0 || r1[i-1] != keySep {
		return ""
	}
	return string(r1[:i])
}

// Lookup retrieves the value behind these keys.
// If the key is present, the value behind it is returned and the boolean is true.
func (d D) Lookup(keys ...string) (interface{}, error) {
	if len(keys) == 0 {
		return nil, ErrNotFound
	}
	var (
		v  interface{} = d.D
		m  map[string]interface{}
		ok bool
	)
	for i := 0; i < len(keys); i++ {
		m, ok = v.(map[string]interface{})
		if !ok {
			return nil, ErrNotFound
		}
		v, ok = m[keys[i]]
		if !ok {
			return nil, ErrNotFound
		}
	}
	return v, nil
}

// JSONEncode JSON encodes D into w.
func (d D) JSONEncode(w io.Writer) error {
	return json.NewEncoder(w).Encode(d)
}

// MarshalJSON implements the json.Marshaler interface.
func (d D) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.D)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *D) UnmarshalJSON(b []byte) (err error) {
	if b == nil {
		d.D = nil
		return
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	return dec.Decode(&d.D)
}

// XMLEncode XML encodes D into w.
func (d D) XMLEncode(w io.Writer) error {
	return xml.NewEncoder(w).Encode(d)
}

// MarshalXML implements the xml.Marshaler interface.
func (d D) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if len(d.D) == 0 {
		return nil
	}
	start.Name.Local = d.xmlName
	start.Name.Space = d.xmlns
	start.Attr = d.xmlAttributes
	return marshallXML(d.D, enc, start, d.xmlArraySep)
}

type charData struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func marshallXML(m map[string]interface{}, enc *xml.Encoder, start xml.StartElement, arraySep string) error {
	err := enc.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range m {
		d, ok := v.(map[string]interface{})
		if ok {
			err = marshallXML(d, enc, xml.StartElement{Name: xml.Name{Local: k}}, arraySep)
		} else {
			err = enc.Encode(charData{XMLName: xml.Name{Local: k}, Value: fmtString(v, arraySep)})
		}
		if err != nil {
			return err
		}
	}
	return enc.EncodeToken(start.End())
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (d *D) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var (
		attr = func(list []xml.Attr) map[string]string {
			m := make(map[string]string, len(list))
			for _, v := range list {
				m[v.Value] = v.Name.Local
			}
			return m
		}(start.Attr)
		tree       = []string{xmlName(start.Name, attr)}
		temp       = make(map[string]interface{})
		name, data string
		grow       bool
	)
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		switch t := token.(type) {
		case xml.StartElement:
			tree = append(tree, xmlName(t.Name, attr))
			grow = true
		case xml.CharData:
			data = string(t)
		case xml.EndElement:
			name, tree = tree[len(tree)-1], tree[:len(tree)-1]
			if !grow {
				continue
			}
			temp[strings.Join(append(tree, name), xmlLevelSep)] = data
			grow = false
		}
	}
	d.D = make(map[string]interface{})
	return expanded(temp, d.D)
}

func expanded(in, out map[string]interface{}) error {
	var (
		a  []string
		mv = func(m map[string]interface{}, to []string) map[string]interface{} {
			for i := 0; i < len(to)-1; i++ {
				_, ok := m[to[i]]
				if !ok {
					m[to[i]] = make(map[string]interface{})
				}
				m = m[to[i]].(map[string]interface{})
			}
			return m
		}
	)
	for k, v := range in {
		a = strings.Split(k, xmlLevelSep)
		mv(out, a[1:])[a[len(a)-1]] = v
	}
	return nil
}

const (
	xmlNSSep    = ":"
	xmlLevelSep = ">"
)

func xmlName(name xml.Name, space map[string]string) string {
	if ns, ok := space[name.Space]; ok {
		return ns + xmlNSSep + name.Local
	}
	return name.Local
}

// Bool forces the returned value behind these keys as a bool.
// An error is returned if the key does not exist or if the requested type is wrong.
func (d D) Bool(keys ...string) (bool, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return false, err
	}
	return toBool(m)
}

// Float64 forces the returned value behind these keys as a float64.
// An error is returned if the key does not exist or if the requested type is wrong.
func (d D) Float64(keys ...string) (float64, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return 0, err
	}
	return toFloat64(m)
}

// Int64 forces the returned value behind these keys as an int64.
// An error is returned if the key does not exist or if the requested type is wrong.
func (d D) Int64(keys ...string) (int64, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return 0, err
	}
	return toInt64(m)
}

// String forces the returned value behind these keys as a string.
// An error is returned if the key does not exist or if the requested type is wrong.
func (d D) String(keys ...string) (string, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return "", err
	}
	return toString(m)
}

// Strings returns if exists, the content of the given key as a slice of strings.
func (d D) Strings(keys ...string) ([]string, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return nil, err
	}
	v, ok := m.([]interface{})
	if !ok {
		var x []string
		return nil, newErrOutOfRange(x, v)
	}
	a := make([]string, len(v))
	for k2, v2 := range v {
		a[k2], err = toString(v2)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// Time tries to return the value behind the key as a time.Time matching the given time layout.
func (d D) Time(layout string, keys ...string) (time.Time, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return time.Time{}, err
	}
	s, err := toString(m)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(layout, s)
}

// Uint64 forces the returned value behind these keys as an uint64.
// An error is returned if the key does not exist or if the requested type is wrong.
func (d D) Uint64(keys ...string) (uint64, error) {
	m, err := d.Lookup(keys...)
	if err != nil {
		return 0, err
	}
	return toUint64(m)
}
