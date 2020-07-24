# Compare Go JSON

Not all JSON tools cover the same features which make it difficult to
select a set of tools for a project. Here is an attempt to compare
features and benchmarks for a few of the JSON tools for Go.

## Features

| Feature                         | [go/json](https://golang.org/pkg/encoding/json/) | [fastjson](https://github.com/valyala/fastjson) | [jsoniter](https://github.com/json-iterator/go) | [OjG](https://github.com/ohler55/ojg) | [simdjson](https://github.com/minio/simdjson-go) | [gjson](https://github.com/tidwall/gjson)
| ------------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| Parse []byte to simple go types | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark:* |
| Validate                        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Parse - io.Reader (large file)  | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| Parse from file                 | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Parse to structs                | :white_check_mark: | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Parse to interface types        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Multiple JSON file/stream       | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| ndjson (newline separated)      | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Marshal/Write                   | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                | :x:                |
| JSON Builder                    | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| JSONPath                        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:**              |
| Data type converters            | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Simple Encoding Notation        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                |
| Parser Test coverage            | --                 | 98.1%              | 21.2%              | 100%               | 59.4%              | 91.8%              |

 * _gjson does not validate while parsing (try a number of 1.2e3e4)_
 ** _gjson has an alternative search feature_

[_Details of each feature listed are at the bottom of the page_](#Feature-Explanations)

# Benchmarks

```

Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
       json.Unmarshal                 45508 ns/op             17983 B/op               336 allocs/op
         oj.Parse                     18494 ns/op              5984 B/op               366 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Unmarshal                 40247 ns/op             19800 B/op               451 allocs/op
   simdjson.Parse                     87288 ns/op            136900 B/op               370 allocs/op
      gjson.ParseBytes                57299 ns/op             20176 B/op               175 allocs/op

         oj ████████████████████████▌ 2.46
   jsoniter ███████████▎ 1.13
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
      gjson ███████▉ 0.79
   simdjson █████▏ 0.52
   fastjson >>> not supported <<<

Validate string/[]byte
       json.Valid                     10931 ns/op                 0 B/op                 0 allocs/op
         oj.Validate                   3728 ns/op                 0 B/op                 0 allocs/op
   fastjson.Validate                   4605 ns/op                 0 B/op                 0 allocs/op
   jsoniter.Valid                      8800 ns/op              2192 B/op               100 allocs/op
   simdjson.Validate                  25959 ns/op            114241 B/op                18 allocs/op
      gjson.Validate                   3957 ns/op                 0 B/op                 0 allocs/op

         oj █████████████████████████████▎ 2.93
      gjson ███████████████████████████▌ 2.76
   fastjson ███████████████████████▋ 2.37
   jsoniter ████████████▍ 1.24
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson ████▏ 0.42

Marshal to string/[]byte
       json.Marshal                   89291 ns/op             27333 B/op               352 allocs/op
         oj.JSON                      12707 ns/op              4096 B/op                 1 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Marshal                   21749 ns/op              7291 B/op                94 allocs/op
   simdjson                 >>> not supported <<<
      gjson                 >>> not supported <<<

         oj ██████████████████████████████████████████████████████████████████████▎ 7.03
   jsoniter █████████████████████████████████████████  4.11
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<
      gjson >>> not supported <<<

Read from single JSON file
       json.Decode                    71909 ns/op             32625 B/op               344 allocs/op
         oj.ParseReader               29370 ns/op             10080 B/op               367 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode                    58657 ns/op             20471 B/op               456 allocs/op
   simdjson                 >>> not supported <<<
      gjson                 >>> not supported <<<

         oj ████████████████████████▍ 2.45
   jsoniter ████████████▎ 1.23
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> not supported <<<
      gjson >>> not supported <<<

Read multiple JSON in a small log file (100MB)
       json.Decode               1315178539 ns/op         855642232 B/op          14810405 allocs/op
         oj.ParseReader           636970161 ns/op         518047424 B/op          15995225 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode               1061236313 ns/op         927208492 B/op          19390064 allocs/op
   simdjson.ParseReader           640959154 ns/op        1285476984 B/op          15416465 allocs/op
      gjson                 >>> not supported <<<

         oj ████████████████████▋ 2.06
   simdjson ████████████████████▌ 2.05
   jsoniter ████████████▍ 1.24
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
      gjson >>> not supported <<<

Read multiple JSON in a semi large log file (5GB)
       json.Decode              86674796158 ns/op       29360262568 B/op         740519910 allocs/op
         oj.ParseReader         38972170247 ns/op       12480775848 B/op         799761415 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Decode              85041581614 ns/op       32939363048 B/op         969507128 allocs/op
   simdjson.ParseReader    >>> out of memory <<<
      gjson                 >>> not supported <<<

         oj ██████████████████████▏ 2.22
   jsoniter ██████████▏ 1.02
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   fastjson >>> not supported <<<
   simdjson >>> out of memory <<<
      gjson >>> not supported <<<

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

 - **Parser Test coverage** percent unit test coverage of the parser
   package. It does not include coverage of other package in the
   offering.
