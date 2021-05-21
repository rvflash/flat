# Flat

[![GoDoc](https://godoc.org/github.com/rvflash/flat?status.svg)](https://godoc.org/github.com/rvflash/flat)
[![Build Status](https://api.travis-ci.com/rvflash/flat.svg?branch=main)](https://travis-ci.com/rvflash/flat?branch=main)
[![Code Coverage](https://codecov.io/gh/rvflash/flat/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/flat)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/flat?)](https://goreportcard.com/report/github.com/rvflash/flat)


The package `flat` provides methods to handle XML or JSON data as a `map[string]interface{}`, 
useful to manipulate unknown structures or to flatten them into a single dimension.


### Installation

To install it, you need to install Go and set your Go workspace first.
Then, download and install it:

```bash
$ go get -u github.com/rvflash/flat
```    
Import it in your code:

```go
import "github.com/rvflash/flat"
```


### Prerequisite

`flat` uses the Go modules that required Go 1.11 or later.


### XML Samples (see the example tests)

> Errors are ignored just for the demo.

#### Marshal

```go
var d = map[string]interface{}{
    "languages": map[string]interface{}{
        "en": "English",
        "fr": "French",
    },
}
var (
    res   = bytes.Buffer{}
    attrs = []xml.Attr{
        {Name: xml.Name{Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
    }
    _ = flat.New(
        d,
        flat.XMLName("custom"),
        flat.XMLNS("http://schemas.xmlsoap.org/soap/envelope/"),
        flat.XMLAttributes(attrs),
    ).XMLEncode(&res)
)
fmt.Println(res.String())
// Output:
// <custom xmlns="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><languages><en>English</en><fr>French</fr></languages></custom>
```

#### Unmarshal

```go
var (
    d   = flat.D{}
    _ = xml.Unmarshal([]byte(`<custom><languages><en>English</en><fr>French</fr></languages></custom>`), &d)
)
fmt.Printf("%#v", d.Flatten())
// Output:
// map[string]interface {}{"languages_en":"English", "languages_fr":"French"}
```
