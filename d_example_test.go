// Copyright (c) 2020 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat_test

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/rvflash/flat"
)

func ExampleD_XMLEncode() {
	var d = map[string]interface{}{
		"languages": map[string]interface{}{
			"fr": "French",
		},
	}
	var (
		res   = bytes.Buffer{}
		attrs = []xml.Attr{
			{Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
		}
		err = flat.New(
			d,
			flat.XMLName("custom"),
			flat.XMLNS("http://schemas.xmlsoap.org/soap/envelope/"),
			flat.XMLAttributes(attrs),
		).XMLEncode(&res)
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res.String())
	// Output:
	// <custom xmlns="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><languages><fr>French</fr></languages></custom>
}
