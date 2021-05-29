# Compare Go JSON

Not all JSON tools cover the same features which make it difficult to
select a set of tools for a project. Here is an attempt to compare
features and benchmarks for a few of the JSON tools for Go.

## Features

| Feature                         | [go/json](https://golang.org/pkg/encoding/json/) | [fastjson](https://github.com/valyala/fastjson) | [jsoniter](https://github.com/json-iterator/go) | [OjG](https://github.com/ohler55/ojg) | [simdjson](https://github.com/minio/simdjson-go) | [gjson](https://github.com/tidwall/gjson)
| ------------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| Parse []byte to simple go types | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark::boom: |
| Validate                        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Parse - io.Reader (large file)  | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| Parse from file                 | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Parse to structs                | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| Parse to interface types        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Multiple JSON file/stream       | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| ndjson (newline separated)      | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Marshal/Write                   | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| JSON Builder                    | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| JSONPath                        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x::small_blue_diamond: |
| Data type converters            | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Simple Encoding Notation        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Test coverage                   | --                 | 93%                | 21%                | 100%               | 57.4%              | 91.5%              |

 :boom: _gjson does not validate while parsing (try a number of 1.2e3e4) although it does catch that error in validation._

 :small_blue_diamond: _gjson has an alternative search feature_

[_Details of each feature listed are at the bottom of the page_](#Feature-Explanations)

# Benchmarks

```
Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
     json.Unmarshal          50142 ns/op        17776 B/op          334 allocs/op
       oj.Parse              19453 ns/op         5691 B/op          364 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Unmarshal          40849 ns/op        19656 B/op          451 allocs/op
 simdjson.Parse             119196 ns/op       141345 B/op          370 allocs/op
    gjson.ParseBytes         53440 ns/op        20039 B/op          175 allocs/op

       oj ██████████████████  2.58
 jsoniter ████████▌ 1.23
     json ▓▓▓▓▓▓▓ 1.00
    gjson ██████▌ 0.94
 simdjson ██▉ 0.42
 fastjson >>> not supported <<<

Validate string/[]byte
     json.Valid              12034 ns/op            0 B/op            0 allocs/op
       oj.Validate            3802 ns/op            0 B/op            0 allocs/op
 fastjson.Validate            4521 ns/op            0 B/op            0 allocs/op
 jsoniter.Valid               8699 ns/op         2184 B/op          100 allocs/op
 simdjson.Validate           26498 ns/op       118842 B/op           20 allocs/op
    gjson.Validate            3904 ns/op            0 B/op            0 allocs/op

       oj ██████████████████████▏ 3.17
    gjson █████████████████████▌ 3.08
 fastjson ██████████████████▋ 2.66
 jsoniter █████████▋ 1.38
     json ▓▓▓▓▓▓▓ 1.00
 simdjson ███▏ 0.45

Iterate tokens in a string/[]byte
     json.Decode             90142 ns/op        22600 B/op         1175 allocs/op
       oj.Tokenize            8049 ns/op         1976 B/op          156 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode             48138 ns/op        20360 B/op          456 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj ██████████████████████████████████████████████████████████████████████████████▍ 11.20
 jsoniter █████████████  1.87
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Unmarshal string/[]byte to a struct
     json.Unmarshal          39226 ns/op         2560 B/op           73 allocs/op
       oj.Unmarshal          41857 ns/op         9355 B/op          456 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Unmarshal          13005 ns/op         3137 B/op          170 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

 jsoniter █████████████████████  3.02
     json ▓▓▓▓▓▓▓ 1.00
       oj ██████▌ 0.94
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Marshal simple types to string/[]byte
     json.Marshal            51471 ns/op        17559 B/op          345 allocs/op
       oj.JSON                6615 ns/op            0 B/op            0 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Marshal            18855 ns/op         6299 B/op           63 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj ██████████████████████████████████████████████████████▍ 7.78
 jsoniter ███████████████████  2.73
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Marshal a struct to string/[]byte
     json.Marshal            11353 ns/op         3457 B/op            1 allocs/op
       oj.Marshal             8453 ns/op         1712 B/op           44 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Marshal            11315 ns/op         3466 B/op            2 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj █████████▍ 1.34
 jsoniter ███████  1.00
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Read from single JSON file
     json.Decode             85118 ns/op        32416 B/op          342 allocs/op
       oj.ParseReader        28191 ns/op         9787 B/op          365 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode             59635 ns/op        20327 B/op          456 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj █████████████████████▏ 3.02
 jsoniter █████████▉ 1.43
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Read multiple JSON in a small log file (100MB)
     json.Decode        1354784239 ns/op   1102470120 B/op     14810436 allocs/op
       oj.ParseReader    668270061 ns/op    761294040 B/op     15402846 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode        1131730647 ns/op   1178775816 B/op     19390099 allocs/op
 simdjson.ParseReader    679189142 ns/op   1398403912 B/op     13639254 allocs/op
    gjson >>> not supported <<<

       oj ██████████████▏ 2.03
 simdjson █████████████▉ 1.99
 jsoniter ████████▍ 1.20
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
    gjson >>> not supported <<<

Read multiple JSON in a semi large log file (5GB)
     json.Decode       88433748107 ns/op  28649383232 B/op    740520188 allocs/op
       oj.ParseReader  38353089960 ns/op  11590577208 B/op    770140750 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode       86478665854 ns/op  32465480800 B/op    969507711 allocs/op
 simdjson.ParseReader >>> out of memory <<<
    gjson >>> not supported <<<

       oj ████████████████▏ 2.31
 jsoniter ███████▏ 1.02
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> out of memory <<<
    gjson >>> not supported <<<

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, the go json
 package.

Tests run on:
 OS:              Ubuntu 20.04.2 LTS
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

 - **Parser Test coverage** percent unit test coverage of the parser
   package. It does not include coverage of other package in the
   offering.
