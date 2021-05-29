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
	"strconv"
	"strings"
	"testing"
	"time"

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

	smallLogFile = "data/log-small.json"
	smallSize    = 100

	largeLogFile = "data/log-large.json"
	largeSize    = 5000
)

type specs struct {
	os        string
	model     string
	processor string
	cores     string
	speed     string
	memory    string
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
		&gjsonPkg,
	}
	for _, s := range []*suite{
		{fun: "parse", title: "Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)", ref: "json"},
		{fun: "validate", title: "Validate string/[]byte", ref: "json"},
		{fun: "decode", title: "Iterate tokens in a string/[]byte", ref: "json"},
		{fun: "unmarshal-struct", title: "Unmarshal string/[]byte to a struct", ref: "json"},
		{fun: "marshal", title: "Marshal simple types to string/[]byte", ref: "json"},
		{fun: "marshal-struct", title: "Marshal a struct to string/[]byte", ref: "json"},
		{fun: "file1", title: "Read from single JSON file", ref: "json"},
		{fun: "small-file", title: "Read multiple JSON in a small log file (100MB)", ref: "json"},
		{fun: "large-file", title: "Read multiple JSON in a semi large log file (5GB)", ref: "json"},
	} {
		s.exec(pkgs)
	}
	// TBD read multiple json, indented small, maybe a few patients in one file
	// TBD validate io.Reader

	fmt.Println()
	fmt.Println(" Higher values (longer bars) are better in all cases. The bar graph compares the")
	fmt.Println(" parsing performance. The lighter colored bar is the reference, the go json")
	fmt.Println(" package.")
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
		fmt.Printf(" Memory:          %s\n", s.memory)
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
			fmt.Printf(" %8s >>> not supported <<<\n", p.name)
			continue
		}
		if r.ref {
			ref = c
		}
		c.res = testing.Benchmark(c.fun)
		if benchErr != nil {
			c.err = benchErr
			c.ns = math.MaxInt64
			fmt.Printf(" %8s.%-11s >>> %s <<<\n", p.name, c.name, benchErr)
			continue
		}
		c.ns = c.res.NsPerOp()
		c.bytes = c.res.AllocedBytesPerOp()
		c.allocs = c.res.AllocsPerOp()
		fmt.Printf(" %8s.%-11s %12d ns/op %12d B/op %12d allocs/op\n",
			p.name, c.name, c.ns, c.bytes, c.allocs)
	}
	fmt.Println()
	scale := 7 // TBD adjust to fit screen better?
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
				fmt.Printf(" %8s >>> %s <<<\n", r.pkg, c.err)
				continue
			}
		}
		fmt.Printf(" %8s %s %3.2f\n", r.pkg, bar, x)
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

func openSmallLogFile() *os.File {
	f, err := os.Open(smallLogFile)
	if err != nil {
		if err = createLogFile(smallLogFile, smallSize); err != nil {
			log.Fatalf("Failed to create %s. %s\n", smallLogFile, err)
		}
		if f, err = os.Open(smallLogFile); err != nil {
			log.Fatalf("Failed to open %s. %s\n", smallLogFile, err)
		}
	}
	return f
}

func openLargeLogFile() *os.File {
	f, err := os.Open(largeLogFile)
	if err != nil {
		if err = createLogFile(largeLogFile, largeSize); err != nil {
			log.Fatalf("Failed to create %s. %s\n", largeLogFile, err)
		}
		if f, err = os.Open(largeLogFile); err != nil {
			log.Fatalf("Failed to open %s. %s\n", largeLogFile, err)
		}
	}
	return f
}

// size is in MB.
func createLogFile(filename string, size int) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Build a log entry.
	var b oj.Builder
	_ = b.Object()
	_ = b.Value(time.Now().UnixNano(), "when")
	_ = b.Value("Just some fake log entry for a generated log file.", "what")
	_ = b.Array("where")
	_ = b.Object()
	_ = b.Value("example.go", "file")
	_ = b.Value(123, "line")
	b.Pop()
	b.Pop()
	_ = b.Value("benchmark-application", "who")
	_ = b.Value("INFO", "level")
	b.PopAll()
	entry := b.Result()

	var whenX jp.Expr
	if whenX, err = jp.Parse([]byte("when")); err != nil {
		return err
	}
	j := oj.JSON(entry)
	cnt := size * 1024 * 1024 / (len(j) + 1)
	for i := 0; i < cnt; i++ {
		// Update entry.
		if err = whenX.Set(entry, time.Now().UnixNano()); err != nil {
			return err
		}
		if err = oj.Write(f, entry); err != nil {
			return err
		}
		_, _ = f.Write([]byte{'\n'})
	}
	return nil
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
				memory:    alt.String(jp.C("SPHardwareDataType").N(0).C("physical_memory").First(js)),
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
		if out, err = ioutil.ReadFile("/proc/meminfo"); err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.Contains(line, "MemTotal") {
					parts := strings.Split(line, ":")
					if 1 < len(parts) {
						s.memory = strings.TrimSpace(parts[1])
						if strings.HasSuffix(s.memory, "kB") {
							if i, err := strconv.Atoi(strings.Split(s.memory, " ")[0]); err == nil {
								s.memory = fmt.Sprintf("%d GB", i/1000000)
							}
						}
					}
				}
			}
		}
	}
	return
}
