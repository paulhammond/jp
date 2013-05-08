// +build linux darwin,go1.1

package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"os"
)

func isTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
