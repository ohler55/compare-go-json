// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"io/ioutil"
	"testing"

	"github.com/valyala/fastjson"
)

var fastjsonPkg = pkg{
	name: "fastjson",
	calls: map[string]*call{
		"validate": {name: "Validate", fun: fastjsonValidate},
	},
}

func fastjsonValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if benchErr = fastjson.ValidateBytes(sample); benchErr != nil {
			b.Fail()
		}
	}
}
