// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"io/ioutil"
	"testing"

	"github.com/ohler55/ojg/oj"
)

var ojPkg = pkg{
	name: "oj",
	calls: map[string]*call{
		"parse":    {name: "Parse", fun: ojParse},
		"validate": {name: "Validate", fun: ojValidate},
		"marshal":  {name: "JSON", fun: ojJSON},
	},
}

func ojParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		if _, benchErr := p.Parse(sample); benchErr != nil {
			b.Fail()
		}
	}
}

func ojValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var v oj.Validator
	for n := 0; n < b.N; n++ {
		if benchErr := v.Validate(sample); benchErr != nil {
			b.Fail()
		}
	}
}

func ojJSON(b *testing.B) {
	data := loadSample()
	opt := oj.Options{OmitNil: true, Indent: 2}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = oj.JSON(data, &opt)
	}
}
