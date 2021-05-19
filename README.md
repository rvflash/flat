# Flat

[![GoDoc](https://godoc.org/github.com/rvflash/flat?status.svg)](https://godoc.org/github.com/rvflash/flat)
[![Build Status](https://api.travis-ci.org/rvflash/flat.svg?branch=main)](https://travis-ci.org/rvflash/flat?branch=main)
[![Code Coverage](https://codecov.io/gh/rvflash/flat/branch/main/graph/badge.svg)](https://codecov.io/gh/rvflash/flat)
[![Go Report Card](https://goreportcard.com/badge/github.com/rvflash/flat?)](https://goreportcard.com/report/github.com/rvflash/flat)


The package `flat` provides methods to handle XML or JSON data as a `map[string]interface`, 
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