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
}`

var pretty16 = strings.ReplaceAll(`\033[0;32m{
  \033[0m"\033[1mfoo\033[0m"\033[0;32m: \033[0m"\033[1mIñtërnâtiônàlizætiøn\033[0m"\033[0;32m,
  \033[0m"\033[1mempty\033[0m"\033[0;32m: \033[0;32m{ }\033[0;32m,
  \033[0m"\033[1msub\033[0m"\033[0;32m: \033[0;32m{
    \033[0m"\033[1mñ\033[0m"\033[0;32m: \033[0m"\033[1m\u00F1\033[0m"\033[0;32m,
    \033[0m"\033[1mn˜\033[0m"\033[0;32m: \033[0m"\033[1mn\u0303\033[0m"\033[0;32m,
    \033[0m"\033[1marray\033[0m"\033[0;32m: \033[0;32m[
      \033[1;33m1\033[0;32m,
      \033[1;33m2\033[0;32m,
      \033[1;33m3\033[0;32m
    ]\033[0;32m,
    \033[0m"\033[1marray\033[0m"\033[0;32m: \033[0;32m[ ]\033[0;32m
  }\033[0;32m
}\033[0m`, `\033`, "\033")

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
		{compact, pretty16, "pretty16"},
		{extraspaces, pretty, "pretty"},
		{pretty, pretty, "pretty"},
		{compact, compact, "compact"},
		// this checks for an edge case in strings
		{`{"slash\\" : "foo" }`, `{"slash\\":"foo"}`, "compact"},
		{`{"" : "foo" }`, `{"":"foo"}`, "compact"},
		// check json inside string isn't expanded
		{`{"foo":"{\"a\":\"b\",\"b\":{},\"c\":{\"a\":\"b\",\"b\":[1,2],\"c\":[]}}"}`, `{
  "foo": "{\"a\":\"b\",\"b\":{},\"c\":{\"a\":\"b\",\"b\":[1,2],\"c\":[]}}"
}`, "pretty"},
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
		want := test.out + "\n"
		if w.String() != want {
			t.Errorf("unexpected JSON\ngot  %q\nwant %q", escape(w.String()), escape(want))
		}
	}
}

func escape(str string) string {
	return strings.ReplaceAll(str, "\033", "\\033")
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
