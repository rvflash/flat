// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat

import "fmt"

type errFlat string

// Error implements the error interface.
func (s errFlat) Error() string {
	return "flat: " + string(s)
}

const (
	// ErrNotFound is returned when the key is unknown.
	ErrNotFound = errFlat("not found")
	// ErrOutFoRange is returned when the type of data requested does not correspond to that of the data.
	ErrOutOfRange = errFlat("wrong data type")
)

func newErrOutOfRange(exp, got interface{}) error {
	return fmt.Errorf("%w: %T expected, got %T", ErrOutOfRange, exp, got)
}
