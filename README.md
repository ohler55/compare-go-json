# Compare Go JSON

Not all JSON tools cover the same features which make it difficult to
select a set of tools for a project. Here is an attempt to compare
features and benchmarks for a few of the JSON tools for Go.

## Features

| Feature                         | [go/json](https://golang.org/pkg/encoding/json/) | [fastjson](https://github.com/valyala/fastjson) | [jsoniter](https://github.com/json-iterator/go) | [OjG](https://github.com/ohler55/ojg) | [simdjson](https://github.com/minio/simdjson-go) |
| ------------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| Parse []byte to simple go types | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Validate                        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Parse - io.Reader (large file)  | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                 |
| Parse from file                 | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Parse to structs                | :white_check_mark: | :x:                | :x:                | :white_check_mark: | :x:                |
| Parse to interface types        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |
| Multiple JSON file/stream       | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                |
| ndjson (newline separated)      | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Marshal/Write                   | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                |
| JSON Builder                    | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |
| JSONPath                        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |
| Data type converters            | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |
| Simple Encoding Notation        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |

[_Details of each feature listed are at the bottom of the page_](#Feature-Explanations)

# Benchmarks

```
Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
       json.Unmarshal                 48892 ns/op             17984 B/op               336 allocs/op
         oj.Parse                     19196 ns/op              5984 B/op               366 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Unmarshal                 40239 ns/op             19800 B/op               451 allocs/op
   simdjson.Parse                     90506 ns/op            136901 B/op               370 allocs/op

         oj █████████████████████████▍ 2.55
   jsoniter ████████████▏ 1.22
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson █████▍ 0.54
   fastjson >>> not supported <<<

Validate string/[]byte
       json.Valid                     10981 ns/op                 0 B/op                 0 allocs/op
         oj.Validate                   3681 ns/op                 0 B/op                 0 allocs/op
   fastjson.Validate                   4671 ns/op                 0 B/op                 0 allocs/op
   jsoniter.Valid                      9105 ns/op              2192 B/op               100 allocs/op
   simdjson.Validate                  26411 ns/op            114240 B/op                18 allocs/op

         oj █████████████████████████████▊ 2.98
   fastjson ███████████████████████▌ 2.35
   jsoniter ████████████  1.21
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson ████▏ 0.42

Marshal to string/[]byte
       json.Marshal                   74607 ns/op             27342 B/op               352 allocs/op
         oj.JSON                      11939 ns/op              4096 B/op                 1 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Marshal                   19446 ns/op              7291 B/op                94 allocs/op
   simdjson                 >>> not supported <<<

         oj ██████████████████████████████████████████████████████████████▍ 6.25
   jsoniter ██████████████████████████████████████▎ 3.84
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<

Read from single JSON file
       json.Decode                    72469 ns/op             32624 B/op               344 allocs/op
         oj.ParseReader               23420 ns/op             10080 B/op               367 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode                    46320 ns/op             20472 B/op               456 allocs/op
   simdjson                 >>> not supported <<<

         oj ██████████████████████████████▉ 3.09
   jsoniter ███████████████▋ 1.56
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<

Read multiple JSON in a small log file (100MB)
       json.Decode               1348331757 ns/op         855642218 B/op          14810403 allocs/op
         oj.ParseReader           649038041 ns/op         518047360 B/op          15995224 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode               1079811228 ns/op         927208456 B/op          19390064 allocs/op
   simdjson.ParseReader           652181597 ns/op        1285476585 B/op          15416464 allocs/op

         oj ████████████████████▊ 2.08
   simdjson ████████████████████▋ 2.07
   jsoniter ████████████▍ 1.25
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<

Read multiple JSON in a semi large log file (5GB)
       json.Decode              80493161639 ns/op       29360272024 B/op         740519974 allocs/op
         oj.ParseReader         36165605980 ns/op       12480784192 B/op         799761453 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode              77476413417 ns/op       32939373944 B/op         969507185 allocs/op
   simdjson.ParseReader    >>> out of memory <<<

         oj ██████████████████████▎ 2.23
   jsoniter ██████████▍ 1.04
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> out of memory <<<

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, the go json
 package.

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

 - **JSON Builder** is the ability to create new data structures suitable for JSON encoding.

 - **[JSONPath](https://goessner.net/articles/JsonPath)** is the
   ability to navigate data using JSONPath expressions.

 - **Data type converters** tools for converting from type to simple
   data types. Basically marshalling and unmarshalling to simple types
   instead of to JSON.

 - **[Simple Encoding Notation](https://github.com/ohler55/ojg/blob/develop/sen.md)** is
   a lazy JSON format where quotes and commas are optional in most
   cases. A merge of JSON and GraphQL formats for those of us that
   don't want to be bothered with strict syntax checking.
