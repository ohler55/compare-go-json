# Compare Go JSON

Not all JSON tools cover the same features which make it difficult to
select a set of tools for a project. Here is an attempt to compare
feature and benchmarks for a few of the JSON tools for Go.

## Features

| Feature                         | go/json            | fastjson           | jsoniter           | OjG                | simdjson           |
| ------------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| Parse []byte to simple go types | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Validate                        | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Read from io.Reader             | :white_check_mark: |                    |                    | :white_check_mark: | :x:                |
| Read from file                  | :white_check_mark: |                    |                    | :white_check_mark: | :white_check_mark: |
| Multiple JSON file/stream       | :x:                |                    |                    | :white_check_mark: | :x:                |
| ndjson (newline separated)      | :x:                |                    |                    | :white_check_mark: | :white_check_mark: |
| Marshal/Write                   | :white_check_mark: | :x:                | :white_check_mark: | :white_check_mark: | :x:                |
| JSONPath                        | :x:                | :x:                | :x:                | :white_check_mark: | :x:                |

# Benchmarks

```
Parse string/[]byte to simple go types ([]interface{}, int64, string, etc)
       json.Unmarshal       42225 ns/op   17983 B/op     336 allocs/op
         oj.Parse           18500 ns/op    5984 B/op     366 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Unmarshal       32111 ns/op   19797 B/op     451 allocs/op
   simdjson.Parse           45285 ns/op  136899 B/op     370 allocs/op

         oj ██████████████████████▊ 2.28
   jsoniter █████████████▏ 1.31
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson █████████▎ 0.93
   fastjson >>> not supported <<<

Validate string/[]byte
       json.Valid           12533 ns/op       0 B/op       0 allocs/op
         oj.Validate         4292 ns/op       0 B/op       0 allocs/op
   fastjson.Validate         4705 ns/op       0 B/op       0 allocs/op
   jsoniter.Valid            9860 ns/op    2192 B/op     100 allocs/op
   simdjson.Validate        24649 ns/op  114241 B/op      18 allocs/op

         oj █████████████████████████████▏ 2.92
   fastjson ██████████████████████████▋ 2.66
   jsoniter ████████████▋ 1.27
       json ▓▓▓▓▓▓▓▓▓▓ 1.00
   simdjson █████  0.51

Marshal to string/[]byte
       json.Marshal         71818 ns/op   27334 B/op     352 allocs/op
         oj.JSON            13205 ns/op    4096 B/op       1 allocs/op
   fastjson                 >>> not supported <<<
   jsoniter.Marshal         17435 ns/op    7290 B/op      94 allocs/op
   simdjson                 >>> not supported <<<

         oj ██████████████████████████████████████████████████████▍ 5.44
   jsoniter █████████████████████████████████████████▏ 4.12
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
 Machine:         MacBookPro15,2
 OS:              Mac OS X 10.15.5
 Processor:       Quad-Core Intel Core i7
 Cores:           4
 Processor Speed: 2.8 GHz
```
