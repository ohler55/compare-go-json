// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/tidwall/gjson"
)

var gjsonPkg = pkg{
	name: "gjson",
	calls: map[string]*call{
		"parse":    {name: "ParseBytes", fun: gjsonParse},
		"validate": {name: "Validate", fun: gjsonValid},
	},
}

func gjsonParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = gjson.ParseBytes(sample).Value()
	}
}

func gjsonValid(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if !gjson.ValidBytes(sample) {
			benchErr = errors.New("JSON not valid")
			b.Fail()
		}
	}
}
