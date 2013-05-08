// This exists because go.crypto/ssh/terminal doesn't build on darwin with go 1.0.
// +build darwin,!go1.1 !linux,!darwin

package main

import "os"

func isTerminal(f *os.File) bool {
	return false
}