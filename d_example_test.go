// Copyright (c) 2021 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package flat_test

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/rvflash/flat"

	"gopkg.in/yaml.v3"
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

func ExampleD_UnmarshalXML() {
	var (
		d   = flat.D{}
		err = xml.Unmarshal([]byte(`<custom><languages><fr>French</fr></languages></custom>`), &d)
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%#v", d.Flatten())
	// Output:
	// map[string]interface {}{"languages_fr":"French"}
}

func ExampleD_UnmarshalYAML() {
	var (
		d   = flat.D{}
		err = yaml.Unmarshal([]byte(`db:
    host: localhost
    name: database
    user:
        login: root
        pass: "insecure"
http:
    timeout: 0`), &d)
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%#v", d.Flatten())
	// Output:
	// map[string]interface {}{"db_host":"localhost", "db_name":"database", "db_user_login":"root", "db_user_pass":"insecure", "http_timeout":0}
}
