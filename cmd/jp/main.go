package main

import (
	"fmt"
	"os"

	flag "github.com/ogier/pflag"
	"github.com/paulhammond/jp"
	"golang.org/x/term"
)

func main() {
	os.Exit(run())
}

func run() int {
	usage := `jp: a JSON reformatter

usage: jp [options] [file]

options:
      --color: force colored output (default autodetects)
      --compact: compact format
  -h  --help: show this help text
`

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	compact := flag.Bool("compact", false, "compact format")
	colors := flag.Bool("color", false, "force colored output")
	help := flag.BoolP("help", "h", false, "show help text")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 || *help {
		flag.Usage()
		if *help {
			return 0
		}
		return 2
	}

	format := "pretty"
	if *compact {
		format = "compact"
	}
	if useColors(*colors) {
		format += "16"
	}

	var fd *os.File
	var e error
	if args[0] == "-" {
		fd = os.Stdin
	} else {
		fd, e = os.Open(args[0])
		if e != nil {
			fmt.Fprintln(os.Stderr, "Error:", e)
			return 1
		}
	}

	e = jp.Expand(fd, os.Stdout, format)
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		return 1
	}
	return 0
}

func useColors(override bool) bool {
	if override {
		return true
	}
	// https://no-color.org
	_, noColor := os.LookupEnv("NO_COLOR")
	if noColor {
		return false
	}
	// https://bixense.com/clicolors/
	_, forceColor := os.LookupEnv("CLICOLOR_FORCE")
	if forceColor {
		return true
	}

	return term.IsTerminal(int(os.Stdout.Fd()))
}
