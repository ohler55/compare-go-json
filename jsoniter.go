// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"errors"
	"io/ioutil"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var jsoniterPkg = pkg{
	name: "jsoniter",
	calls: map[string]*call{
		"parse":    {name: "Unmarshal", fun: jsoniterUnmarshal},
		"validate": {name: "Valid", fun: jsoniterValid},
		"marshal":  {name: "Marshal", fun: jsoniterMarshal},
	},
}

func jsoniterUnmarshal(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	var result interface{}
	for n := 0; n < b.N; n++ {
		if benchErr := jsoniter.Unmarshal(sample, &result); benchErr != nil {
			b.Fail()
		}
	}
}

func jsoniterValid(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if !jsoniter.Valid(sample) {
			benchErr = errors.New("JSON not valid")
			b.Fail()
		}
	}
}

func jsoniterMarshal(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, benchErr = jsoniter.Marshal(data); benchErr != nil {
			b.Fail()
		}
	}
}
