// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
)

var jsonPkg = pkg{
	name: "json",
	calls: map[string]*call{
		"parse":    {name: "Unmarshal", fun: goParse},
		"validate": {name: "Valid", fun: goValidate},
		"marshal":  {name: "Marshal", fun: goMarshal},
	},
}

func goParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var result interface{}
	for n := 0; n < b.N; n++ {
		if benchErr := json.Unmarshal(sample, &result); benchErr != nil {
			b.Fail()
		}
	}
}

func goValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !json.Valid(sample) {
			benchErr = errors.New("JSON not valid")
			b.Fail()
		}
	}
}

func goMarshal(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, benchErr := json.MarshalIndent(data, "", "  "); benchErr != nil {
			b.Fail()
		}
	}
}
