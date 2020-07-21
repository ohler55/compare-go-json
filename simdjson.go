// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/minio/simdjson-go"
)

var simdjsonPkg = pkg{
	name: "simdjson",
	calls: map[string]*call{
		"parse":    {name: "Parse", fun: simdjsonParse},
		"validate": {name: "Validate", fun: simdjsonValidate},
	},
}

func simdjsonParse(b *testing.B) {
	if !simdjson.SupportedCPU() {
		benchErr = errors.New("Unsupported CPU by simdjson")
		b.Fail()
	}
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	var pj simdjson.ParsedJson
	tmp := &simdjson.Iter{}
	obj := &simdjson.Object{}
	ary := &simdjson.Array{}
	for n := 0; n < b.N; n++ {
		parsed, err := simdjson.Parse(sample, &pj)
		if err != nil {
			benchErr = err
			b.Fail()
		}
		iter := parsed.Iter()
		typ := iter.Advance()
		switch typ {
		case simdjson.TypeRoot:
			if typ, tmp, benchErr = iter.Root(tmp); benchErr != nil {
				b.Fail()
			}
			switch typ {
			case simdjson.TypeArray:
				if ary, benchErr = tmp.Array(ary); benchErr != nil {
					b.Fail()
				}
				if _, benchErr = ary.Interface(); benchErr != nil {
					b.Fail()
				}
			case simdjson.TypeObject:
				if obj, benchErr = tmp.Object(obj); benchErr != nil {
					b.Fail()
				}
				var m map[string]interface{}
				if m, benchErr = obj.Map(m); benchErr != nil {
					b.Fail()
				}
			}
		}
	}
}

func simdjsonValidate(b *testing.B) {
	if !simdjson.SupportedCPU() {
		benchErr = errors.New("Unsupported CPU by simdjson")
		b.Fail()
	}
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	var pj simdjson.ParsedJson
	for n := 0; n < b.N; n++ {
		if _, benchErr = simdjson.Parse(sample, &pj); benchErr != nil {
			b.Fail()
		}
	}
}
