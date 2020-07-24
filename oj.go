// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/ohler55/ojg/oj"
)

var ojPkg = pkg{
	name: "oj",
	calls: map[string]*call{
		"parse":      {name: "Parse", fun: ojParse},
		"validate":   {name: "Validate", fun: ojValidate},
		"marshal":    {name: "JSON", fun: ojJSON},
		"file1":      {name: "ParseReader", fun: ojFile1},
		"small-file": {name: "ParseReader", fun: ojFileManySmallLoad},
		"large-file": {name: "ParseReader", fun: ojFileManyLarge},
	},
}

func ojParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		if _, benchErr = p.Parse(sample); benchErr != nil {
			b.Fail()
		}
	}
}

func ojValidate(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var v oj.Validator
	for n := 0; n < b.N; n++ {
		if benchErr = v.Validate(sample); benchErr != nil {
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

func ojFile1(b *testing.B) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, benchErr = p.ParseReader(f); benchErr != nil {
			b.Fail()
		}
	}
}

func ojFileManySmallChan(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()

	rc := make(chan interface{}, 1000)
	ready := make(chan bool)
	go func() {
		ready <- true
		for {
			if v := <-rc; v == nil {
				break
			}
		}
	}()
	<-ready
	b.ResetTimer()

	var p oj.Parser
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, benchErr = p.ParseReader(f, rc); benchErr != nil {
			b.Fail()
		}
	}
	rc <- nil
}

func ojCb(_ interface{}) bool {
	return false
}

func ojFileManySmallReader(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, benchErr = p.ParseReader(f, ojCb); benchErr != nil {
			b.Fail()
		}
	}
}

func ojFileManySmallLoad(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		j, _ := ioutil.ReadAll(f)
		if _, benchErr = p.Parse(j, ojCb); benchErr != nil {
			b.Fail()
		}
	}
}

func ojFileManyLarge(b *testing.B) {
	f := openLargeLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	p := &oj.Parser{Reuse: true}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		if _, benchErr = p.ParseReader(f, ojCb); benchErr != nil {
			b.Fail()
		}
	}
}
