package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestRun(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"unescape": cmdUnescape,
		},
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"jp": run,
	}))
}

func cmdUnescape(ts *testscript.TestScript, neg bool, args []string) {
	if neg {
		ts.Fatalf("unsupported: ! unescape")
	}
	for _, arg := range args {
		file := ts.MkAbs(arg)
		data, err := os.ReadFile(file)
		ts.Check(err)
		data = bytes.ReplaceAll(data, []byte(`\033`), []byte{033})
		err = os.WriteFile(file, data, 0666)
		ts.Check(err)
	}
}
