package jp

import (
	"bytes"
	"strings"
	"testing"
)

var compact = `{"foo":"bar","empty":{},"sub":{"a":"b","c":"d","array":[1,2,3],"array":[]}}`
var pretty = `{
  "foo": "bar",
  "empty": { },
  "sub": {
    "a": "b",
    "c": "d",
    "array": [
      1,
      2,
      3
    ],
    "array": [ ]
  }
}`
var extraspaces = `{
	"foo":  "bar"  ,     "empty"    : {  }
,"sub":  {"a"  :  "b"  ,"c":"d","array" : [  1,   2, 3 ]  
,"array":[
]  }
  }  
`

func TestExpand(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{compact, pretty},
		{extraspaces, pretty},
	}
	for _, test := range tests {
		r := strings.NewReader(test.in)
		w := &bytes.Buffer{}
		err := Expand(r, w)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if w.String() != test.out {
			t.Errorf("unexpected JSON, got\n%s\nexpected\n%s", w.String(), test.out)
		}
	}
}
