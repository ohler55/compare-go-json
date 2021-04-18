// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var jsoniterPkg = pkg{
	name: "jsoniter",
	calls: map[string]*call{
		"parse":      {name: "Unmarshal", fun: jsoniterUnmarshal},
		"validate":   {name: "Valid", fun: jsoniterValid},
		"decode":     {name: "Decode", fun: jsoniterDecode},
		"marshal":    {name: "Marshal", fun: jsoniterMarshal},
		"file1":      {name: "Decode", fun: jsoniterFile1},
		"small-file": {name: "Decode", fun: jsoniterFileManySmall},
		"large-file": {name: "Decode", fun: jsoniterFileManyLarge},
	},
}

func jsoniterUnmarshal(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	var result interface{}
	for n := 0; n < b.N; n++ {
		if benchErr = jsoniter.Unmarshal(sample, &result); benchErr != nil {
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

func jsoniterDecode(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()

	var data interface{}
	for n := 0; n < b.N; n++ {
		dec := jsoniter.NewDecoder(bytes.NewReader(sample))
		if err := dec.Decode(&data); err != nil {
			benchErr = err
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

func jsoniterFile1(b *testing.B) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := jsoniter.NewDecoder(f)
		var data interface{}
		if err := dec.Decode(&data); err != nil {
			benchErr = err
			b.Fail()
		}
	}
}

func jsoniterFileManySmall(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		j, _ := ioutil.ReadAll(f)
		dec := jsoniter.NewDecoder(bytes.NewReader(j))
		for {
			var data interface{}
			if !dec.More() {
				break
			}
			if err := dec.Decode(&data); err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}

func jsoniterFileManyLarge(b *testing.B) {
	f := openLargeLogFile()
	defer func() { _ = f.Close() }()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := jsoniter.NewDecoder(f)
		for {
			var data interface{}
			if !dec.More() {
				break
			}
			if err := dec.Decode(&data); err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}
