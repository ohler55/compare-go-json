// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var jsonPkg = pkg{
	name: "json",
	calls: map[string]*call{
		"parse":      {name: "Unmarshal", fun: goParse},
		"validate":   {name: "Valid", fun: goValidate},
		"decode":     {name: "Decode", fun: goDecode},
		"marshal":    {name: "Marshal", fun: goMarshal},
		"file1":      {name: "Decode", fun: goFile1},
		"small-file": {name: "Decode", fun: goFileManySmallLoad},
		"large-file": {name: "Decode", fun: goFileManyLarge},
	},
}

func goParse(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	var result interface{}
	for n := 0; n < b.N; n++ {
		if benchErr = json.Unmarshal(sample, &result); benchErr != nil {
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

func goDecode(b *testing.B) {
	sample, _ := ioutil.ReadFile(filename)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		dec := json.NewDecoder(bytes.NewReader(sample))
		for {
			_, err := dec.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}

func goMarshal(b *testing.B) {
	data := loadSample()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, benchErr = json.MarshalIndent(data, "", "  "); benchErr != nil {
			b.Fail()
		}
	}
}

func goFile1(b *testing.B) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	var data interface{}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := json.NewDecoder(f)
		if err := dec.Decode(&data); err != nil && err != io.EOF {
			benchErr = err
			b.Fail()
		}
	}
}

func goFileManySmall(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	var data interface{}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := json.NewDecoder(f)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}

func goFileManySmallLoad(b *testing.B) {
	f := openSmallLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	var data interface{}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		j, _ := ioutil.ReadAll(f)
		dec := json.NewDecoder(bytes.NewReader(j))
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}

func goFileManyLarge(b *testing.B) {
	f := openLargeLogFile()
	defer func() { _ = f.Close() }()
	b.ResetTimer()
	var data interface{}
	for n := 0; n < b.N; n++ {
		_, _ = f.Seek(0, 0)
		dec := json.NewDecoder(f)
		for {
			if err := dec.Decode(&data); err == io.EOF {
				break
			} else if err != nil {
				benchErr = err
				b.Fail()
			}
		}
	}
}
