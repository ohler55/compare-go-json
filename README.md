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
| Parser Test coverage            | --                 | 93%                | 21%                | 100%               | 57.4%              | 91.5%              |

 :boom: _gjson does not validate while parsing (try a number of 1.2e3e4) although it does catch that error in validation._

 :small_blue_diamond: _gjson has an alternative search feature_

[_Details of each feature listed are at the bottom of the page_](#Feature-Explanations)

# Benchmarks

```
Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
     json.Unmarshal          42561 ns/op        17778 B/op          334 allocs/op
       oj.Parse              18878 ns/op         5691 B/op          364 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Unmarshal          31988 ns/op        19654 B/op          451 allocs/op
 simdjson.Parse              51125 ns/op       141341 B/op          370 allocs/op
    gjson.ParseBytes         33225 ns/op        20039 B/op          175 allocs/op

       oj ███████████████▊ 2.25
 jsoniter █████████▎ 1.33
    gjson ████████▉ 1.28
     json ▓▓▓▓▓▓▓ 1.00
 simdjson █████▊ 0.83
 fastjson >>> not supported <<<

Validate string/[]byte
     json.Valid              12926 ns/op            0 B/op            0 allocs/op
       oj.Validate            3992 ns/op            0 B/op            0 allocs/op
 fastjson.Validate            4831 ns/op            0 B/op            0 allocs/op
 jsoniter.Valid               9334 ns/op         2184 B/op          100 allocs/op
 simdjson.Validate           28691 ns/op       118841 B/op           20 allocs/op
    gjson.Validate            4118 ns/op            0 B/op            0 allocs/op

       oj ██████████████████████▋ 3.24
    gjson █████████████████████▉ 3.14
 fastjson ██████████████████▋ 2.68
 jsoniter █████████▋ 1.38
     json ▓▓▓▓▓▓▓ 1.00
 simdjson ███▏ 0.45

Iterate tokens in a string/[]byte
     json.Decode             82299 ns/op        22600 B/op         1175 allocs/op
       oj.Tokenize            8826 ns/op         1976 B/op          156 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode             33946 ns/op        20360 B/op          456 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj █████████████████████████████████████████████████████████████████▎ 9.32
 jsoniter ████████████████▉ 2.42
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Marshal to string/[]byte
     json.Marshal            69970 ns/op        26977 B/op          352 allocs/op
       oj.JSON               11278 ns/op         4096 B/op            1 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Marshal            14886 ns/op         6298 B/op           63 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj ███████████████████████████████████████████▍ 6.20
 jsoniter ████████████████████████████████▉ 4.70
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Read from single JSON file
     json.Decode             51254 ns/op        32417 B/op          342 allocs/op
       oj.ParseReader        21984 ns/op         9788 B/op          365 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode             40161 ns/op        20330 B/op          456 allocs/op
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

       oj ████████████████▎ 2.33
 jsoniter ████████▉ 1.28
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> not supported <<<
    gjson >>> not supported <<<

Read multiple JSON in a small log file (100MB)
     json.Decode        1606279141 ns/op   1102468376 B/op     14810435 allocs/op
       oj.ParseReader    687538598 ns/op    761293120 B/op     15402844 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode        1234408088 ns/op   1178778576 B/op     19390009 allocs/op
 simdjson.ParseReader    735606094 ns/op   1414514684 B/op     14824078 allocs/op
    gjson >>> not supported <<<

       oj ████████████████▎ 2.34
 simdjson ███████████████▎ 2.18
 jsoniter █████████  1.30
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
    gjson >>> not supported <<<

Read multiple JSON in a semi large log file (5GB)
     json.Decode       75877623348 ns/op  28649462840 B/op    740520932 allocs/op
       oj.ParseReader  34969385723 ns/op  11590640008 B/op    770141203 allocs/op
 fastjson >>> not supported <<<
 jsoniter.Decode       67916024904 ns/op  32465353640 B/op    969505876 allocs/op
 simdjson.ParseReader >>> out of memory <<<
    gjson >>> not supported <<<

       oj ███████████████▏ 2.17
 jsoniter ███████▊ 1.12
     json ▓▓▓▓▓▓▓ 1.00
 fastjson >>> not supported <<<
 simdjson >>> out of memory <<<
    gjson >>> not supported <<<

 Higher values (longer bars) are better in all cases. The bar graph compares the
 parsing performance. The lighter colored bar is the reference, the go json
 package.

Tests run on:
 Machine:         MacBookPro15,2
 OS:              macOS 11.2.3
 Processor:       Quad-Core Intel Core i7
 Cores:           4
 Processor Speed: 2.8 GHz
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
