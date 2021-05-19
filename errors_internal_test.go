// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat

import (
	"testing"

	"github.com/matryer/is"
)

func TestErrFlat_Error(t *testing.T) {
	is.New(t).Equal("flat: not found", ErrNotFound.Error())
}

func TestNewErrOutOfRange(t *testing.T) {
	var (
		x bool
		g float64
	)
	is.New(t).Equal("flat: wrong data type: bool expected, got float64", newErrOutOfRange(x, g).Error())
}
