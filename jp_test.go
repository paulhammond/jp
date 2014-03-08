// Copyright 2013 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package jp

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

var compact = `{"foo":"Iñtërnâtiônàlizætiøn","empty":{},"sub":{"ñ":"\u00F1","n˜":"n\u0303","array":[1,2,3],"array":[]}}`
var pretty = `{
  "foo": "Iñtërnâtiônàlizætiøn",
  "empty": { },
  "sub": {
    "ñ": "\u00F1",
    "n˜": "n\u0303",
    "array": [
      1,
      2,
      3
    ],
    "array": [ ]
  }
}
`
var extraspaces = `{
	"foo":  "Iñtërnâtiônàlizætiøn"  ,     "empty"    : {  }
,"sub":  {"ñ"  :  "\u00F1"  ,"n˜":"n\u0303","array" : [  1,   2, 3 ]  
,"array":[
]  }
  }  
`

func TestExpand(t *testing.T) {
	tests := []struct {
		in     string
		out    string
		format string
	}{
		{compact, pretty, "pretty"},
		{extraspaces, pretty, "pretty"},
		{pretty, pretty, "pretty"},
		{compact, compact, "compact"},
		// this checks for an edge case in strings
		{`{"slash\\" : "foo" }`, `{"slash\\":"foo"}`, "compact"},
		{`{"" : "foo" }`, `{"":"foo"}`, "compact"},
		// check invalid spaces between elements are preserved
		{`{"foo":[1 2, 3]}`, `{"foo":[1 2,3]}`, "compact"},
		{`{"foo":[1 true, 3]}`, `{"foo":[1 true,3]}`, "compact"},
		{`{"foo":[true false, 3]}`, `{"foo":[true false,3]}`, "compact"},
		{`{"foo":[true false]}`, `{"foo":[true false]}`, "compact"},
		{`{"foo":[true {}]}`, `{"foo":[true{}]}`, "compact"},
	}
	for _, test := range tests {
		r := strings.NewReader(test.in)
		w := &bytes.Buffer{}
		err := Expand(r, w, test.format)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if w.String() != test.out {
			t.Errorf("unexpected JSON, got\n%s\nexpected\n%s", w.String(), test.out)
		}
	}
}

// This benchmark isn't run by default. To run it, create "bench.json", then:
// go test -c .
// ./jp.test -test.bench="Expand" -test.cpuprofile cpu.out -test.benchtime=5 2> tmp.out
// go tool pprof jp.test cpu.out
func BenchmarkExpand(b *testing.B) {
	b.StopTimer()
	var r, _ = os.Open("./bench.json")
	var w = os.Stderr
	b.StartTimer()
	var err error
	for i := 0; i < b.N; i++ {
		err = Expand(r, w, "pretty")
	}
	b.StopTimer()
	fmt.Println(err)
	b.StartTimer()
}
