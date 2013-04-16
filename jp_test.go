package jp

import (
	"bytes"
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
		{extraspaces, compact, "compact"},
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
