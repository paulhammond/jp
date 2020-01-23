// Copyright 2013-2014 Paul Hammond.
// This software is licensed under the MIT license, see LICENSE.txt for details.

package main

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	flag "github.com/ogier/pflag"
	"github.com/paulhammond/jp"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: jp [file]\n")
		flag.PrintDefaults()
	}

	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	compact := flag.Bool("compact", false, "compact format")
	colors := flag.Bool("color", isTerminal, "colored format")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	format := "pretty"
	if *compact {
		format = "compact"
	}
	if *colors {
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
			os.Exit(1)
		}
	}

	e = jp.Expand(fd, os.Stdout, format)
	if e != nil {
		fmt.Fprintln(os.Stderr, "Error:", e)
		os.Exit(1)
	}
}
