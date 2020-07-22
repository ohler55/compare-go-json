# Compare Go JSON

Not all JSON tools cover the same features which make it difficult to
select a set of tools for a project. Here is an attempt to compare
feature and benchmarks for a few of the JSON tools for Go.

## Features

| Feature                         | [go/json](https://golang.org/pkg/encoding/json/) | [fastjson](https://github.com/valyala/fastjson) | [jsoniter](https://github.com/json-iterator/go) | [OjG](https://github.com/ohler55/compare-go-json) | [simdjson](https://github.com/minio/simdjson-go) |
| ------------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| Parse []byte to simple go types | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Validate                        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Read from io.Reader             | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | ??                 |
| Read from file                  | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Parse to structs                | :white_check_mark: | :x:                | :x:                | :white_check_mark: | :x:                |
| Parse to interface types        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |
| Multiple JSON file/stream       | :white_check_mark: | :x:                | ??                 | :white_check_mark: | :x:                |
| ndjson (newline separated)      | :white_check_mark: | :x:                | ??                 | :white_check_mark: | :white_check_mark: |
| Marshal/Write                   | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                |
| JSONPath                        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |

[_Details of each feature listed are at the bottom of the page_](#Feature-Explanations)

# Benchmarks

```
Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
       json.Unmarshal       45071 ns/op   17985 B/op     336 allocs/op
         oj.Parse           18157 ns/op    5984 B/op     366 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Unmarshal       41758 ns/op   19799 B/op     451 allocs/op
   simdjson.Parse           64753 ns/op  136901 B/op     370 allocs/op

         oj ████████████████████████▊ 2.48
   jsoniter ██████████▊ 1.08
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson ██████▉ 0.70
   fastjson >>> not supported <<<

Validate string/[]byte
       json.Valid           11299 ns/op       0 B/op       0 allocs/op
         oj.Validate         3858 ns/op       0 B/op       0 allocs/op
   fastjson.Validate         4815 ns/op       0 B/op       0 allocs/op
   jsoniter.Valid            8737 ns/op    2192 B/op     100 allocs/op
   simdjson.Validate        27724 ns/op  114241 B/op      18 allocs/op

         oj █████████████████████████████▎ 2.93
   fastjson ███████████████████████▍ 2.35
   jsoniter ████████████▉ 1.29
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson ████  0.41

Marshal to string/[]byte
       json.Marshal         72439 ns/op   27344 B/op     352 allocs/op
         oj.JSON            11870 ns/op    4096 B/op       1 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Marshal         18504 ns/op    7291 B/op      94 allocs/op
   simdjson                 >>> not supported <<<

         oj █████████████████████████████████████████████████████████████  6.10
   jsoniter ███████████████████████████████████████▏ 3.91
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<

Read from single JSON file
       json.Decode          56546 ns/op   32655 B/op     346 allocs/op
         oj.ParseReader     21717 ns/op   10080 B/op     367 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode          37733 ns/op   20472 B/op     456 allocs/op
   simdjson                 >>> not supported <<<

         oj ██████████████████████████  2.60
   jsoniter ██████████████▉ 1.50
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, usually the go
 json package.

 The Benchmarks reflect a use case where JSON is either provided as a string or
 read from a file (io.Reader) then parsed into simple go types of nil, bool, int64
 float64, string, []interface{}, or map[string]interface{}. When supported, an
 io.Writer benchmark is also included along with some miscellaneous operations.

Tests run on:
 OS:              Ubuntu 18.04.4 LTS
 Processor:       Intel(R) Core(TM) i7-8700 CPU
 Cores:           12
 Processor Speed: 3.20GHz
 Memory:          16 GB
```

## Feature Explanations

 - **Parse** parse a string to []byte slice in simple go types of
   `[]interface`, `map[string]interface{}`, `string`, `float64`,
   `int64`, `bool`, or `nil`. This support the use case of extracting
   data from a JSON suitable for natigating as well as handing off to
   other packages such as a database for storage.

 - **Validate** a string or []byte slice without extracting values.

 - **Read from io.Reader** indicates a source such as a socket or file
   larger than will fit into memory can be parsed.

 - **Read from file** indicates a parser can read from a file if not
   directly then using ioutils.

 - **Parse to structs** is the ability to reconstitute a struct type
   from JSON.

 - **Parse to interface types** is the ability to reconstitutes types
   even if they are included as interfaces in a containing struct or
   slice.

 - **Multiple JSON** indicates a file or stream with multiple JSON
   documents can be parsed. This is no restricted to the limited case
   of exactly one JSON element per line. Encountered in database dumps
   and load files.

 - **ndjson** is a multiple document JSON where each JSON document
   must be on exactly one line. Found in log files.

 - **Marshal/Write** is the ability of the package to marshal go types
   in JSON.

 - **[JSONPath](https://goessner.net/articles/JsonPath)** is the
   ability to navigate data using JSONPath expressions.
