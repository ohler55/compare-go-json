// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/minio/simdjson-go"
)

var simdjsonPkg = pkg{
	name: "simdjson",
	calls: map[string]*call{
		"parse":      {name: "Parse", fun: simdjsonParse},
		"validate":   {name: "Validate", fun: simdjsonValidate},
		"small-file": {name: "ParseReader", fun: simdjsonFileManySmall},
		"large-file": {name: "ParseReader", fun: simdjsonFileManyLarge},
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
	for n := 0; n < b.N; n++ {
		parsed, err := simdjson.Parse(sample, &pj)
		if err != nil {
			benchErr = err
			b.Fail()
			break
		}
		if benchErr = simdjsonExtract(parsed); benchErr != nil {
			b.Fail()
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

func simdjsonFileManySmall(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// simdjson closes the chan when done parsing so a new one has to be
		// created on each parse.
		done := make(chan bool)
		rc := make(chan simdjson.Stream, 1000)
		go func() {
			cnt := 0
			for {
				v := <-rc
				cnt++
				if v.Error != nil {
					if v.Error != io.EOF {
						benchErr = v.Error
						b.Fail()
					}
					break
				}
				if benchErr = simdjsonExtract(v.Value); benchErr != nil {
					b.Fail()
				}
			}
			done <- true
		}()
		_, _ = f.Seek(0, 0)
		simdjson.ParseNDStream(f, rc, nil)
		<-done
	}
}

func simdjsonFileManyLarge(b *testing.B) {
	// On larger files such as the 5GB file used for a large file (not that
	// large really) simdjson apparently attempts to pull the whole file into
	// memory and which causes an out of memory error or kills the
	// application.
	benchErr = errors.New("out of memory")
	b.Fail()
	/*
		f := openLargeLogFile()
		defer func() { _ = f.Close() }()

		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			// simdjson closes the chan when done parsing so a new one has to be
			// created on each parse.
			done := make(chan bool)
			rc := make(chan simdjson.Stream, 1000)
			go func() {
				cnt := 0
				for {
					v := <-rc
					cnt++
					if v.Error != nil {
						if v.Error != io.EOF {
							benchErr = v.Error
							b.Fail()
						}
						break
					}
					if benchErr = simdjsonExtract(v.Value); benchErr != nil {
						b.Fail()
					}
				}
				done <- true
			}()
			_, _ = f.Seek(0, 0)
			simdjson.ParseNDStream(f, rc, nil)
			<-done
		}
	*/
}

func simdjsonExtract(pj *simdjson.ParsedJson) (err error) {
	tmp := &simdjson.Iter{}

	iter := pj.Iter()
	for {
		typ := iter.Advance()
		switch typ {
		case simdjson.TypeRoot:
			if typ, tmp, err = iter.Root(tmp); err != nil {
				return
			}
			switch typ {
			case simdjson.TypeArray:
				ary := &simdjson.Array{}
				if ary, err = tmp.Array(ary); err != nil {
					return
				}
				if _, err = ary.Interface(); err != nil {
					return
				}
			case simdjson.TypeObject:
				obj := &simdjson.Object{}
				if obj, err = tmp.Object(obj); err != nil {
					return
				}
				var m map[string]interface{}
				if m, err = obj.Map(m); err != nil {
					return
				}
			}
		default:
			return
		}
	}
	return
}
