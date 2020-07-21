// Copyright (c) 2020, Peter Ohler, All rights reserved.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strings"
	"testing"

	"github.com/ohler55/ojg/alt"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

const (
	blocks    = " ▏▎▍▌▋▊▉█"
	darkBlock = "▓"
)

var (
	filename = "data/patient.json"
	benchErr error
)

type specs struct {
	os        string
	model     string
	processor string
	cores     string
	speed     string
}

type call struct {
	name   string
	fun    func(b *testing.B)
	res    testing.BenchmarkResult
	ns     int64 // base adjusted
	bytes  int64 // base adjusted
	allocs int64 // base adjusted
	err    error
}

type pkg struct {
	name  string
	calls map[string]*call
}

type result struct {
	pkg  string
	call *call
	ref  bool
}

type suite struct {
	title string
	fun   string // key into the pkg calls
	ref   string // reference package for the suite
}

type noWriter int

func (w noWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func main() {
	testing.Init()
	flag.Parse()
	if 0 < len(flag.Args()) {
		filename = flag.Args()[0]
	}

	pkgs := []*pkg{
		&jsonPkg,
		&ojPkg,
		&fastjsonPkg,
		&jsoniterPkg,
		&simdjsonPkg,
	}
	for _, s := range []*suite{
		{fun: "parse", title: "Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)", ref: "json"},
		{fun: "validate", title: "Validate string/[]byte", ref: "json"},
		{fun: "marshal", title: "Marshal to string/[]byte", ref: "json"},
	} {
		s.exec(pkgs)
	}
	// TBD read from file (single json)
	// TBD read multiple json, single line (small, medium, large)
	// TBD read multiple json, indented small
	// TBD io.Reader multiple json, indented small
	// TBD write
	// TBD validate io.Reader

	fmt.Println()
	fmt.Println(" Higher values (longer bars) are better in all cases. The bar graph compares the")
	fmt.Println(" parsing performance. The lighter colored bar is the reference, usually the go")
	fmt.Println(" json package.")
	fmt.Println()
	fmt.Println(" The Benchmarks reflect a use case where JSON is either provided as a string or")
	fmt.Println(" read from a file (io.Reader) then parsed into simple go types of nil, bool, int64")
	fmt.Println(" float64, string, []interface{}, or map[string]interface{}. When supported, an")
	fmt.Println(" io.Writer benchmark is also included along with some miscellaneous operations.")
	fmt.Println()
	if s := getSpecs(); s != nil {
		fmt.Println("Tests run on:")
		if 0 < len(s.model) {
			fmt.Printf(" Machine:         %s\n", s.model)
		}
		fmt.Printf(" OS:              %s\n", s.os)
		fmt.Printf(" Processor:       %s\n", s.processor)
		fmt.Printf(" Cores:           %s\n", s.cores)
		fmt.Printf(" Processor Speed: %s\n", s.speed)
		// TBD add memory
	}
	fmt.Println()
}

func (s *suite) exec(pkgs []*pkg) {
	fmt.Println()
	fmt.Println(s.title)
	var results []*result
	var ref *call
	for _, p := range pkgs {
		benchErr = nil
		c := p.calls[s.fun]
		r := result{pkg: p.name, call: c, ref: s.ref == p.name}
		results = append(results, &r)
		if c == nil {
			r.call = &call{ns: math.MaxInt64, err: fmt.Errorf("not supported")}
			fmt.Printf(" %10s                 >>> not supported <<<\n", p.name)
			continue
		}
		if r.ref {
			ref = c
		}
		c.res = testing.Benchmark(c.fun)
		if benchErr != nil {
			c.err = benchErr
			c.ns = math.MaxInt64
			fmt.Printf(" %10s.%-14s >>> %s <<<\n", p.name, c.name, benchErr)
			continue
		}
		c.ns = c.res.NsPerOp()
		c.bytes = c.res.AllocedBytesPerOp()
		c.allocs = c.res.AllocsPerOp()
		fmt.Printf(" %10s.%-14s %6d ns/op  %6d B/op  %6d allocs/op\n",
			p.name, c.name, c.ns, c.bytes, c.allocs)
	}
	fmt.Println()
	scale := 10 // TBD adjust to fit screen better?
	sort.Slice(results, func(i, j int) bool { return results[i].call.ns < results[j].call.ns })
	for _, r := range results {
		c := r.call
		x := 1.0
		var bar string
		if r.pkg == s.ref {
			bar = strings.Repeat(darkBlock, scale)
		} else {
			if c.err == nil {
				x = float64(ref.ns) / float64(c.ns)
				size := x * float64(scale)
				bar = strings.Repeat(string([]rune(blocks)[8:]), int(size))
				frac := int(size*8.0) - (int(size) * 8)
				bar += string([]rune(blocks)[frac : frac+1])
			} else {
				fmt.Printf(" %10s >>> %s <<<\n", r.pkg, c.err)
				continue
			}
		}
		fmt.Printf(" %10s %s %3.2f\n", r.pkg, bar, x)
	}
}

func loadSample() (data interface{}) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to load %s. %s\n", filename, err)
	}
	defer func() { _ = f.Close() }()

	var p oj.Parser
	if data, err = p.ParseReader(f); err != nil {
		log.Fatalf("Failed to parse %s. %s\n", filename, err)
	}
	return
}

func getSpecs() (s *specs) {
	// Assume MacOS and try system_profiler. If that fails assume linux and check /proc.
	out, err := exec.Command("system_profiler", "-json", "SPHardwareDataType").Output()
	if err == nil {
		var js interface{}
		if js, err = oj.Parse(out); err == nil {
			s = &specs{
				model:     alt.String(jp.C("SPHardwareDataType").N(0).C("machine_model").First(js)),
				processor: alt.String(jp.C("SPHardwareDataType").N(0).C("cpu_type").First(js)),
				cores:     alt.String(jp.C("SPHardwareDataType").N(0).C("number_processors").First(js)),
				speed:     alt.String(jp.C("SPHardwareDataType").N(0).C("current_processor_speed").First(js)),
			}
			var b []byte
			if out, err = exec.Command("sw_vers", "-productName").Output(); err == nil {
				b = append(b, bytes.TrimSpace(out)...)
				b = append(b, ' ')
			}
			if out, err = exec.Command("sw_vers", "-productVersion").Output(); err == nil {
				b = append(b, bytes.TrimSpace(out)...)
			}
			s.os = string(b)
		}
		return
	}
	// Try Ubuntu next.
	if out, err = exec.Command("lsb_release", "-d").Output(); err == nil {
		s = &specs{}
		parts := strings.Split(string(out), ":")
		if 1 < len(parts) {
			s.os = string(strings.TrimSpace(parts[1]))
		}
		if out, err = ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			cnt := 0
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "processor") {
					cnt++
				} else if strings.Contains(line, "model name") {
					parts := strings.Split(line, ":")
					if 1 < len(parts) {
						parts = strings.Split(parts[1], "@")
						s.processor = strings.TrimSpace(parts[0])
						if 1 < len(parts) {
							s.speed = strings.TrimSpace(parts[1])
						}
					}
				}
				s.cores = fmt.Sprintf("%d", cnt)
			}
		}
	}
	return
}
